package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	clientResponse "bitbucket.org/cmps128gofour/homework3/response"
	"github.com/gorilla/mux"
)

var mainIP string
var client = http.Client{}

func proxySubjectGET(w http.ResponseWriter, r *http.Request) {
	// Parse the key from url variable and (store) value from the request
	vars := mux.Vars(r)
	key := vars["subject"]

	// make request
	requestString := fmt.Sprintf("%s/keyValue-store/%s", mainIP, key)
	request, err := http.NewRequest(http.MethodGet, requestString, nil)

	if err != nil {
		log.Println("could not make request:", err)
		return
	}

	// Send the request
	response, err := client.Do(request)

	if err != nil {
		log.Println("could not get response:", err)
		// Main server is down
		respondError501(w)
		return
	}
	log.Println("response status code: ", response.StatusCode)

	// Write the header
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(response.StatusCode)
	// Read the body

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Println("could not read body:", err)
		return
	}

	// Write the body
	_, err = w.Write(body)
	if err != nil {
		log.Println("could not write body:", err)
	}

}

func proxySubjectSEARCH(w http.ResponseWriter, r *http.Request) {
	// Parse the key from url variable and (store) value from the request
	vars := mux.Vars(r)
	key := vars["subject"]

	// make request
	requestString := fmt.Sprintf("%s/keyValue-store/search/%s", mainIP, key)
	request, err := http.NewRequest(http.MethodGet, requestString, nil)
	if err != nil {
		log.Println("could not make request:", err)
		return
	}

	// Send the request
	response, err := client.Do(request)
	if err != nil {
		log.Println("could not get response:", err)
		// Main server is down
		respondError501(w)
		return
	}
	log.Println("response status code: ", response.StatusCode)

	// Write the header
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(response.StatusCode)

	// Read the body
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("could not read body:", err)
		return
	}

	// Write the body
	_, err = w.Write(body)
	if err != nil {
		log.Println("could not write body:", err)
	}

}

func proxySubjectDEL(w http.ResponseWriter, r *http.Request) {
	// Parse the key from url variable and (store) value from the request
	vars := mux.Vars(r)
	key := vars["subject"]

	// make request
	requestString := fmt.Sprintf("%s/keyValue-store/%s", mainIP, key)
	request, err := http.NewRequest(http.MethodDelete, requestString, nil)
	if err != nil {
		log.Println("could not make request:", err)
		return
	}

	// Send the request
	response, err := client.Do(request)
	if err != nil {
		log.Println("could not get response:", err)
		// Main server is down
		respondError501(w)
		return
	}

	// Write the header
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(response.StatusCode)

	// Read the body
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("could not read body:", err)
		return
	}

	// Write the body
	_, err = w.Write(body)
	if err != nil {
		log.Println("could not write body:", err)
	}
}

func proxySubjectPUT(w http.ResponseWriter, r *http.Request) {
	// Parse the key from url variable and (store) value from the request
	vars := mux.Vars(r)
	key := vars["subject"]

	// Create form data
	form := url.Values{}
	form.Add("val", r.PostFormValue("val"))

	// Make request
	requestString := fmt.Sprintf("%s/keyValue-store/%s", mainIP, key)
	request, err := http.NewRequest(http.MethodPut, requestString, strings.NewReader(form.Encode()))
	if err != nil {
		log.Println("could not create request:", err)
		return
	}

	request.Header.Add("content-type", "application/x-www-form-urlencoded")

	// Send the request
	response, err := client.Do(request)
	if err != nil {
		log.Println("could not get response:", err)
		// Main server is down
		respondError501(w)
		return
	}

	// Write the header
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(response.StatusCode)

	// Read the body
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("could not read body:", err)
		return
	}

	// Write the body
	_, err = w.Write(body)
	if err != nil {
		log.Println("could not write body:", err)
	}
}

// Have the ResponseWriter write response 501 in JSON format.
func respondError501(w http.ResponseWriter) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	var resp *clientResponse.Response
	resp = &clientResponse.Response{
		Result: "Error",
		Msg:    "Server unavailable",
	}

	// Convert response into json structure and then into bytes
	data, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Unable to marshal response: %v\n", err)
		http.Error(w, "Unable to marshal response", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(data)
	if err != nil {
		log.Println("could not write body:", err)
	}
}
