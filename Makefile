.PHONY: bootstrap
bootstrap:
	go mod download && go mod vendor

.PHONY: test
test: bootstrap
	CGO_ENABLED=0 go test -mod=vendor -v -covermode=atomic -coverprofile=.coverage.out .
	go tool cover -html=.coverage.out -o=.coverage.html

.PHONY: covhtml
covhtml:
	open .coverage.html
