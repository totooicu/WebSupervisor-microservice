module WebSupervisor

go 1.23.0

toolchain go1.23.12

require (
	github.com/go-redis/redis/v8 v8.11.5
	golang.org/x/net v0.37.0
	streams-communication v0.0.0
)

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
)

replace streams-communication => ./streams-library
