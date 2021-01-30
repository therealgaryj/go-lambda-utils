package httpclient

import "encoding/base64"

func GenerateBasicAuthHeader(username string, password string) string {
	value := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	return "Basic " + value
}