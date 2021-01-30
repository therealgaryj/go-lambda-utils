# go-lambda-utils

A collection of golang packages used to accelerate serverless functions deployed to AWS 

## lambdatohttp

Used with AWS Lambda to allow running a golang lambda locally with a HTTP handler.

Where no AWS services are being used other than basic API Gateway (REST) with Lambda, this lib takes the incoming API 
Gateway event, transforms to a http.HttpRequest and routes to the correct handler using Gorilla's mux.

### Usage

To run an HTTP endpoint (used for local development or to run behind, say, a load balancer):

```golang

func main() {
    r := mux.NewRouter()
    
    r.Path("foo").
    Methods(http.MethodGet).
    HandlerFunc(handlerFunc)

	http.ListenAndServe(":8080", r)
}

```

To run on AWS Lambda:

```golang

func HandleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    r := mux.NewRouter()
    
    r.Path("foo").
    Methods(http.MethodGet).
    HandlerFunc(handlerFunc)
    
    return aws.ServeRequest(r, ctx, req), nil
}

func main() {
    lambda.Start(HandleRequest)
}

```

This lib may not work for you, but that's ok. I've created it for a specific use case I required.
``