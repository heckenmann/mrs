FROM golang:1.14.2-alpine as build
RUN apk update && apk upgrade && apk add git
RUN mkdir -p /tmp/build
WORKDIR /tmp/build
COPY ./multiregexsuche.go .
COPY ./mrs.yml .
RUN go get "gopkg.in/yaml.v2"
RUN go get "github.com/gorilla/mux"
RUN go build multiregexsuche.go

####################################################
FROM alpine:3.11
EXPOSE 8080
RUN mkdir -p /opt/mrs
WORKDIR /opt/mrs
COPY --from=build /tmp/build/multiregexsuche .
COPY --from=build /tmp/build/mrs.yml .
CMD ./multiregexsuche