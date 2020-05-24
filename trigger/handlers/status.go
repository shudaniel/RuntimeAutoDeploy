package handlers

import (
	"RuntimeAutoDeploy/common"
	"net/http"

	log "github.com/sirupsen/logrus"
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
	log.Info("received status request")
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
	statusList, err = common.GetStatusList(traceId)
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
	common.ConnectToRedis()
	StatusRoutine = &Status{
		redisConn:  common.RedisConn,
		statusList: make([]string, 0),
	}
}
