package main

import (
	"github.com/KasumiMercury/alchemark/token"
	"reflect"
	"testing"
)

func TestCountIndent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  int
	}{
		{
			name:  "No indent",
			input: "No indent",
			want:  0,
		},
		{
			name:  "4 spaces indent",
			input: "    4 spaces indent",
			want:  1,
		},
		{
			name:  "8 spaces indent",
			input: "        8 spaces indent",
			want:  2,
		},
		{
			name:  "1 tab indent",
			input: "\t1 tab indent",
			want:  1,
		},
		{
			name:  "2 tabs indent",
			input: "\t\t2 tabs indent",
			want:  2,
		},
		{
			name:  "1 tab 4 spaces indent",
			input: "\t    1 tab 4 spaces indent",
			want:  2,
		},
		{
			name:  "4 spaces 1 tab indent",
			input: "    \t4 spaces 1 tab indent",
			want:  2,
		},
		{
			name:  "1 tab 3 spaces indent",
			input: "\t   1 tab 3 spaces indent",
			want:  1,
		},
		{
			name:  "3 spaces 1 tab indent",
			input: "   \t3 spaces 1 tab indent",
			want:  1,
		},
		{
			name:  "indent sandwiched by spaces",
			input: "  \t  mix 4 spaces and 1 tab",
			want:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := countIndent([]rune(tt.input)); got != tt.want {
				t.Errorf("countIndent() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
			want: token.NewCodeBlock("go", ""),
		},
		{
			name: "Horizontal by ---",
			args: args{input: "---"},
			want: token.NewHorizontal(),
		},
		{
			name: "Horizontal by ***",
			args: args{input: "***"},
			want: token.NewHorizontal(),
		},
		{
			name: "Horizontal by ___",
			args: args{input: "___"},
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
		{
			name:  "Paragraph",
			input: "Paragraph",
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
		{
			name: "Horizontal shortage3",
			args: args{input: "**"},
			want: token.NewParagraphBlock("**", 0),
		},
		{
			name: "Horizontal shortage4",
			args: args{input: "*"},
			want: token.NewParagraphBlock("*", 0),
		},
		{
			name: "Horizontal shortage5",
			args: args{input: "__"},
			want: token.NewParagraphBlock("__", 0),
		},
		{
			name: "Horizontal shortage6",
			args: args{input: "_"},
			want: token.NewParagraphBlock("_", 0),
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
