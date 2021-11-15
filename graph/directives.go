package graph

import (
	"context"
	"errors"

	"github.com/99designs/gqlgen/graphql"
	"github.com/wisdommatt/ecommerce-microservice-public-api/grpc/proto"
)

type ContextKey string

var (
	userContextKey ContextKey = "user-graphql-key"
	JwtContextKey  ContextKey = "jwt-context-key"
)

type DirectiveFunc func(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error)

func IsAuthenticated(userService proto.UserServiceClient) DirectiveFunc {
	return func(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
		jwtTokenI := ctx.Value(JwtContextKey)
		if jwtTokenI == nil {
			return nil, errors.New("you are not authenticated")
		}
		jwtToken := jwtTokenI.(string)
		authUser, err := userService.GetUserFromJWT(ctx, &proto.GetUserFromJWTInput{JwtToken: jwtToken})
		if err != nil {
			return nil, parseGrpcError(err)
		}
		usr := ProtoUserToGql(authUser.User)
		ctx = context.WithValue(ctx, userContextKey, usr)
		return next(ctx)
	}
}
