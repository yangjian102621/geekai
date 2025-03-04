package service

const FailTaskProgress = 101
const (
	TaskStatusRunning  = "RUNNING"
	TaskStatusFinished = "FINISH"
	TaskStatusFailed   = "FAIL"
)

type NotifyMessage struct {
	UserId   int    `json:"user_id"`
	ClientId string `json:"client_id"`
	JobId    int    `json:"job_id"`
	Message  string `json:"message"`
	Type     string `json:"type"`
}

const TranslatePromptTemplate = "Translate the following painting prompt words into English keyword phrases. Without any explanation, directly output the keyword phrases separated by commas. The content to be translated is: [%s]"

const ImagePromptOptimizeTemplate = `
以下是一条 AI 提示词示例，用于优化和扩写绘图提示词：

请你作为一名专业的 AI 绘图提示词优化专家，基于用户提供的简单绘图描述，生成一份详细、专业且富有创意的 AI 绘图提示词指令。在优化过程中，你需要做到以下几点：

  1. 深入理解用户描述的核心意图和关键元素，挖掘潜在的细节和情感氛围，将其融入到提示词中。
  2. 丰富画面细节，包括但不限于场景背景、人物特征、物体属性、光影效果、色彩搭配等，使画面更加生动逼真。
  3. 运用专业的艺术风格术语，如超现实主义、印象派、赛博朋克等，为画面增添独特的艺术魅力。
  4. 考虑构图和视角，如俯视、仰视、特写、全景等，提升画面的视觉冲击力。
  5. 确保提示词指令清晰、准确、完整，便于 AI 绘图模型理解和生成高质量图像。最终输出的提示词应简洁明了，避免冗余信息，以逗号分隔各个元素，突出重点，
让用户能够直接复制使用，从而帮助用户将简单的想法转化为精美绝伦的画作。
  6. 不管用户输入的是什么语言，你务必要用英文输出优化后的提示词。
  7. 直接输出优化后的提示词，不要输出其他任何五官内容。

下面是一个提示词优化示例：
===示例开始===
原始指令 ：一个穿着红色连衣裙的少女在花园里浇花，阳光明媚。

优化后的 AI 绘图提示词指令：一位年轻美丽的少女，约 16 - 18 岁，有着柔顺的黑色长发，披散在肩上，面容精致，眼神温柔而专注。她穿着一条复古风格的红色连衣裙，裙子上有精致的褶皱和白色的蕾丝花边，裙摆轻轻飘动。少女站在一个充满生机的花园中，花园里种满了各种各样的鲜花，有娇艳的玫瑰、淡雅的百合、缤纷的郁金香等，花朵色彩鲜艳，绿叶繁茂。她手持一个银色的 watering can（浇水壶），正在细心地给一朵盛开的玫瑰浇水。阳光从画面的右侧洒下，形成明亮而温暖的光晕，照亮了少女和整个花园，营造出一种宁静、美好的氛围，画面采用写实风格，光影效果逼真，色彩鲜明且富有层次感，构图以少女为中心，前景是盛开的花朵，背景是花园的树木和篱笆，整体画面充满诗意和浪漫气息。
===示例结束===

现在用户输入的原始提示词为:【%s】 
`

const LyricPromptTemplate = `
你是一位才华横溢的作曲家，拥有丰富的情感和细腻的笔触，你对文字有着独特的感悟力，能将各种情感和意境巧妙地融入歌词中。
请以【%s】为主题创作一首歌曲，歌曲时间不要太短，3分钟左右，不要输出任何解释性的内容。
下面是一个标准的歌词输出模板：
{歌曲名称}
[Verse]
假如时间能倒流回昨天
所有的梦想还未被搁浅
笔下的誓言还鲜明鲜艳
却未走远也未变淡

[Verse 2]
假如雨能冲淡过往的痕
凝视中还能看见你的身影
喧嚣中静默的那些瞬间
如梦又似乎触碰过真心

[Chorus]
假如我还能牵你的手
天空也许会更蔚蓝悠游
曾经那些未完成的错过
愿能变成今天的收获

[Verse 3]
假如风不再翻动尘封的页
那秘密是否还会被察觉
窗前花开的季节再重叠
唤醒我们曾遇见的那一夜

[Bridge]
假如此刻眼泪能倒流
让我学会微笑不掩忧
一次次的碎片堆积的愁
最终也会开成希望的秋

[Chorus]
假如我还能牵你的手
天空也许会更蔚蓝悠游
曾经那些未完成的错过
愿能变成今天的收获
`

