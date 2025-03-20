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
				break
			} else {
				return nil, false
			}
		}
	}

	if level == 0 {
		return nil, false
	}

	return token.NewHeadingBlock(string(input[level+1:]), level), true
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

func BlockQuoteDetector(input []rune) (token.BlockToken, bool) {
	if len(input) < 2 {
		return nil, false
	}

	level := 0
	contents := make([]rune, 0, len(input))

	for _, char := range input {
		if char == '>' {
			level++
		} else if char == ' ' {
			contents = input[level+1:]
			break
		} else {
			return nil, false
		}
	}

	contentBlock := DetectBlockType(string(contents))

	return token.NewBlockQuote(level, contentBlock), true
}

func ListItemDetector(input []rune) (token.BlockToken, bool) {
	if input[0] != '-' && input[0] != '+' && input[0] != '*' {
		return nil, false
	}

	if input[1] != ' ' {
		return nil, false
	}

	// skip space
	pos := 2
	for _, char := range input[2:] {
		if char == ' ' {
			pos++
		} else {
			break
		}
	}

	contentBlock := DetectBlockType(string(input[pos:]))

	return token.NewListItem(input[0], 0, contentBlock), true
}

func HyphenDetector(input []rune) (token.BlockToken, bool) {
	if input[0] != '-' {
		return nil, false
	}

	if len(input) > 1 && input[1] == ' ' {
		// TODO: when the line can be a horizontal line, it should be a horizontal line
		tk, ok := ListItemDetector(input)
		return tk, ok
	}

	_, ok := HorizontalDetector(input)

	return token.NewHyphen(ok, input), true
}

func AsteriskDetector(input []rune) (token.BlockToken, bool) {
	if input[0] != '*' {
		return nil, false
	}

	if len(input) > 1 && input[1] == ' ' {
		// TODO: when the line can be a horizontal line, it should be a horizontal line
		tk, ok := ListItemDetector(input)
		return tk, ok
	}

	if tk, ok := HorizontalDetector(input); ok {
		return tk, ok
	}

	return nil, false
}

type IndentInfo struct {
	Depth       int
	SeekPos     int
	RemainSpace int
}

func countIndent(input []rune) IndentInfo {
	spaceCount := 0
	pos := 0

	for _, char := range input {

		if char == '\t' {
			spaceCount += 4
			pos++
			continue
		}

		if char == ' ' {
			spaceCount++
			pos++
			continue
		}

		break
	}

	return IndentInfo{
		Depth:       spaceCount / 4,
		SeekPos:     pos,
		RemainSpace: spaceCount % 4,
	}
}

func DetectBlockType(line string) token.BlockToken {
	input := []rune(line)

	indentInfo := countIndent(input)
	input = input[indentInfo.SeekPos:]

	if len(input) == 0 {
		return token.NewBlank()
	}

	firstChar := input[0]

	if indentInfo.Depth > 0 {
		// TODO: handle remaining space

		switch firstChar {
		case '-':
			// TODO: If there is remaining space, can the line be a list item?
			if tk, ok := ListItemDetector(input); ok {
				return tk.(token.ListItem).Indent(indentInfo.Depth)
			}
		default:
			return token.NewIndentedBlock(indentInfo.Depth, input)
		}
	}

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
		if tk, ok := HyphenDetector(input); ok {
			return tk
		}
	case '*':
		if tk, ok := AsteriskDetector(input); ok {
			return tk
		}
	case '>':
		if tk, ok := BlockQuoteDetector(input); ok {
			return tk
		}
	case '+':
		if tk, ok := ListItemDetector(input); ok {
			return tk
		}
	case '_':
		if tk, ok := HorizontalDetector(input); ok {
			return tk
		}
	case '=':
		return token.NewEqual(input)
	default:
		return token.NewParagraphBlock(line, 0)
	}

	return token.NewParagraphBlock(line, 0)
}
