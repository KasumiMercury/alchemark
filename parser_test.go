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
			want:  []token.Token{},
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

func BenchmarkParser_ParseToBlock(b *testing.B) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "case1",
			input: "# Heading\nParagraph",
		},
		{
			name:  "case2",
			input: "# Heading\nParagraph\n# Heading2\nParagraph2",
		},
		{
			name:  "case3",
			input: "# Heading\n---\nParagraph",
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				p := NewParser(tt.input)
				p.ParseToBlocks()
			}
		})
	}
}
