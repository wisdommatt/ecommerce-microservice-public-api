protoc:
	protoc user.proto --go-grpc_out=. --go_out=.
	protoc product.proto --go-grpc_out=. --go_out=.
	protoc cart.proto --go-grpc_out=. --go_out=.
	
run:
	go run main.go

tests:
	go test ./... -race -cover

watch:
	go install github.com/cespare/reflex@latest
	reflex -s -- sh -c 'clear && go run main.go'