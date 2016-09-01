install:
	@cd app; go get
server:
	revel run github.com/huacnlee/mediom
release:
	GOOS=linux GOARCH=amd64 revel package github.com/huacnlee/mediom
test:
	@cd tests; go test