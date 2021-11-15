package graph

import (
	"html"
	"strings"

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

func GqlNewUserToProto(user *model.NewUser) *proto.NewUser {
	return &proto.NewUser{
		Email:    html.EscapeString(strings.ToLower(user.Email)),
		FullName: html.EscapeString(user.FullName),
		Password: user.Password,
		Country:  html.EscapeString(user.Country),
	}
}

func GqlNewProductToProto(product *model.NewProduct) *proto.NewProduct {
	return &proto.NewProduct{
		Name:        html.EscapeString(product.Name),
		Category:    html.EscapeString(product.Category),
		Description: html.EscapeString(product.Description),
		Brand:       html.EscapeString(unpointStr(product.Brand)),
		Price:       product.Price,
		ImageUrl:    html.EscapeString(product.ImageURL),
	}
}

func ProtoProductToGql(product *proto.Product) *model.Product {
	return &model.Product{
		Sku:         product.Sku,
		Name:        product.Name,
		Category:    product.Category,
		Description: product.Description,
		Brand:       product.Brand,
		Price:       product.Price,
		ImageURL:    product.ImageUrl,
	}
}

func GqlNewCartItemToProto(item *model.NewCartItem) *proto.NewCartItem {
	return &proto.NewCartItem{
		ProductSku: html.EscapeString(item.ProductSku),
		Quantity:   int32(item.Quantity),
	}
}

func ProtoCartItemToGql(item *proto.CartItem) *model.CartItem {
	return &model.CartItem{
		ID:         item.Id,
		ProductSku: item.ProductSku,
		Quantity:   int(item.Quantity),
	}
}

func unpointStr(str *string) string {
	if str == nil {
		return ""
	}
	return *str
}
