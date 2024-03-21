package main

import (
	"fmt"
	"regexp"
	"strings"
)

type Person struct {
	Name string
	Age  int
}

type Student struct {
	Person
	School string
}

func main() {

	text := `
开始创作您的LOGO了！基于您提供的描述，我将设计一个主题内容为一个小人推着一个购物车的卡通风格LOGO。稍等片刻，您的创意即将变为现实。
	{"prompt":"Create a cartoon-style logo featuring a character pushing a shopping cart. The design should be colorful, vibrant, and engaging, showcasing the character in a dynamic and cheerful pose. The shopping cart should be noticeable but not overpower the character. Use bright and welcoming colors to make the logo inviting and fun. The character can be stylized in a cute and approachable manner, suitable for a wide range of audiences. Ensure the logo is clear and easily recognizable at small sizes.","size":"1024x1024"}


![image1](https://filesystem.site/cdn/20240320/JQIliW99JMPZRjMkgS2PlWNlfUtqDu.webp)

[下载1](https://filesystem.site/cdn/download/20240320/JQIliW99JMPZRjMkgS2PlWNlfUtqDu.webp)

这是为您设计的卡通风格LOGO，主题是一个小人推着购物车。请查看图像，看它是否满足您的需求。如果您对这个设计满意，请访问 [Vectorizer.ai](https://vectorizer.ai/) 将其转换为矢量图，以便在不同大小和格式下保持清晰度。如果您觉得需要进一步的调整或改进，请告诉我，我们可以继续优化设计。`
	pattern := `!\[([^\]]*)]\(([^)]+)\)`
	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(text, -1)

	// 下载图片并替换链接地址
	for _, match := range matches {
		imageURL := match[2]
		fmt.Println(imageURL)
		// 对于相同地址的图片，已经被替换了，就不再重复下载了
		if !strings.Contains(text, imageURL) {
			continue
		}
	}
}
