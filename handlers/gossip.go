package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var replicationFactor = 2
var unvisitedNodes []string

func selectNodesToPropagateTo() []string {
	count := 0
	rand.Seed(time.Now().Unix())
	var result []string

	for count < replicationFactor && len(unvisitedNodes) > 0 {
		i := rand.Int() % len(unvisitedNodes)
		//Cut a random ipPort from the set.
		ipPort := unvisitedNodes[i]
		unvisitedNodes = append(unvisitedNodes[:i], unvisitedNodes[i+1:]...)
		result = append(result, ipPort)
		count++
	}
	return result
}

func onReceiveGossip() {
	if len(unvisitedNodes) > 0 {
		//ipPortToPropagate := selectNodesToPropagateTo()
		/*TODO: Send Gossip to ipPort1, ipPort2 here along with set unvisitedNodes:*/
	}
}

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
