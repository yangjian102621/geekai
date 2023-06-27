package utils

import (
	"fmt"
	"github.com/pkoukk/tiktoken-go"
)

func CalcTokens(text string, model string) (int, error) {
	encoding, ok := tiktoken.MODEL_TO_ENCODING[model]
	if !ok {
		encoding = "cl100k_base"
	}
	tke, err := tiktoken.GetEncoding(encoding)
	if err != nil {
		return 0, fmt.Errorf("getEncoding: %v", err)
	}

	token := tke.Encode(text, nil, nil)
	return len(token), nil
}
