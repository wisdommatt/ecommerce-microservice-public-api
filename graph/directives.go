package graph

import (
	"context"
	"errors"

	"github.com/99designs/gqlgen/graphql"
	"github.com/wisdommatt/ecommerce-microservice-public-api/grpc/proto"
)

type contextKey string

var userContextKey contextKey = "user-graphql-key"

type DirectiveFunc func(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error)

func IsAuthenticated(jwtContextKey interface{}, userService proto.UserServiceClient) DirectiveFunc {
	return func(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
		jwtTokenI := ctx.Value(jwtContextKey)
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
