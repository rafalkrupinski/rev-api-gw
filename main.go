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
	verbose := flag.Bool("v", false, "should every proxy request be logged to stdout")
	addr := flag.String("addr", ":8080", "proxy listen address")
	usage := flag.Bool("h", false, "print help")
	configPath := flag.String("config", "application.yaml", "path to configuration file")

	flag.Parse()

	if *usage {
		flag.Usage()
		os.Exit(0)
	}

	cfg, err := config.ReadEndpointConfig(*configPath)
	if err != nil {
		panic(err)
	}

	if *verbose {
		log.Printf("Listening on %s", *addr)
	}

	serveMux := http.NewServeMux()
	handlers.Configure(serveMux, cfg, http.DefaultTransport, *verbose)

	server := &http.Server{Addr: *addr, Handler: serveMux}
	log.Fatal(server.ListenAndServe())
}
