package main

import (
	"github.com/KasumiMercury/alchemark/token"
	"sort"
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

	blocks := make([]lineData, 0, len(p.lines))

	for data := range lineDataCh {
		if data.token == nil {
			continue
		}

		blocks = append(blocks, data)
	}

	sort.Slice(blocks, func(i, j int) bool {
		return blocks[i].index < blocks[j].index
	})

	tokens := make([]token.BlockToken, 0, len(p.lines))

	var openingCodeBlockFence *token.CodeBlockFence
	codeBuffer := make([]string, 0)

	for _, block := range blocks {
		if block.token.Type() == token.CodeBlockFenceType {
			if openingCodeBlockFence == nil {
				openingCodeBlockFence = block.token.(*token.CodeBlockFence)
				continue
			}

			fenceToken := block.token.(*token.CodeBlockFence)
			if fenceToken.InfoString() == "" && fenceToken.FenceChar() == openingCodeBlockFence.FenceChar() {
				tokens = append(tokens, token.NewCodeBlock(openingCodeBlockFence.InfoString(), codeBuffer))
				openingCodeBlockFence = nil
				codeBuffer = codeBuffer[:0]
				continue
			}
		}

		if openingCodeBlockFence != nil {
			codeBuffer = append(codeBuffer, p.lines[block.index])
			continue
		}

		tokens = append(tokens, block.token)
	}

	return tokens
}
