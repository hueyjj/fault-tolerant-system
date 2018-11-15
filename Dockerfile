FROM golang:latest AS build
WORKDIR /go/src/bitbucket.org/cmps128gofour/homework3
RUN go get -u github.com/gorilla/mux
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=build /go/src/bitbucket.org/cmps128gofour/homework3/app .
CMD ["./app"]
