package shard

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/serialx/hashring"
)

func FixLonelyNodes(shardTable map[int][]string) error {

	// Make an arr to store our lonely nodes
	// and our keys from the hash table
	lonelyNodes := []string{}
	keys := []string{}

	// Run through the map
	for key, value := range shardTable {
		// If the arr has less than 2
		// nodes its a lonely node
		if len(value) <= 1 {
			// push this node into the lonely nodes array
			lonelyNodes = append(lonelyNodes, value[0])
			// and delete the key from the table
			delete(shardTable, key)
			continue
		}
		// otherwise append the key to our key arr
		keys = append(keys, strconv.Itoa(key))
	}

	// Make a hashring to distribute our nodes equally
	ring := hashring.New(keys)

	// for each lonely node
	for _, node := range lonelyNodes {
		// figure out where to put it
		index, _ := ring.GetNode(node)
		// convert to int since the key was a string and we want an int
		conv, err := strconv.Atoi(index)
		if err != nil {
			return err
		}
		// add the node to map
		shardTable[conv] = append(shardTable[conv], node)
	}

	return nil
}

func ShardIt(views string, nodesPerShard int) (map[int][]string, error) {

	// you can't have less than one node per shard
	if nodesPerShard <= 1 {
		return nil, fmt.Errorf("cannot have less than 2 nodes per shard")
	}

	// Turn the views into an arr
	tviews := strings.Split(views, ",")

	// Group up the nodes
	arrCount := -1
	shards := []([]string){}

	for i, str := range tviews {

		if i%nodesPerShard == 0 {
			arrCount++
			shards = append(shards, []string{})
		}

		shards[arrCount] = append(shards[arrCount], str)

	}

	// Make the groups into a map,
	// where the key is the id
	// and the value is the group
	idToShards := map[int][]string{}
	for i, shard := range shards {
		idToShards[i] = shard
	}

	// TODO: Fix any groups that have only one node
	err := FixLonelyNodes(idToShards)
	if err != nil {
		return nil, err
	}

	return idToShards, nil

}
