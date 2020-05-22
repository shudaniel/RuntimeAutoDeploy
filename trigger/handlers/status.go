package handlers

import (
	"RuntimeAutoDeploy/common"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/redis.v5"
)

var (
	StatusRoutine *Status
)

type Status struct {
	redisConn  *redis.Client
	statusList []string
}

func RADStatusHandler(w http.ResponseWriter, r *http.Request) {
	// When the user pings this, return the status
	var (
		err        error
		statusList []string
	)
	if r.Method != "GET" {
		log.Error("error. Received incorrect HTTP method. Expecting GET")
		return
	}
	keys, OK := r.URL.Query()[common.TRACE_ID]
	if !OK || len(keys[0]) < 1 {
		log.Println("Url Param 'TRACE_ID' is missing")
		return
	}
	traceId := r.FormValue(common.TRACE_ID)

	// Fetch the status list from REDIS
	statusList, err = StatusRoutine.getStatusList(traceId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("internal error, please try again"))
		return
	}

	for i := range statusList {
		_, _ = w.Write([]byte(statusList[i]))
		_, _ = w.Write([]byte("\n"))
	}
	return
}

func StartStatusService() {
	StatusRoutine = &Status{
		redisConn:  nil,
		statusList: make([]string, 0),
	}
	StatusRoutine.redisConn = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}
