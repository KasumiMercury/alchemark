package main

import (
	"github.com/KasumiMercury/alchemark/token"
	"reflect"
	"testing"
)

func TestDetectBlockType(t *testing.T) {
	t.Parallel()

	type args struct {
		input string
	}

	tests := []struct {
		name string
		args args
		want []token.Token
	}{
		{
			name: "Heading",
			args: args{input: "# Heading"},
			want: []token.Token{{Type: token.Heading}},
		},
		{
			name: "Heading2",
			args: args{input: "## Heading"},
			want: []token.Token{{Type: token.Heading}},
		},
		{
			name: "CodeBlock",
			args: args{input: "```"},
			want: []token.Token{{Type: token.CodeBlock}},
		},
		{
			name: "CodeBlock with language",
			args: args{input: "````go"},
			want: []token.Token{{Type: token.CodeBlock}},
		},
		{
			name: "Horizontal by ---",
			args: args{input: "---"},
			want: []token.Token{{Type: token.Horizontal}},
		},
		{
			name: "Paragraph",
			args: args{input: "Paragraph"},
			want: []token.Token{{Type: token.Paragraph}},
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
