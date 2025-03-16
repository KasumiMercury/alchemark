package main

import (
	"github.com/KasumiMercury/alchemark/token"
)

func HeadingDetector(input []rune) (token.BlockToken, bool) {
	level := 0
	for _, char := range input {
		if char == '#' && level < 6 {
			level++
		} else {
			if char == ' ' {
				return token.NewHeadingBlock(string(input[level+1:]), level), true
			} else {
				return nil, false
			}
		}
	}

	return nil, false
}

func CodeBlockDetector(input []rune) (token.BlockToken, bool) {
	if len(input) < 3 {
		return nil, false
	}

	if input[0] == '`' && input[1] == '`' && input[2] == '`' {
		// TODO: infoStringのスペース取り扱い
		infoString := ""
		if len(input) > 3 {
			infoString = string(input[4:])
		}
		return token.NewCodeBlockFence(input[0], infoString), true
	}

	return nil, false
}

func HorizontalDetector(input []rune) (token.BlockToken, bool) {
	if len(input) < 3 {
		return nil, false
	}

	hChar := input[0]
	hCharCount := 1

	// first char is already checked
	for _, char := range input[1:] {
		if char != ' ' {
			continue
		}

		if char == hChar {
			hCharCount++
		} else {
			return nil, false
		}

		if hCharCount > 2 {
			break
		}
	}

	return token.NewHorizontal(), true
}

func countIndent(input []rune) int {
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

func DetectBlockType(line string) token.BlockToken {
	input := []rune(line)

	indent := countIndent(input)
	input = input[indent:]

	if len(input) == 0 {
		return token.NewEmpty()
	}

	firstChar := input[0]

	// TODO: インデントの取り扱いを修正

	switch firstChar {
	case '#':
		if tk, ok := HeadingDetector(input); ok {
			return tk
		}
	case '`':
		if tk, ok := CodeBlockDetector(input); ok {
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
		return token.NewParagraphBlock(line, indent)
	}

	return token.NewParagraphBlock(line, indent)
}
