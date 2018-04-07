package handlers

import (
	"github.com/rafalkrupinski/rev-api-gw/morego/morehttp"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func CleanupHandler(req *http.Request) (*http.Request, *http.Response, error) {
	addr := req.URL

	if req.RequestURI != "" {
		addr, _ = url.Parse(req.RequestURI)
		req.RequestURI = ""
	}

	if addr.Host == "" && req.Host != "" {
		addr.Host = req.Host
	}

	if addr.Scheme == morehttp.SCHEME_HTTP {
		strings.TrimSuffix(addr.Host, ":"+strconv.Itoa(morehttp.PORT_HTTP))
	} else if addr.Scheme == morehttp.SCHEME_HTTPS {
		strings.TrimSuffix(addr.Host, ":"+strconv.Itoa(morehttp.PORT_HTTPS))
	}
	return nil, nil, nil
}

type RequestHandler interface {
	HandleRequest(*http.Request) (*http.Request, *http.Response, error)
}

type RequestHandlerFunc func(*http.Request) (*http.Request, *http.Response, error)

func (f RequestHandlerFunc) HandleRequest(r *http.Request) (*http.Request, *http.Response, error) {
	return f(r)
}

type ResponseHandler interface {
	HandleResponse(*http.Request, *http.Response) *http.Response
}

type ResponseHandlerFunc func(*http.Request, *http.Response) *http.Response

func (f ResponseHandlerFunc) HandleResponse(req *http.Request, resp *http.Response) *http.Response {
	return f(req, resp)
}

type HandlerChain struct {
	RequestHandlers  []RequestHandler
	ResponseHandlers []ResponseHandler
}

func (c *HandlerChain) AddRequestHandlerFunc(handlerFunc RequestHandlerFunc) {
	c.RequestHandlers = append(c.RequestHandlers, handlerFunc)
}

func (c *HandlerChain) AddRequestHandler(handler RequestHandler) {
	c.RequestHandlers = append(c.RequestHandlers, handler)
}

func (c *HandlerChain) AddResponseHandlerFunc(handlerFunc ResponseHandlerFunc) {
	c.ResponseHandlers = append(c.ResponseHandlers, handlerFunc)
}

func (c *HandlerChain) AddResponseHandler(handler ResponseHandler) {
	c.ResponseHandlers = append(c.ResponseHandlers, handler)
}

func (c *HandlerChain) HandleRequest(req *http.Request) (*http.Request, *http.Response, error) {
	currentRequest := req
	for _, h := range c.RequestHandlers {
		request, response, err := h.HandleRequest(currentRequest)

		if request != nil {
			currentRequest = request
		}
		if err != nil || response != nil {
			return currentRequest, response, err
		}
	}
	return currentRequest, nil, nil
}

func (c *HandlerChain) HandleResponse(req *http.Request, resp *http.Response) *http.Response {
	currentResponse := resp
	for _, h := range c.ResponseHandlers {
		response := h.HandleResponse(req, currentResponse)
		if response != nil {
			currentResponse = response
		}
	}
	return currentResponse
}

func (c *HandlerChain) ServeHTTP(to http.ResponseWriter, req *http.Request) {
	request, response, err := c.HandleRequest(req)

	if err != nil {
		http.Error(to, err.Error(), 500)
		return
	}

	if response == nil {
		log.Printf("request handler chain didn't return response for %s", req.URL)
		http.Error(to, "nil response", 500)
	}

	if request == nil {
		request = req
	}

	response = c.HandleResponse(req, response)
	if response == nil {
		http.Error(to, "response handler chain returned nil response", 500)
		return
	}

	writeResponse(response, to)
}

func writeResponse(resp *http.Response, w http.ResponseWriter) {
	log.Printf("Copying response to client %v [%d]", resp.Status, resp.StatusCode)
	copyHeaders(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	nr, err := io.Copy(w, resp.Body)
	if err := resp.Body.Close(); err != nil {
		log.Printf("Can't close response body %v", err)
	}
	log.Printf("Copied %v bytes to client error=%v", nr, err)
}

func copyHeaders(dst, src http.Header) {
	for k := range dst {
		dst.Del(k)
	}
	for k, vs := range src {
		for _, v := range vs {
			dst.Add(k, v)
		}
	}
}
