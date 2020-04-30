package main

import (
    "fmt"
    "log"
    "net/http"
)

// AppPreferencesHandler handles the requests from the apps
func AppPreferencesHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        // ctx, _ := context.WithCancel(r.Context())
        fmt.Println("GET")
    }
    w.Write([]byte("hello"))
    return
}

func main() {

    fmt.Println("User Preference Service Started")
    http.HandleFunc("/trigger",http.HandlerFunc(AppPreferencesHandler))
	err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatal(err.Error())
    }
}








/////////////////////////////////////////




