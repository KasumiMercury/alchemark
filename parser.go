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

func (p *Parser) ParseToBlocks() []token.Token {
	if len(p.lines) == 0 {
		return nil
	}

	type lineData struct {
		index int
		token token.Token
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

	tokens := make([]token.Token, 0, len(p.lines))
	for data := range lineDataCh {
		if data.token == nil {
			continue
		}

		tokens = append(tokens, data.token)
	}

	return tokens
}
