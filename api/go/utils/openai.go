package utils

import (
	"fmt"
	"github.com/pkoukk/tiktoken-go"
)

func CalcTokens(text string, model string) (int, error) {
	encoding := tiktoken.MODEL_TO_ENCODING[model]
	tke, err := tiktoken.GetEncoding(encoding)
	if err != nil {
		return 0, fmt.Errorf("getEncoding: %v", err)
	}

	token := tke.Encode(text, nil, nil)
	return len(token), nil
}
