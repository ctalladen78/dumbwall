package main

import (
	"flag"
	"net/http"
)

var (
	version string
	commit  string
	branch  string
)

func init() {
	if commit == "" {
		commit = "unknown"
	}

	if branch == "" {
		branch = "unknown"
	}
}

var (
	etcPath    = flag.String("conf", "/etc/dumbwall", "configuration directory")
	listenAddr = flag.String("listen", ":80", "web server address")
)

func main() {
	flag.Parse()

	server := new(httpd)

	err := server.init(*etcPath)
	if err != nil {
		panic(err)
	}

	http.Handle("/", server.mux)

	if err = http.ListenAndServe(*listenAddr, nil); err != nil {
		panic(err)
	}
}
