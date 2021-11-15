//go:generate go run github.com/99designs/gqlgen
package graph

import (
	"errors"

	"github.com/opentracing/opentracing-go"
	"github.com/wisdommatt/ecommerce-microservice-public-api/grpc/proto"
	"google.golang.org/grpc/status"
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

func parseGrpcError(err error) error {
	return errors.New(status.Convert(err).Message())
}
