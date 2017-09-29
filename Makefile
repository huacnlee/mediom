install:
	@npm install -g coffee-script@1.6.2 node-sass
	@go get github.com/huacnlee/train
	@go get github.com/revel/cmd/revel
	@go build -o $GOPATH/bin/train github.com/huacnlee/train/cmd
	@dep ensure
server:
	revel run github.com/huacnlee/mediom
release: assets
	GOOS=linux GOARCH=amd64 revel package github.com/huacnlee/mediom prod
assets:
	@train --source app/assets --out public
test:
	@cd app; go test
	@cd app/controllers; go test
	@cd app/models; go test
	@cd tests; go test