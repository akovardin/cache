SERVER=server
CLIENT=client

# ${CACHE_BIN}: install cahce.go src/4gophers.com/*/*.go
# 	go build -o ./bin/${BIN}

server: server.go src/4gophers.com/cache/server/*.go src/4gophers.com/cache/safemap/*.go
	rm -rf ./pkg/*
	go build -o ./bin/${SERVER} server.go

client: client.go src/4gophers.com/cache/client/*.go
	rm -rf ./pkg/*
	go build -o ./bin/${CLIENT} client.go

install:
	gpm install

clear:
	rm bin/${SERVER}
	rm -rf ./pkg/*