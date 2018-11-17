package vectorclock

import (
	"reflect"
	"testing"
)

// Some maps that we can cast later
var l = map[string]int{
	"1": 2,
	"2": 1,
	"3": 5,
}

var m = map[string]int{
	"1": 2,
	"2": 2,
	"3": 3,
}

var n = map[string]int{
	"1": 1,
	"2": 1,
	"3": 1,
}

var y = map[string]int{
	"1": 1,
	"2": 1,
}

var z = map[string]int{
	"3": 1,
	"4": 1,
	"5": 1,
}

var alpha = map[string]int{
	"3": 1,
	"4": 2,
	"5": 1,
}

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want VectorClock
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVectorClock_HappenedBefore(t *testing.T) {
	type args struct {
		u VectorClock
	}
	tests := []struct {
		name string
		v    VectorClock
		args args
		want bool
	}{
		{
			name: "Did Happen Before",
			v:    VectorClock(n),
			args: args{u: VectorClock(l)},
			want: true,
		},
		{
			name: "Did Not Happen Before",
			v:    VectorClock(l),
			args: args{u: VectorClock(n)},
			want: false,
		},
		{
			name: "Differnt Keys",
			v:    VectorClock(l),
			args: args{u: VectorClock(z)},
			want: false,
		},
		{
			name: "Differnt lengths",
			v:    VectorClock(l),
			args: args{u: VectorClock(y)},
			want: false,
		},
		{
			name: "Concurrent Vectors",
			v:    VectorClock(l),
			args: args{u: VectorClock(m)},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.HappenedBefore(tt.args.u); got != tt.want {
				t.Errorf("VectorClock.HappenedBefore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVectorClock_ConcurrentWith(t *testing.T) {
	type args struct {
		u VectorClock
	}
	tests := []struct {
		name string
		v    VectorClock
		args args
		want bool
	}{
		{
			name: "Concurrent Vectors",
			v:    VectorClock(l),
			args: args{u: VectorClock(m)},
			want: true,
		},
		{
			name: "Non - Concurrent Vectors",
			v:    VectorClock(l),
			args: args{u: VectorClock(n)},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.ConcurrentWith(tt.args.u); got != tt.want {
				t.Errorf("VectorClock.ConcurrentWith() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVectorClock_Equals(t *testing.T) {
	type args struct {
		u VectorClock
	}
	tests := []struct {
		name string
		v    VectorClock
		args args
		want bool
	}{
		// TODO: Add test cases.
		{
			name: "Equal Vectors",
			v:    VectorClock(l),
			args: args{u: VectorClock(l)},
			want: true,
		},
		{
			name: "Differnt Keys",
			v:    VectorClock(l),
			args: args{u: VectorClock(z)},
			want: false,
		},
		{
			name: "Differnt lengths",
			v:    VectorClock(l),
			args: args{u: VectorClock(y)},
			want: false,
		},
		{
			name: "Differnt values",
			v:    VectorClock(z),
			args: args{u: VectorClock(alpha)},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.Equals(tt.args.u); got != tt.want {
				t.Errorf("VectorClock.Equals() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestVectorClock_MergeE(t *testing.T) {

	// Test 1
	// Merge l and m
	t.Run("Valid Merge", func(t *testing.T) {
		v := VectorClock(l)
		u := VectorClock(m)

		err := v.MergeE(u)

		if err != nil {
			t.Errorf("VectorClock.MergeE() = %v, want %v", err, nil)
		}
		want := VectorClock(map[string]int{"1": 2, "2": 2, "3": 5})
		result := v.Equals(want)
		if !result {
			t.Errorf("VectorClock.MergeE() = %v, want %v", v, want)
		}
	})

	// Test 2
	t.Run("Invalid Merge Due To Length", func(t *testing.T) {
		v := VectorClock(l)
		u := VectorClock(y)

		err := v.MergeE(u)

		if err == nil {
			t.Errorf("VectorClock.MergeE() = %v, want %v", err, nil)
		}

		want := VectorClock(map[string]int{"1": 2, "2": 2, "3": 5})

		result := v.Equals(want)
		if !result {
			t.Errorf("VectorClock.MergeE() = %v, want %v", v, want)
		}
	})

	// Test 3
	t.Run("Invalid Merge Due To Differnt Keys", func(t *testing.T) {
		v := VectorClock(l)
		u := VectorClock(z)

		err := v.MergeE(u)

		if err == nil {
			t.Errorf("VectorClock.MergeE() = %v, want %v", err, nil)
		}

		want := VectorClock(map[string]int{"1": 2, "2": 2, "3": 5})

		result := v.Equals(want)
		if !result {
			t.Errorf("VectorClock.MergeE() = %v, want %v", v, want)
		}
	})

}
