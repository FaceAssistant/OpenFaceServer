package main

import (
    "net/http"
    "github.com/gorilla/mux"
    "log"
    "fa/web"
)

func main() {
    r := mux.NewRouter().StrictSlash(true)
    r.HandleFunc("/train", web.TrainingHandler()).Methods("Post")
    log.Fatal(http.ListenAndServe(":8082", r))
}
