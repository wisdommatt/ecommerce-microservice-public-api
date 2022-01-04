# Ecommerce Public Graphql API

This is the public graphql api for interacting with other internal services in the ecommerce application.

## Applications

* ## [Jaeger](https://www.jaegertracing.io/)

  * Jaeger is an **open source software for tracing transactions between distributed services**.
* ## [GRPC](https://grpc.io)

  * GRPC is an open source remote procedure call system.
  * The public api interacts with other internal services using grpc.

### Usage

To install / run the user microservice run the command below:

```bash
docker-compose up
```

## Requirements

The application requires the following:

* Go (v 1.5+)
* Docker (v3+)
* Docker Compose

## Other Services / Resources

* #### [Product Service](https://github.com/wisdommatt/ecommerce-microservice-product-service)
* #### [Notification Service](https://github.com/wisdommatt/ecommerce-microservice-notification-service)
* #### [Cart Service](https://github.com/wisdommatt/ecommerce-microservice-cart-service)
* #### [Shared](https://github.com/wisdommatt/ecommerce-microservice-shared)
* #### [User Service](https://github.com/wisdommatt/ecommerce-microservice-cart-service)
