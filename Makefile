install:
	@npm install -g node-sass@3.3.2
	@npm install -g coffee-script@1.6.2
	@go get github.com/huacnlee/train
	@go build -o $GOPATH/bin/train github.com/huacnlee/train/cmd
	@cd app; go get
server:
	revel run github.com/huacnlee/mediom
release:
	@make assets
	GOOS=linux GOARCH=amd64 revel package github.com/huacnlee/mediom prod
assets:
	@train --source app/assets --out public
test:
	@cd app; go test
	@cd app/controllers; go test
	@cd app/models; go test
	@cd tests; go test