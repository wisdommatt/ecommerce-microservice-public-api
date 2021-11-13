package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/joho/godotenv"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/wisdommatt/ecommerce-microservice-public-api/graph"
	"github.com/wisdommatt/ecommerce-microservice-public-api/graph/generated"
)

const defaultPort = "1212"

func main() {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{PrettyPrint: true})
	log.SetReportCaller(true)
	log.SetOutput(os.Stdout)

	mustLoadDotenv(log)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{
		Tracer: initTracer("graphql-api"),
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
