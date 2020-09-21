GOFILES := $(wildcard *.go)

all: pluto-vs-relay.linux64

pluto-vs-relay.linux64: $(GOFILES)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o pluto-vs-relay.linux64
	chmod a+x pluto-vs-relay.linux64

docker: pluto-vs-relay.linux64
	docker build . -t guardianmultimedia/pluto-vs-relay:DEV
