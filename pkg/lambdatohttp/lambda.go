package lambdatohttp

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func ServeRequest(router *mux.Router, ctx context.Context, req events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	customHttpResponse := customHttpResponse{
		Status:  0,
		Headers: http.Header{},
		Body:    nil,
	}

	router.ServeHTTP(customHttpResponseWriter{
		response: &customHttpResponse,
		}, createHttpRequest(req))


	return events.APIGatewayProxyResponse{
		StatusCode:        determineResponseStatus(customHttpResponse),
		Headers:           flattenResponseHeaders(customHttpResponse),
		MultiValueHeaders: customHttpResponse.Headers,
		Body:              string(customHttpResponse.Body),
		IsBase64Encoded:   false,
	}
}

func determineResponseStatus(res customHttpResponse) int {
	status := 200
	if res.Status > 0 {
		status = res.Status
	}

	return status
}

func flattenResponseHeaders(res customHttpResponse) map[string]string {
	singleValueHeaders := make(map[string]string)

	for header, value := range res.Headers {
		if len(value) > 1 {
			singleValueHeaders[header] = strings.Join(value, ",")
		}
	}

	return singleValueHeaders
}

func createHttpRequest(req events.APIGatewayProxyRequest) *http.Request {
	return &http.Request{
		Method:           req.HTTPMethod,
		URL:              generateUrl(req),
		Proto:            req.RequestContext.Protocol,
		ProtoMajor:       1,
		ProtoMinor:       1,
		Header:           req.MultiValueHeaders,
		Body:             requestBodyReader{body: []byte(req.Body)}, // reader?!
		GetBody:          nil,
		ContentLength:    int64(len([]byte(req.Body))),
		TransferEncoding: nil,
		Close:            false,
		Host:             req.RequestContext.DomainName,
		Form:             nil,
		PostForm:         nil,
		MultipartForm:    nil,
		Trailer:          nil,
		RemoteAddr:       "",
		RequestURI:       "",
		TLS:              nil,
		Response:         nil,
	}
}

func generateQueryString(req events.APIGatewayProxyRequest) string {
	queryParameters := make([]string, 0)

	for param, val := range req.MultiValueQueryStringParameters {
		for _, nextVal := range val {
			queryParameters = append(queryParameters, param+"="+nextVal)
		}
	}

	return strings.Join(queryParameters, "&")
}

func generateUrl(req events.APIGatewayProxyRequest) *url.URL {
	return &url.URL{
		Scheme:      req.Headers["X-Forwarded-Proto"],
		Opaque:      "",
		User:        nil,
		Host:        req.RequestContext.DomainName,
		Path:        req.Path,
		RawPath:     req.Path,
		ForceQuery:  false,
		RawQuery:    generateQueryString(req),
		Fragment:    "",
		RawFragment: "",
	}
}

type customHttpResponse struct {
	Status int
	Headers http.Header
	Body []byte
}
type customHttpResponseWriter struct {
	response *customHttpResponse
}
func (w customHttpResponseWriter) Header() http.Header {
	return w.response.Headers
}

func (w customHttpResponseWriter) Write(content []byte) (int, error) {
	w.response.Body = append(w.response.Body, content...)

	return len(content), nil
}
func (w customHttpResponseWriter) WriteHeader(statusCode int) {
	w.response.Status = statusCode
}

type requestBodyReader struct {
	body []byte
}

func (r requestBodyReader) Read(p []byte) (int, error) {
	remaining := len(r.body)

	if remaining == 0 {
		return 0, io.EOF
	}

	if remaining > len(p) {
		for x:=0;x<len(p);x++ {
			p[x] = r.body[x]
		}
		r.body = r.body[len(p) - 1:]
		return len(p), nil
	} else {
		for x:=0;x<remaining;x++ {
			p[x] = r.body[x]
		}
		return remaining, nil
	}
}
func (r requestBodyReader) Close() error {
	return nil
}
