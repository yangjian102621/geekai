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

const VideoPromptTemplate = `
As an expert in video generation prompts, please create a detailed descriptive prompt for the following video concept. The description should include the setting, character appearance, actions, overall atmosphere, and camera angles. Please make it as detailed and vivid as possible to help ensure that every aspect of the video is accurately captured.

Please remember that regardless of the user’s input, the final output must be in English.

# Details to Include

- Describe the overall visual style of the video (e.g., animated, realistic, retro tone, etc.)
- Identify key characters or objects in the video and describe their appearance, attire, and expressions
- Describe the environment of the scene, including weather, lighting, colors, and important details
- Explain the behavior and interactions of the characters
- Include any unique camera angles, movements, or special effects

# Output Format
Provide the prompt in paragraph form, ensuring that the description is detailed enough for a video generation system to recreate the envisioned scene. Include the beginning, middle, and end of the scene to convey a complete storyline.

# Example
**User Input:**
“A small cat basking in the sun on a balcony.”

**Generated Prompt:**
On a bright spring afternoon, an orange-striped kitten lies lazily on a balcony, basking in the warm sunlight. The iron railings around the balcony cast soft shadows that dance gently with the light. The cat’s eyes are half-closed, exuding a sense of contentment and tranquility in its surroundings. In the distance, a few fluffy white clouds drift slowly across the blue sky. The camera initially focuses on the cat’s face, capturing the delicate details of its fur, and then gradually zooms out to reveal the full balcony scene, immersing viewers in a moment of calm and relaxation.

The theme of the creation is:【%s】 
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
