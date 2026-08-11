package ppt

import (
	"context"
	"fmt"
	"geekai/core/types"
	"geekai/log"
	"net/http"
	"time"

	"github.com/imroc/req/v3"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
	"golang.org/x/time/rate"
)

var imageLogger = log.GetLogger()

// ImageGenerator 图片生成适配器接口
type ImageGenerator interface {
	Provider() string
	Generate(ctx context.Context, prompt string) (string, error)
	// GenerateWithReference 图生图；referenceImages 为公网 URL 或 data:image/...;base64,...（本地文件应在调用前经 PrepareReferenceInputsForImg2Img 转换）
	GenerateWithReference(ctx context.Context, prompt string, referenceImages []string) (string, error)
}

// Nano Banana 适配器（OpenAI DALL-E 风格 API）
type nanoBananaImageGenerator struct {
	client  *req.Client
	cfg     types.PPTConfig
	limiter *rate.Limiter
}

// Seedream 适配器（火山引擎 arkruntime SDK）
type seedreamImageGenerator struct {
	cfg     types.PPTConfig
	limiter *rate.Limiter
}

// NewImageGenerator 根据配置创建对应的图片生成适配器
func NewImageGenerator(cfg types.PPTConfig) (ImageGenerator, error) {
	qps := cfg.QPSLimit
	if qps <= 0 {
		qps = 1
	}
	limiter := rate.NewLimiter(rate.Limit(qps), 1)

	switch cfg.ActiveImageProvider {
	case types.PPTImageProviderNanoBanana:
		if cfg.NanoBananaApiURL == "" || cfg.NanoBananaApiKey == "" {
			return nil, fmt.Errorf("nano banana api not configured")
		}
		return &nanoBananaImageGenerator{
			client:  req.C().SetTimeout(3 * time.Minute),
			cfg:     cfg,
			limiter: limiter,
		}, nil
	case types.PPTImageProviderSeedream:
		if cfg.SeedreamBaseURL == "" || cfg.SeedreamApiKey == "" || cfg.SeedreamModel == "" {
			return nil, fmt.Errorf("seedream api not configured")
		}
		return &seedreamImageGenerator{
			cfg:     cfg,
			limiter: limiter,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported image provider: %s", cfg.ActiveImageProvider)
	}
}

func (g *nanoBananaImageGenerator) Provider() string {
	return string(types.PPTImageProviderNanoBanana)
}

// nanoBananaReq 按 OpenAI DALL-E 风格 / Nano-banana API 文档
type nanoBananaReq struct {
	Model          string   `json:"model"`
	Prompt         string   `json:"prompt"`
	ResponseFormat string   `json:"response_format,omitempty"` // url 或 b64_json
	AspectRatio    string   `json:"aspect_ratio,omitempty"`    // 1:1, 4:3, 3:4, 16:9, 9:16, 2:3, 3:2, 4:5, 5:4, 21:9
	Image          []string `json:"image,omitempty"`           // 参考图 url 或 b64
}

// nanoBananaRes 响应为 data[].url（DALL-E 风格）
type nanoBananaRes struct {
	Data []struct {
		URL     string `json:"url,omitempty"`
		B64JSON string `json:"b64_json,omitempty"`
	} `json:"data"`
}

type nanoBananaErr struct {
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

func (g *nanoBananaImageGenerator) buildReqBody(prompt string, referenceImages []string) nanoBananaReq {
	modelName := g.cfg.NanoBananaModel
	if modelName == "" {
		modelName = "nano-banana"
	}
	reqBody := nanoBananaReq{
		Model:  modelName,
		Prompt: prompt,
	}
	if len(referenceImages) > 0 {
		reqBody.Image = referenceImages
	}
	if g.cfg.NanoBananaResponseFormat != "" {
		reqBody.ResponseFormat = g.cfg.NanoBananaResponseFormat
	} else {
		reqBody.ResponseFormat = "url"
	}
	if g.cfg.NanoBananaAspectRatio != "" {
		reqBody.AspectRatio = g.cfg.NanoBananaAspectRatio
	} else {
		reqBody.AspectRatio = "16:9"
	}
	return reqBody
}

func (g *nanoBananaImageGenerator) Generate(ctx context.Context, prompt string) (string, error) {
	return g.GenerateWithReference(ctx, prompt, nil)
}

func (g *nanoBananaImageGenerator) GenerateWithReference(ctx context.Context, prompt string, referenceImages []string) (string, error) {
	reqBody := g.buildReqBody(prompt, referenceImages)

	var (
		result nanoBananaRes
		errRes nanoBananaErr
	)

	do := func() (int, error) {
		if err := g.limiter.Wait(ctx); err != nil {
			return 0, err
		}

		imageLogger.Infof("nano banana generate image, api: %s", g.cfg.NanoBananaApiURL)
		r, err := g.client.R().
			SetContext(ctx).
			SetHeader("Content-Type", "application/json").
			SetHeader("Authorization", "Bearer "+g.cfg.NanoBananaApiKey).
			SetBody(reqBody).
			SetSuccessResult(&result).
			SetErrorResult(&errRes).
			Post(g.cfg.NanoBananaApiURL)
		if err != nil {
			return 0, err
		}
		if r.IsErrorState() {
			return r.StatusCode, fmt.Errorf("nano banana error: %s, %s", r.Status, errRes.Error.Message)
		}
		if len(result.Data) == 0 || result.Data[0].URL == "" {
			return r.StatusCode, fmt.Errorf("nano banana returned empty data")
		}
		return r.StatusCode, nil
	}

	if err := callWithRetry(ctx, do); err != nil {
		return "", err
	}

	return result.Data[0].URL, nil
}

func (g *seedreamImageGenerator) Provider() string {
	return string(types.PPTImageProviderSeedream)
}

func (g *seedreamImageGenerator) Generate(ctx context.Context, prompt string) (string, error) {
	return g.GenerateWithReference(ctx, prompt, nil)
}

func (g *seedreamImageGenerator) GenerateWithReference(ctx context.Context, prompt string, referenceImages []string) (string, error) {
	if err := g.limiter.Wait(ctx); err != nil {
		return "", err
	}

	client := arkruntime.NewClientWithApiKey(g.cfg.SeedreamApiKey, arkruntime.WithBaseUrl(g.cfg.SeedreamBaseURL))

	size := g.cfg.SeedreamSize
	if size == "" {
		size = "1920x1080"
	}
	responseFormat := g.cfg.SeedreamResponseType
	if responseFormat == "" {
		responseFormat = "url"
	}

	generateReq := model.GenerateImagesRequest{
		Model:          g.cfg.SeedreamModel,
		Prompt:         prompt,
		Size:           volcengine.String(size),
		ResponseFormat: volcengine.String(responseFormat),
		Watermark:      volcengine.Bool(g.cfg.SeedreamWatermark),
	}
	if len(referenceImages) > 0 {
		generateReq.Image = referenceImages
	}

	var lastErr error
	for attempt := 0; attempt < 3; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return "", ctx.Err()
			case <-time.After(time.Duration(attempt) * 2 * time.Second):
			}
		}
		imageLogger.Infof("seedream generate image, api: %s", g.cfg.SeedreamBaseURL)
		if err := generateReq.NormalizeImages(); err != nil {
			return "", fmt.Errorf("seedream normalize images: %w", err)
		}
		resp, err := client.GenerateImages(ctx, generateReq)
		if err != nil {
			lastErr = fmt.Errorf("seedream error: %w", err)
			continue
		}
		if resp.Data == nil || len(resp.Data) == 0 {
			lastErr = fmt.Errorf("seedream returned empty data")
			continue
		}
		if resp.Data[0].Url == nil || *resp.Data[0].Url == "" {
			lastErr = fmt.Errorf("seedream returned empty url")
			continue
		}
		return *resp.Data[0].Url, nil
	}
	return "", lastErr
}

// callWithRetry 对 429 错误做指数退避重试
func callWithRetry(ctx context.Context, fn func() (int, error)) error {
	var (
		retries   = 3
		backoffs  = []time.Duration{2 * time.Second, 4 * time.Second, 8 * time.Second}
		lastError error
	)

	for i := 0; i < retries; i++ {
		status, err := fn()
		if err == nil {
			return nil
		}
		lastError = err

		// 仅对 429 做指数退避重试
		if status != http.StatusTooManyRequests || i == retries-1 {
			break
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(backoffs[i]):
		}
	}

	return lastError
}
