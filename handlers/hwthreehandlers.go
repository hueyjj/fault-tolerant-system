package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"bitbucket.org/cmps128gofour/homework4/shard"

	"bitbucket.org/cmps128gofour/homework4/response"
	"bitbucket.org/cmps128gofour/homework4/vectorclock"
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
	data := r.PostFormValue("payload")

	payload := new(response.Payload)
	payload.VectorClocks = make(map[string]vectorclock.Unit)
	payload.IPTable = make(map[string]int)
	if data != "" {
		log.Printf("subjectPUT: data=%+v\n", data)
		if err := json.Unmarshal([]byte(data), &payload); err != nil {
			log.Printf("subjectPUT: Unable to unmarshal payload: %v\n", err)
			log.Printf("subjectPUT: payload=%+v\n", payload)
			log.Printf("subjectPUT: r.Body=%+v\n", r.Body)
		}
	} else {
		log.Printf("subjectPUT: payload is empty")
	}

	//msg := new(response.Response)
	//if err := json.Unmarshal([]byte(payload), &msg); err != nil {
	//	log.Printf("subjectPUT: Unable to unmarshal payload: %v\n", err)
	//	log.Printf("subjectPUT: payload=%+v\n", payload)
	//	log.Printf("subjectPUT: r.Body=%+v\n", r.Body)
	//	//http.Error(w, "subjectPUT: Unable to unmarshal payload", http.StatusInternalServerError)
	//	//return
	//	if payload := r.PostFormValue("vectorclocks"); payload != "" {
	//		msg = getBody(r.Body)
	//		log.Printf("subjectPUT: Attempt to parse entire body instead: %+v\n", msg)
	//	}
	//}

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
		if isGreaterEqual(key, payload.VectorClocks, vectorClocks) {
			KVStore.Put(key, value)
			vectorClocks[key] = mergeClock(key, vectorClocks, payload.VectorClocks)
			*replaced = true
			resp = &response.Response{
				Replaced: replaced,
				Msg:      "Updated successfully",
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
	} else {
		// Put key into store if it doesn't exists
		if isGreaterEqual(key, payload.VectorClocks, vectorClocks) {
			KVStore.Put(key, value)
			vectorClocks[key] = mergeClock(key, vectorClocks, payload.VectorClocks)
			*replaced = false
			resp = &response.Response{
				Replaced: replaced,
				Msg:      "Added successfully",
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

	//iptableValue := r.PostFormValue("iptable")
	//iptable := make(map[string]int)
	//if iptableValue != "" {
	//	if err := json.Unmarshal([]byte(iptableValue), &iptable); err != nil {
	//		log.Printf("Unable to unmarshal iptable: %v\n", err)
	//		//http.Error(w, "Unable to unmarshal iptable", http.StatusInternalServerError)
	//		//return
	//	}
	//}

	log.Printf("subjectPUT: len(iptable)=%d\n", len(payload.IPTable))
	if len(payload.IPTable) <= 0 {
		// Start a new gossip
		for _, view := range views {
			payload.IPTable[view] = 0
		}
		payload.IPTable[myIP] = 1

		nextNodeURL, err := findNextNode(payload.IPTable)
		if err == nil {
			log.Printf("subjectPUT: nextNodeURL=%s", nextNodeURL)
			gossipSubjectPUT(nextNodeURL, key, value, payload)
		}
	} else {
		// Gossip if there's an ip that hasn't seen the message
		payload.IPTable[myIP] = 1
		nextNodeURL, err := findNextNode(payload.IPTable)
		if err == nil {
			log.Printf("subjectPUT: nextNodeURL=%s", nextNodeURL)
			gossipSubjectPUT(nextNodeURL, key, value, payload)
		}
	}

	log.Printf("subjectPUT: views: %+v iptable: %+v\n", views, payload.IPTable)
	log.Printf("subjectPUT: vectorClocks=%+v\n", vectorClocks)
	log.Printf("subjectPUT: store=%+v\n", KVStore)
	// Convert response into json structure and then into bytes
	dataResp, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Unable to marshal response: %v\n", err)
		//http.Error(w, "Unable to marshal response", http.StatusInternalServerError)
		//return
	}

	// Send response
	w.Write(dataResp)
}

func subjectGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	log.Printf("subjectGET: GET request received\n")

	// Parse the key from url variable and (store) value from the request
	vars := mux.Vars(r)
	key := vars["subject"]
	//payload := r.PostFormValue("payload")
	//if payload == "" {
	//	r.ParseForm()
	//	payload = r.Form.Get("payload")
	//	log.Printf("subjectGET: payload=%s\n", payload)
	//}

	data := r.PostFormValue("payload")

	payload := new(response.Payload)
	payload.VectorClocks = make(map[string]vectorclock.Unit)
	payload.IPTable = make(map[string]int)
	if data != "" {
		log.Printf("subjectPUT: data=%+v\n", data)
		if err := json.Unmarshal([]byte(data), &payload); err != nil {
			log.Printf("subjectPUT: Unable to unmarshal payload: %v\n", err)
			log.Printf("subjectPUT: payload=%+v\n", payload)
			log.Printf("subjectPUT: r.Body=%+v\n", r.Body)
		}
	} else {
		log.Printf("subjectPUT: payload is empty")
	}

	//msg := new(response.Response)
	//if err := json.Unmarshal([]byte(payload), &msg); err != nil {
	//	log.Printf("subjectGET: Unable to unmarshal payload: %v\n", err)
	//	log.Printf("subjectGET: payload=%+v\n", payload)
	//	log.Printf("subjectGET: r.Body=%+v\n", r.Body)
	//	//http.Error(w, "subjectGET: Unable to unmarshal payload", http.StatusInternalServerError)
	//	//return
	//}

	var resp *response.Response
	if KVStore.Exists(key) {
		if isGreaterEqual(key, vectorClocks, payload.VectorClocks) {
			vectorClocks[key] = mergeClock(key, vectorClocks, payload.VectorClocks)
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
		if isGreaterEqual(key, vectorClocks, payload.VectorClocks) {
			vectorClocks[key] = mergeClock(key, vectorClocks, payload.VectorClocks)
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

	//iptableValue := r.PostFormValue("iptable")
	//iptable := make(map[string]int)
	//if iptableValue != "" {
	//	if err := json.Unmarshal([]byte(iptableValue), &iptable); err != nil {
	//		log.Printf("Unable to unmarshal iptable: %v\n", err)
	//		//http.Error(w, "Unable to unmarshal iptable", http.StatusInternalServerError)
	//		//return
	//	}
	//}

	log.Printf("subjectGET: len(iptable)=%d\n", len(payload.IPTable))
	if len(payload.IPTable) <= 0 {
		// Start a new gossip
		for _, view := range views {
			payload.IPTable[view] = 0
		}
		payload.IPTable[myIP] = 1
		nextNodeURL, err := findNextNode(payload.IPTable)
		if err == nil {
			log.Printf("subjectGET: nextNodeURL=%s", nextNodeURL)
			gossipSubjectGET(nextNodeURL, key, payload)
		}
	} else {
		// Gossip if there's an ip that hasn't seen the message
		payload.IPTable[myIP] = 1
		nextNodeURL, err := findNextNode(payload.IPTable)
		if err == nil {
			log.Printf("subjectGET: nextNodeURL=%s", nextNodeURL)
			gossipSubjectGET(nextNodeURL, key, payload)
		}
	}

	log.Printf("subjectGET: views: %+v iptable: %+v\n", views, payload.IPTable)
	log.Printf("subjectGET: vectorClocks=%+v\n", vectorClocks)
	log.Printf("subjectGET: store=%+v\n", KVStore)

	// Convert response into json structure and then into bytes
	dataResp, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Unable to marshal response: %v\n", err)
		//http.Error(w, "Unable to marshal response", http.StatusInternalServerError)
		//return
	}

	// Send response
	w.Write(dataResp)
}

func subjectSEARCH(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Parse the key from url variable and (store) value from the request
	vars := mux.Vars(r)
	key := vars["subject"]

	//payload := r.PostFormValue("payload")

	data := r.PostFormValue("payload")

	payload := new(response.Payload)
	payload.VectorClocks = make(map[string]vectorclock.Unit)
	payload.IPTable = make(map[string]int)
	if data != "" {
		log.Printf("subjectPUT: data=%+v\n", data)
		if err := json.Unmarshal([]byte(data), &payload); err != nil {
			log.Printf("subjectPUT: Unable to unmarshal payload: %v\n", err)
			log.Printf("subjectPUT: payload=%+v\n", payload)
			log.Printf("subjectPUT: r.Body=%+v\n", r.Body)
		}
	} else {
		log.Printf("subjectPUT: payload is empty")
	}

	//msg := new(response.Response)
	//if err := json.Unmarshal([]byte(payload), &msg); err != nil {
	//	log.Printf("subjectSEARCH: Unable to unmarshal payload: %v\n", err)
	//	log.Printf("subjectSEARCH: payload=%+v\n", payload)
	//	//http.Error(w, "subjectSEARCH: Unable to unmarshal payload", http.StatusInternalServerError)
	//	//return
	//}

	var resp *response.Response
	isExists := new(bool)
	if KVStore.Exists(key) {
		*isExists = true
		if isGreaterEqual(key, vectorClocks, payload.VectorClocks) {
			vectorClocks[key] = mergeClock(key, vectorClocks, payload.VectorClocks)
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
		} else if len(payload.VectorClocks) <= 0 {
			resp = &response.Response{
				Result: "Error",
				Msg:    "Payload out of date",
				Payload: response.Payload{
					VectorClocks: vectorClocks,
				},
			}
			w.WriteHeader(http.StatusBadRequest)
		} else if isGreaterEqual(key, vectorClocks, payload.VectorClocks) {
			vectorClocks[key] = mergeClock(key, vectorClocks, payload.VectorClocks)
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

	//iptableValue := r.PostFormValue("iptable")
	//iptable := make(map[string]int)
	//if iptableValue != "" {
	//	if err := json.Unmarshal([]byte(iptableValue), &iptable); err != nil {
	//		log.Printf("Unable to unmarshal iptable: %v\n", err)
	//		//http.Error(w, "Unable to unmarshal iptable", http.StatusInternalServerError)
	//		//return
	//	}
	//}

	log.Printf("subjectSEARCH: len(iptable)=%d\n", len(payload.IPTable))
	if len(payload.IPTable) <= 0 {
		// Start a new gossip
		for _, view := range views {
			payload.IPTable[view] = 0
		}
		payload.IPTable[myIP] = 1
		nextNodeURL, err := findNextNode(payload.IPTable)
		if err == nil {
			log.Printf("subjectSEARCH: nextNodeURL=%s", nextNodeURL)
			gossipSubjectSEARCH(nextNodeURL, key, payload)
		}
	} else {
		// Gossip if there's an ip that hasn't seen the message
		payload.IPTable[myIP] = 1
		nextNodeURL, err := findNextNode(payload.IPTable)
		if err == nil {
			log.Printf("subjectSEARCH: nextNodeURL=%s", nextNodeURL)
			gossipSubjectSEARCH(nextNodeURL, key, payload)
		}
	}

	log.Printf("subjectSEARCH: views: %+v iptable: %+v\n", views, payload.IPTable)
	log.Printf("subjectSEARCH: vectorClocks=%+v\n", vectorClocks)
	log.Printf("subjectSEARCH: store=%+v\n", KVStore)

	// Convert response into json structure and then into bytes
	dataResp, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Unable to marshal response: %v\n", err)
		//http.Error(w, "Unable to marshal response", http.StatusInternalServerError)
		//return
	}

	// Send response
	w.Write(dataResp)
}

func subjectDEL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse the key from url variable and (store) value from the request
	vars := mux.Vars(r)
	key := vars["subject"]

	r.Header.Add("content-type", "application/x-www-form-urlencoded")
	data := r.PostFormValue("payload")

	payload := new(response.Payload)
	payload.VectorClocks = make(map[string]vectorclock.Unit)
	payload.IPTable = make(map[string]int)
	if data != "" {
		log.Printf("subjectPUT: data=%+v\n", data)
		if err := json.Unmarshal([]byte(data), &payload); err != nil {
			log.Printf("subjectPUT: Unable to unmarshal payload: %v\n", err)
			log.Printf("subjectPUT: payload=%+v\n", payload)
			log.Printf("subjectPUT: r.Body=%+v\n", r.Body)
		}
	} else {
		log.Printf("subjectPUT: payload is empty")
	}

	//msg := new(response.Response)
	//if err := json.Unmarshal([]byte(payload), &msg); err != nil {
	//	log.Printf("subjectDEL: Unable to unmarshal payload: %v\n", err)
	//	log.Printf("subjectDEL: payload=%+v\n", payload)
	//	//http.Error(w, "subjectDEL: Unable to unmarshal payload", http.StatusInternalServerError)
	//	//return
	//}

	var resp *response.Response
	if KVStore.Exists(key) {
		if isGreaterEqual(key, payload.VectorClocks, vectorClocks) {
			err := KVStore.Delete(key)
			if err != nil {
				log.Printf("Unable to delete key: %v\n", err)
				http.Error(w, "Unable to delete key", http.StatusBadRequest)
			}
			vectorClocks[key] = mergeClock(key, vectorClocks, payload.VectorClocks)
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
		if isGreaterEqual(key, payload.VectorClocks, vectorClocks) {
			vectorClocks[key] = mergeClock(key, vectorClocks, payload.VectorClocks)
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

	//iptableValue := r.PostFormValue("iptable")
	//iptable := make(map[string]int)
	//if iptableValue != "" {
	//	if err := json.Unmarshal([]byte(iptableValue), &iptable); err != nil {
	//		log.Printf("Unable to unmarshal iptable: %v\n", err)
	//		//http.Error(w, "Unable to unmarshal iptable", http.StatusInternalServerError)
	//		//return
	//	}
	//}

	log.Printf("subjectDEL: len(iptable)=%d\n", len(payload.IPTable))
	if len(payload.IPTable) <= 0 {
		// Start a new gossip
		for _, view := range views {
			payload.IPTable[view] = 0
		}
		payload.IPTable[myIP] = 1
		nextNodeURL, err := findNextNode(payload.IPTable)
		if err == nil {
			log.Printf("subjectDEL: nextNodeURL=%s", nextNodeURL)
			gossipSubjectDEL(nextNodeURL, key, payload)
		}
	} else {
		// Gossip if there's an ip that hasn't seen the message
		payload.IPTable[myIP] = 1
		nextNodeURL, err := findNextNode(payload.IPTable)
		if err == nil {
			log.Printf("subjectDEL: nextNodeURL=%s", nextNodeURL)
			gossipSubjectDEL(nextNodeURL, key, payload)
		}
	}

	log.Printf("subjectDEL: views: %+v iptable: %+v\n", views, payload.IPTable)
	log.Printf("subjectDEL: vectorClocks=%+v\n", vectorClocks)
	log.Printf("subjectDEL: store=%+v\n", KVStore)

	// Convert response into json structure and then into bytes
	dataResp, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Unable to marshal response: %v\n", err)
		//http.Error(w, "Unable to marshal response", http.StatusInternalServerError)
		//return
	}

	// Send response
	w.Write(dataResp)
}

var views []string
var myIP string
var shardMap map[int][]string

func GetShardMap() map[int][]string {
	return shardMap
}

func SetShardMap(s string, v string) {

	myViews := strings.Split(v, ",")
	sInt, _ := strconv.Atoi(s)
	shardMap, _ = shard.Shard(myViews, sInt)
}

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
		//http.Error(w, "Unable to marshal response", http.StatusInternalServerError)
		//return
	}

	// Send response
	w.Write(data)
}

func viewPUT(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse the key from url variable and (store) value from the request
	ipport := r.PostFormValue("ip_port") // This doesn't work
	log.Printf("ipport=%s", ipport)
	if ipport == "" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			//http.Error(w, "can't read body", http.StatusBadRequest)
			//return
		}
		ipport = strings.Split(string(body), "=")[1]
		ipport = strings.Replace(ipport, "%3A", ":", -1)
		log.Printf("ipport=%s\n", ipport)
		//return
	}
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
	iptable := make(map[string]int)
	if iptableValue != "" {
		if err := json.Unmarshal([]byte(iptableValue), &iptable); err != nil {
			log.Printf("Unable to unmarshal iptable: %v\n", err)
			//http.Error(w, "Unable to unmarshal iptable", http.StatusInternalServerError)
			//return
		}
	}

	log.Printf("viewPUT: ipport=%s\n", ipport)
	log.Printf("viewPUT: len(iptable)=%d\n", len(iptable))
	if len(iptable) <= 0 {
		// Start a new gossip
		for _, view := range views {
			iptable[view] = 0
		}
		iptable[myIP] = 1
		iptable[ipport] = 1
		nextNodeURL, err := findNextNode(iptable)
		if err == nil {
			log.Printf("viewPUT: nextNodeURL=%s", nextNodeURL)
			gossipViewPUT(nextNodeURL, ipport, iptable)
		}
		// Update new node to latest KVS and vectorclocks with this node
		updateNode(ipport, vectorClocks, KVStore)
	} else {
		// Gossip if there's an ip that hasn't seen the message
		iptable[myIP] = 1
		nextNodeURL, err := findNextNode(iptable)
		if err == nil {
			log.Printf("viewPUT: nextNodeURL=%s", nextNodeURL)
			gossipViewPUT(nextNodeURL, ipport, iptable)
		}
	}

	log.Printf("viewPUT: views: %+v iptable: %+v\n", views, iptable)

	// Convert response into json structure and then into bytes
	data, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Unable to marshal response: %v\n", err)
		//http.Error(w, "Unable to marshal response", http.StatusInternalServerError)
		//return
	}

	// Send response
	w.Write(data)
}

