package main

import (
	"github.com/KasumiMercury/alchemark/token"
	"strings"
)

func DetectBlockType(input string) []token.Token {
	if strings.HasPrefix(input, "#") {
		return []token.Token{{Type: token.Heading}}
	}

	if strings.HasPrefix(input, "```") {
		return []token.Token{{Type: token.CodeBlock}}
	}

	if strings.HasPrefix(input, "---") {
		return []token.Token{{Type: token.Horizontal}}
	}

	return []token.Token{{Type: token.Paragraph}}
}
