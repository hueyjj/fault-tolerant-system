package vectorclock

import (
	"encoding/json"
	"fmt"
)

type Unit struct {
	Tick      int
	Timestamp int64 // Unix timestamp
}

// VectorClock is a map of the id of a process to it's current clock value
type VectorClock map[string]int

// New Creates a new vector clock
func New() VectorClock {
	return VectorClock(make(map[string]int))
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

// InitProcessE is similar to InitProcess, only that
// it will return an error if the process already exists
func (v VectorClock) InitProcessE(id string) error {
	_, ok := v[id]

	if ok {
		return fmt.Errorf("cannot init process %s: process already exists", id)
	}

	v[id] = 0
	return nil
}

// SetClockE takes in a process id and sets its clock to the given value.
// if the process id does not exist an error is returned
func (v VectorClock) SetClockE(id string, clock int) error {

	if _, ok := v[id]; !ok {
		return fmt.Errorf("no such process exists with id: %s", id)
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

// TickE is similar to Tick. Only that it returns an error
// if the id is not found
func (v VectorClock) TickE(id string) error {
	val, ok := v[id]
	if !ok {
		return fmt.Errorf("no such process exists with id: %s", id)
	}
	v[id] = val + 1
	return nil
}

// Merge takes the max of both vectors and updates the value
// value of the calling vector
// Only updates keys that are shared between two vectors
func (v VectorClock) Merge(u VectorClock) {
	// Loop through v
	for key, value := range v {
		// If the key exists in the u vector and
		// the value in v is less than the value in u, then
		// assign it
		if other, ok := u[key]; (ok) && (value < other) {
			v[key] = other
		}
	}

}

// MergeE is similar to Merge, only that it returns
// an error if there non similar keys, or the vectors
// are of differnt lengths. If an error is encounted the
// calling vector is reverted to it's original pre-calling
// state
func (v *VectorClock) MergeE(u VectorClock) error {

	if len(*v) != len(u) {
		return fmt.Errorf("could not merge vectors: vectors are different lengths - len(u)->%d , len(v)->%d", len(*v), len(u))
	}

	copy := v.clone()

	// Loop through v
	for key, value := range *v {
		// If the key exists in the u vector and
		// the value in v is less than the value in u, then
		// assign it
		other, ok := u[key]

		if !ok {
			*v = copy
			return fmt.Errorf("could not merge vectors: non-shared key - %s", key)
		}

		if (ok) && (value < other) {
			(*v)[key] = other
		}
	}

	return nil
}

// Equal returns true if v == u
// that is that v and u are of the same length &&
// v and u contain the same keys that map to the same values
func (v VectorClock) Equals(u VectorClock) bool {
	// len has to be the same
	if len(u) != len(v) {
		return false
	}

	// loop through v keys
	for key, value := range v {
		otherVal, ok := u[key]
		if !ok {
			return false
		}
		if otherVal != value {
			return false
		}
	}

	return true
}

// HappenedBefore returns true if v->u
// see https://en.wikipedia.org/wiki/Vector_clock#Partial_ordering_property
func (v VectorClock) HappenedBefore(u VectorClock) bool {
	return v.lessthan(u)
}

// ConcurrentWith returns true if v || u
// (i.e.) !(v < u) && !(u < v)
func (v VectorClock) ConcurrentWith(u VectorClock) bool {
	return !v.lessthan(u) && !u.lessthan(v)
}

// lessthan returns true if v < u
// see https://en.wikipedia.org/wiki/Vector_clock#Partial_ordering_property
func (v VectorClock) lessthan(u VectorClock) bool {
	// if the len doesn't match up then we're comparing apples and oranges
	if len(v) != len(u) {
		return false
	}

	// Start a strictly less than counter
	strict := 0

	for key, value := range v {

		// make sure that both vectors contain the same keys
		if _, ok := u[key]; !ok {
			return false
		}

		// if our value is greater than we can't be less than u
		if value > u[key] {
			return false
		}

		// here we must be less than or equal to, here we only need to check
		// for strictly less than, so that we can increment our counter
		if value < u[key] {
			strict++
		}

	}

	// return true if the counter is at least 1
	return strict >= 1
}

// clone clones the calling vector
func (v VectorClock) clone() VectorClock {
	c := New()
	for key, value := range v {
		c[key] = value
	}
	return c
}