func viewDELETE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//body, err := ioutil.ReadAll(r.Body)
	//if err != nil {
	//	log.Printf("Error reading body: %v", err)
	//	http.Error(w, "can't read body", http.StatusBadRequest)
	//	return
	//}
	//log.Printf("viewDELETE: body=%s", string(body))
	// Parse the key from url variable and (store) value from the request

	//name := r.PostFormValue("name")
	//fmt.Fprintf(w, "Hello, %s!", name)
	ipport := r.PostFormValue("ip_port") // This doesn't work
	log.Printf("ipport=%s", ipport)
	if ipport == "" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			//http.Error(w, "can't read body", http.StatusBadRequest)
			//return
		}
		ipport = strings.Split(string(body), "=")[1]
		ipport = strings.Replace(ipport, "%3A", ":", -1)
		log.Printf("ipport=%s\n", ipport)
		//return
	}
	//ipport := strings.Split(string(body), "=")[1]
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

	iptableValue := r.PostFormValue("iptable")
	iptable := make(map[string]int)
	if iptableValue != "" {
		if err := json.Unmarshal([]byte(iptableValue), &iptable); err != nil {
			log.Printf("Unable to unmarshal iptable: %v\n", err)
			//http.Error(w, "Unable to unmarshal iptable", http.StatusInternalServerError)
			//return
		}
	}

	log.Printf("viewPUT: ipport=%s\n", ipport)
	log.Printf("viewPUT: len(iptable)=%d\n", len(iptable))
	if len(iptable) <= 0 {
		// Start a new gossip
		for _, view := range views {
			iptable[view] = 0
		}
		iptable[myIP] = 1
		nextNodeURL, err := findNextNode(iptable)
		if err == nil {
			log.Printf("viewPUT: nextNodeURL=%s", nextNodeURL)
			gossipViewDELETE(nextNodeURL, ipport, iptable)
		}
	} else {
		// Gossip if there's an ip that hasn't seen the message
		iptable[myIP] = 1
		nextNodeURL, err := findNextNode(iptable)
		if err == nil {
			log.Printf("viewPUT: nextNodeURL=%s", nextNodeURL)
			gossipViewDELETE(nextNodeURL, ipport, iptable)
		}
	}

	log.Printf("viewPUT: views: %+v iptable: %+v\n", views, iptable)

	// Convert response into json structure and then into bytes
	data, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Unable to marshal response: %v\n", err)
		//http.Error(w, "Unable to marshal response", http.StatusInternalServerError)
		//return
	}

	// Send response
	w.Write(data)
}

