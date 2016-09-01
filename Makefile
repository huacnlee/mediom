install:
	@cd app; go get
server:
	revel run github.com/huacnlee/mediom
release:
	GOOS=linux GOARCH=amd64 revel package github.com/huacnlee/mediom
test:
	@cd app; go test
	@cd app/controllers; go test
	@cd app/models; go test
	@cd tests; go test