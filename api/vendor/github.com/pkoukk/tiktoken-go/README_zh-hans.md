# tiktoken-go
Go 语言版本的 OpenAI 的 tiktoken。  
帮你把文本转换成 OpenAI 的模型可以识别的 token。  
tiktoken的原项目地址[tiktoken](https://github.com/openai/tiktoken).  

# 用法

## 安装


```bash
go get github.com/pkoukk/tiktoken-go
# 默认的tiktoken需要从openai下载token字典，如果想要离线使用，可以使用以下分支
go get github.com/pkoukk/tiktoken-go@embed
```
## 缓存
Tiktoken-go 和原始的 Tiktoken 库一样，具有相同的缓存机制。  
您可以使用环境变量 TIKTOKEN_CACHE_DIR 来设置缓存目录。  
一旦设置了该变量，tiktoken-go 将使用该目录来缓存令牌字典。  
如果您未设置此环境变量，则 tiktoken-go 将在每次首次初始化编码时下载字典。  

## 例子

### get token by encoding

```go
package main

import (
    "fmt"
    "github.com/pkoukk/tiktoken-go"
)

func main() (num_tokens int) {
    text = "Hello, world!"
    encoding = "r50k_base"

	tke, err := tiktoken.GetEncoding(encoding)
	if err != nil {
		err = fmt.Errorf("GetEncoding: %v", err)
		return
	}

    // encode
	token := tke.Encode(text, nil, nil)

    // num_tokens
    num_tokens = len(token)
}
```

### get token by Model

```go
package main

import (
    "fmt"
    "github.com/pkoukk/tiktoken-go"
)

func main() (num_tokens int) {
    text = "Hello, world!"
    encoding = "davinci"

   tkm, err := tiktoken.EncodingForModel(model)
	if err != nil {
		err = fmt.Errorf(": %v", err)
		return
	}

	 // encode
	token := tke.Encode(text, nil, nil)

    // num_tokens
    num_tokens = len(token)
}
```

### 计算chat API消息当中的token消耗
这段代码由@nasa1024根据[官方示例](https://github.com/openai/openai-cookbook/blob/main/examples/How_to_count_tokens_with_tiktoken.ipynb)编写

请注意，消息的token计算方式可能随时会发生改变，以下代码并不一定在将来适用，如果您需要精确的计算，请关注官方文档。

如果您发现这段代码不再适用，欢迎您提PR或Issue。

```go
package main

import (
	"fmt"

	"github.com/pkoukk/tiktoken-go"
	"github.com/sashabaranov/go-openai"
)

func NumTokensFromMessages(messages []openai.ChatCompletionMessage, model string) (num_tokens int) {
	tkm, err := tiktoken.EncodingForModel(model)
	if err != nil {
		err = fmt.Errorf("EncodingForModel: %v", err)
		fmt.Println(err)
		return
	}

	var tokens_per_message int
	var tokens_per_name int
	if model == "gpt-3.5-turbo-0301" || model == "gpt-3.5-turbo" {
		tokens_per_message = 4
		tokens_per_name = -1
	} else if model == "gpt-4-0314" || model == "gpt-4" {
		tokens_per_message = 3
		tokens_per_name = 1
	} else {
		fmt.Println("Warning: model not found. Using cl100k_base encoding.")
		tokens_per_message = 3
		tokens_per_name = 1
	}

	for _, message := range messages {
		num_tokens += tokens_per_message
		num_tokens += len(tkm.Encode(message.Content, nil, nil))
		num_tokens += len(tkm.Encode(message.Role, nil, nil))
		if message.Name != "" {
			num_tokens += tokens_per_name
		}
	}
	num_tokens += 3
	return num_tokens
}
```

# available encodings
 | Encoding name           | OpenAI models                                        |
 | ----------------------- | ---------------------------------------------------- |
 | `cl100k_base`           | `gpt-4`, `gpt-3.5-turbo`, `text-embedding-ada-002`   |
 | `p50k_base`             | Codex models, `text-davinci-002`, `text-davinci-003` |
 | `r50k_base` (or `gpt2`) | GPT-3 models like `davinci`                          |



# available models
| Model name                   | OpenAI models |
| ---------------------------- | ------------- |
| gpt-4                        | cl100k_base   |
| gpt-4-*                      | cl100k_base   |
| gpt-3.5-turbo                | cl100k_base   |
| gpt-3.5-turbo-*              | cl100k_base   |
| text-davinci-003             | p50k_base     |
| text-davinci-002             | p50k_base     |
| text-davinci-001             | r50k_base     |
| text-curie-001               | r50k_base     |
| text-babbage-001             | r50k_base     |
| text-ada-001                 | r50k_base     |
| davinci                      | r50k_base     |
| curie                        | r50k_base     |
| babbage                      | r50k_base     |
| ada                          | r50k_base     |
| code-davinci-002             | p50k_base     |
| code-davinci-001             | p50k_base     |
| code-cushman-002             | p50k_base     |
| code-cushman-001             | p50k_base     |
| davinci-codex                | p50k_base     |
| cushman-codex                | p50k_base     |
| text-davinci-edit-001        | p50k_edit     |
| code-davinci-edit-001        | p50k_edit     |
| text-embedding-ada-002       | cl100k_base   |
| text-similarity-davinci-001  | r50k_base     |
| text-similarity-curie-001    | r50k_base     |
| text-similarity-babbage-001  | r50k_base     |
| text-similarity-ada-001      | r50k_base     |
| text-search-davinci-doc-001  | r50k_base     |
| text-search-curie-doc-001    | r50k_base     |
| text-search-babbage-doc-001  | r50k_base     |
| text-search-ada-doc-001      | r50k_base     |
| code-search-babbage-code-001 | r50k_base     |
| code-search-ada-code-001     | r50k_base     |
| gpt2                         | gpt2          |

# 与官方 [tiktoken](https://github.com/openai/tiktoken) 的对比

## get token by encoding

| python tiktoken                                          | golang tiktoken-go                                       |
| :------------------------------------------------------- | :------------------------------------------------------- |
| text: hallo world!, encoding: cl100k_base, token: 4      | text: hallo world!, encoding: cl100k_base, token: 4      |
| text: hallo world!, encoding: p50k_base, token: 4        | text: hallo world!, encoding: p50k_base, token: 4        |
| text: hallo world!, encoding: r50k_base, token: 4        | text: hallo world!, encoding: r50k_base, token: 4        |
| text: 你好世界！, encoding: cl100k_base, token: 6        | text: 你好世界！, encoding: cl100k_base, token: 6        |
| text: 你好世界！, encoding: p50k_base, token: 11         | text: 你好世界！, encoding: p50k_base, token: 11         |
| text: 你好世界！, encoding: r50k_base, token: 11         | text: 你好世界！, encoding: r50k_base, token: 11         |
| text: こんにちは世界！, encoding: cl100k_base, token: 5  | text: こんにちは世界！, encoding: cl100k_base, token: 5  |
| text: こんにちは世界！, encoding: p50k_base, token: 13   | text: こんにちは世界！, encoding: p50k_base, token: 13   |
| text: こんにちは世界！, encoding: r50k_base, token: 13   | text: こんにちは世界！, encoding: r50k_base, token: 13   |
| text: 안녕하세요 세계!, encoding: cl100k_base, token: 10 | text: 안녕하세요 세계!, encoding: cl100k_base, token: 10 |
| text: 안녕하세요 세계!, encoding: p50k_base, token: 21   | text: 안녕하세요 세계!, encoding: p50k_base, token: 21   |
| text: 안녕하세요 세계!, encoding: r50k_base, token: 21   | text: 안녕하세요 세계!, encoding: r50k_base, token: 21   |
| text: Привет мир!, encoding: cl100k_base, token: 6       | text: Привет мир!, encoding: cl100k_base, token: 6       |
| text: Привет мир!, encoding: p50k_base, token: 12        | text: Привет мир!, encoding: p50k_base, token: 12        |
| text: Привет мир!, encoding: r50k_base, token: 12        | text: Привет мир!, encoding: r50k_base, token: 12        |
| text: ¡Hola mundo!, encoding: cl100k_base, token: 4      | text: ¡Hola mundo!, encoding: cl100k_base, token: 4      |
| text: ¡Hola mundo!, encoding: p50k_base, token: 7        | text: ¡Hola mundo!, encoding: p50k_base, token: 7        |
| text: ¡Hola mundo!, encoding: r50k_base, token: 7        | text: ¡Hola mundo!, encoding: r50k_base, token: 7        |
| text: Hallo Welt!, encoding: cl100k_base, token: 3       | text: Hallo Welt!, encoding: cl100k_base, token: 3       |
| text: Hallo Welt!, encoding: p50k_base, token: 5         | text: Hallo Welt!, encoding: p50k_base, token: 5         |
| text: Hallo Welt!, encoding: r50k_base, token: 5         | text: Hallo Welt!, encoding: r50k_base, token: 5         |
| text: Bonjour le monde!, encoding: cl100k_base, token: 4 | text: Bonjour le monde!, encoding: cl100k_base, token: 4 |
| text: Bonjour le monde!, encoding: p50k_base, token: 7   | text: Bonjour le monde!, encoding: p50k_base, token: 7   |
| text: Bonjour le monde!, encoding: r50k_base, token: 7   | text: Bonjour le monde!, encoding: r50k_base, token: 7   |
| text: Ciao mondo!, encoding: cl100k_base, token: 4       | text: Ciao mondo!, encoding: cl100k_base, token: 4       |
| text: Ciao mondo!, encoding: p50k_base, token: 5         | text: Ciao mondo!, encoding: p50k_base, token: 5         |
| text: Ciao mondo!, encoding: r50k_base, token: 5         | text: Ciao mondo!, encoding: r50k_base, token: 5         |
| text: Hej världen!, encoding: cl100k_base, token: 7      | text: Hej världen!, encoding: cl100k_base, token: 7      |
| text: Hej världen!, encoding: p50k_base, token: 8        | text: Hej världen!, encoding: p50k_base, token: 8        |
| text: Hej världen!, encoding: r50k_base, token: 8        | text: Hej världen!, encoding: r50k_base, token: 8        |
| text: Hallo wereld!, encoding: cl100k_base, token: 3     | text: Hallo wereld!, encoding: cl100k_base, token: 3     |
| text: Hallo wereld!, encoding: p50k_base, token: 5       | text: Hallo wereld!, encoding: p50k_base, token: 5       |
| text: Hallo wereld!, encoding: r50k_base, token: 5       | text: Hallo wereld!, encoding: r50k_base, token: 5       |
| text: Hallo verden!, encoding: cl100k_base, token: 4     | text: Hallo verden!, encoding: cl100k_base, token: 4     |
| text: Hallo verden!, encoding: p50k_base, token: 5       | text: Hallo verden!, encoding: p50k_base, token: 5       |
| text: Hallo verden!, encoding: r50k_base, token: 5       | text: Hallo verden!, encoding: r50k_base, token: 5       |
| text: Hallo wereld!, encoding: cl100k_base, token: 3     | text: Hallo wereld!, encoding: cl100k_base, token: 3     |
| text: Hallo wereld!, encoding: p50k_base, token: 5       | text: Hallo wereld!, encoding: p50k_base, token: 5       |
| text: Hallo wereld!, encoding: r50k_base, token: 5       | text: Hallo wereld!, encoding: r50k_base, token: 5       |
| text: Hallo verden!, encoding: cl100k_base, token: 4     | text: Hallo verden!, encoding: cl100k_base, token: 4     |
| text: Hallo verden!, encoding: p50k_base, token: 5       | text: Hallo verden!, encoding: p50k_base, token: 5       |
| text: Hallo verden!, encoding: r50k_base, token: 5       | text: Hallo verden!, encoding: r50k_base, token: 5       |


## get token by model  

| python tiktoken                                                       | golang tiktoken-go                                                    |
| --------------------------------------------------------------------- | --------------------------------------------------------------------- |
| text: hallo world!, model: gpt-4, token: 4                            | text: hallo world!, model: gpt-4, token: 4                            |
| text: hallo world!, model: gpt-3.5-turbo, token: 4                    | text: hallo world!, model: gpt-3.5-turbo, token: 4                    |
| text: hallo world!, model: text-davinci-003, token: 4                 | text: hallo world!, model: text-davinci-003, token: 4                 |
| text: hallo world!, model: text-davinci-002, token: 4                 | text: hallo world!, model: text-davinci-002, token: 4                 |
| text: hallo world!, model: text-davinci-001, token: 4                 | text: hallo world!, model: text-davinci-001, token: 4                 |
| text: hallo world!, model: text-curie-001, token: 4                   | text: hallo world!, model: text-curie-001, token: 4                   |
| text: hallo world!, model: text-babbage-001, token: 4                 | text: hallo world!, model: text-babbage-001, token: 4                 |
| text: hallo world!, model: text-ada-001, token: 4                     | text: hallo world!, model: text-ada-001, token: 4                     |
| text: hallo world!, model: davinci, token: 4                          | text: hallo world!, model: davinci, token: 4                          |
| text: hallo world!, model: curie, token: 4                            | text: hallo world!, model: curie, token: 4                            |
| text: hallo world!, model: babbage, token: 4                          | text: hallo world!, model: babbage, token: 4                          |
| text: hallo world!, model: ada, token: 4                              | text: hallo world!, model: ada, token: 4                              |
| text: hallo world!, model: code-davinci-002, token: 4                 | text: hallo world!, model: code-davinci-002, token: 4                 |
| text: hallo world!, model: code-davinci-001, token: 4                 | text: hallo world!, model: code-davinci-001, token: 4                 |
| text: hallo world!, model: code-cushman-002, token: 4                 | text: hallo world!, model: code-cushman-002, token: 4                 |
| text: hallo world!, model: code-cushman-001, token: 4                 | text: hallo world!, model: code-cushman-001, token: 4                 |
| text: hallo world!, model: davinci-codex, token: 4                    | text: hallo world!, model: davinci-codex, token: 4                    |
| text: hallo world!, model: cushman-codex, token: 4                    | text: hallo world!, model: cushman-codex, token: 4                    |
| text: hallo world!, model: text-davinci-edit-001, token: 4            | text: hallo world!, model: text-davinci-edit-001, token: 4            |
| text: hallo world!, model: code-davinci-edit-001, token: 4            | text: hallo world!, model: code-davinci-edit-001, token: 4            |
| text: hallo world!, model: text-embedding-ada-002, token: 4           | text: hallo world!, model: text-embedding-ada-002, token: 4           |
| text: hallo world!, model: text-similarity-davinci-001, token: 4      | text: hallo world!, model: text-similarity-davinci-001, token: 4      |
| text: 你好世界！, model: gpt-4, token: 6                              | text: 你好世界！, model: gpt-4, token: 6                              |
| text: 你好世界！, model: gpt-3.5-turbo, token: 6                      | text: 你好世界！, model: gpt-3.5-turbo, token: 6                      |
| text: 你好世界！, model: text-davinci-003, token: 11                  | text: 你好世界！, model: text-davinci-003, token: 11                  |
| text: 你好世界！, model: text-davinci-002, token: 11                  | text: 你好世界！, model: text-davinci-002, token: 11                  |
| text: 你好世界！, model: text-davinci-001, token: 11                  | text: 你好世界！, model: text-davinci-001, token: 11                  |
| text: 你好世界！, model: text-curie-001, token: 11                    | text: 你好世界！, model: text-curie-001, token: 11                    |
| text: 你好世界！, model: text-babbage-001, token: 11                  | text: 你好世界！, model: text-babbage-001, token: 11                  |
| text: 你好世界！, model: text-ada-001, token: 11                      | text: 你好世界！, model: text-ada-001, token: 11                      |
| text: 你好世界！, model: davinci, token: 11                           | text: 你好世界！, model: davinci, token: 11                           |
| text: 你好世界！, model: curie, token: 11                             | text: 你好世界！, model: curie, token: 11                             |
| text: 你好世界！, model: babbage, token: 11                           | text: 你好世界！, model: babbage, token: 11                           |
| text: 你好世界！, model: ada, token: 11                               | text: 你好世界！, model: ada, token: 11                               |
| text: 你好世界！, model: code-davinci-002, token: 11                  | text: 你好世界！, model: code-davinci-002, token: 11                  |
| text: 你好世界！, model: code-davinci-001, token: 11                  | text: 你好世界！, model: code-davinci-001, token: 11                  |
| text: 你好世界！, model: code-cushman-002, token: 11                  | text: 你好世界！, model: code-cushman-002, token: 11                  |
| text: 你好世界！, model: code-cushman-001, token: 11                  | text: 你好世界！, model: code-cushman-001, token: 11                  |
| text: 你好世界！, model: davinci-codex, token: 11                     | text: 你好世界！, model: davinci-codex, token: 11                     |
| text: 你好世界！, model: cushman-codex, token: 11                     | text: 你好世界！, model: cushman-codex, token: 11                     |
| text: 你好世界！, model: text-davinci-edit-001, token: 11             | text: 你好世界！, model: text-davinci-edit-001, token: 11             |
| text: 你好世界！, model: code-davinci-edit-001, token: 11             | text: 你好世界！, model: code-davinci-edit-001, token: 11             |
| text: 你好世界！, model: text-embedding-ada-002, token: 6             | text: 你好世界！, model: text-embedding-ada-002, token: 6             |
| text: 你好世界！, model: text-similarity-davinci-001, token: 11       | text: 你好世界！, model: text-similarity-davinci-001, token: 11       |
| text: こんにちは世界！, model: gpt-4, token: 5                        | text: こんにちは世界！, model: gpt-4, token: 5                        |
| text: こんにちは世界！, model: gpt-3.5-turbo, token: 5                | text: こんにちは世界！, model: gpt-3.5-turbo, token: 5                |
| text: こんにちは世界！, model: text-davinci-003, token: 13            | text: こんにちは世界！, model: text-davinci-003, token: 13            |
| text: こんにちは世界！, model: text-davinci-002, token: 13            | text: こんにちは世界！, model: text-davinci-002, token: 13            |
| text: こんにちは世界！, model: text-davinci-001, token: 13            | text: こんにちは世界！, model: text-davinci-001, token: 13            |
| text: こんにちは世界！, model: text-curie-001, token: 13              | text: こんにちは世界！, model: text-curie-001, token: 13              |
| text: こんにちは世界！, model: text-babbage-001, token: 13            | text: こんにちは世界！, model: text-babbage-001, token: 13            |
| text: こんにちは世界！, model: text-ada-001, token: 13                | text: こんにちは世界！, model: text-ada-001, token: 13                |
| text: こんにちは世界！, model: davinci, token: 13                     | text: こんにちは世界！, model: davinci, token: 13                     |
| text: こんにちは世界！, model: curie, token: 13                       | text: こんにちは世界！, model: curie, token: 13                       |
| text: こんにちは世界！, model: babbage, token: 13                     | text: こんにちは世界！, model: babbage, token: 13                     |
| text: こんにちは世界！, model: ada, token: 13                         | text: こんにちは世界！, model: ada, token: 13                         |
| text: こんにちは世界！, model: code-davinci-002, token: 13            | text: こんにちは世界！, model: code-davinci-002, token: 13            |
| text: こんにちは世界！, model: code-davinci-001, token: 13            | text: こんにちは世界！, model: code-davinci-001, token: 13            |
| text: こんにちは世界！, model: code-cushman-002, token: 13            | text: こんにちは世界！, model: code-cushman-002, token: 13            |
| text: こんにちは世界！, model: code-cushman-001, token: 13            | text: こんにちは世界！, model: code-cushman-001, token: 13            |
| text: こんにちは世界！, model: davinci-codex, token: 13               | text: こんにちは世界！, model: davinci-codex, token: 13               |
| text: こんにちは世界！, model: cushman-codex, token: 13               | text: こんにちは世界！, model: cushman-codex, token: 13               |
| text: こんにちは世界！, model: text-davinci-edit-001, token: 13       | text: こんにちは世界！, model: text-davinci-edit-001, token: 13       |
| text: こんにちは世界！, model: code-davinci-edit-001, token: 13       | text: こんにちは世界！, model: code-davinci-edit-001, token: 13       |
| text: こんにちは世界！, model: text-embedding-ada-002, token: 5       | text: こんにちは世界！, model: text-embedding-ada-002, token: 5       |
| text: こんにちは世界！, model: text-similarity-davinci-001, token: 13 | text: こんにちは世界！, model: text-similarity-davinci-001, token: 13 |
| text: 안녕하세요 세계!, model: gpt-4, token: 10                       | text: 안녕하세요 세계!, model: gpt-4, token: 10                       |
| text: 안녕하세요 세계!, model: gpt-3.5-turbo, token: 10               | text: 안녕하세요 세계!, model: gpt-3.5-turbo, token: 10               |
| text: 안녕하세요 세계!, model: text-davinci-003, token: 21            | text: 안녕하세요 세계!, model: text-davinci-003, token: 21            |
| text: 안녕하세요 세계!, model: text-davinci-002, token: 21            | text: 안녕하세요 세계!, model: text-davinci-002, token: 21            |
| text: 안녕하세요 세계!, model: text-davinci-001, token: 21            | text: 안녕하세요 세계!, model: text-davinci-001, token: 21            |
| text: 안녕하세요 세계!, model: text-curie-001, token: 21              | text: 안녕하세요 세계!, model: text-curie-001, token: 21              |
| text: 안녕하세요 세계!, model: text-babbage-001, token: 21            | text: 안녕하세요 세계!, model: text-babbage-001, token: 21            |
| text: 안녕하세요 세계!, model: text-ada-001, token: 21                | text: 안녕하세요 세계!, model: text-ada-001, token: 21                |
| text: 안녕하세요 세계!, model: davinci, token: 21                     | text: 안녕하세요 세계!, model: davinci, token: 21                     |
| text: 안녕하세요 세계!, model: curie, token: 21                       | text: 안녕하세요 세계!, model: curie, token: 21                       |
| text: 안녕하세요 세계!, model: babbage, token: 21                     | text: 안녕하세요 세계!, model: babbage, token: 21                     |
| text: 안녕하세요 세계!, model: ada, token: 21                         | text: 안녕하세요 세계!, model: ada, token: 21                         |
| text: 안녕하세요 세계!, model: code-davinci-002, token: 21            | text: 안녕하세요 세계!, model: code-davinci-002, token: 21            |
| text: 안녕하세요 세계!, model: code-davinci-001, token: 21            | text: 안녕하세요 세계!, model: code-davinci-001, token: 21            |
| text: 안녕하세요 세계!, model: code-cushman-002, token: 21            | text: 안녕하세요 세계!, model: code-cushman-002, token: 21            |
| text: 안녕하세요 세계!, model: code-cushman-001, token: 21            | text: 안녕하세요 세계!, model: code-cushman-001, token: 21            |
| text: 안녕하세요 세계!, model: davinci-codex, token: 21               | text: 안녕하세요 세계!, model: davinci-codex, token: 21               |
| text: 안녕하세요 세계!, model: cushman-codex, token: 21               | text: 안녕하세요 세계!, model: cushman-codex, token: 21               |
| text: 안녕하세요 세계!, model: text-davinci-edit-001, token: 21       | text: 안녕하세요 세계!, model: text-davinci-edit-001, token: 21       |
| text: 안녕하세요 세계!, model: code-davinci-edit-001, token: 21       | text: 안녕하세요 세계!, model: code-davinci-edit-001, token: 21       |
| text: 안녕하세요 세계!, model: text-embedding-ada-002, token: 10      | text: 안녕하세요 세계!, model: text-embedding-ada-002, token: 10      |
| text: 안녕하세요 세계!, model: text-similarity-davinci-001, token: 21 | text: 안녕하세요 세계!, model: text-similarity-davinci-001, token: 21 |
| text: Привет мир!, model: gpt-4, token: 6                             | text: Привет мир!, model: gpt-4, token: 6                             |
| text: Привет мир!, model: gpt-3.5-turbo, token: 6                     | text: Привет мир!, model: gpt-3.5-turbo, token: 6                     |
| text: Привет мир!, model: text-davinci-003, token: 12                 | text: Привет мир!, model: text-davinci-003, token: 12                 |
| text: Привет мир!, model: text-davinci-002, token: 12                 | text: Привет мир!, model: text-davinci-002, token: 12                 |
| text: Привет мир!, model: text-davinci-001, token: 12                 | text: Привет мир!, model: text-davinci-001, token: 12                 |
| text: Привет мир!, model: text-curie-001, token: 12                   | text: Привет мир!, model: text-curie-001, token: 12                   |
| text: Привет мир!, model: text-babbage-001, token: 12                 | text: Привет мир!, model: text-babbage-001, token: 12                 |
| text: Привет мир!, model: text-ada-001, token: 12                     | text: Привет мир!, model: text-ada-001, token: 12                     |
| text: Привет мир!, model: davinci, token: 12                          | text: Привет мир!, model: davinci, token: 12                          |
| text: Привет мир!, model: curie, token: 12                            | text: Привет мир!, model: curie, token: 12                            |
| text: Привет мир!, model: babbage, token: 12                          | text: Привет мир!, model: babbage, token: 12                          |
| text: Привет мир!, model: ada, token: 12                              | text: Привет мир!, model: ada, token: 12                              |
| text: Привет мир!, model: code-davinci-002, token: 12                 | text: Привет мир!, model: code-davinci-002, token: 12                 |
| text: Привет мир!, model: code-davinci-001, token: 12                 | text: Привет мир!, model: code-davinci-001, token: 12                 |
| text: Привет мир!, model: code-cushman-002, token: 12                 | text: Привет мир!, model: code-cushman-002, token: 12                 |
| text: Привет мир!, model: code-cushman-001, token: 12                 | text: Привет мир!, model: code-cushman-001, token: 12                 |
| text: Привет мир!, model: davinci-codex, token: 12                    | text: Привет мир!, model: davinci-codex, token: 12                    |
| text: Привет мир!, model: cushman-codex, token: 12                    | text: Привет мир!, model: cushman-codex, token: 12                    |
| text: Привет мир!, model: text-davinci-edit-001, token: 12            | text: Привет мир!, model: text-davinci-edit-001, token: 12            |
| text: Привет мир!, model: code-davinci-edit-001, token: 12            | text: Привет мир!, model: code-davinci-edit-001, token: 12            |
| text: Привет мир!, model: text-embedding-ada-002, token: 6            | text: Привет мир!, model: text-embedding-ada-002, token: 6            |
| text: Привет мир!, model: text-similarity-davinci-001, token: 12      | text: Привет мир!, model: text-similarity-davinci-001, token: 12      |
| text: ¡Hola mundo!, model: gpt-4, token: 4                            | text: ¡Hola mundo!, model: gpt-4, token: 4                            |
| text: ¡Hola mundo!, model: gpt-3.5-turbo, token: 4                    | text: ¡Hola mundo!, model: gpt-3.5-turbo, token: 4                    |
| text: ¡Hola mundo!, model: text-davinci-003, token: 7                 | text: ¡Hola mundo!, model: text-davinci-003, token: 7                 |
| text: ¡Hola mundo!, model: text-davinci-002, token: 7                 | text: ¡Hola mundo!, model: text-davinci-002, token: 7                 |
| text: ¡Hola mundo!, model: text-davinci-001, token: 7                 | text: ¡Hola mundo!, model: text-davinci-001, token: 7                 |
| text: ¡Hola mundo!, model: text-curie-001, token: 7                   | text: ¡Hola mundo!, model: text-curie-001, token: 7                   |
| text: ¡Hola mundo!, model: text-babbage-001, token: 7                 | text: ¡Hola mundo!, model: text-babbage-001, token: 7                 |
| text: ¡Hola mundo!, model: text-ada-001, token: 7                     | text: ¡Hola mundo!, model: text-ada-001, token: 7                     |
| text: ¡Hola mundo!, model: davinci, token: 7                          | text: ¡Hola mundo!, model: davinci, token: 7                          |
| text: ¡Hola mundo!, model: curie, token: 7                            | text: ¡Hola mundo!, model: curie, token: 7                            |
| text: ¡Hola mundo!, model: babbage, token: 7                          | text: ¡Hola mundo!, model: babbage, token: 7                          |
| text: ¡Hola mundo!, model: ada, token: 7                              | text: ¡Hola mundo!, model: ada, token: 7                              |
| text: ¡Hola mundo!, model: code-davinci-002, token: 7                 | text: ¡Hola mundo!, model: code-davinci-002, token: 7                 |
| text: ¡Hola mundo!, model: code-davinci-001, token: 7                 | text: ¡Hola mundo!, model: code-davinci-001, token: 7                 |
| text: ¡Hola mundo!, model: code-cushman-002, token: 7                 | text: ¡Hola mundo!, model: code-cushman-002, token: 7                 |
| text: ¡Hola mundo!, model: code-cushman-001, token: 7                 | text: ¡Hola mundo!, model: code-cushman-001, token: 7                 |
| text: ¡Hola mundo!, model: davinci-codex, token: 7                    | text: ¡Hola mundo!, model: davinci-codex, token: 7                    |
| text: ¡Hola mundo!, model: cushman-codex, token: 7                    | text: ¡Hola mundo!, model: cushman-codex, token: 7                    |
| text: ¡Hola mundo!, model: text-davinci-edit-001, token: 7            | text: ¡Hola mundo!, model: text-davinci-edit-001, token: 7            |
| text: ¡Hola mundo!, model: code-davinci-edit-001, token: 7            | text: ¡Hola mundo!, model: code-davinci-edit-001, token: 7            |
| text: ¡Hola mundo!, model: text-embedding-ada-002, token: 4           | text: ¡Hola mundo!, model: text-embedding-ada-002, token: 4           |
| text: ¡Hola mundo!, model: text-similarity-davinci-001, token: 7      | text: ¡Hola mundo!, model: text-similarity-davinci-001, token: 7      |
| text: Hallo Welt!, model: gpt-4, token: 3                             | text: Hallo Welt!, model: gpt-4, token: 3                             |
| text: Hallo Welt!, model: gpt-3.5-turbo, token: 3                     | text: Hallo Welt!, model: gpt-3.5-turbo, token: 3                     |
| text: Hallo Welt!, model: text-davinci-003, token: 5                  | text: Hallo Welt!, model: text-davinci-003, token: 5                  |
| text: Hallo Welt!, model: text-davinci-002, token: 5                  | text: Hallo Welt!, model: text-davinci-002, token: 5                  |
| text: Hallo Welt!, model: text-davinci-001, token: 5                  | text: Hallo Welt!, model: text-davinci-001, token: 5                  |
| text: Hallo Welt!, model: text-curie-001, token: 5                    | text: Hallo Welt!, model: text-curie-001, token: 5                    |
| text: Hallo Welt!, model: text-babbage-001, token: 5                  | text: Hallo Welt!, model: text-babbage-001, token: 5                  |
| text: Hallo Welt!, model: text-ada-001, token: 5                      | text: Hallo Welt!, model: text-ada-001, token: 5                      |
| text: Hallo Welt!, model: davinci, token: 5                           | text: Hallo Welt!, model: davinci, token: 5                           |
| text: Hallo Welt!, model: curie, token: 5                             | text: Hallo Welt!, model: curie, token: 5                             |
| text: Hallo Welt!, model: babbage, token: 5                           | text: Hallo Welt!, model: babbage, token: 5                           |
| text: Hallo Welt!, model: ada, token: 5                               | text: Hallo Welt!, model: ada, token: 5                               |
| text: Hallo Welt!, model: code-davinci-002, token: 5                  | text: Hallo Welt!, model: code-davinci-002, token: 5                  |
| text: Hallo Welt!, model: code-davinci-001, token: 5                  | text: Hallo Welt!, model: code-davinci-001, token: 5                  |
| text: Hallo Welt!, model: code-cushman-002, token: 5                  | text: Hallo Welt!, model: code-cushman-002, token: 5                  |
| text: Hallo Welt!, model: code-cushman-001, token: 5                  | text: Hallo Welt!, model: code-cushman-001, token: 5                  |
| text: Hallo Welt!, model: davinci-codex, token: 5                     | text: Hallo Welt!, model: davinci-codex, token: 5                     |
| text: Hallo Welt!, model: cushman-codex, token: 5                     | text: Hallo Welt!, model: cushman-codex, token: 5                     |
| text: Hallo Welt!, model: text-davinci-edit-001, token: 5             | text: Hallo Welt!, model: text-davinci-edit-001, token: 5             |
| text: Hallo Welt!, model: code-davinci-edit-001, token: 5             | text: Hallo Welt!, model: code-davinci-edit-001, token: 5             |
| text: Hallo Welt!, model: text-embedding-ada-002, token: 3            | text: Hallo Welt!, model: text-embedding-ada-002, token: 3            |
| text: Hallo Welt!, model: text-similarity-davinci-001, token: 5       | text: Hallo Welt!, model: text-similarity-davinci-001, token: 5       |
| text: Bonjour le monde!, model: gpt-4, token: 4                       | text: Bonjour le monde!, model: gpt-4, token: 4                       |
| text: Bonjour le monde!, model: gpt-3.5-turbo, token: 4               | text: Bonjour le monde!, model: gpt-3.5-turbo, token: 4               |
| text: Bonjour le monde!, model: text-davinci-003, token: 7            | text: Bonjour le monde!, model: text-davinci-003, token: 7            |
| text: Bonjour le monde!, model: text-davinci-002, token: 7            | text: Bonjour le monde!, model: text-davinci-002, token: 7            |
| text: Bonjour le monde!, model: text-davinci-001, token: 7            | text: Bonjour le monde!, model: text-davinci-001, token: 7            |
| text: Bonjour le monde!, model: text-curie-001, token: 7              | text: Bonjour le monde!, model: text-curie-001, token: 7              |
| text: Bonjour le monde!, model: text-babbage-001, token: 7            | text: Bonjour le monde!, model: text-babbage-001, token: 7            |
| text: Bonjour le monde!, model: text-ada-001, token: 7                | text: Bonjour le monde!, model: text-ada-001, token: 7                |
| text: Bonjour le monde!, model: davinci, token: 7                     | text: Bonjour le monde!, model: davinci, token: 7                     |
| text: Bonjour le monde!, model: curie, token: 7                       | text: Bonjour le monde!, model: curie, token: 7                       |
| text: Bonjour le monde!, model: babbage, token: 7                     | text: Bonjour le monde!, model: babbage, token: 7                     |
| text: Bonjour le monde!, model: ada, token: 7                         | text: Bonjour le monde!, model: ada, token: 7                         |
| text: Bonjour le monde!, model: code-davinci-002, token: 7            | text: Bonjour le monde!, model: code-davinci-002, token: 7            |
| text: Bonjour le monde!, model: code-davinci-001, token: 7            | text: Bonjour le monde!, model: code-davinci-001, token: 7            |
| text: Bonjour le monde!, model: code-cushman-002, token: 7            | text: Bonjour le monde!, model: code-cushman-002, token: 7            |
| text: Bonjour le monde!, model: code-cushman-001, token: 7            | text: Bonjour le monde!, model: code-cushman-001, token: 7            |
| text: Bonjour le monde!, model: davinci-codex, token: 7               | text: Bonjour le monde!, model: davinci-codex, token: 7               |
| text: Bonjour le monde!, model: cushman-codex, token: 7               | text: Bonjour le monde!, model: cushman-codex, token: 7               |
| text: Bonjour le monde!, model: text-davinci-edit-001, token: 7       | text: Bonjour le monde!, model: text-davinci-edit-001, token: 7       |
| text: Bonjour le monde!, model: code-davinci-edit-001, token: 7       | text: Bonjour le monde!, model: code-davinci-edit-001, token: 7       |
| text: Bonjour le monde!, model: text-embedding-ada-002, token: 4      | text: Bonjour le monde!, model: text-embedding-ada-002, token: 4      |
| text: Bonjour le monde!, model: text-similarity-davinci-001, token: 7 | text: Bonjour le monde!, model: text-similarity-davinci-001, token: 7 |
| text: Ciao mondo!, model: gpt-4, token: 4                             | text: Ciao mondo!, model: gpt-4, token: 4                             |
| text: Ciao mondo!, model: gpt-3.5-turbo, token: 4                     | text: Ciao mondo!, model: gpt-3.5-turbo, token: 4                     |
| text: Ciao mondo!, model: text-davinci-003, token: 5                  | text: Ciao mondo!, model: text-davinci-003, token: 5                  |
| text: Ciao mondo!, model: text-davinci-002, token: 5                  | text: Ciao mondo!, model: text-davinci-002, token: 5                  |
| text: Ciao mondo!, model: text-davinci-001, token: 5                  | text: Ciao mondo!, model: text-davinci-001, token: 5                  |
| text: Ciao mondo!, model: text-curie-001, token: 5                    | text: Ciao mondo!, model: text-curie-001, token: 5                    |
| text: Ciao mondo!, model: text-babbage-001, token: 5                  | text: Ciao mondo!, model: text-babbage-001, token: 5                  |
| text: Ciao mondo!, model: text-ada-001, token: 5                      | text: Ciao mondo!, model: text-ada-001, token: 5                      |
| text: Ciao mondo!, model: davinci, token: 5                           | text: Ciao mondo!, model: davinci, token: 5                           |
| text: Ciao mondo!, model: curie, token: 5                             | text: Ciao mondo!, model: curie, token: 5                             |
| text: Ciao mondo!, model: babbage, token: 5                           | text: Ciao mondo!, model: babbage, token: 5                           |
| text: Ciao mondo!, model: ada, token: 5                               | text: Ciao mondo!, model: ada, token: 5                               |
| text: Ciao mondo!, model: code-davinci-002, token: 5                  | text: Ciao mondo!, model: code-davinci-002, token: 5                  |
| text: Ciao mondo!, model: code-davinci-001, token: 5                  | text: Ciao mondo!, model: code-davinci-001, token: 5                  |
| text: Ciao mondo!, model: code-cushman-002, token: 5                  | text: Ciao mondo!, model: code-cushman-002, token: 5                  |
| text: Ciao mondo!, model: code-cushman-001, token: 5                  | text: Ciao mondo!, model: code-cushman-001, token: 5                  |
| text: Ciao mondo!, model: davinci-codex, token: 5                     | text: Ciao mondo!, model: davinci-codex, token: 5                     |
| text: Ciao mondo!, model: cushman-codex, token: 5                     | text: Ciao mondo!, model: cushman-codex, token: 5                     |
| text: Ciao mondo!, model: text-davinci-edit-001, token: 5             | text: Ciao mondo!, model: text-davinci-edit-001, token: 5             |
| text: Ciao mondo!, model: code-davinci-edit-001, token: 5             | text: Ciao mondo!, model: code-davinci-edit-001, token: 5             |
| text: Ciao mondo!, model: text-embedding-ada-002, token: 4            | text: Ciao mondo!, model: text-embedding-ada-002, token: 4            |
| text: Ciao mondo!, model: text-similarity-davinci-001, token: 5       | text: Ciao mondo!, model: text-similarity-davinci-001, token: 5       |
| text: Hej världen!, model: gpt-4, token: 7                            | text: Hej världen!, model: gpt-4, token: 7                            |
| text: Hej världen!, model: gpt-3.5-turbo, token: 7                    | text: Hej världen!, model: gpt-3.5-turbo, token: 7                    |
| text: Hej världen!, model: text-davinci-003, token: 8                 | text: Hej världen!, model: text-davinci-003, token: 8                 |
| text: Hej världen!, model: text-davinci-002, token: 8                 | text: Hej världen!, model: text-davinci-002, token: 8                 |
| text: Hej världen!, model: text-davinci-001, token: 8                 | text: Hej världen!, model: text-davinci-001, token: 8                 |
| text: Hej världen!, model: text-curie-001, token: 8                   | text: Hej världen!, model: text-curie-001, token: 8                   |
| text: Hej världen!, model: text-babbage-001, token: 8                 | text: Hej världen!, model: text-babbage-001, token: 8                 |
| text: Hej världen!, model: text-ada-001, token: 8                     | text: Hej världen!, model: text-ada-001, token: 8                     |
| text: Hej världen!, model: davinci, token: 8                          | text: Hej världen!, model: davinci, token: 8                          |
| text: Hej världen!, model: curie, token: 8                            | text: Hej världen!, model: curie, token: 8                            |
| text: Hej världen!, model: babbage, token: 8                          | text: Hej världen!, model: babbage, token: 8                          |
| text: Hej världen!, model: ada, token: 8                              | text: Hej världen!, model: ada, token: 8                              |
| text: Hej världen!, model: code-davinci-002, token: 8                 | text: Hej världen!, model: code-davinci-002, token: 8                 |
| text: Hej världen!, model: code-davinci-001, token: 8                 | text: Hej världen!, model: code-davinci-001, token: 8                 |
| text: Hej världen!, model: code-cushman-002, token: 8                 | text: Hej världen!, model: code-cushman-002, token: 8                 |
| text: Hej världen!, model: code-cushman-001, token: 8                 | text: Hej världen!, model: code-cushman-001, token: 8                 |
| text: Hej världen!, model: davinci-codex, token: 8                    | text: Hej världen!, model: davinci-codex, token: 8                    |
| text: Hej världen!, model: cushman-codex, token: 8                    | text: Hej världen!, model: cushman-codex, token: 8                    |
| text: Hej världen!, model: text-davinci-edit-001, token: 8            | text: Hej världen!, model: text-davinci-edit-001, token: 8            |
| text: Hej världen!, model: code-davinci-edit-001, token: 8            | text: Hej världen!, model: code-davinci-edit-001, token: 8            |
| text: Hej världen!, model: text-embedding-ada-002, token: 7           | text: Hej världen!, model: text-embedding-ada-002, token: 7           |
| text: Hej världen!, model: text-similarity-davinci-001, token: 8      | text: Hej världen!, model: text-similarity-davinci-001, token: 8      |
| text: Hallo wereld!, model: gpt-4, token: 3                           | text: Hallo wereld!, model: gpt-4, token: 3                           |
| text: Hallo wereld!, model: gpt-3.5-turbo, token: 3                   | text: Hallo wereld!, model: gpt-3.5-turbo, token: 3                   |
| text: Hallo wereld!, model: text-davinci-003, token: 5                | text: Hallo wereld!, model: text-davinci-003, token: 5                |
| text: Hallo wereld!, model: text-davinci-002, token: 5                | text: Hallo wereld!, model: text-davinci-002, token: 5                |
| text: Hallo wereld!, model: text-davinci-001, token: 5                | text: Hallo wereld!, model: text-davinci-001, token: 5                |
| text: Hallo wereld!, model: text-curie-001, token: 5                  | text: Hallo wereld!, model: text-curie-001, token: 5                  |
| text: Hallo wereld!, model: text-babbage-001, token: 5                | text: Hallo wereld!, model: text-babbage-001, token: 5                |
| text: Hallo wereld!, model: text-ada-001, token: 5                    | text: Hallo wereld!, model: text-ada-001, token: 5                    |
| text: Hallo wereld!, model: davinci, token: 5                         | text: Hallo wereld!, model: davinci, token: 5                         |
| text: Hallo wereld!, model: curie, token: 5                           | text: Hallo wereld!, model: curie, token: 5                           |
| text: Hallo wereld!, model: babbage, token: 5                         | text: Hallo wereld!, model: babbage, token: 5                         |
| text: Hallo wereld!, model: ada, token: 5                             | text: Hallo wereld!, model: ada, token: 5                             |
| text: Hallo wereld!, model: code-davinci-002, token: 5                | text: Hallo wereld!, model: code-davinci-002, token: 5                |
| text: Hallo wereld!, model: code-davinci-001, token: 5                | text: Hallo wereld!, model: code-davinci-001, token: 5                |
| text: Hallo wereld!, model: code-cushman-002, token: 5                | text: Hallo wereld!, model: code-cushman-002, token: 5                |
| text: Hallo wereld!, model: code-cushman-001, token: 5                | text: Hallo wereld!, model: code-cushman-001, token: 5                |
| text: Hallo wereld!, model: davinci-codex, token: 5                   | text: Hallo wereld!, model: davinci-codex, token: 5                   |
| text: Hallo wereld!, model: cushman-codex, token: 5                   | text: Hallo wereld!, model: cushman-codex, token: 5                   |
| text: Hallo wereld!, model: text-davinci-edit-001, token: 5           | text: Hallo wereld!, model: text-davinci-edit-001, token: 5           |
| text: Hallo wereld!, model: code-davinci-edit-001, token: 5           | text: Hallo wereld!, model: code-davinci-edit-001, token: 5           |
| text: Hallo wereld!, model: text-embedding-ada-002, token: 3          | text: Hallo wereld!, model: text-embedding-ada-002, token: 3          |
| text: Hallo wereld!, model: text-similarity-davinci-001, token: 5     | text: Hallo wereld!, model: text-similarity-davinci-001, token: 5     |
| text: Hallo verden!, model: gpt-4, token: 4                           | text: Hallo verden!, model: gpt-4, token: 4                           |
| text: Hallo verden!, model: gpt-3.5-turbo, token: 4                   | text: Hallo verden!, model: gpt-3.5-turbo, token: 4                   |
| text: Hallo verden!, model: text-davinci-003, token: 5                | text: Hallo verden!, model: text-davinci-003, token: 5                |
| text: Hallo verden!, model: text-davinci-002, token: 5                | text: Hallo verden!, model: text-davinci-002, token: 5                |
| text: Hallo verden!, model: text-davinci-001, token: 5                | text: Hallo verden!, model: text-davinci-001, token: 5                |
| text: Hallo verden!, model: text-curie-001, token: 5                  | text: Hallo verden!, model: text-curie-001, token: 5                  |
| text: Hallo verden!, model: text-babbage-001, token: 5                | text: Hallo verden!, model: text-babbage-001, token: 5                |
| text: Hallo verden!, model: text-ada-001, token: 5                    | text: Hallo verden!, model: text-ada-001, token: 5                    |
| text: Hallo verden!, model: davinci, token: 5                         | text: Hallo verden!, model: davinci, token: 5                         |
| text: Hallo verden!, model: curie, token: 5                           | text: Hallo verden!, model: curie, token: 5                           |
| text: Hallo verden!, model: babbage, token: 5                         | text: Hallo verden!, model: babbage, token: 5                         |
| text: Hallo verden!, model: ada, token: 5                             | text: Hallo verden!, model: ada, token: 5                             |
| text: Hallo verden!, model: code-davinci-002, token: 5                | text: Hallo verden!, model: code-davinci-002, token: 5                |
| text: Hallo verden!, model: code-davinci-001, token: 5                | text: Hallo verden!, model: code-davinci-001, token: 5                |
| text: Hallo verden!, model: code-cushman-002, token: 5                | text: Hallo verden!, model: code-cushman-002, token: 5                |
| text: Hallo verden!, model: code-cushman-001, token: 5                | text: Hallo verden!, model: code-cushman-001, token: 5                |
| text: Hallo verden!, model: davinci-codex, token: 5                   | text: Hallo verden!, model: davinci-codex, token: 5                   |
| text: Hallo verden!, model: cushman-codex, token: 5                   | text: Hallo verden!, model: cushman-codex, token: 5                   |
| text: Hallo verden!, model: text-davinci-edit-001, token: 5           | text: Hallo verden!, model: text-davinci-edit-001, token: 5           |
| text: Hallo verden!, model: code-davinci-edit-001, token: 5           | text: Hallo verden!, model: code-davinci-edit-001, token: 5           |
| text: Hallo verden!, model: text-embedding-ada-002, token: 4          | text: Hallo verden!, model: text-embedding-ada-002, token: 4          |
| text: Hallo verden!, model: text-similarity-davinci-001, token: 5     | text: Hallo verden!, model: text-similarity-davinci-001, token: 5     |
| text: Hallo wereld!, model: gpt-4, token: 3                           | text: Hallo wereld!, model: gpt-4, token: 3                           |
| text: Hallo wereld!, model: gpt-3.5-turbo, token: 3                   | text: Hallo wereld!, model: gpt-3.5-turbo, token: 3                   |
| text: Hallo wereld!, model: text-davinci-003, token: 5                | text: Hallo wereld!, model: text-davinci-003, token: 5                |
| text: Hallo wereld!, model: text-davinci-002, token: 5                | text: Hallo wereld!, model: text-davinci-002, token: 5                |
| text: Hallo wereld!, model: text-davinci-001, token: 5                | text: Hallo wereld!, model: text-davinci-001, token: 5                |
| text: Hallo wereld!, model: text-curie-001, token: 5                  | text: Hallo wereld!, model: text-curie-001, token: 5                  |
| text: Hallo wereld!, model: text-babbage-001, token: 5                | text: Hallo wereld!, model: text-babbage-001, token: 5                |
| text: Hallo wereld!, model: text-ada-001, token: 5                    | text: Hallo wereld!, model: text-ada-001, token: 5                    |
| text: Hallo wereld!, model: davinci, token: 5                         | text: Hallo wereld!, model: davinci, token: 5                         |
| text: Hallo wereld!, model: curie, token: 5                           | text: Hallo wereld!, model: curie, token: 5                           |
| text: Hallo wereld!, model: babbage, token: 5                         | text: Hallo wereld!, model: babbage, token: 5                         |
| text: Hallo wereld!, model: ada, token: 5                             | text: Hallo wereld!, model: ada, token: 5                             |
| text: Hallo wereld!, model: code-davinci-002, token: 5                | text: Hallo wereld!, model: code-davinci-002, token: 5                |
| text: Hallo wereld!, model: code-davinci-001, token: 5                | text: Hallo wereld!, model: code-davinci-001, token: 5                |
| text: Hallo wereld!, model: code-cushman-002, token: 5                | text: Hallo wereld!, model: code-cushman-002, token: 5                |
| text: Hallo wereld!, model: code-cushman-001, token: 5                | text: Hallo wereld!, model: code-cushman-001, token: 5                |
| text: Hallo wereld!, model: davinci-codex, token: 5                   | text: Hallo wereld!, model: davinci-codex, token: 5                   |
| text: Hallo wereld!, model: cushman-codex, token: 5                   | text: Hallo wereld!, model: cushman-codex, token: 5                   |
| text: Hallo wereld!, model: text-davinci-edit-001, token: 5           | text: Hallo wereld!, model: text-davinci-edit-001, token: 5           |
| text: Hallo wereld!, model: code-davinci-edit-001, token: 5           | text: Hallo wereld!, model: code-davinci-edit-001, token: 5           |
| text: Hallo wereld!, model: text-embedding-ada-002, token: 3          | text: Hallo wereld!, model: text-embedding-ada-002, token: 3          |
| text: Hallo wereld!, model: text-similarity-davinci-001, token: 5     | text: Hallo wereld!, model: text-similarity-davinci-001, token: 5     |
| text: Hallo verden!, model: gpt-4, token: 4                           | text: Hallo verden!, model: gpt-4, token: 4                           |
| text: Hallo verden!, model: gpt-3.5-turbo, token: 4                   | text: Hallo verden!, model: gpt-3.5-turbo, token: 4                   |
| text: Hallo verden!, model: text-davinci-003, token: 5                | text: Hallo verden!, model: text-davinci-003, token: 5                |
| text: Hallo verden!, model: text-davinci-002, token: 5                | text: Hallo verden!, model: text-davinci-002, token: 5                |
| text: Hallo verden!, model: text-davinci-001, token: 5                | text: Hallo verden!, model: text-davinci-001, token: 5                |
| text: Hallo verden!, model: text-curie-001, token: 5                  | text: Hallo verden!, model: text-curie-001, token: 5                  |
| text: Hallo verden!, model: text-babbage-001, token: 5                | text: Hallo verden!, model: text-babbage-001, token: 5                |
| text: Hallo verden!, model: text-ada-001, token: 5                    | text: Hallo verden!, model: text-ada-001, token: 5                    |
| text: Hallo verden!, model: davinci, token: 5                         | text: Hallo verden!, model: davinci, token: 5                         |
| text: Hallo verden!, model: curie, token: 5                           | text: Hallo verden!, model: curie, token: 5                           |
| text: Hallo verden!, model: babbage, token: 5                         | text: Hallo verden!, model: babbage, token: 5                         |
| text: Hallo verden!, model: ada, token: 5                             | text: Hallo verden!, model: ada, token: 5                             |
| text: Hallo verden!, model: code-davinci-002, token: 5                | text: Hallo verden!, model: code-davinci-002, token: 5                |
| text: Hallo verden!, model: code-davinci-001, token: 5                | text: Hallo verden!, model: code-davinci-001, token: 5                |
| text: Hallo verden!, model: code-cushman-002, token: 5                | text: Hallo verden!, model: code-cushman-002, token: 5                |
| text: Hallo verden!, model: code-cushman-001, token: 5                | text: Hallo verden!, model: code-cushman-001, token: 5                |
| text: Hallo verden!, model: davinci-codex, token: 5                   | text: Hallo verden!, model: davinci-codex, token: 5                   |
| text: Hallo verden!, model: cushman-codex, token: 5                   | text: Hallo verden!, model: cushman-codex, token: 5                   |
| text: Hallo verden!, model: text-davinci-edit-001, token: 5           | text: Hallo verden!, model: text-davinci-edit-001, token: 5           |
| text: Hallo verden!, model: code-davinci-edit-001, token: 5           | text: Hallo verden!, model: code-davinci-edit-001, token: 5           |
| text: Hallo verden!, model: text-embedding-ada-002, token: 4          | text: Hallo verden!, model: text-embedding-ada-002, token: 4          |
| text: Hallo verden!, model: text-similarity-davinci-001, token: 5     | text: Hallo verden!, model: text-similarity-davinci-001, token: 5     |


# License
[MIT](./LICENSE)
