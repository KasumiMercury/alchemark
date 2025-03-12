package main

import (
	"github.com/KasumiMercury/alchemark/token"
	"reflect"
	"testing"
)

func TestDetectBlockTypeSuccess(t *testing.T) {
	t.Parallel()

	type args struct {
		input string
	}

	tests := []struct {
		name string
		args args
		want token.Token
	}{
		{
			name: "Heading",
			args: args{input: "# Heading"},
			want: token.NewHeadingBlock("Heading", 1),
		},
		{
			name: "Heading2",
			args: args{input: "## Heading"},
			want: token.NewHeadingBlock("Heading", 2),
		},
		{
			name: "CodeBlock",
			args: args{input: "```"},
			want: token.NewCodeBlock("", ""),
		},
		{
			name: "CodeBlock with language",
			args: args{input: "````go"},
			want: token.NewCodeBlock("", ""),
		},
		{
			name: "Horizontal by ---",
			args: args{input: "---"},
			want: token.NewHorizontal(),
		},
		{
			name: "Paragraph",
			args: args{input: "Paragraph"},
			want: token.NewParagraphBlock("Paragraph", 0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := DetectBlockType(tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DetectBlockType() = %v, want %v", got, tt.want)
			}
		})
	}

}

func BenchmarkDetectBlockType(b *testing.B) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "Heading",
			input: "# Heading",
		},
		{
			name:  "Heading2",
			input: "## Heading",
		},
		{
			name:  "CodeBlock",
			input: "```",
		},
		{
			name:  "CodeBlock with language",
			input: "````go",
		},
		{
			name:  "Horizontal by ---",
			input: "---",
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				DetectBlockType(tt.input)
			}
		})
	}
}

func TestDetectBlockTypeNoSpace(t *testing.T) {
	t.Parallel()

	type args struct {
		input string
	}

	tests := []struct {
		name string
		args args
		want token.Token
	}{
		{
			name: "Heading no space",
			args: args{input: "#Heading"},
			want: token.NewParagraphBlock("#Heading", 0),
		},
		{
			name: "Heading2 no space",
			args: args{input: "##Heading"},
			want: token.NewParagraphBlock("##Heading", 0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := DetectBlockType(tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DetectBlockType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDetectBlockTypeShortage(t *testing.T) {
	t.Parallel()

	type args struct {
		input string
	}

	tests := []struct {
		name string
		args args
		want token.Token
	}{
		{
			name: "CodeBlock shortage",
			args: args{input: "`"},
			want: token.NewParagraphBlock("`", 0),
		},
		{
			name: "CodeBlock shortage2",
			args: args{input: "``"},
			want: token.NewParagraphBlock("``", 0),
		},
		{
			name: "Horizontal shortage",
			args: args{input: "-"},
			want: token.NewParagraphBlock("-", 0),
		},
		{
			name: "Horizontal shortage2",
			args: args{input: "--"},
			want: token.NewParagraphBlock("--", 0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := DetectBlockType(tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DetectBlockType() = %v, want %v", got, tt.want)
			}
		})
	}
}
