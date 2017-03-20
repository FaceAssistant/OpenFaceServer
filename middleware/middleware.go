package middleware

import (
    "net/http"
    "net/http/httputil"
    "fmt"
)

func AuthMiddleWare(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request) {
        client := &http.Client{}
        req, err := http.NewRequest("POST", "http://localhost/api/v1/users/auth", nil)
        if err != nil {
            http.Error(w, err.Error(), http.StatusUnauthorized)
            return
        }
        req.Header.Set("Authorization", r.Header.Get("Authorization"))
        resp, err := client.Do(req)
        if err != nil {
            http.Error(w, err.Error(), http.StatusUnauthorized)
            return
        }
        defer resp.Body.Close()

        r.Header.Set("Authorization", resp.Header.Get("Authorization"))
        next.ServeHTTP(w, r)
    })
}

func RequestDump(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        dump, err := httputil.DumpRequest(r, true)
        if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println(string(dump))
		next.ServeHTTP(w, r)
    })
}
