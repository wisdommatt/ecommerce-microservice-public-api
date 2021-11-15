package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"html"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"github.com/wisdommatt/ecommerce-microservice-public-api/graph/generated"
	"github.com/wisdommatt/ecommerce-microservice-public-api/graph/model"
	"github.com/wisdommatt/ecommerce-microservice-public-api/grpc/proto"
	"google.golang.org/grpc/metadata"
)

func (r *cartItemResolver) Product(ctx context.Context, obj *model.CartItem) (*model.Product, error) {
	// TODO(wisdommatt): use dataloader to retrieve products
	return r.Query().GetProduct(ctx, obj.ProductSku)
}

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

	if input.FullName == "" || input.Country == "" || input.Password == "" || input.Email == "" {
		ext.Error.Set(span, true)
		span.LogFields(
			log.Error(errors.New("all fields are required")),
		)
		return nil, errors.New("all fields are required")
	}
	span.LogFields(
		log.Object("param.input", input),
	)
	newUser, err := r.UserServiceClient.CreateUser(ctx, GqlNewUserToProto(&input))
	if err != nil {
		return nil, parseGrpcError(err)
	}
	return ProtoUserToGql(newUser), nil
}

func (r *mutationResolver) AddNewProduct(ctx context.Context, input model.NewProduct) (*model.Product, error) {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, r.Tracer, "AddNewProduct")
	defer span.Finish()

	if input.Name == "" || input.Category == "" || input.Description == "" || input.ImageURL == "" || input.Price == 0 {
		ext.Error.Set(span, true)
		span.LogFields(
			log.Error(errors.New("all fields are required")),
		)
		return nil, errors.New("all fields except brand are required")
	}
	span.LogFields(
		log.Object("param.input", input),
	)
	metaData := metadata.New(map[string]string{
		"Authorization": ctx.Value(JwtContextKey).(string),
	})
	ctx = metadata.NewOutgoingContext(ctx, metaData)
	newProduct, err := r.ProductServiceClient.AddProduct(ctx, GqlNewProductToProto(&input))
	if err != nil {
		return nil, parseGrpcError(err)
	}
	return ProtoProductToGql(newProduct), nil
}

func (r *mutationResolver) AddToCart(ctx context.Context, input model.NewCartItem) (*model.CartItem, error) {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, r.Tracer, "AddToCart")
	defer span.Finish()

	if input.ProductSku == "" || input.Quantity == 0 {
		ext.Error.Set(span, true)
		span.LogFields(
			log.Error(errors.New("all fields are required")),
		)
		return nil, errors.New("all fields are required")
	}
	metaData := metadata.New(map[string]string{
		"Authorization": ctx.Value(JwtContextKey).(string),
	})
	ctx = metadata.NewOutgoingContext(ctx, metaData)
	newCartItem, err := r.CartServiceClient.AddToCart(ctx, GqlNewCartItemToProto(&input))
	if err != nil {
		return nil, parseGrpcError(err)
	}
	return ProtoCartItemToGql(newCartItem), nil
}

func (r *queryResolver) GetProduct(ctx context.Context, sku string) (*model.Product, error) {
	if sku == "" {
		return nil, errors.New("sku cannot be empty")
	}
	product, err := r.ProductServiceClient.GetProduct(ctx, &proto.GetProductInput{Sku: sku})
	if err != nil {
		return nil, parseGrpcError(err)
	}
	return ProtoProductToGql(product), nil
}

// CartItem returns generated.CartItemResolver implementation.
func (r *Resolver) CartItem() generated.CartItemResolver { return &cartItemResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type cartItemResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
