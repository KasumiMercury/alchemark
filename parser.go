package main

import (
	"github.com/KasumiMercury/alchemark/token"
)

func HeadingDetector(input string) []token.Token {
	level := 0
	for _, char := range input {
		if char == '#' && level < 6 {
			level++
		} else {
			if char == ' ' {
				return []token.Token{{Type: token.Heading}}
			} else {
				return []token.Token{{Type: token.Paragraph}}
			}
		}
	}

	return []token.Token{{Type: token.Paragraph}}
}

func CodeBlockDetector(input string) []token.Token {
	if len(input) < 3 {
		return []token.Token{{Type: token.Paragraph}}
	}

	if input[0:3] == "```" {
		return []token.Token{{Type: token.CodeBlock}}
	}

	return []token.Token{{Type: token.Paragraph}}
}

func HorizontalDetector(input string) []token.Token {
	if len(input) < 3 {
		return []token.Token{{Type: token.Paragraph}}
	}

	if input[0:3] == "---" {
		return []token.Token{{Type: token.Horizontal}}
	}

	return []token.Token{{Type: token.Paragraph}}
}

func DetectBlockType(input string) []token.Token {
	firstChar := input[0]

	switch firstChar {
	case '#':
		return HeadingDetector(input)
	case '`':
		return CodeBlockDetector(input)
	case '-':
		return HorizontalDetector(input)
	default:
		return []token.Token{{Type: token.Paragraph}}
	}
}