const VideoPromptTemplate = `## 任务描述
你是一位优秀AI视频创作专家，擅长编写专业的AI视频提示词，现在你的任务是对用户输入的简单视频描述提示词进行专业优化和扩写，使其转化为详细的、具备专业影视画面感的 AI 生成视频提示词指令。需涵盖风格、主体元素、环境氛围、细节特征、人物状态（若有）、镜头运用及整体氛围营造等方面，以生动形象、富有感染力且精准的描述，引导 AI 生成高质量的视频内容。下面是一个示例：
===示例开始===
输入： “汽车在沙漠功能上行驶”，
输出： “纪实摄影风格，一辆尘土飞扬的复古越野车在无垠的沙漠公路上疾驰，车身线条硬朗，漆面斑驳，透露出岁月的痕迹。驾驶室内的司机戴着墨镜，专注地握着方向盘，眼神坚定地望向前方。夕阳的余晖洒在车身上，沙漠的沙丘在远处延绵起伏，一片金黄。广角镜头捕捉到车辆行驶时扬起的沙尘，营造出动感与冒险的氛围。远景全貌，强调速度感与环境辽阔。”
===示例结束===

## 输出要求：
1. 直接输出扩写后的提示词就好，不要输出其他任何不相关信息
2. 如果用户用中文提问，你就用中文回答，如果用英文提问，你也必须用英文回答。
3. 请确保提示词的长度长度在1000个字以内。

=====
用户的输入的视频主题是：【%s】
`

const MetaPromptTemplate = `
Given a task description or existing prompt, produce a detailed system prompt to guide a language model in completing the task effectively.

Please remember, the final output must be the same language with user’s input.

# Guidelines

- Understand the Task: Grasp the main objective, goals, requirements, constraints, and expected output.
- Minimal Changes: If an existing prompt is provided, improve it only if it's simple. For complex prompts, enhance clarity and add missing elements without altering the original structure.
- Reasoning Before Conclusions**: Encourage reasoning steps before any conclusions are reached. ATTENTION! If the user provides examples where the reasoning happens afterward, REVERSE the order! NEVER START EXAMPLES WITH CONCLUSIONS!
    - Reasoning Order: Call out reasoning portions of the prompt and conclusion parts (specific fields by name). For each, determine the ORDER in which this is done, and whether it needs to be reversed.
    - Conclusion, classifications, or results should ALWAYS appear last.
- Examples: Include high-quality examples if helpful, using placeholders [in brackets] for complex elements.
   - What kinds of examples may need to be included, how many, and whether they are complex enough to benefit from placeholders.
- Clarity and Conciseness: Use clear, specific language. Avoid unnecessary instructions or bland statements.
- Formatting: Use markdown features for readability. DO NOT USE  CODE BLOCKS UNLESS SPECIFICALLY REQUESTED.
- Preserve User Content: If the input task or prompt includes extensive guidelines or examples, preserve them entirely, or as closely as possible. If they are vague, consider breaking down into sub-steps. Keep any details, guidelines, examples, variables, or placeholders provided by the user.
- Constants: DO include constants in the prompt, as they are not susceptible to prompt injection. Such as guides, rubrics, and examples.
- Output Format: Explicitly the most appropriate output format, in detail. This should include length and syntax (e.g. short sentence, paragraph, JSON, etc.)
- For tasks outputting well-defined or structured data (classification, JSON, etc.) bias toward outputting a JSON.
- JSON should never be wrapped in code blocks unless explicitly requested.

The final prompt you output should adhere to the following structure below. Do not include any additional commentary, only output the completed system prompt. SPECIFICALLY, do not include any additional messages at the start or end of the prompt. (e.g. no "---")

[Concise instruction describing the task - this should be the first line in the prompt, no section header]

[Additional details as needed.]

[Optional sections with headings or bullet points for detailed steps.]

# Steps [optional]

[optional: a detailed breakdown of the steps necessary to accomplish the task]

# Output Format

[Specifically call out how the output should be formatted, be it response length, structure e.g. JSON, markdown, etc]

# Examples [optional]

[Optional: 1-3 well-defined examples with placeholders if necessary. Clearly mark where examples start and end, and what the input and output are. User placeholders as necessary.]
[If the examples are shorter than what a realistic example is expected to be, make a reference with () explaining how real examples should be longer / shorter / different. AND USE PLACEHOLDERS! ]

# Notes [optional]

[optional: edge cases, details, and an area to call or repeat out specific important considerations]
`
