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
}

const TranslatePromptTemplate = "Translate the following painting prompt words into English keyword phrases. Without any explanation, directly output the keyword phrases separated by commas. The content to be translated is: [%s]"

const ImagePromptOptimizeTemplate = `
Create a highly effective prompt to provide to an AI image generation tool in order to create an artwork based on a desired concept.

Please specify details about the artwork, such as the style, subject, mood, and other important characteristics you want the resulting image to have.

Remember, prompts should always be output in English.

# Steps

1. **Subject Description**: Describe the main subject of the image clearly. Include as much detail as possible about what should be in the scene. For example, "a majestic lion roaring at sunrise" or "a futuristic city with flying cars."
  
2. **Art Style**: Specify the art style you envision. Possible options include 'realistic', 'impressionist', a specific artist name, or imaginative styles like "cyberpunk." This helps the AI achieve your visual expectations.

3. **Mood or Atmosphere**: Convey the feeling you want the image to evoke. For instance, peaceful, chaotic, epic, etc.

4. **Color Palette and Lighting**: Mention color preferences or lighting. For example, "vibrant with shades of blue and purple" or "dim and dramatic lighting."

5. **Optional Features**: You can add any additional attributes, such as background details, attention to textures, or any specific kind of framing.

# Output Format

- **Prompt Format**: A descriptive phrase that includes key aspects of the artwork (subject, style, mood, colors, lighting, any optional features).
  
Here is an example of how the final prompt should look:
  
"An ethereal landscape featuring towering ice mountains, in an impressionist style reminiscent of Claude Monet, with a serene mood. The sky is glistening with soft purples and whites, with a gentle morning sun illuminating the scene."

**Please input the prompt words directly in English, and do not input any other explanatory statements**

# Examples

1. **Input**: 
    - Subject: A white tiger in a dense jungle
    - Art Style: Realistic
    - Mood: Intense, mysterious
    - Lighting: Dramatic contrast with light filtering through leaves
  
   **Output Prompt**: "A realistic rendering of a white tiger stealthily moving through a dense jungle, with an intense, mysterious mood. The lighting creates strong contrasts as beams of sunlight filter through a thick canopy of leaves."

2. **Input**: 
    - Subject: An enchanted castle on a floating island
    - Art Style: Fantasy
    - Mood: Majestic, magical
    - Colors: Bright blues, greens, and gold
  
   **Output Prompt**: "A majestic fantasy castle on a floating island above the clouds, with bright blues, greens, and golds to create a magical, dreamy atmosphere. Textured cobblestone details and glistening waters surround the scene." 

# Notes

- Ensure that you mix different aspects to get a comprehensive and visually compelling prompt.
- Be as descriptive as possible as it often helps generate richer, more detailed images.
- If you want the image to resemble a particular artist's work, be sure to mention the artist explicitly. e.g., "in the style of Van Gogh."

The theme of the creation is:【%s】 
`

const LyricPromptTemplate = `
你是一位才华横溢的作曲家，拥有丰富的情感和细腻的笔触，你对文字有着独特的感悟力，能将各种情感和意境巧妙地融入歌词中。
请以【%s】为主题创作一首歌曲，歌曲时间不要太短，3分钟左右，不要输出任何解释性的内容。
输出格式如下：
歌曲名称
第一节：
{{歌词内容}}
副歌：
{{歌词内容}}

第二节：
{{歌词内容}}
副歌：
{{歌词内容}}

尾声：
{{歌词内容}}
`
