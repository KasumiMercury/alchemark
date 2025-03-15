package main

import (
	"github.com/KasumiMercury/alchemark/token"
	"strings"
	"sync"
)

type Parser struct {
	lines []string
}

func NewParser(text string) *Parser {
	return &Parser{
		lines: strings.Split(text, "\n"),
	}
}

func (p *Parser) ParseToBlocks() []token.BlockToken {
	if len(p.lines) == 0 {
		return nil
	}

	type lineData struct {
		index int
		token token.BlockToken
	}

	lineDataCh := make(chan lineData, len(p.lines))
	var wg sync.WaitGroup

	for i, line := range p.lines {
		wg.Add(1)
		go func(index int, line string) {
			defer wg.Done()
			lineDataCh <- lineData{index, DetectBlockType(line)}
		}(i, line)
	}

	go func() {
		wg.Wait()
		close(lineDataCh)
	}()

	tokens := make([]token.BlockToken, 0, len(p.lines))

	var openingCodeBlockFence *token.CodeBlockFence
	codeBuffer := make([]string, 0)

	for data := range lineDataCh {
		if data.token == nil {
			continue
		}

		if data.token.Type() == token.CodeBlockFenceType {
			if openingCodeBlockFence == nil {
				openingCodeBlockFence = data.token.(*token.CodeBlockFence)
				continue
			}

			fenceToken := data.token.(*token.CodeBlockFence)
			if fenceToken.InfoString() == "" && fenceToken.FenceChar() == openingCodeBlockFence.FenceChar() {
				tokens = append(tokens, token.NewCodeBlock(openingCodeBlockFence.InfoString(), codeBuffer))
				openingCodeBlockFence = nil
				codeBuffer = codeBuffer[:0]
				continue
			}
		}

		if openingCodeBlockFence != nil {
			codeBuffer = append(codeBuffer, p.lines[data.index])
			continue
		}

		tokens = append(tokens, data.token)
	}

	return tokens
}
