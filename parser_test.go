package main

import (
	"github.com/KasumiMercury/alchemark/token"
	"reflect"
	"testing"
)

func TestParser_ParseToBlock(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input string
		want  []token.Token
	}{
		{
			input: "",
			want:  nil,
		},
		{
			input: "# Heading\nParagraph",
			want: []token.Token{
				token.NewHeadingBlock("Heading", 1),
				token.NewParagraphBlock("Paragraph", 0),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()

			p := NewParser(tt.input)
			if got := p.ParseToBlocks(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.ParseBlocks() = %v, want %v", got, tt.want)
			}
		})
	}
}
