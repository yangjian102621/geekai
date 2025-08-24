package midware

import (
	"context"
	"fmt"
	"geekai/core/types"
	"geekai/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

// RateLimitEvery 使用 Redis 做固定间隔限流：在 interval 内仅允许一次请求
// Key 优先使用登录用户ID，若没有则退化为 route + IP
func RateLimitEvery(redisClient *redis.Client, interval time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		keyID := ""
		if userID, ok := c.Get(types.LoginUserID); ok {
			keyID = fmt.Sprintf("user:%s", utils.InterfaceToString(userID))
		} else {
			keyID = fmt.Sprintf("ip:%s", c.ClientIP())
		}

		fullPath := c.FullPath()
		if fullPath == "" {
			fullPath = c.Request.URL.Path
		}
		key := fmt.Sprintf("rl:%s:%s", fullPath, keyID)

		okSet, err := redisClient.SetNX(context.Background(), key, 1, interval).Result()
		if err != nil {
			// Redis 异常时放行，避免误伤可用性
			return
		}
		if !okSet {
			c.JSON(http.StatusTooManyRequests, types.BizVo{Code: types.Failed, Message: "请求过于频繁，请稍后重试"})
			c.Abort()
			return
		}
	}
}
