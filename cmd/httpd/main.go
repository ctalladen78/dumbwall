package main

import "net/http"

func main() {
	server := new(httpd)

	err := server.init()
	if err != nil {
		panic(err)
	}

	http.Handle("/", server.mux)

	if err = http.ListenAndServe(":1234", nil); err != nil {
		panic(err)
	}
}
