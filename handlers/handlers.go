package handlers

import (
	ht "github.com/rafalkrupinski/rev-api-gw/http"
	"github.com/rafalkrupinski/rev-api-gw/httplog"
	"io"
	"log"
	"net/http"
)

func CleanupHandler(req *http.Request) (*http.Request, *http.Response, error) {
	ht.CleanupRequest(req)
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
		response = errorToResponse(err)
	}

	if response == nil {
		log.Panicf("Handler chain didn't return response for %s", req.URL)
	}

	if request == nil {
		request = req
	}

	response = c.HandleResponse(req, response)

	if err != nil {
		http.Error(to, err.Error(), 500)
	} else {
		writeResponse(response, to)
	}
}

func errorToResponse(err error) (resp *http.Response) {
	resp = &http.Response{Header: http.Header{}, StatusCode: 500}
	resp.Header.Set("Content-Type", "text/plain; charset=utf-8")
	resp.Header.Set("X-Content-Type-Options", "nosniff")
	resp.Body = httplog.NewReadCloserFromString(err.Error())
	return
}

func writeResponse(resp *http.Response, w http.ResponseWriter) {
	origBody := resp.Body
	defer origBody.Close()
	log.Printf("Copying response to client %v [%d]", resp.Status, resp.StatusCode)
	// http.ResponseWriter will take care of filling the correct response length
	// Setting it now, might impose wrong value, contradicting the actual new
	// body the user returned.
	// We keep the original body to remove the header only if things changed.
	// This will prevent problems with HEAD requests where there's no body, yet,
	// the Content-Length header should be set.
	if origBody != resp.Body {
		resp.Header.Del("Content-Length")
	}
	copyHeaders(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	nr, err := io.Copy(w, resp.Body)
	if err := resp.Body.Close(); err != nil {
		log.Printf("Can't close response body %v", err)
	}
	log.Printf("Copied %v bytes to client error=%v", nr, err)
}

func copyHeaders(dst, src http.Header) {
	for k, _ := range dst {
		dst.Del(k)
	}
	for k, vs := range src {
		for _, v := range vs {
			dst.Add(k, v)
		}
	}
}
