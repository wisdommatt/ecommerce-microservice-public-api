//go:generate go run github.com/99designs/gqlgen
package graph

import (
	"github.com/opentracing/opentracing-go"
	"github.com/wisdommatt/ecommerce-microservice-public-api/grpc/proto"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Tracer               opentracing.Tracer
	UserServiceClient    proto.UserServiceClient
	ProductServiceClient proto.ProductServiceClient
	CartServiceClient    proto.CartServiceClient
}
