package main

import (
	"fmt"
	"geekai/service/sora"
)

func main() {
	text := ` **Analyzing input prompt** > ***- The input provides a JSON with an empty history and a prompt in Chinese: "女孩在雪山上滑雪，欢快的节奏，姿势优雅，伴随女孩的欢声笑语".*** > ID: task_01k8pye324ef7t6heq6jyaxbfe >[数据预览](https://asyncdata.net/web/task_01k8pye324ef7t6heq6jyaxbfe) | [原始数据](https://asyncdata.net/source/task_01k8pye324ef7t6heq6jyaxbfe) > 排队中... > 生成中. >🏃‍ 进度 36..47..57..61..69..76..79..81..`

	soraService := sora.NewSoraService(nil)
	url, err := soraService.DownloadVideoURL(text)
	if err != nil {
		fmt.Printf("extract video URL failed: %v", err)
		return
	}
	fmt.Println(url)
}
