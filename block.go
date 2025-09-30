package main

import (
	"github.com/KasumiMercury/alchemark/token"
)

func HeadingDetector(input []rune) (token.BlockToken, bool) {
	if len(input) == 0 {
		return nil, false
	}

	level := 0
	idx := 0

	for idx < len(input) && input[idx] == '#' {
		level++
		idx++
	}

	if level == 0 || level > 6 {
		return nil, false
	}

	if idx < len(input) {
		if input[idx] != ' ' && input[idx] != '\t' {
			return nil, false
		}

		for idx < len(input) && (input[idx] == ' ' || input[idx] == '\t') {
			idx++
		}
	}

	content := input[idx:]

	if len(content) == 0 {
		return token.NewHeadingBlock("", level), true
	}

	end := len(content) - 1
	for end >= 0 && (content[end] == ' ' || content[end] == '\t') {
		end--
	}

	if end < 0 {
		return token.NewHeadingBlock("", level), true
	}

	if content[end] != '#' {
		return token.NewHeadingBlock(string(content[:end+1]), level), true
	}

	closingStart := end
	for closingStart >= 0 && content[closingStart] == '#' {
		closingStart--
	}

	if closingStart >= 0 && content[closingStart] != ' ' && content[closingStart] != '\t' {
		return token.NewHeadingBlock(string(content[:end+1]), level), true
	}

	trimPos := closingStart
	for trimPos >= 0 && (content[trimPos] == ' ' || content[trimPos] == '\t') {
		trimPos--
	}

	if trimPos < 0 {
		content = content[:0]
	} else {
		content = content[:trimPos+1]
	}

	return token.NewHeadingBlock(string(content), level), true
}

func CodeBlockDetector(input []rune) (token.BlockToken, bool) {
	if len(input) < 3 {
		return nil, false
	}

	fenceChar := input[0]

	if input[1] == fenceChar && input[2] == fenceChar {
		// TODO: infoStringのスペース取り扱い
		infoString := ""

		if len(input) > 3 {
			for i, char := range input[3:] {
				if char == fenceChar {
					continue
				}

				infoString += string(input[i+3:])
				break
			}
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
		if char == ' ' {
			continue
		}

		if char == hChar {
			hCharCount++
		} else {
			return nil, false
		}
	}

	if hCharCount < 3 {
		return nil, false
	}

	return token.NewHorizontal(), true
}

func BlockQuoteDetector(input []rune) (token.BlockToken, bool) {
	if len(input) < 2 {
		return nil, false
	}

	level := 0
	contents := make([]rune, 0, len(input))
	spaceCount := 0

	for _, char := range input {
		if char == '>' {
			level++
		} else if char == ' ' {
			spaceCount++
			continue
		} else {
			contents = input[level+spaceCount:]
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
		_, ok := HorizontalDetector(input)
		if ok {
			return token.NewHyphen(ok, input), true
		}
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
		hTk, ok := HorizontalDetector(input)
		if ok {
			return hTk, true
		}
		lTk, ok := ListItemDetector(input)
		return lTk, ok
	}

	hTk, ok := HorizontalDetector(input)
	if ok {
		return hTk, true
	}

	return token.NewParagraphBlock(string(input), 0), true
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
			// Add remaining space to input
			remainingSpaceRunes := make([]rune, indentInfo.RemainSpace)
			for i := 0; i < indentInfo.RemainSpace; i++ {
				remainingSpaceRunes[i] = ' '
			}

			selfRunes := append(remainingSpaceRunes, input...)

			return token.NewIndentedBlock(indentInfo.Depth, selfRunes)
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
	case '~':
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
