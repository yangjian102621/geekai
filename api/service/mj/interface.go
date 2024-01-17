package mj

type MjClientInterface interface {
	Imagine(prompt string) error
	Upscale(index int, messageId string, hash string) error
	Variation(index int, messageId string, hash string) error
}

type MjServiceInterface interface {
	Run()
}
