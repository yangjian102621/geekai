package handler

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"geekai/core"
	"geekai/core/middleware"
	"geekai/service"
	"geekai/store/model"
	"geekai/utils/resp"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ModelAPIHandler struct {
	BaseHandler
	userService *service.UserService
}

const modelTokenPrefix = "sk-thinkbuddy-"

func NewModelAPIHandler(app *core.AppServer, db *gorm.DB, users *service.UserService) *ModelAPIHandler {
	return &ModelAPIHandler{BaseHandler: BaseHandler{App: app, DB: db}, userService: users}
}

func (h *ModelAPIHandler) RegisterRoutes() {
	account := h.App.Engine.Group("/api/user/model-token", middleware.UserAuthMiddleware(h.App.Config.Session.SecretKey, h.App.Redis))
	account.GET("", h.TokenStatus)
	account.POST("", h.RotateToken)
	account.DELETE("", h.RevokeToken)

	api := h.App.Engine.Group("/v1", h.authenticate)
	api.GET("/models", h.Models)
}

func modelAPIError(c *gin.Context, status int, code, message string) {
	c.AbortWithStatusJSON(status, gin.H{"error": gin.H{"message": message, "type": code, "code": code, "param": nil}})
}

func modelTokenHash(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func modelUserAvailable(user model.User) bool {
	return user.Status && (user.ExpiredTime <= 0 || user.ExpiredTime > time.Now().Unix())
}

func (h *ModelAPIHandler) tokenUser(c *gin.Context) (model.User, bool) {
	c.Header("Cache-Control", "no-store")
	user, err := h.GetLoginUser(c)
	if err != nil || !modelUserAvailable(user) {
		resp.ERROR(c, "账号不可用，请重新登录或联系管理员")
		return user, false
	}
	return user, true
}

func (h *ModelAPIHandler) TokenStatus(c *gin.Context) {
	user, ok := h.tokenUser(c)
	if !ok {
		return
	}
	var token model.ModelToken
	err := h.DB.Where("user_id = ?", user.Id).First(&token).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		resp.ERROR(c, "查询模型令牌失败")
		return
	}
	resp.SUCCESS(c, gin.H{"enabled": err == nil, "created_at": token.CreatedAt})
}

func (h *ModelAPIHandler) RotateToken(c *gin.Context) {
	user, ok := h.tokenUser(c)
	if !ok {
		return
	}
	var random [32]byte
	if _, err := rand.Read(random[:]); err != nil {
		resp.ERROR(c, "生成模型令牌失败")
		return
	}
	token := modelTokenPrefix + base64.RawURLEncoding.EncodeToString(random[:])
	item := model.ModelToken{UserId: user.Id, TokenHash: modelTokenHash(token), CreatedAt: time.Now()}
	err := h.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"token_hash", "created_at"}),
	}).Create(&item).Error
	if err != nil {
		resp.ERROR(c, "保存模型令牌失败")
		return
	}
	resp.SUCCESS(c, gin.H{"token": token, "created_at": item.CreatedAt})
}

func (h *ModelAPIHandler) RevokeToken(c *gin.Context) {
	user, ok := h.tokenUser(c)
	if !ok {
		return
	}
	if err := h.DB.Where("user_id = ?", user.Id).Delete(&model.ModelToken{}).Error; err != nil {
		resp.ERROR(c, "停用模型令牌失败")
		return
	}
	resp.SUCCESS(c, "模型令牌已停用")
}

func (h *ModelAPIHandler) authenticate(c *gin.Context) {
	c.Header("Cache-Control", "no-store")
	parts := strings.Fields(c.GetHeader("Authorization"))
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || !strings.HasPrefix(parts[1], modelTokenPrefix) || len(parts[1]) != len(modelTokenPrefix)+43 {
		modelAPIError(c, http.StatusUnauthorized, "invalid_api_key", "请使用 ThinkBuddy 模型令牌进行 Bearer 认证")
		return
	}
	var token model.ModelToken
	err := h.DB.Where("token_hash = ?", modelTokenHash(parts[1])).First(&token).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			modelAPIError(c, http.StatusUnauthorized, "invalid_api_key", "模型令牌无效或已停用")
		} else {
			modelAPIError(c, http.StatusInternalServerError, "server_error", "验证模型令牌失败")
		}
		return
	}
	var user model.User
	err = h.DB.First(&user, token.UserId).Error
	if err != nil || !modelUserAvailable(user) {
		modelAPIError(c, http.StatusForbidden, "account_unavailable", "账号不可用")
		return
	}
	c.Set("model_api_user", user)
}

func (h *ModelAPIHandler) availableModels(user model.User) *gorm.DB {
	var ids []uint
	// 无效的个人模型配置不会授予额外权限。
	if json.Unmarshal([]byte(user.ChatModels), &ids) != nil {
		ids = nil
	}
	return h.DB.Where("type = ? AND enabled = ?", "chat", true).
		Where(h.DB.Where("open = ?", true).Or("id IN ?", ids))
}

func (h *ModelAPIHandler) Models(c *gin.Context) {
	user := c.MustGet("model_api_user").(model.User)
	var models []model.ChatModel
	if err := h.availableModels(user).Order("sort_num ASC, id ASC").Find(&models).Error; err != nil {
		modelAPIError(c, http.StatusInternalServerError, "server_error", "读取模型列表失败")
		return
	}
	data := make([]gin.H, 0, len(models))
	seen := make(map[string]bool)
	for _, item := range models {
		if item.Value == "" || seen[item.Value] {
			continue
		}
		seen[item.Value] = true
		data = append(data, gin.H{"id": item.Value, "object": "model", "created": item.CreatedAt.Unix(), "owned_by": "thinkbuddy"})
	}
	c.JSON(http.StatusOK, gin.H{"object": "list", "data": data})
}
