package httpclient

type HttpResponse struct {
	StatusCode int
	StatusText string
	Headers map[string][]string
	Body string
}