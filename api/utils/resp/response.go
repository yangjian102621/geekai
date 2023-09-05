package resp

import (
	"chatplus/core/types"
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
		c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: messages[0]})
	} else {
		c.JSON(http.StatusOK, types.BizVo{Code: types.Failed})
	}
}

func HACKER(c *gin.Context) {
	c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: "Hacker attempt!!!"})
}

func NotAuth(c *gin.Context, messages ...string) {
	if messages != nil {
		c.JSON(http.StatusOK, types.BizVo{Code: types.NotAuthorized, Message: messages[0]})
	} else {
		c.JSON(http.StatusOK, types.BizVo{Code: types.NotAuthorized, Message: "Not Authorized"})
	}
}
