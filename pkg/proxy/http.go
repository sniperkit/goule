package proxy

import (
	"io"
	"net/http"
)

func ProxyHTTP(context *Context, client *http.Client) {
	if code, err := proxyHTTPInternal(context, client); err != nil {
		context.Response.Header().Set("Content-Type", "text/plain")
		context.Response.WriteHeader(code)
		context.Response.Write([]byte(err.Error()))
	}
}

func proxyHTTPInternal(context *Context, client *http.Client) (int, error) {
	req, err := http.NewRequest(context.Request.Method,
		context.DestURL.String(), context.Request.Body)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	for header, value := range RequestHeaders(context) {
		req.Header[header] = value
	}
	res, err := client.Do(req)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	for header, value := range ResponseHeaders(context, res.Header) {
		context.Response.Header()[header] = value
	}
	context.Response.WriteHeader(res.StatusCode)
	io.Copy(context.Response, res.Body)
	res.Body.Close()

	// I don't think I need to close context.Response because it is
	// automatically closed by the HTTP server.
	return 0, nil
}
