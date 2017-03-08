package main

import (
    "net/http"
    "github.com/gorilla/mux"
    "log"
)

func main() {
    r := mux.NewRouter().StrictSlash(true)
    r.HandleFUnc("/train",).Methods("Post")
    log.Fatal(http.ListenAndServe(":8082", r))
}
