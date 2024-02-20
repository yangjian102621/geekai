package main

import (
	"fmt"
	"regexp"
)

func main() {
	text := `
> search("Shenzhen weather January 15, 2024") 

> mclick([0, 9, 16]) 

> **end-searching**

今天深圳的天气情况如下：

- 白天气温预计在21°C至24°C之间，天气晴朗。
- 晚上气温预计在21°C左右，云量较多，可能会有间断性小雨。
- 风向主要是东南风，风速大约在6至12公里每小时之间。

这些信息表明深圳今天的天气相对舒适，适合户外活动。晚上可能需要带伞以应对间断性小雨。温度较为宜人，早晚可能稍微凉爽一些【[Shenzhen weather in January 2024 | Shenzhen 14 day weather](https://www.weather25.com/asia/china/guangdong/shenzhen?page=month&month=January)】【[Hourly forecast for Shenzhen, Guangdong, China](https://www.timeanddate.com/weather/china/shenzhen/hourly)】【[Shenzhen Guangdong China 15 Day Weather Forecast](https://www.weatheravenue.com/en/asia/cn/guangdong/shenzhen-weather-15-days.html)】。

我将根据这些信息生成一张气象图，展示深圳今天的天气情况。

	{"prompt":"A detailed weather map for Shenzhen, China, on January 15, 2024. The map shows a sunny day with clear skies during the day and partly cloudy skies at night. Temperatures range from 21\u00b0C to 24\u00b0C during the day and around 21\u00b0C at night. There are indications of light southeast winds during the day and evening, with wind speeds ranging from 6 to 12 km/h. The map includes symbols for sunshine, light clouds, and wind direction arrows, along with temperature readings for different times of the day. The layout is clear, with a focus on Shenzhen's geographical location and the surrounding region.","size":"1024x1024"}


![image1](https://filesystem.site/cdn/20240115/XD6EjyPDGCD4X3AQt3h3FijRmSb6fB.webp)

![下载1](https://filesystem.site/cdn/download/20240115/XD6EjyPDGCD4X3AQt3h3FijRmSb6fB.webp)

And here is another image link: ![another image](https://example.com/another-image.png).


这是根据今天深圳的天气情况制作的气象图。图中展示了白天晴朗、夜间部分多云的天气，以及相关的温度和风向信息。`
	pattern := `!\[([^\]]*)]\(([^)]+)\)`

	// 编译正则表达式
	re := regexp.MustCompile(pattern)

	// 查找匹配的字符串
	matches := re.FindAllStringSubmatch(text, -1)

	// 提取链接并打印
	for _, match := range matches {
		fmt.Println(match[2])
	}
}
