package main

import (
    "net/http"
    "github.com/gorilla/mux"
    "github.com/justinas/alice"
    "log"
    "fa/web"
    "fa/middleware"
)

func main() {
    r := mux.NewRouter().StrictSlash(true)
    r.HandleFunc("/train", web.TrainingHandler()).Methods("Post")
    r.HandleFunc("/infer", web.FaceRecogHandler()).Methods("Post")
    a := alice.New(middleware.AuthMiddleWare).Then(r)
    log.Fatal(http.ListenAndServe(":8081", a))
    //log.Fatal(http.ListenAndServe(":8081", r))
}
