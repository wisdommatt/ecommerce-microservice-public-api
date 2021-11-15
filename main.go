package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
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

type contextKey string

const defaultPort = "1212"

var jwtContextKey contextKey = "jwt-context-key"

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

	config := generated.Config{Resolvers: &graph.Resolver{
		Tracer:               initTracer("graphql-api"),
		UserServiceClient:    userServiceClient,
		ProductServiceClient: productServiceClient,
		CartServiceClient:    cartServiceClient,
	}}
	config.Directives.IsAuthenticated = graph.IsAuthenticated(jwtContextKey, userServiceClient)
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(config))

	router := chi.NewRouter()
	router.Use(addJwtToHTTPContext)
	router.Handle("/graphql/", playground.Handler("GraphQL playground", "/graphql/query"))
	router.Handle("/graphql/query", srv)
	httpServer := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(httpServer.ListenAndServe())
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

func addJwtToHTTPContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")
		jwtToken := strings.ReplaceAll(authorizationHeader, "Bearer ", "")
		ctx := context.WithValue(r.Context(), jwtContextKey, jwtToken)
		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}
