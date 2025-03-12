package main

import (
	"github.com/KasumiMercury/alchemark/token"
)

func HeadingDetector(input string) (token.Token, bool) {
	level := 0
	for _, char := range input {
		if char == '#' && level < 6 {
			level++
		} else {
			if char == ' ' {
				return token.NewHeadingBlock(input[level+1:], level), true
			} else {
				return nil, false
			}
		}
	}

	return nil, false
}

func CodeBlockDetector(input string) (token.Token, bool) {
	if len(input) < 3 {
		return nil, false
	}

	if input[0:3] == "```" {
		// TODO: lang, code
		return token.NewCodeBlock("", ""), true
	}

	return nil, false
}

func HorizontalDetector(input string) (token.Token, bool) {
	if len(input) < 3 {
		return nil, false
	}

	if input[0:3] == "---" {
		return token.NewHorizontal(), true
	}

	return nil, false
}

func DetectBlockType(input string) token.Token {
	firstChar := input[0]

	switch firstChar {
	case '#':
		if tk, ok := HeadingDetector(input); ok {
			return tk
		}
	case '`':
		if tk, ok := CodeBlockDetector(input); ok {
			return tk
		}
	case '-':
		if tk, ok := HorizontalDetector(input); ok {
			return tk
		}
	default:
		// TODO: check depth
		return token.NewParagraphBlock(input, 0)
	}

	// TODO: check depth
	return token.NewParagraphBlock(input, 0)
}
