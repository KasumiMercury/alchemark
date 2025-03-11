package main

import (
	"github.com/KasumiMercury/alchemark/token"
)

func DetectBlockType(input string) []token.Token {
	currentPos := 0
	currentChar := input[currentPos]

	if currentChar == '#' {
		currentPos++
		currentChar = input[currentPos]

		for currentChar == '#' {
			currentPos++
			currentChar = input[currentPos]
		}

		if currentChar != ' ' {
			return []token.Token{{Type: token.Paragraph}}
		}

		return []token.Token{{Type: token.Heading}}
	}

	if currentChar == '`' {
		for i := 0; i < 2; i++ {
			if currentChar != '`' {
				return []token.Token{{Type: token.Paragraph}}
			}

			currentPos++
			currentChar = input[currentPos]
		}

		return []token.Token{
			{Type: token.CodeBlock},
		}
	}

	if currentChar == '-' {
		for i := 0; i < 2; i++ {
			if currentChar != '-' {
				return []token.Token{{Type: token.Paragraph}}
			}

			currentPos++
			currentChar = input[currentPos]
		}

		return []token.Token{
			{Type: token.Horizontal},
		}
	}

	return []token.Token{{Type: token.Paragraph}}
}
