// Package handlers defines the http handler functions
// that will be used in the homework 1 program
package handlers

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func Test_helloGET(t *testing.T) {
	type want struct {
		code int
		body []byte
	}
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Regular GET hello",
			args: args{
				r: httptest.NewRequest("GET", "localhost:8080/hello", nil),
				w: httptest.NewRecorder(),
			},
			want: want{
				code: http.StatusOK,
				body: []byte("Hello world!"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			helloGET(tt.args.w, tt.args.r)
			res := tt.args.w.Result()
			defer res.Body.Close()
			bb, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Errorf("could not read body: %v", err)
			}
			got := want{
				code: res.StatusCode,
				body: bb,
			}
			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("helloGET() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_helloPOST(t *testing.T) {
	type want struct {
		code int
	}
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Regular POST hello",
			args: args{
				r: httptest.NewRequest("POST", "localhost:8080/hello", nil),
				w: httptest.NewRecorder(),
			},
			want: want{
				code: http.StatusMethodNotAllowed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			helloPOST(tt.args.w, tt.args.r)
			res := tt.args.w.Result()
			got := want{
				code: res.StatusCode,
			}
			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("helloPOST() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_testGET(t *testing.T) {
	type want struct {
		code int
		body []byte
	}
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Regular GET test",
			args: args{
				r: httptest.NewRequest("GET", "localhost:8080/test", nil),
				w: httptest.NewRecorder(),
			},
			want: want{
				code: http.StatusOK,
				body: []byte("GET request received"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testGET(tt.args.w, tt.args.r)
			res := tt.args.w.Result()
			defer res.Body.Close()
			bb, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Errorf("could not read body: %v", err)
			}
			got := want{
				code: res.StatusCode,
				body: bb,
			}
			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("testGET() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_testPOST(t *testing.T) {
	type want struct {
		code int
		body []byte
	}
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Regular POST test",
			args: args{
				r: httptest.NewRequest("POST", "localhost:8080/test?msg=cs128", nil),
				w: httptest.NewRecorder(),
			},
			want: want{
				code: http.StatusOK,
				body: []byte("POST message received: cs128"),
			},
		},
		{
			name: "Regular POST test 2",
			args: args{
				r: httptest.NewRequest("POST", "localhost:8080/test?msg=jorge", nil),
				w: httptest.NewRecorder(),
			},
			want: want{
				code: http.StatusOK,
				body: []byte("POST message received: jorge"),
			},
		},
		{
			name: "No forms POST test",
			args: args{
				r: httptest.NewRequest("POST", "localhost:8080/test", nil),
				w: httptest.NewRecorder(),
			},
			want: want{
				code: http.StatusOK,
				body: []byte("POST message received: "),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testPOST(tt.args.w, tt.args.r)
			res := tt.args.w.Result()
			defer res.Body.Close()
			bb, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Errorf("could not read body: %v", err)
			}
			got := want{
				code: res.StatusCode,
				body: bb,
			}
			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("testPOST() = %v, want %v", got, tt.want)
			}
		})
	}
}
