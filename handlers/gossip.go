package handlers

import (
	"math/rand"
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

func gossipViewGET() {

}

func gossipViewPUT() {

}

func gossipViewDELETE() {

}
