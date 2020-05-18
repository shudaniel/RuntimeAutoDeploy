package main

import (
	trigger "RuntimeAutoDeploy/Trigger/handlers"
	"RuntimeAutoDeploy/common"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

func main() {
	log.Info("RAD Service Started")
	err := trigger.Cleanup(common.GIT_BUILD_FOLDER)
	if err != nil {
		log.Error("Exiting now.")
		return
	}
	http.HandleFunc("/trigger", trigger.RADTriggerHandler)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err.Error())
	}
}
