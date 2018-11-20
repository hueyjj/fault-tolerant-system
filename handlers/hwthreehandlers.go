package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"bitbucket.org/cmps128gofour/homework3/response"
	"bitbucket.org/cmps128gofour/homework3/vectorclock"
	"github.com/gorilla/mux"
)

var (
	vectorClocks = make(map[string]vectorclock.Unit)
)

func subjectPUT(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse the key from url variable and (store) value from the request
	vars := mux.Vars(r)
	key := vars["subject"]
	value := r.PostFormValue("val")
	payload := r.PostFormValue("payload")

	msg := new(response.Response)
	if err := json.Unmarshal([]byte(payload), &msg); err != nil {
		log.Printf("Unable to unmarshal payload: %v\n", err)
		http.Error(w, "Unable to unmarshal payload", http.StatusInternalServerError)
		return
	}

	var resp *response.Response
	replaced := new(bool)
	// Return error message if key is 1 and 200 characters
	if len(key) < 1 || len(key) > 200 {
		resp = &response.Response{
			Msg:    "Key not valid",
			Result: "Error",
		}
		w.WriteHeader(http.StatusBadRequest)
	} else if len(value) > 1e6 {
		// Return error message if value is greater than 1 MB
		resp = &response.Response{
			Msg:    "Object too large. Size limit is 1MB",
			Result: "Error",
		}
		w.WriteHeader(http.StatusBadRequest)
	} else if KVStore.Exists(key) {
		// Replace value in store
		if isGreaterEqual(key, msg.Payload.VectorClocks, vectorClocks) {
			KVStore.Put(key, value)
			vectorClocks[key] = mergeClock(key, vectorClocks, msg.Payload.VectorClocks)
			*replaced = true
			resp = &response.Response{
				Replaced: replaced,
				Msg:      "Updated successfully",
				Payload: response.Payload{
					VectorClocks: vectorClocks,
				},
			}
			w.WriteHeader(http.StatusOK)
		} else {
			resp = &response.Response{
				Result: "Error",
				Msg:    "Payload out of date",
				Payload: response.Payload{
					VectorClocks: vectorClocks,
				},
			}
			w.WriteHeader(http.StatusBadRequest)
		}
	} else {
		// Put key into store if it doesn't exists
		if isGreaterEqual(key, msg.Payload.VectorClocks, vectorClocks) {
			KVStore.Put(key, value)
			vectorClocks[key] = mergeClock(key, vectorClocks, msg.Payload.VectorClocks)
			*replaced = false
			resp = &response.Response{
				Replaced: replaced,
				Msg:      "Added successfully",
				Payload: response.Payload{
					VectorClocks: vectorClocks,
				},
			}
			w.WriteHeader(http.StatusCreated)
		} else {
			resp = &response.Response{
				Result: "Error",
				Msg:    "Payload out of date",
				Payload: response.Payload{
					VectorClocks: vectorClocks,
				},
			}
			w.WriteHeader(http.StatusBadRequest)
		}
	}

	// Convert response into json structure and then into bytes
	data, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Unable to marshal response: %v\n", err)
		http.Error(w, "Unable to marshal response", http.StatusInternalServerError)
		return
	}

	// Send response
	w.Write(data)
}

func subjectGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse the key from url variable and (store) value from the request
	vars := mux.Vars(r)
	key := vars["subject"]
	payload := r.PostFormValue("payload")

	msg := new(response.Response)
	if err := json.Unmarshal([]byte(payload), &msg); err != nil {
		log.Printf("Unable to unmarshal payload: %v\n", err)
		http.Error(w, "Unable to unmarshal payload", http.StatusInternalServerError)
		return
	}

	var resp *response.Response
	if KVStore.Exists(key) {
		if isGreaterEqual(key, vectorClocks, msg.Payload.VectorClocks) {
			vectorClocks[key] = mergeClock(key, vectorClocks, msg.Payload.VectorClocks)
			value, _ := KVStore.Get(key)
			resp = &response.Response{
				Result: "Success",
				Value:  value,
				Payload: response.Payload{
					VectorClocks: vectorClocks,
				},
			}
			w.WriteHeader(http.StatusOK)
		} else {
			resp = &response.Response{
				Result: "Error",
				Msg:    "Payload out of date",
				Payload: response.Payload{
					VectorClocks: vectorClocks,
				},
			}
			w.WriteHeader(http.StatusBadRequest)
		}
	} else {
		if isGreaterEqual(key, vectorClocks, msg.Payload.VectorClocks) {
			vectorClocks[key] = mergeClock(key, vectorClocks, msg.Payload.VectorClocks)
			resp = &response.Response{
				Result: "Error",
				Msg:    "Key does not exist",
				Payload: response.Payload{
					VectorClocks: vectorClocks,
				},
			}
			w.WriteHeader(http.StatusNotFound)
		} else {
			resp = &response.Response{
				Result: "Error",
				Msg:    "Payload out of date",
				Payload: response.Payload{
					VectorClocks: vectorClocks,
				},
			}
			w.WriteHeader(http.StatusBadRequest)
		}
	}

	// Convert response into json structure and then into bytes
	data, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Unable to marshal response: %v\n", err)
		http.Error(w, "Unable to marshal response", http.StatusInternalServerError)
		return
	}

	// Send response
	w.Write(data)
}

func subjectSEARCH(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Parse the key from url variable and (store) value from the request
	vars := mux.Vars(r)
	key := vars["subject"]

	payload := r.PostFormValue("payload")

	msg := new(response.Response)
	if err := json.Unmarshal([]byte(payload), &msg); err != nil {
		log.Printf("Unable to unmarshal payload: %v\n", err)
		http.Error(w, "Unable to unmarshal payload", http.StatusInternalServerError)
		return
	}

	var resp *response.Response
	isExists := new(bool)
	if KVStore.Exists(key) {
		*isExists = true
		if isGreaterEqual(key, vectorClocks, msg.Payload.VectorClocks) {
			vectorClocks[key] = mergeClock(key, vectorClocks, msg.Payload.VectorClocks)
			resp = &response.Response{
				Result:   "Success",
				IsExists: isExists,
				Payload: response.Payload{
					VectorClocks: vectorClocks,
				},
			}
			w.WriteHeader(http.StatusOK)
		} else {
			resp = &response.Response{
				Result: "Error",
				Msg:    "Payload out of date",
				Payload: response.Payload{
					VectorClocks: vectorClocks,
				},
			}
			w.WriteHeader(http.StatusBadRequest)
		}
	} else {
		if !clockExists(key) {
			*isExists = false
			newClock(key) // creates new clock
			resp = &response.Response{
				Result:   "Success",
				IsExists: isExists,
				Payload: response.Payload{
					VectorClocks: vectorClocks,
				},
			}
			w.WriteHeader(http.StatusOK)
		} else if len(msg.Payload.VectorClocks) <= 0 {
			resp = &response.Response{
				Result: "Error",
				Msg:    "Payload out of date",
				Payload: response.Payload{
					VectorClocks: vectorClocks,
				},
			}
			w.WriteHeader(http.StatusBadRequest)
		} else if isGreaterEqual(key, vectorClocks, msg.Payload.VectorClocks) {
			vectorClocks[key] = mergeClock(key, vectorClocks, msg.Payload.VectorClocks)
			resp = &response.Response{
				Result:   "Success",
				IsExists: isExists,
				Payload: response.Payload{
					VectorClocks: vectorClocks,
				},
			}
			w.WriteHeader(http.StatusOK)
		} else {
			resp = &response.Response{
				Result: "Error",
				Msg:    "Payload out of date",
				Payload: response.Payload{
					VectorClocks: vectorClocks,
				},
			}
			w.WriteHeader(http.StatusBadRequest)
		}
	}

	// Convert response into json structure and then into bytes
	data, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Unable to marshal response: %v\n", err)
		http.Error(w, "Unable to marshal response", http.StatusInternalServerError)
		return
	}

	// Send response
	w.Write(data)
}

