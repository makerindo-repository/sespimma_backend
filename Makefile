BINARY  := api-go
BIN_DIR := ./bin
CMD     := ./cmd/main.go

.PHONY: build run test clean deploy lint tidy

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	  go build -ldflags="-w -s" -o $(BIN_DIR)/$(BINARY) $(CMD)

run:
	go run $(CMD)

test:
	go test ./...

tidy:
	go mod tidy

lint:
	golangci-lint run ./...

clean:
	rm -rf $(BIN_DIR)

# Deploy: build then copy binary to VPS and restart PM2
deploy: build
	scp $(BIN_DIR)/$(BINARY) maker@103.172.205.183:/var/www/sespimma-api-go/$(BINARY)
	scp -r migrations/ maker@103.172.205.183:/var/www/sespimma-api-go/
	ssh maker@103.172.205.183 " \
	  export NVM_DIR=\$$HOME/.nvm && \
	  [ -s \$$NVM_DIR/nvm.sh ] && . \$$NVM_DIR/nvm.sh; \
	  chmod +x /var/www/sespimma-api-go/$(BINARY) && \
	  pm2 restart backend-golang || pm2 start /var/www/sespimma-api-go/$(BINARY) \
	    --name backend-golang --cwd /var/www/sespimma-api-go && \
	  pm2 save"
