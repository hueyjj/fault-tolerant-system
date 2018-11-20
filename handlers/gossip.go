package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func gossipSubjectPUT() {
}

func gossipSubjectGET() {

}

func gossipSubjectSEARCH() {

}

func gossipSubjectDEL() {

}

// Not needed
//func gossipViewGET() {
//}

func gossipViewPUT(nodeURL, ipport string, iptable map[string]bool) {
	form := url.Values{}
	form.Add("ip_port", ipport)
	data, err := json.Marshal(iptable)
	if err != nil {
		log.Printf("Unable to marshal iptable: %v\n", err)
		return
	}
	form.Add("iptable", string(data))
	req, err := http.NewRequest(http.MethodPut, nodeURL, strings.NewReader(form.Encode()))
	if err != nil {
		log.Printf("Unable to create PUT request: nodeURL=%s %v\n", nodeURL, err)
	}

	// Don't care about response here, just do it
	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		log.Printf("Unable to do PUT request: %v\n", err)
	}
}

func gossipViewDELETE() {

}

func findNextNode(iptable map[string]bool) (string, error) {
	nodeURL := ""
	for key, value := range iptable {
		if value == false {
			nodeURL = key
		}
	}
	if nodeURL == "" {
		return "", errors.New("Unable to find valid node url")
	}
	return "http://" + nodeURL + "/view", nil
}
