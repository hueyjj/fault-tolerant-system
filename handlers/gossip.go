package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"bitbucket.org/cmps128gofour/homework3/response"
	"bitbucket.org/cmps128gofour/homework3/store"
	"bitbucket.org/cmps128gofour/homework3/vectorclock"
)

const replicationFactor int = 2

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

func onReceiveGossip(incomingUnvisitedSet []string) {
	unvisitedNodes = incomingUnvisitedSet
	if len(unvisitedNodes) > 0 {
		//ipPortsToPropagate := selectNodesToPropagateTo()
		/*TODO: Send Gossip to ipPort1, ipPort2 here along with set unvisitedNodes:*/
	} else {
		terminateGossip()
	}
}

func startGossip(iptable []string, myIPport string) {
	//removing IP port of the veiw that recieved the direct request from the client
	for index, currentIP := range iptable {
		if currentIP == myIPport {
			iptable = append(iptable[:index], iptable[index+1:]...)
		}
	}
	onReceiveGossip(iptable)
}

func terminateGossip() {
	unvisitedNodes = nil
}

//https://stackoverflow.com/questions/44956031/how-to-get-intersection-of-two-slice-in-golang
//Time Complexity : O(m+n)
func intersection(s1, s2 []string) (inter []string) {
	hash := make(map[string]bool)
	for _, e := range s1 {
		hash[e] = true
	}
	for _, e := range s2 {
		// If elements present in the hashmap then append intersection list.
		if hash[e] {
			inter = append(inter, e)
		}
	}
	//Remove dups from slice.
	inter = removeDups(inter)
	return
}

//Remove dups from slice.
func removeDups(elements []string) (nodups []string) {
	encountered := make(map[string]bool)
	for _, element := range elements {
		if !encountered[element] {
			nodups = append(nodups, element)
			encountered[element] = true
		}
	}
	return
}

func updateNode(nodeURL string, vectorClocks map[string]vectorclock.Unit, kvs *store.Store) {
	updatePayload := &response.Update{
		VectorClocks: vectorClocks,
		KVS:          kvs.KeyvalMap,
	}

	data, err := json.Marshal(updatePayload)
	if err != nil {
		log.Printf("Unable to marshal iptable: %v\n", err)
		return
	}

	form := url.Values{}
	form.Add("payload", string(data))

	nodeURL = fmt.Sprintf("http://%s/view/update", nodeURL)
	log.Printf("updateNODE: nodeURL=%s", nodeURL)

	req, err := http.NewRequest(http.MethodPost, nodeURL, strings.NewReader(form.Encode()))
	if err != nil {
		log.Printf("Unable to create POST (update) request: nodeURL=%s %v\n", nodeURL, err)
	}
	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	// Don't care about response here, just do it
	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		log.Printf("Unable to do POST (update) request: %v\n", err)
	}
}

func gossipSubjectPUT(nodeURL, key, value, payload string, iptable map[string]int) {
	form := url.Values{}
	form.Add("val", value)
	data, err := json.Marshal(iptable)
	if err != nil {
		log.Printf("Unable to marshal iptable: %v\n", err)
		return
	}
	form.Add("iptable", string(data))
	form.Add("payload", payload)

	nodeURL = fmt.Sprintf("%s/keyValue-store/%s", nodeURL, key)
	log.Printf("gossipSubjectPUT: nodeURL=%s", nodeURL)

	req, err := http.NewRequest(http.MethodPut, nodeURL, strings.NewReader(form.Encode()))
	//log.Printf("iptable:string(data)=%s\n", string(data))

	req.Header.Add("content-type", "application/x-www-form-urlencoded")
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

func gossipSubjectGET(nodeURL, key, payload string, iptable map[string]int) {
	form := url.Values{}
	data, err := json.Marshal(iptable)
	if err != nil {
		log.Printf("Unable to marshal iptable: %v\n", err)
		return
	}
	form.Add("iptable", string(data))
	form.Add("payload", payload)

	nodeURL = fmt.Sprintf("%s/hackedroute/%s", nodeURL, key)
	log.Printf("gossipSubjectGET: nodeURL=%s", nodeURL)

	req, err := http.NewRequest(http.MethodPost, nodeURL, strings.NewReader(form.Encode()))
	//log.Printf("iptable:string(data)=%s\n", string(data))

	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	if err != nil {
		log.Printf("Unable to create GET request: nodeURL=%s %v\n", nodeURL, err)
	}

	// Don't care about response here, just do it
	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		log.Printf("Unable to do GET request: %v\n", err)
	}
}

