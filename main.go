package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/joho/godotenv"
	otgrpc "github.com/opentracing-contrib/go-grpc"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/wisdommatt/ecommerce-microservice-public-api/graph"
	"github.com/wisdommatt/ecommerce-microservice-public-api/graph/generated"
	"github.com/wisdommatt/ecommerce-microservice-public-api/grpc/proto"
	"google.golang.org/grpc"
)

const defaultPort = "1212"

func main() {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{PrettyPrint: true})
	log.SetReportCaller(true)
	log.SetOutput(os.Stdout)

	mustLoadDotenv(log)

	tracer := initTracer("graphql-api")
	opentracing.SetGlobalTracer(tracer)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	userServiceClientConn, err := grpc.Dial(
		os.Getenv("USER_SERVICE_ADDR"),
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(tracer)),
		grpc.WithStreamInterceptor(otgrpc.OpenTracingStreamClientInterceptor(tracer)),
	)
	if err != nil {
		log.WithError(err).WithField("userServiceAddr", os.Getenv("USER_SERVICE_ADDR")).
			Fatal("an error occured while connecting to user service")
	}
	userServiceClient := proto.NewUserServiceClient(userServiceClientConn)

	productServiceClientConn, err := grpc.Dial(
		os.Getenv("PRODUCT_SERVICE_ADDR"),
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(tracer)),
		grpc.WithStreamInterceptor(otgrpc.OpenTracingStreamClientInterceptor(tracer)),
	)
	if err != nil {
		log.WithError(err).WithField("productServiceAddr", os.Getenv("PRODUCT_SERVICE_ADDR")).
			Fatal("an error occured while connecting to product service")
	}
	productServiceClient := proto.NewProductServiceClient(productServiceClientConn)

	cartServiceClientConn, err := grpc.Dial(
		os.Getenv("CART_SERVICE_ADDR"),
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(tracer)),
		grpc.WithStreamInterceptor(otgrpc.OpenTracingStreamClientInterceptor(tracer)),
	)
	if err != nil {
		log.WithError(err).WithField("cartServiceAddr", os.Getenv("CART_SERVICE_ADDR")).
			Fatal("an error occured while connecting to cart service")
	}
	cartServiceClient := proto.NewCartServiceClient(cartServiceClientConn)

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{
		Tracer:               initTracer("graphql-api"),
		UserServiceClient:    userServiceClient,
		ProductServiceClient: productServiceClient,
		CartServiceClient:    cartServiceClient,
	}}))

	http.Handle("/graphql/", playground.Handler("GraphQL playground", "/graphql/query"))
	http.Handle("/graphql/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func mustLoadDotenv(log *logrus.Logger) {
	err := godotenv.Load(".env", ".env-defaults")
	if err != nil {
		log.WithError(err).Fatal("Unable to load env files")
	}
}

func initTracer(serviceName string) opentracing.Tracer {
	return initJaegerTracer(serviceName)
}

func initJaegerTracer(serviceName string) opentracing.Tracer {
	cfg := &config.Configuration{
		ServiceName: serviceName,
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
	}
	tracer, _, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
	if err != nil {
		log.Fatal("ERROR: cannot init Jaeger", err)
	}
	return tracer
}
