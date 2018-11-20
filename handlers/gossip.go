package handlers

import (
	"math/rand"
	"time"
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

func gossipSubjectPUT() {
}

func gossipSubjectGET() {

}

func gossipSubjectSEARCH() {

}

func gossipSubjectDEL() {

}

func gossipViewGET() {

}

func gossipViewPUT() {

}

func gossipViewDELETE() {

}
