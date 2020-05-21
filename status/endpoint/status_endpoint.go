package status

import (
	"RuntimeAutoDeploy/common"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

func RADStatusHandler(w http.ResponseWriter, r *http.Request) {
	// When the user pings this, return the status
	var (
		err   error
		statusList  []string
	)
	if r.Method != "GET" {
		log.Error("error. Received incorrect HTTP method. Expecting GET")
		return
	}

	keys, ok := r.URL.Query()["traceid"]    
    if !ok || len(keys[0]) < 1 {
        log.Println("Url Param 'traceid' is missing")
        return
	}

	traceId := keys[0]

	// Fetch the redis key
	res, err = routine.redisConn.Get(fmt.Sprintf("%s-%s", common.TRACE_ID, traceId)).Result()
	err = json.Unmarshal([]byte(res), &statusList)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Error("error unmarshalling the retrieved value from redis")
		return
	}

	for i := 0; i < len(statusList); i++ {
		w.Write([]byte(statusList[i]))
	}

}

