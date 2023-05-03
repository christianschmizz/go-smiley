build:
	@CGO_ENABLED=0 go build -v ./cmd/smileyd

check:
	@go get -u -a golang.org/x/tools/cmd/stringer

all: check
	@stringer -type=PollEventResponseCode
	@stringer -type=ResponseCode
	@stringer -type=CommandCode
	@stringer -type=UnitType

.PHONY: test
test:
	@go test ./...
