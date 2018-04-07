package main

import (
	"flag"
	"github.com/rafalkrupinski/rev-api-gw/config"
	"github.com/rafalkrupinski/rev-api-gw/handlers"
	"log"
	"net/http"
	"os"
)

func main() {
	appConfig, usage, err := config.ParseFlags()
	if err != nil {
		panic(err)
	}

	if usage {
		flag.Usage()
		os.Exit(0)
	}

	appConfig.Dump()

	serveMux := http.NewServeMux()
	handlers.Configure(appConfig, serveMux, http.DefaultTransport)

	server := &http.Server{Addr: appConfig.Address, Handler: serveMux}
	log.Fatal(server.ListenAndServe())
}
