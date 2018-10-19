package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var hostURL = "http://localhost:8080"
var subURLhello = "/hello"
var subURLtest = "/test"

func TestGet_hello(t *testing.T) {

	assert := assert.New(t)
	resBody := makeGetRequest(subURLhello)

	// assert equality
	assert.Equal("Hello world!", resBody, "they should be equal")

	// assert inequality
	assert.NotEqual("123fail", resBody, "they should not be equal")
}

func TestGet_test(t *testing.T) {

	assert := assert.New(t)
	resBody := makeGetRequest(subURLtest)

	// assert equality
	assert.Equal("GET request received", resBody, "they should be equal")

	// assert inequality
	assert.NotEqual("123fail", resBody, "they should not be equal")
}

func TestPost_hello(t *testing.T) {

	assert := assert.New(t)
	pmr := "POST message received: "
	// assert equality
	assert.Equal("", makePostRequest(subURLhello, "This is a post123"), "they should be equal")
	assert.Equal("", makePostRequest(subURLhello, "ACoolMessage"), "they should be equal")
	assert.Equal("", makePostRequest(subURLhello, ""), "they should be equal")

	// assert inequality
	assert.NotEqual(pmr+"123fail", makePostRequest(subURLhello, "ACoolMessage"), "they should not be equal")
}

func TestPost_test(t *testing.T) {

	assert := assert.New(t)
	pmr := "POST message received: "
	// assert equality
	assert.Equal(pmr+"This is a post123", makePostRequest(subURLtest, "This is a post123"), "they should be equal")
	assert.Equal(pmr+"ACoolMessage", makePostRequest(subURLtest, "ACoolMessage"), "they should be equal")
	assert.Equal(pmr+"", makePostRequest(subURLtest, ""), "they should be equal")

	// assert inequality
	assert.NotEqual(pmr+"123fail", makePostRequest(subURLtest, "ACoolMessage"), "they should not be equal")
}

// Make a simple GET request
// Returns response message-body
func makeGetRequest(subURL string) string {

	response, err := http.Get(hostURL + subURL)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
	return string(body)
}

// Make a simple POST request that takes a message
// Returns response message-body
func makePostRequest(subURL string, msg string) string {

	response, err := http.PostForm(hostURL+subURL,
		url.Values{"msg": {msg}})

	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
	return string(body)
}
