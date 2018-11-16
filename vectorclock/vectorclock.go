package vectorclock

import (
	"encoding/json"
	"fmt"
)

// VectorClock is a map of the id of a process to it's current clock value
type VectorClock map[string]int

// New Creates a new vector clock
func New() VectorClock {
	return VectorClock{}
}

// Marshall creates a json encoded vectorclock
func (v VectorClock) Marshall() ([]byte, error) {
	m := make(map[string]interface{})

	for k, val := range v {
		m[k] = val
	}

	return json.Marshal(m)
}

// Unmarshall creates a new VectorClock from json-encoded data
func Unmarshall(clock []byte) (VectorClock, error) {
	// From here: https://blog.golang.org/json-and-go

	// Create an empty interface and unmarshall data into it
	var f interface{}
	err := json.Unmarshal(clock, &f)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshall data: %v", err)
	}

	// assert that m is a map
	m := f.(map[string]interface{})

	// Create a new clock
	vclock := New()

	// loop through the map
	for k, v := range m {
		switch vv := v.(type) {
		case float64:
			vclock.SetClock(k, int(vv))
		default:
			return nil, fmt.Errorf("unexpected value: %v", v)
		}
	}

	return vclock, nil
}

// ProcessValue takes an process's id, and returns its clock value
// returns an error if the process is not found
func (v VectorClock) ProcessValue(id string) (int, error) {
	val, ok := v[id]

	if !ok {
		return 0, fmt.Errorf("no such process id: %s", id)
	}

	return val, nil
}

// InitProcess takes in an id and makes it clock 0. Does nothing if
// the process already exists
func (v VectorClock) InitProcess(id string) {
	_, ok := v[id]

	if ok {
		return
	}

	v[id] = 0
}

// SetClockE takes in a process id and sets its clock to the given value.
// if the process id does not exist an error is returned
func (v VectorClock) SetClockE(id string, clock int) error {

	if _, ok := v[id]; !ok {
		return fmt.Errorf("no such process exist with id: %s", id)
	}

	v[id] = clock
	return nil
}

// SetClock is similar to SetClockE, except that it creates a process if the id
// is not found.
func (v VectorClock) SetClock(id string, clock int) {
	v[id] = clock
}

// Tick updates the value of the process by one, does nothing
// if the id does not exist
func (v VectorClock) Tick(id string) {
	val, ok := v[id]
	if ok {
		v[id] = val + 1
	}
}
