package mappers

import (
	"github.com/wisdommatt/ecommerce-microservice-public-api/graph/model"
	"github.com/wisdommatt/ecommerce-microservice-public-api/grpc/proto"
)

func ProtoUserToGql(user *proto.User) *model.User {
	return &model.User{
		ID:       user.Id,
		FullName: user.FullName,
		Email:    user.Email,
		Country:  user.Country,
	}
}

func ProtoLoginResponseToGql(res *proto.LoginResponse) *model.LoginResponse {
	return &model.LoginResponse{
		JwtToken: res.JwtToken,
		User:     ProtoUserToGql(res.User),
	}
}
