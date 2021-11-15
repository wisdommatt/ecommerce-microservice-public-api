package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"html"
	"strings"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"github.com/wisdommatt/ecommerce-microservice-public-api/graph/generated"
	"github.com/wisdommatt/ecommerce-microservice-public-api/graph/model"
	"github.com/wisdommatt/ecommerce-microservice-public-api/grpc/proto"
)

func (r *mutationResolver) AuthLogin(ctx context.Context, email string, password string) (*model.LoginResponse, error) {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, r.Tracer, "AuthLogin")
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
	return ProtoLoginResponseToGql(authResponse), nil
}

func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (*model.User, error) {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, r.Tracer, "CreateUser")
	defer span.Finish()
	span.LogFields(
		log.Object("param.input", input),
	)
	// validating the user input.
	input.FullName = html.EscapeString(input.FullName)
	input.Country = html.EscapeString(input.Country)
	input.Password = html.EscapeString(input.Password)
	input.Email = html.EscapeString(strings.ToLower(input.Email))
	if input.FullName == "" || input.Country == "" || input.Password == "" || input.Email == "" {
		ext.Error.Set(span, true)
		span.LogFields(
			log.Error(errors.New("all fields are required")),
		)
		return nil, errors.New("all fields are required")
	}
	newUser, err := r.UserServiceClient.CreateUser(ctx, GqlNewUserToProto(&input))
	if err != nil {
		return nil, parseGrpcError(err)
	}
	return ProtoUserToGql(newUser), nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
