package main

import (
	"flag"
	"log"
	"os"

	"github.com/skatiyar/byor"
)

func main() {
	var port string
	flag.StringVar(&port, "port", ":7654", "Give port in format `:7654` to connect to server")
	var help bool
	flag.BoolVar(&help, "help", false, "Show help")

	flag.Parse()

	if help {
		flag.Usage()
		os.Exit(0)
	}

	if err := byor.Client(port); err != nil {
		log.Fatalln(err)
	}
}
