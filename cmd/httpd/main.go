package main

import (
	"flag"
	"net/http"
)

var (
	etcPath = flag.String("conf", "/etc/dumbwall", "configuration directory")
)

func main() {
	flag.Parse()

	server := new(httpd)

	err := server.init(*etcPath)
	if err != nil {
		panic(err)
	}

	http.Handle("/", server.mux)
	println("starting web server")

	if err = http.ListenAndServe(":1234", nil); err != nil {
		panic(err)
	}
}
