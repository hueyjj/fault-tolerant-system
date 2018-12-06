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
		// If the arr has only one
		// node its a lonely node
		if len(value) == 1 {
			// push this node into the lonely nodes array
			lonelyNodes = append(lonelyNodes, value[0])
			// and delete the key from the table
			delete(shardTable, key)
			continue
		}
		if len(value) == 0 {
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

func manualPlace(views []string, S int) (map[int][]string, error) {
	idToShards := map[int][]string{}

	for i := 0; i < S; i++ {
		idToShards[i] = []string{}
	}

	count := 0

	for i, view := range views {

		if i != 0 && i%(len(views)/S) == 0 {
			count++
		}

		idToShards[count] = append(idToShards[count], view)
	}

	return idToShards, nil
}

// Shard creates a map of ids to shards , the map will have a size that is <= S
func Shard(views string, S int) (map[int][]string, error) {

	// Turn the views into an arr
	tviews := strings.Split(views, ",")

	// If len(views) is divisible by S then
	// manually place items in the map
	if len(tviews)%S == 0 {
		return manualPlace(tviews, S)
	}

	// Make the groups into a map,
	// where the key is the id
	// and the value is the group
	idToShards := map[int][]string{}

	shardIDs := []string{}
	// Prepopulate the map w empty shards
	for i := 0; i < S; i++ {
		idToShards[i] = []string{}
		conv := strconv.Itoa(i)
		shardIDs = append(shardIDs, conv)
	}

	// Make a hashring using the shardIDS
	ring := hashring.New(shardIDs)

	// Uniformly distribute values in tviews
	for _, view := range tviews {
		asciiShardNumber, _ := ring.GetNode(view)
		shardNumber, _ := strconv.Atoi(asciiShardNumber)
		idToShards[shardNumber] = append(idToShards[shardNumber], view)
	}

	FixLonelyNodes(idToShards)
	return idToShards, nil
}

// ShardN makes a map of shards where each shard has
// nodePerShard nodes in it
func ShardN(views string, nodesPerShard int) (map[int][]string, error) {

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

func GetShardID(iptable map[int][]string, key string) int {
	// Get the keys from the table
	keys := []string{}
	for key := range iptable {
		conv := strconv.Itoa(key)
		keys = append(keys, conv)
	}

	// use the keys to make a hash ring
	ring := hashring.New(keys)
	// Hash the key, and find the shard it corr. to
	node, _ := ring.GetNode(key)
	// convert the string node to an int
	shardID, _ := strconv.Atoi(node)

	return shardID
}
