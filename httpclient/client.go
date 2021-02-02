package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/therealgaryj/go-lambda-utils/log"
	"io/ioutil"
	"net/http"
	"strings"
)

func HttpGet(url string, headers map[string]string) (HttpResponse, error) {

	request, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Error(fmt.Sprintf("%v", err))

		return HttpResponse{}, err
	}

	request.Header.Set("Accept", "application/json; charset=utf-8")

	addHeadersToRequest(request, headers)

	return executeRequest(request)
}

func HttpPut(url string, headers map[string]string, payload interface{}) (HttpResponse, error) {
	return httpSendWithBody("PUT", url, headers, payload)
}
func HttpPost(url string, headers map[string]string, payload interface{}) (HttpResponse, error) {
	return httpSendWithBody("POST", url, headers, payload)
}
func httpSendWithBody(method string, url string, headers map[string]string, payload interface{}) (HttpResponse, error) {

	requestBody, error := json.Marshal(payload)

	if error != nil {
		log.Error(error.Error())

		return HttpResponse{}, error
	}

	request, createRequestError := http.NewRequest(method, url, bytes.NewBuffer(requestBody))

	if createRequestError != nil {
		log.Error(createRequestError.Error())
	}

	request.Header.Set("Accept", "application/json;charset=utf-8")
	request.Header.Set("Content-Type", "application/json; charset=utf-8")

	addHeadersToRequest(request, headers)

	return executeRequest(request)
}

func addHeadersToRequest(request *http.Request, headersToAdd map[string]string) {
	for header, value := range headersToAdd {
		request.Header.Add(header, value)
	}
}

func executeRequest(request *http.Request) (HttpResponse, error) {
	log.Info("-------")
	log.Info(fmt.Sprintf("Sending %s request to: %s", request.Method, request.URL.Path))
	debugLogRequest(request)

	client := &http.Client{}

	response, executeRequestError := client.Do(request)
	if executeRequestError != nil {
		return HttpResponse{}, executeRequestError
	}

	defer response.Body.Close()

	rawBody, readResponseError := ioutil.ReadAll(response.Body)
	if readResponseError != nil {
		return HttpResponse{}, readResponseError
	}

	debugLogResponse(response)

	log.Info("-------")

	return HttpResponse{
		StatusCode: response.StatusCode,
		StatusText: response.Status,
		Headers:    response.Header,
		Body:       string(rawBody),
	}, nil
}

// DebugLog Log the full request
func debugLogRequest(request *http.Request) {

	log.Debug(fmt.Sprintf("%s %s", request.Method, request.URL.String()))

	for headerName, headerValue := range request.Header {
		formattedHeaderValue := strings.Join(headerValue, ", ")
		prefix := ""
		sanitizedHeaderValue := formattedHeaderValue

		if strings.ToLower(headerName) == "authorization" {
			if authType := string([]rune(strings.ToLower(formattedHeaderValue))[0:5]); authType == "basic" {
				prefix = "Basic "
			} else if authType := string([]rune(strings.ToLower(formattedHeaderValue))[0:6]); authType =="bearer" {
				prefix = "Bearer "
			}

			sanitizedHeaderValue = "<sanitized>"
		}

		log.Debug(fmt.Sprintf("%s: %s%s", headerName, prefix, sanitizedHeaderValue))
	}

	log.Debug("<Body Omitted>")
}

// DebugLog Log the full response
func debugLogResponse(response *http.Response) {

	log.Debug(response.Status)

	for headerName, headerValue := range response.Header {
		log.Debug(fmt.Sprintf("%s: %s", headerName, headerValue))
	}

	log.Debug("<Body Omitted>")
}
