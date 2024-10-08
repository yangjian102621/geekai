package resp

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"geekai/core/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SUCCESS(c *gin.Context, values ...interface{}) {
	if values != nil {
		c.JSON(http.StatusOK, types.BizVo{Code: types.Success, Data: values[0]})
	} else {
		c.JSON(http.StatusOK, types.BizVo{Code: types.Success})
	}

}

func ERROR(c *gin.Context, messages ...string) {
	if messages != nil {
		c.JSON(http.StatusBadRequest, types.BizVo{Code: types.Failed, Message: messages[0]})
	} else {
		c.JSON(http.StatusBadRequest, types.BizVo{Code: types.Failed})
	}
}

func HACKER(c *gin.Context) {
	c.JSON(http.StatusBadRequest, types.BizVo{Code: types.Failed, Message: "Hacker attempt!!!"})
}

func NotAuth(c *gin.Context, messages ...string) {
	if messages != nil {
		c.JSON(http.StatusUnauthorized, types.BizVo{Code: types.NotAuthorized, Message: messages[0]})
	} else {
		c.JSON(http.StatusUnauthorized, types.BizVo{Code: types.NotAuthorized, Message: "Not Authorized"})
	}
}