func gossipSubjectSEARCH(nodeURL, key, payload string, iptable map[string]int) {
	form := url.Values{}
	data, err := json.Marshal(iptable)
	if err != nil {
		log.Printf("Unable to marshal iptable: %v\n", err)
		return
	}
	form.Add("iptable", string(data))
	form.Add("payload", payload)

	nodeURL = fmt.Sprintf("%s/whatevenisthisroute/%s", nodeURL, key)
	log.Printf("gossipSubjectSEARCH: nodeURL=%s", nodeURL)

	req, err := http.NewRequest(http.MethodPost, nodeURL, strings.NewReader(form.Encode()))
	//log.Printf("iptable:string(data)=%s\n", string(data))

	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	if err != nil {
		log.Printf("Unable to create GET (search) request: nodeURL=%s %v\n", nodeURL, err)
	}

	// Don't care about response here, just do it
	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		log.Printf("Unable to do GET (search) request: %v\n", err)
	}
}

func gossipSubjectDEL(nodeURL, key, payload string, iptable map[string]int) {
	form := url.Values{}
	data, err := json.Marshal(iptable)
	if err != nil {
		log.Printf("Unable to marshal iptable: %v\n", err)
		return
	}
	form.Add("iptable", string(data))
	form.Add("payload", payload)

	nodeURL = fmt.Sprintf("%s/keyValue-store/%s", nodeURL, key)
	log.Printf("gossipSubjectDEL: nodeURL=%s", nodeURL)

	req, err := http.NewRequest(http.MethodPost, nodeURL, strings.NewReader(form.Encode()))
	//log.Printf("iptable:string(data)=%s\n", string(data))

	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	if err != nil {
		log.Printf("Unable to create DEL request: nodeURL=%s %v\n", nodeURL, err)
	}

	// Don't care about response here, just do it
	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		log.Printf("Unable to do DEL request: %v\n", err)
	}
}

// Not needed
//func gossipViewGET() {
//}

func gossipViewPUT(nodeURL, ipport string, iptable map[string]int) {
	form := url.Values{}
	form.Add("ip_port", ipport)
	data, err := json.Marshal(iptable)
	if err != nil {
		log.Printf("Unable to marshal iptable: %v\n", err)
		return
	}
	form.Add("iptable", string(data))
	nodeURL = fmt.Sprintf("%s/view", nodeURL)
	req, err := http.NewRequest(http.MethodPut, nodeURL, strings.NewReader(form.Encode()))
	log.Printf("iptable:string(data)=%s\n", string(data))
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
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

func gossipViewDELETE(nodeURL, ipport string, iptable map[string]int) {
	form := url.Values{}
	form.Add("ip_port", ipport)
	data, err := json.Marshal(iptable)
	if err != nil {
		log.Printf("Unable to marshal iptable: %v\n", err)
		return
	}
	form.Add("iptable", string(data))
	//nodeURL = fmt.Sprintf("%s/?ip_port=%s&ip_table=%s", nodeURL, ipport, string(data))
	nodeURL = fmt.Sprintf("%s/view", nodeURL)
	log.Printf("nodeURL=%s", nodeURL)
	req, err := http.NewRequest(http.MethodPost, nodeURL, strings.NewReader(form.Encode()))
	log.Printf("iptable:string(data)=%s\n", string(data))
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	//req.Header.Add("content-type", "multipart/form-data")
	if err != nil {
		log.Printf("Unable to create DELETE request: nodeURL=%s %v\n", nodeURL, err)
	}

	// Don't care about response here, just do it
	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		log.Printf("Unable to do DELETE request: %v\n", err)
	}
}

func findNextNode(iptable map[string]int) (string, error) {
	nodeURL := ""
	for key, value := range iptable {
		if value == 0 {
			nodeURL = key
		}
	}
	if nodeURL == "" {
		return "", errors.New("Unable to find valid node url")
	}
	return "http://" + nodeURL, nil
}
