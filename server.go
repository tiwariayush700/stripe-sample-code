package main

import (
	`log`
	`net/http`
)

func main() {
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)

	addr := "localhost:4242"
	log.Printf("Listening on %s ...", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