func subjectDEL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse the key from url variable and (store) value from the request
	vars := mux.Vars(r)
	key := vars["subject"]

	payload := r.PostFormValue("payload")

	msg := new(response.Response)
	if err := json.Unmarshal([]byte(payload), &msg); err != nil {
		log.Printf("Unable to unmarshal payload: %v\n", err)
		http.Error(w, "Unable to unmarshal payload", http.StatusInternalServerError)
		return
	}

	var resp *response.Response
	if KVStore.Exists(key) {
		if isGreaterEqual(key, msg.Payload.VectorClocks, vectorClocks) {
			err := KVStore.Delete(key)
			if err != nil {
				log.Printf("Unable to delete key: %v\n", err)
				http.Error(w, "Unable to delete key", http.StatusBadRequest)
			}
			vectorClocks[key] = mergeClock(key, vectorClocks, msg.Payload.VectorClocks)
			resp = &response.Response{
				Result: "Success",
				Msg:    "Key deleted",
				Payload: response.Payload{
					VectorClocks: vectorClocks,
				},
			}
			w.WriteHeader(http.StatusOK)
		} else {
			resp = &response.Response{
				Result: "Error",
				Msg:    "Payload out of date",
				Payload: response.Payload{
					VectorClocks: vectorClocks,
				},
			}
			w.WriteHeader(http.StatusBadRequest)
		}
	} else {
		if isGreaterEqual(key, msg.Payload.VectorClocks, vectorClocks) {
			vectorClocks[key] = mergeClock(key, vectorClocks, msg.Payload.VectorClocks)
			resp = &response.Response{
				Result: "Error",
				Msg:    "Key does not exist",
				Payload: response.Payload{
					VectorClocks: vectorClocks,
				},
			}
			w.WriteHeader(http.StatusNotFound)
		} else {
			resp = &response.Response{
				Result: "Error",
				Msg:    "Payload out of date",
				Payload: response.Payload{
					VectorClocks: vectorClocks,
				},
			}
			w.WriteHeader(http.StatusBadRequest)
		}
	}

	// Convert response into json structure and then into bytes
	data, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Unable to marshal response: %v\n", err)
		http.Error(w, "Unable to marshal response", http.StatusInternalServerError)
		return
	}

	// Send response
	w.Write(data)
}

var views []string
var myIP string

func SetViews(v string) {
	views = strings.Split(v, ",")
}

func SetMyIpPort(ip string) {
	myIP = ip
}

func viewGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	iplist := strings.Join(views, ",")
	var resp *response.IPTableResponse
	resp = &response.IPTableResponse{
		View: iplist,
	}
	w.WriteHeader(http.StatusOK)

	// Convert response into json structure and then into bytes
	data, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Unable to marshal response: %v\n", err)
		http.Error(w, "Unable to marshal response", http.StatusInternalServerError)
		return
	}

	// Send response
	w.Write(data)
}

func viewPUT(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse the key from url variable and (store) value from the request
	ipport := r.PostFormValue("ip_port")
	isIpportExist := false
	for _, ip := range views {
		if ip == ipport {
			isIpportExist = true
		}
	}

	var resp *response.ViewResponse
	if isIpportExist {
		resp = &response.ViewResponse{
			Result: "Error",
			Msg:    fmt.Sprintf("%s is already in view", ipport),
		}
		w.WriteHeader(http.StatusNotFound)
	} else {
		views = append(views, ipport)
		resp = &response.ViewResponse{
			Result: "Success",
			Msg:    fmt.Sprintf("Successfully added %s to view", ipport),
		}
		w.WriteHeader(http.StatusOK)
	}

	iptableValue := r.PostFormValue("iptable")
	iptable := make(map[string]bool)
	if iptableValue != "" {
		if err := json.Unmarshal([]byte(iptableValue), &iptable); err != nil {
			log.Printf("Unable to unmarshal iptable: %v\n", err)
			http.Error(w, "Unable to unmarshal iptable", http.StatusInternalServerError)
			return
		}
	}

	if len(iptable) <= 0 {
		// Start a new gossip
		for _, view := range views {
			iptable[view] = false
		}
		iptable[myIP] = true
		nextNodeURL, err := findNextNode(iptable)
		if err == nil {
			gossipViewPUT(nextNodeURL, ipport, iptable)
		}
	} else {
		// Gossip if there's an ip that hasn't seen the message
		iptable[myIP] = true
		nextNodeURL, err := findNextNode(iptable)
		if err == nil {
			gossipViewPUT(nextNodeURL, ipport, iptable)
		}
	}

	log.Printf("viewPUT: views: %+v iptable: %+v\n", views, iptable)

	// Convert response into json structure and then into bytes
	data, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Unable to marshal response: %v\n", err)
		http.Error(w, "Unable to marshal response", http.StatusInternalServerError)
		return
	}

	// Send response
	w.Write(data)
}

