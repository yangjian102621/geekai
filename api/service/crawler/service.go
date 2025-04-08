package crawler

import (
	"context"
	"errors"
	"fmt"
	"geekai/logger"
	"net/url"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

// Service 网络爬虫服务
type Service struct {
	browser *rod.Browser
}

// NewService 创建一个新的爬虫服务
func NewService() (*Service, error) {
	// 启动浏览器
	path, _ := launcher.LookPath()
	u := launcher.New().Bin(path).
		Headless(true).                   // 无头模式
		Set("disable-web-security", "").  // 禁用网络安全限制
		Set("disable-gpu", "").           // 禁用 GPU 加速
		Set("no-sandbox", "").            // 禁用沙箱模式
		Set("disable-setuid-sandbox", "").// 禁用 setuid 沙箱
		MustLaunch()

	browser := rod.New().ControlURL(u).MustConnect()
	
	return &Service{
		browser: browser,
	}, nil
}

// SearchResult 搜索结果
type SearchResult struct {
	Title   string `json:"title"`   // 标题
	URL     string `json:"url"`     // 链接
	Content string `json:"content"` // 内容摘要
}

// WebSearch 网络搜索
func (s *Service) WebSearch(keyword string, maxPages int) ([]SearchResult, error) {
	if keyword == "" {
		return nil, errors.New("搜索关键词不能为空")
	}
	
	if maxPages <= 0 {
		maxPages = 1
	}
	if maxPages > 10 {
		maxPages = 10 // 最多搜索 10 页
	}

	results := make([]SearchResult, 0)
	
	// 使用百度搜索
	searchURL := fmt.Sprintf("https://www.baidu.com/s?wd=%s", url.QueryEscape(keyword))
	
	// 设置页面超时
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	// 创建页面
	page := s.browser.MustPage()
	defer page.MustClose()
	
	// 设置视口大小
	err := page.SetViewport(&proto.EmulationSetDeviceMetricsOverride{
		Width:  1280,
		Height: 800,
	})
	if err != nil {
		return nil, fmt.Errorf("设置视口失败: %v", err)
	}
	
	// 导航到搜索页面
	err = page.Context(ctx).Navigate(searchURL)
	if err != nil {
		return nil, fmt.Errorf("导航到搜索页面失败: %v", err)
	}
	
	// 等待搜索结果加载完成
	err = page.WaitLoad()
	if err != nil {
		return nil, fmt.Errorf("等待页面加载完成失败: %v", err)
	}
	
	// 分析当前页面的搜索结果
	for i := 0; i < maxPages; i++ {
		if i > 0 {
			// 点击下一页按钮
			nextPage, err := page.Element("a.n")
			if err != nil || nextPage == nil {
				break // 没有下一页
			}
			
			err = nextPage.Click(proto.InputMouseButtonLeft, 1)
			if err != nil {
				break // 点击下一页失败
			}
			
			// 等待新页面加载
			err = page.WaitLoad()
			if err != nil {
				break
			}
		}
		
		// 提取搜索结果
		resultElements, err := page.Elements(".result, .c-container")
		if err != nil || resultElements == nil {
			continue
		}
		
		for _, result := range resultElements {
			// 获取标题
			titleElement, err := result.Element("h3, .t")
			if err != nil || titleElement == nil {
				continue
			}
			
			title, err := titleElement.Text()
			if err != nil {
				continue
			}
			
			// 获取 URL
			linkElement, err := titleElement.Element("a")
			if err != nil || linkElement == nil {
				continue
			}
			
			href, err := linkElement.Attribute("href")
			if err != nil || href == nil {
				continue
			}
			
			// 获取内容摘要 - 尝试多个可能的选择器
			var contentElement *rod.Element
			var content string
			
			// 尝试多个可能的选择器来适应不同版本的百度搜索结果
			selectors := []string{".content-right_8Zs40", ".c-abstract", ".content_LJ0WN", ".content"}
			for _, selector := range selectors {
				contentElement, err = result.Element(selector)
				if err == nil && contentElement != nil {
					content, _ = contentElement.Text()
					if content != "" {
						break
					}
				}
			}
			
			// 如果所有选择器都失败，尝试直接从结果块中提取文本
			if content == "" {
				// 获取结果元素的所有文本
				fullText, err := result.Text()
				if err == nil && fullText != "" {
					// 简单处理：从全文中移除标题，剩下的可能是摘要
					fullText = strings.Replace(fullText, title, "", 1)
					// 清理文本
					content = strings.TrimSpace(fullText)
					// 限制内容长度
					if len(content) > 200 {
						content = content[:200] + "..."
					}
				}
			}
			
			// 添加到结果集
			results = append(results, SearchResult{
				Title:   title,
				URL:     *href,
				Content: content,
			})
			
			// 限制结果数量，每页最多 10 条
			if len(results) >= 10*maxPages {
				break
			}
		}
	}
	
	// 获取真实 URL（百度搜索结果中的 URL 是短链接，需要跳转获取真实 URL）
	for i, result := range results {
		realURL, err := s.getRedirectURL(result.URL)
		if err == nil && realURL != "" {
			results[i].URL = realURL
		}
	}
	
	return results, nil
}

// 获取真实 URL
func (s *Service) getRedirectURL(shortURL string) (string, error) {
	// 创建页面
	page, err := s.browser.Page(proto.TargetCreateTarget{URL: ""})
	if err != nil {
		return shortURL, err // 返回原始URL
	}
	defer func() {
		_ = page.Close()
	}()
	
	// 导航到短链接
	err = page.Navigate(shortURL)
	if err != nil {
		return shortURL, err // 返回原始URL
	}
	
	// 等待重定向完成
	time.Sleep(2 * time.Second)
	
	// 获取当前 URL
	info, err := page.Info()
	if err != nil {
		return shortURL, err // 返回原始URL
	}
	
	return info.URL, nil
}

// Close 关闭浏览器
func (s *Service) Close() error {
	if s.browser != nil {
		err := s.browser.Close()
		s.browser = nil
		return err
	}
	return nil
}

// SearchWeb 封装的搜索方法
func SearchWeb(keyword string, maxPages int) (string, error) {
	// 添加panic恢复机制
	defer func() {
		if r := recover(); r != nil {
			log := logger.GetLogger()
			log.Errorf("爬虫服务崩溃: %v", r)
		}
	}()
	
	service, err := NewService()
	if err != nil {
		return "", fmt.Errorf("创建爬虫服务失败: %v", err)
	}
	defer service.Close()
	
	// 设置超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	
	// 使用goroutine和通道来处理超时
	resultChan := make(chan []SearchResult, 1)
	errChan := make(chan error, 1)
	
	go func() {
		results, err := service.WebSearch(keyword, maxPages)
		if err != nil {
			errChan <- err
			return
		}
		resultChan <- results
	}()
	
	// 等待结果或超时
	select {
	case <-ctx.Done():
		return "", fmt.Errorf("搜索超时: %v", ctx.Err())
	case err := <-errChan:
		return "", fmt.Errorf("搜索失败: %v", err)
	case results := <-resultChan:
		if len(results) == 0 {
			return "未找到关于 \"" + keyword + "\" 的相关搜索结果", nil
		}
		
		// 格式化结果
		var builder strings.Builder
		builder.WriteString(fmt.Sprintf("为您找到关于 \"%s\" 的 %d 条搜索结果：\n\n", keyword, len(results)))
		
		for i, result := range results {
			// // 尝试打开链接获取实际内容
			// page := service.browser.MustPage()
			// defer page.MustClose()
			
			// // 设置页面超时
			// pageCtx, pageCancel := context.WithTimeout(context.Background(), 10*time.Second)
			// defer pageCancel()
			
			// // 导航到目标页面
			// err := page.Context(pageCtx).Navigate(result.URL)
			// if err == nil {
			// 	// 等待页面加载
			// 	_ = page.WaitLoad()
				
			// 	// 获取页面标题
			// 	title, err := page.Eval("() => document.title")
			// 	if err == nil && title.Value.String() != "" {
			// 		result.Title = title.Value.String()
			// 	}
				
			// 	// 获取页面主要内容
			// 	if content, err := page.Element("body"); err == nil {
			// 		if text, err := content.Text(); err == nil {
			// 			// 清理并截取内容
			// 			text = strings.TrimSpace(text)
			// 			if len(text) > 200 {
			// 				text = text[:200] + "..."
			// 			}
			// 			result.Content = text
			// 		}
			// 	}
			// }
			
			builder.WriteString(fmt.Sprintf("%d. **%s**\n", i+1, result.Title))
			builder.WriteString(fmt.Sprintf("   链接: %s\n", result.URL))
			if result.Content != "" {
				builder.WriteString(fmt.Sprintf("   摘要: %s\n", result.Content))
			}
			builder.WriteString("\n")
		}
		
		return builder.String(), nil
	}
} 