func viewUPDATE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	payload := r.PostFormValue("payload")
	msg := new(response.Update)
	if err := json.Unmarshal([]byte(payload), &msg); err != nil {
		log.Printf("viewUPDATE: Unable to unmarshal payload: %v\n", err)
		log.Printf("viewUPDATE: payload=%+v\n", payload)
		log.Printf("viewUPDATE: r.Body=%+v\n", r.Body)
		//http.Error(w, "viewUPDATE: Unable to unmarshal payload", http.StatusInternalServerError)
		//return
	}

	log.Printf("%+v\n", msg)

	vectorClocks = msg.VectorClocks
	KVStore.KeyvalMap = msg.KVS

	thankyou := struct {
		Msg string
	}{
		Msg: "Thank you, I have received the update",
	}

	data, err := json.Marshal(thankyou)
	if err != nil {
		log.Printf("viewUPDATE: Unable to marshal response: %v\n", err)
		//http.Error(w, "viewUPDATE: Unable to marshal response", http.StatusInternalServerError)
		//return
	}
	w.WriteHeader(http.StatusOK)
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

type StupidDelNotWorking struct {
	Ipport  string
	Iptable map[string]int
}

func getViewBody(body io.ReadCloser) *StupidDelNotWorking {
	resp := new(StupidDelNotWorking)
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

// Checks if a key (hashed) belongs in a node
//func keyBelongsInNode(key, ip string) bool {
//	// Hash key to get shard id
//	shardID := hashhere()
//	for _, nodeIP := range placeHolderMap[shardID] {
//		if nodeIP == myIP {
//			return true
//		}
//	}
//	return false
//}

func shardGET(w http.ResponseWriter, r *http.Request) {
	myShardID := -1
	for shardID, nodes := range placeHolderMap {
		for _, node := range nodes {
			if myIP == node {
				myShardID = shardID
				// Could break early here?
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	var resp *response.ShardResponse
	resp = &response.ShardResponse{
		ID: myShardID,
	}
	w.WriteHeader(http.StatusOK)

	dataResp, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Unable to marshal response: %v\n", err)
	}
	w.Write(dataResp)
}

func shardAllGET(w http.ResponseWriter, r *http.Request) {
	var shardIDs []string
	for shardID := range placeHolderMap {
		shardIDs = append(shardIDs, strconv.Itoa(shardID))
	}
	allShardIDs := strings.Join(shardIDS, ",")
	var resp *response.ShardResponse
	resp = &response.ShardResponse{
		Result:   "Success",
		ShardIds: allShardIDs,
	}
	w.WriteHeader(http.StatusOK)

	dataResp, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Unable to marshal response: %v\n", err)
	}
	w.Write(dataResp)
}

func shardMemberGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shardID := vars["id"]

	var resp *response.ShardResponse
	if val, ok := placeHolderMap[shardID]; ok {
		members := strings.Join(placeHolderMap[shardID], ",")
		resp = &response.ShardResponse{
			Result:  "Success",
			Members: members,
		}
		w.WriteHeader(http.StatusOK)
	} else {
		resp = &response.ShardResponse{
			Result: "Error",
			Msg:    fmt.Sprintf("No shard with id %s", shardID),
		}
		w.WriteHeader(http.StatusNotFound)
	}

	dataResp, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Unable to marshal response: %v\n", err)
	}
	w.Write(dataResp)
}

func shardCountGET(w http.ResponseWriter, r *http.Request) {
}

func shardChangePUT(w http.ResponseWriter, r *http.Request) {
}
