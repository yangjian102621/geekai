package sd

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import logger2 "geekai/logger"

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
