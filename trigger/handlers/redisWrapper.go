package handlers

import (
	"RuntimeAutoDeploy/common"
	"encoding/json"
	"fmt"

	log "github.com/Sirupsen/logrus"
)

func (routine *Status) AddToStatusList(traceId string, status string, firstAdd bool) {
	var (
		res         string
		err         error
		statusList  []string
		jStatusList []byte
	)
	if firstAdd {
		statusList = make([]string, 0)
		statusList = append(statusList, status)
		jStatusList, _ = json.Marshal(statusList)
		routine.redisConn.Set(fmt.Sprintf("%s-%s", common.TRACE_ID, traceId), jStatusList, 0)
		return
	}
	res, err = routine.redisConn.Get(fmt.Sprintf("%s-%s", common.TRACE_ID, traceId)).Result()
	if err != nil {
		log.WithFields(log.Fields{
			"err":     err.Error(),
			"traceId": traceId,
			"status":  status,
		}).Error("error fetching value from redis")
		return
	}
	err = json.Unmarshal([]byte(res), &statusList)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Error("error unmarshalling the retrieved value from redis")
		return
	}
	statusList = append(statusList, status)
	jStatusList, _ = json.Marshal(statusList)
	routine.redisConn.Set(fmt.Sprintf("%s-%s", common.TRACE_ID, traceId), jStatusList, 0)
}

func (routine *Status) GetStatusList(traceId string) ([]string, error) {
	var (
		res        string
		err        error
		statusList []string
	)
	res, err = routine.redisConn.Get(fmt.Sprintf("%s-%s", common.TRACE_ID, traceId)).Result()
	if err != nil {
		log.WithFields(log.Fields{
			"err":     err.Error(),
			"traceId": traceId,
		}).Error("error fetching value from redis")
		return nil, fmt.Errorf("%s", "error fetching value from redis")
	}
	err = json.Unmarshal([]byte(res), &statusList)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Error("error unmarshalling the retrieved value from redis")
		return nil, fmt.Errorf("%s", "error unmarshalling retreived value from redis")
	}
	return statusList, nil
}
