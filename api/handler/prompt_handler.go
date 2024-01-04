package handler

import (
	"chatplus/core"
	"chatplus/core/types"
	"chatplus/utils"
	"chatplus/utils/resp"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const rewritePromptTemplate = "Please rewrite the following text into AI painting prompt words, and please try to add detailed description of the picture, painting style, scene, rendering effect, picture light and other elements. Please output directly in English without any explanation, within 150 words. The text to be rewritten is: [%s]"
const translatePromptTemplate = "Translate the following painting prompt words into English keyword phrases. Without any explanation, directly output the keyword phrases separated by commas. The content to be translated is: [%s]"

type PromptHandler struct {
	BaseHandler
	db *gorm.DB
}

func NewPromptHandler(app *core.AppServer, db *gorm.DB) *PromptHandler {
	h := &PromptHandler{db: db}
	h.App = app
	return h
}

// Rewrite translate and rewrite prompt with ChatGPT
func (h *PromptHandler) Rewrite(c *gin.Context) {
	var data struct {
		Prompt string `json:"prompt"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	content, err := utils.OpenAIRequest(h.db, data.Prompt, rewritePromptTemplate)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c, content)
}

func (h *PromptHandler) Translate(c *gin.Context) {
	var data struct {
		Prompt string `json:"prompt"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	content, err := utils.OpenAIRequest(h.db, data.Prompt, translatePromptTemplate)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c, content)
}
