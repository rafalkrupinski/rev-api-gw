package main

import (
	"golang.org/x/net/context"
	"flag"
	"github.com/elazarl/goproxy"
	"./handlers"
	"net/http"
	"./httplog"
	"log"
	"./handlers/oauth"
	"github.com/dghubble/oauth1"
	"os"
)

func main() {
	verbose := flag.Bool("v", false, "should every proxy request be logged to stdout")
	addr := flag.String("addr", ":8080", "proxy listen address")

	consumerKey := flag.String("ck", "", "OAuth1 consuer key [required]")
	consumerSecret := flag.String("cs", "", "OAuth1 consumer secret [required]")
	tokenKey := flag.String("t", "", "OAuth token")
	tokenSecret := flag.String("ts", "", "OAuth token secret")
	usage := flag.Bool("h", false, "print help")

	flag.Parse()

	if *usage {
		flag.Usage()
		os.Exit(0)
	}

	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = *verbose

	if *consumerKey == "" {
		panic("No consumerKey")
	}
	if *consumerSecret == "" {
		panic("No consumer secret")
	}

	allRequests := proxy.OnRequest()
	allRequests.DoFunc(handlers.CleanupHandler)

	allRequests.DoFunc(handlers.ViaIn)
	proxy.OnResponse().DoFunc(handlers.ViaOut)

	isEtsy := goproxy.ReqHostIs("openapi.etsy.com:443", "openapi.etsy.com")

	proxy.OnRequest(Not(isEtsy)).DoFunc(handlers.RespondWith(goproxy.ContentTypeText, http.StatusForbidden, "Forbidden Host"))

	etsyReq := proxy.OnRequest(isEtsy)
	// handle incoming https connections with MITM
	allRequests.HandleConnect(goproxy.AlwaysMitm)

	etsyReq.DoFunc(handlers.UpgradeToHttps)

	config := oauth1.NewConfig(*consumerKey, *consumerSecret)
	token := oauth1.NewToken(*tokenKey, *tokenSecret)

	etsyReq.Do(NewClient(config, token, *verbose))

	//http.Handle("/override/", myserver{})
	//http.Handle("/", proxy)
	//log.Fatal(http.ListenAndServe(*addr, nil))

	// required to handle CONNECT
	log.Fatal(http.ListenAndServe(*addr, proxy))
}

func NewClient(config*oauth1.Config, token *oauth1.Token, logging bool) (handler goproxy.ReqHandler) {
	ctx := context.TODO()

	if logging {
		roundTripper := &httplog.LoggingRoundTripper{Super:http.DefaultTransport}
		actualClient := &http.Client{Transport:roundTripper}
		ctx = context.WithValue(ctx, oauth1.HTTPClient, actualClient)
	}

	client := config.Client(ctx, token)
	handler = oauth.New(client)

	return
}

func Not(f goproxy.ReqConditionFunc) goproxy.ReqConditionFunc {
	return func(req *http.Request, ctx *goproxy.ProxyCtx) bool {
		return !f(req, ctx)
	}
}
