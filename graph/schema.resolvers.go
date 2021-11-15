package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"html"

	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"github.com/wisdommatt/ecommerce-microservice-public-api/graph/generated"
	"github.com/wisdommatt/ecommerce-microservice-public-api/graph/mappers"
	"github.com/wisdommatt/ecommerce-microservice-public-api/graph/model"
	"github.com/wisdommatt/ecommerce-microservice-public-api/grpc/proto"
	"google.golang.org/grpc/status"
)

func (r *mutationResolver) AuthLogin(ctx context.Context, email string, password string) (*model.LoginResponse, error) {
	span := r.Tracer.StartSpan("AuthLogin")
	defer span.Finish()
	span.SetTag("param.email", email)

	email = html.EscapeString(email)
	password = html.EscapeString(password)
	if email == "" || password == "" {
		ext.Error.Set(span, true)
		span.LogFields(
			log.Error(errors.New("all fields are required")),
		)
		return nil, errors.New("all fields are required")
	}
	authResponse, err := r.UserServiceClient.LoginUser(ctx, &proto.LoginInput{Email: email, Password: password})
	if err != nil {
		return nil, parseGrpcError(err)
	}
	return mappers.ProtoLoginResponseToGql(authResponse), nil
}

func parseGrpcError(err error) error {
	return errors.New(status.Convert(err).Message())
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
