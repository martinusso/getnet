COVER_FILE="/tmp/go-cover.$$.tmp"

.PHONY: examples
examples:
	go run examples/main.go

.PHONY: test
test:
	go test ./... -race -cover -count=1

.PHONY: cover 
cover:
	go test ./... -race -count=1 -coverprofile=$(COVER_FILE)
	go tool cover -html=$(COVER_FILE)
	rm $(COVER_FILE)
