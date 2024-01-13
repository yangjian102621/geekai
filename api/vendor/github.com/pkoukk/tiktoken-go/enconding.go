package tiktoken

import (
	"errors"
)

const ENDOFTEXT string = "<|endoftext|>"
const FIM_PREFIX string = "<|fim_prefix|>"
const FIM_MIDDLE string = "<|fim_middle|>"
const FIM_SUFFIX string = "<|fim_suffix|>"
const ENDOFPROMPT string = "<|endofprompt|>"

var MODEL_TO_ENCODING = map[string]string{
	// chat
	"gpt-4":         "cl100k_base",
	"gpt-3.5-turbo": "cl100k_base",
	// text
	"text-davinci-003": "p50k_base",
	"text-davinci-002": "p50k_base",
	"text-davinci-001": "r50k_base",
	"text-curie-001":   "r50k_base",
	"text-babbage-001": "r50k_base",
	"text-ada-001":     "r50k_base",
	"davinci":          "r50k_base",
	"curie":            "r50k_base",
	"babbage":          "r50k_base",
	"ada":              "r50k_base",
	// code
	"code-davinci-002": "p50k_base",
	"code-davinci-001": "p50k_base",
	"code-cushman-002": "p50k_base",
	"code-cushman-001": "p50k_base",
	"davinci-codex":    "p50k_base",
	"cushman-codex":    "p50k_base",
	// edit
	"text-davinci-edit-001": "p50k_edit",
	"code-davinci-edit-001": "p50k_edit",
	// embeddings
	"text-embedding-ada-002": "cl100k_base",
	// old embeddings
	"text-similarity-davinci-001":  "r50k_base",
	"text-similarity-curie-001":    "r50k_base",
	"text-similarity-babbage-001":  "r50k_base",
	"text-similarity-ada-001":      "r50k_base",
	"text-search-davinci-doc-001":  "r50k_base",
	"text-search-curie-doc-001":    "r50k_base",
	"text-search-babbage-doc-001":  "r50k_base",
	"text-search-ada-doc-001":      "r50k_base",
	"code-search-babbage-code-001": "r50k_base",
	"code-search-ada-code-001":     "r50k_base",
	// open source
	"gpt2": "gpt2",
}

var MODEL_PREFIX_TO_ENCODING = map[string]string{
	// chat
	"gpt-4-":         "cl100k_base", // e.g., gpt-4-0314, etc., plus gpt-4-32k
	"gpt-3.5-turbo-": "cl100k_base", // e.g, gpt-3.5-turbo-0301, -0401, etc.
}

type Encoding struct {
	Name           string
	PatStr         string
	MergeableRanks map[string]int
	SpecialTokens  map[string]int
	ExplicitNVocab int
}

func getEncoding(encodingName string) (*Encoding, error) {
	encoding, ok := ENCODING_MAP[encodingName]
	if !ok {
		initEncoding, err := initEncoding(encodingName)
		if err != nil {
			return nil, err
		}
		encoding = initEncoding
		ENCODING_MAP[encodingName] = encoding
	}
	return encoding, nil
}

func initEncoding(encodingName string) (*Encoding, error) {
	switch encodingName {
	case "cl100k_base":
		return cl100k_base()
	case "p50k_base":
		return p50k_base()
	case "r50k_base":
		return r50k_base()
	case "p50k_edit":
		return p50k_edit()
	default:
		return nil, errors.New("Unknown encoding: " + encodingName)
	}
}

func cl100k_base() (*Encoding, error) {
	ranks, err := loadTiktokenBpe("cl100k_base.tiktoken")
	if err != nil {
		return nil, err
	}
	special_tokens := map[string]int{
		ENDOFTEXT:   100257,
		FIM_PREFIX:  100258,
		FIM_MIDDLE:  100259,
		FIM_SUFFIX:  100260,
		ENDOFPROMPT: 100276,
	}
	return &Encoding{
		Name:           "cl100k_base",
		PatStr:         `(?i:'s|'t|'re|'ve|'m|'ll|'d)|[^\r\n\p{L}\p{N}]?\p{L}+|\p{N}{1,3}| ?[^\s\p{L}\p{N}]+[\r\n]*|\s*[\r\n]+|\s+(?!\S)|\s+`,
		MergeableRanks: ranks,
		SpecialTokens:  special_tokens,
	}, nil
}

func p50k_edit() (*Encoding, error) {
	ranks, err := loadTiktokenBpe("p50k_base.tiktoken")
	if err != nil {
		return nil, err
	}
	special_tokens := map[string]int{ENDOFTEXT: 50256, FIM_PREFIX: 50281, FIM_MIDDLE: 50282, FIM_SUFFIX: 50283}
	return &Encoding{
		Name:           "p50k_edit",
		PatStr:         `'s|'t|'re|'ve|'m|'ll|'d| ?\p{L}+| ?\p{N}+| ?[^\s\p{L}\p{N}]+|\s+(?!\S)|\s+`,
		MergeableRanks: ranks,
		SpecialTokens:  special_tokens,
	}, nil
}

func p50k_base() (*Encoding, error) {
	ranks, err := loadTiktokenBpe("p50k_base.tiktoken")
	if err != nil {
		return nil, err
	}
	special_tokens := map[string]int{ENDOFTEXT: 50256}

	// ExplicitNVocab := 50281
	// max_tokens := int(math.Max(float64(len(special_tokens)), float64(len(ranks))))

	// if len(special_tokens)+len(ranks) != max_tokens {
	// 	return nil, errors.New("special_tokens and ranks must be disjoint")
	// }

	return &Encoding{
		Name:           "p50k_base",
		PatStr:         `'s|'t|'re|'ve|'m|'ll|'d| ?\p{L}+| ?\p{N}+| ?[^\s\p{L}\p{N}]+|\s+(?!\S)|\s+`,
		MergeableRanks: ranks,
		SpecialTokens:  special_tokens,
		ExplicitNVocab: 50281,
	}, nil
}

func r50k_base() (*Encoding, error) {
	ranks, err := loadTiktokenBpe("r50k_base.tiktoken")
	if err != nil {
		return nil, err
	}
	special_tokens := map[string]int{ENDOFTEXT: 50256}
	return &Encoding{
		Name:           "r50k_base",
		MergeableRanks: ranks,
		PatStr:         `'s|'t|'re|'ve|'m|'ll|'d| ?\p{L}+| ?\p{N}+| ?[^\s\p{L}\p{N}]+|\s+(?!\S)|\s+`,
		SpecialTokens:  special_tokens,
		ExplicitNVocab: 50257,
	}, nil
}

var ENCODING_MAP = map[string]*Encoding{}
