package main

import (
    "fmt"
    "log"
    "net/http"
    "RuntimeAutoDeploy/Trigger/handlers"
)


func main() {

    fmt.Println("User Preference Service Started")
    http.HandleFunc("/trigger",http.HandlerFunc(trigger.AppPreferencesHandler))
	err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatal(err.Error())
    }
}

