# CMPS 128 Project

Class project for UCSC CMPS 128 Distributed Systems.

## Setup

```sh
mkdir -p $GOPATH/src/bitbucket.org/cmps128gofour
cd $GOPATH/src/bitbucket.org/cmps128gofour
git clone https://$USERNAME@bitbucket.org/cmps128gofour/homework4.git
```

## Install dependencies

```sh
go get -u github.com/gorilla/mux
go get -u github.com/stretchr/testify/assert
```

## Building the server

```sh
# Using golang
go build

# Or using docker
docker build -t cmps128gofour/homework4 .
```

## Running the server

```sh
# If built with go
./homework4 # To stop the server, press ctrl-c
# if built with docker
docker run -p 8082:8080 -e VIEW="176.32.164.10:8082,176.32.164.10:8083" -e IP_PORT="176.32.164.10:8082" testing
docker run -p 8083:8080 -e VIEW="176.32.164.10:8082,176.32.164.10:8083" -e IP_PORT="176.32.164.10:8083" testing
```

## Testing the output of the server

Once the server is running, use the `curl` command to test the server

```sh
curl -v $URL -X $METHOD -d $FORMVALUES
# e.g.
curl -v localhost:8080/test -X POST -d "msg=Hello World"
```

Example output (<-- signify area of importance):

```
Note: Unnecessary use of -X or --request, POST is already inferred.
*   Trying 127.0.0.1...
* TCP_NODELAY set
* Connected to localhost (127.0.0.1) port 8080 (#0)
> POST /test HTTP/1.1 								<-- This is the route that was used to send the listed method
> Host: localhost:8080
> User-Agent: curl/7.58.0
> Accept: */*
> Content-Length: 15
> Content-Type: application/x-www-form-urlencoded
>
* upload completely sent off: 15 out of 15 bytes
< HTTP/1.1 200 OK 									<-- This is the response code
< Date: Wed, 10 Oct 2018 07:01:42 GMT
< Content-Length: 34
< Content-Type: text/plain; charset=utf-8
<
* Connection #0 to host localhost left intact
POST message received: Hello World 				<-- This is the output of the sever
```

## Running unit tests

Start the server, and use `go test` to run unit tests on the handlers

```sh
go test -v ./...
```

Example output:

```
?   	bitbucket.org/cmps128gofour/homework1	[no test files]
=== RUN   Test_helloGET
=== RUN   Test_helloGET/Regular_GET_hello
--- PASS: Test_helloGET (0.00s)
    --- PASS: Test_helloGET/Regular_GET_hello (0.00s)
=== RUN   Test_helloPOST
=== RUN   Test_helloPOST/Regular_POST_hello
2018/10/12 01:00:03 POST /hello not supported
--- PASS: Test_helloPOST (0.00s)
    --- PASS: Test_helloPOST/Regular_POST_hello (0.00s)
=== RUN   Test_testGET
=== RUN   Test_testGET/Regular_GET_test
--- PASS: Test_testGET (0.00s)
    --- PASS: Test_testGET/Regular_GET_test (0.00s)
=== RUN   Test_testPOST
=== RUN   Test_testPOST/Regular_POST_test
=== RUN   Test_testPOST/Regular_POST_test_2
=== RUN   Test_testPOST/No_forms_POST_test
2018/10/12 01:00:03 received request with missing key "msg"
--- PASS: Test_testPOST (0.00s)
    --- PASS: Test_testPOST/Regular_POST_test (0.00s)
    --- PASS: Test_testPOST/Regular_POST_test_2 (0.00s)
    --- PASS: Test_testPOST/No_forms_POST_test (0.00s)
PASS
ok  	bitbucket.org/cmps128gofour/homework1/handlers	(cached)
>>>>>>> c0c19f959ea3c6f91dbc7ac988f677559ac826d4
```