func viewDELETE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	// Parse the key from url variable and (store) value from the request
	//ipport := r.PostFormValue("ip_port") // This doesn't work
	ipport := strings.Split(string(body), "=")[1]
	target := -1
	for index, ip := range views {
		if ip == ipport {
			target = index
		}
	}

	var resp *response.ViewResponse
	if target != -1 {
		views = append(views[:target], views[target+1:]...)
		resp = &response.ViewResponse{
			Result: "Success",
			Msg:    fmt.Sprintf("Successfully removed %s from view", ipport),
		}
		w.WriteHeader(http.StatusOK)
	} else {
		resp = &response.ViewResponse{
			Result: "Error",
			Msg:    fmt.Sprintf("%s is not in current view", ipport),
		}
		w.WriteHeader(http.StatusNotFound)
	}

	// Convert response into json structure and then into bytes
	data, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Unable to marshal response: %v\n", err)
		http.Error(w, "Unable to marshal response", http.StatusInternalServerError)
		return
	}

	// Send response
	w.Write(data)
}

func unixNow() int64 {
	return time.Now().Unix()
}

func getBody(body io.ReadCloser) *response.Response {
	resp := new(response.Response)
	if err := json.NewDecoder(body).Decode(resp); err != nil {
		return nil
	}
	return resp
}

//func incClock(key string) {
//	if _, ok := vectorClocks[key]; ok {
//		vectorClocks[key] = vectorclock.Unit{
//			Tick:      vectorClocks[key].Tick + 1,
//			Timestamp: vectorClocks[key].Timestamp,
//		}
//	} else {
//		vectorClocks[key] = vectorclock.Unit{
//			Tick:      1,
//			Timestamp: unixNow(),
//		}
//	}
//}

func newClock(key string) {
	vectorClocks[key] = vectorclock.Unit{
		Tick:      1,
		Timestamp: unixNow(),
	}
}

func clockExists(key string) bool {
	if _, ok := vectorClocks[key]; ok {
		return true
	}
	return false
}

func isGreaterEqual(key string, v1, v2 map[string]vectorclock.Unit) bool {
	log.Printf("isGreaterEqual: v1[%s].Tick=%d v2[%s].Tick=%d v1[%s].Timestamp=%d v2[%s].Timestamp=%d",
		key, v1[key].Tick, key, v2[key].Tick, key, v1[key].Timestamp, key, v2[key].Timestamp)

	v1Val := v1[key].Tick
	v2Val := v2[key].Tick
	if v1Val > v2Val {
		return true
	} else if v1Val < v2Val {
		return false
	} else {
		v1Time := v1[key].Timestamp
		v2Time := v2[key].Timestamp
		if v1Time > v2Time || v1Time == v2Time {
			return true
		}
		return false
	}
}

func mergeClock(key string, v1, v2 map[string]vectorclock.Unit) vectorclock.Unit {
	v1Val := v1[key].Tick
	v2Val := v2[key].Tick
	v1Time := v1[key].Timestamp
	v2Time := v2[key].Timestamp
	var tick int
	var timestamp int64
	if v1Val > v2Val {
		tick = v1Val + 1
		timestamp = v1Time
	} else if v1Val < v2Val {
		tick = v2Val + 1
		timestamp = v2Time
	} else {
		if v1Time > v2Time {
			tick = v1Val + 1
			timestamp = v1Time
		} else {
			tick = v2Val + 1
			timestamp = v2Time
		}
	}
	log.Printf("mergeClock: v1.Tick=%d v2.Tick=%d tick=%d time=%d\n", v1Val, v2Val, tick, timestamp)
	return vectorclock.Unit{
		Tick:      tick,
		Timestamp: timestamp,
	}
}
