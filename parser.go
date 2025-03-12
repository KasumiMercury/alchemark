package main

import (
	"github.com/KasumiMercury/alchemark/token"
)

func HeadingDetector(input string) token.Token {
	level := 0
	for _, char := range input {
		if char == '#' && level < 6 {
			level++
		} else {
			if char == ' ' {
				return token.NewHeadingBlock(input[level:], level)
			} else {
				// TODO: check depth
				return token.NewParagraphBlock(input, 0)
			}
		}
	}

	// TODO: check depth
	return token.NewParagraphBlock(input, 0)
}

func CodeBlockDetector(input string) token.Token {
	if len(input) < 3 {
		// TODO: check depth
		return token.NewParagraphBlock(input, 0)
	}

	if input[0:3] == "```" {
		// TODO: lang, code
		return token.NewCodeBlock("", "")
	}

	// TODO: check depth
	return token.NewParagraphBlock(input, 0)
}

func HorizontalDetector(input string) token.Token {
	if len(input) < 3 {
		// TODO: check depth
		return token.NewParagraphBlock(input, 0)
	}

	if input[0:3] == "---" {
		return token.NewHorizontal()
	}

	// TODO: check depth
	return token.NewParagraphBlock(input, 0)
}

func DetectBlockType(input string) token.Token {
	firstChar := input[0]

	switch firstChar {
	case '#':
		return HeadingDetector(input)
	case '`':
		return CodeBlockDetector(input)
	case '-':
		return HorizontalDetector(input)
	default:
		return token.NewParagraphBlock(input, 0)
	}
}
