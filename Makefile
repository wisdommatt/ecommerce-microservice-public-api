protoc:
	protoc proto-files/user.proto --go-grpc_out=. --go_out=.
	protoc proto-files/product.proto --go-grpc_out=. --go_out=.
	protoc proto-files/cart.proto --go-grpc_out=. --go_out=.
	
run:
	go run main.go

tests:
	go test ./... -race -cover

watch:
	go install github.com/cespare/reflex@latest
	reflex -s -- sh -c 'clear && go run main.go'

generate:
	go get github.com/99designs/gqlgen
	go generate ./...