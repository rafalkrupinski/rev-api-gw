package main

import (
	"flag"
	"log"
	"net/http"
	"github.com/drone/go-bitbucket/oauth1"
	"github.com/elazarl/goproxy"
	"bitbucket.org/mattesilver/etsygw/handlers"
	"bitbucket.org/mattesilver/etsygw/handlers/oauth"
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

	// handle incoming https connections with MITM
	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)

	proxy.OnRequest().DoFunc(handlers.ViaIn)
	proxy.OnResponse().DoFunc(handlers.ViaOut)

	etsyReq := proxy.OnRequest(goproxy.ReqHostIs("openapi.etsy.com:443", "openapi.etsy.com"))

	// handle incoming http connections
	etsyReq.DoFunc(handlers.UpgradeToHttps)

	consumer := oauth.NewConsumer(*consumerKey, *consumerSecret)
	token := oauth1.NewAccessToken(*tokenKey, *tokenSecret, nil)
	etsyReq.Do(oauth.New(token, consumer))

	if *verbose {
		proxy.OnRequest().DoFunc(handlers.LogIn)
		proxy.OnResponse().DoFunc(handlers.LogOut)
	}

	//http.Handle("/override/", myserver{})
	http.Handle("/", proxy)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
