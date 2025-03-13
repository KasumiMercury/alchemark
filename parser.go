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

	if input[1] != input[0] || input[2] != input[0] {
		return nil, false
	}

	return token.NewHorizontal(), true
}

func countIndent(input string) int {
	indent := 0
	spaceCount := 0

	for _, char := range input {
		switch char {
		case '\t':
			indent++
			spaceCount = 0
		case ' ':
			spaceCount++
			if spaceCount == 4 {
				indent++
				spaceCount = 0
			}
		default:
			return indent
		}
	}

	return indent
}

func DetectBlockType(input string) token.Token {
	indent := countIndent(input)
	input = input[indent:]

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
	case '*':
		if tk, ok := HorizontalDetector(input); ok {
			return tk
		}
	case '_':
		if tk, ok := HorizontalDetector(input); ok {
			return tk
		}
	default:
		return token.NewParagraphBlock(input, indent)
	}

	return token.NewParagraphBlock(input, indent)
}
