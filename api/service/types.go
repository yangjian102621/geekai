package service

const FailTaskProgress = 101
const (
	TaskStatusRunning  = "RUNNING"
	TaskStatusFinished = "FINISH"
	TaskStatusFailed   = "FAIL"
)

type NotifyMessage struct {
	UserId  int    `json:"user_id"`
	JobId   int    `json:"job_id"`
	Message string `json:"message"`
}

const RewritePromptTemplate = "Please rewrite the following text into AI painting prompt words, and please try to add detailed description of the picture, painting style, scene, rendering effect, picture light and other creative elements. Just output the final prompt word directly. Do not output any explanation lines. The text to be rewritten is: [%s]"
const TranslatePromptTemplate = "Translate the following painting prompt words into English keyword phrases. Without any explanation, directly output the keyword phrases separated by commas. The content to be translated is: [%s]"
