package sd

import logger2 "chatplus/logger"

var logger = logger2.GetLogger()

type NotifyMessage struct {
	UserId  int    `json:"user_id"`
	JobId   int    `json:"job_id"`
	Message string `json:"message"`
}

const (
	Running  = "RUNNING"
	Finished = "FINISH"
	Failed   = "FAIL"
)
