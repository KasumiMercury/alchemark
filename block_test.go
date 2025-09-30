package main

import (
	"reflect"
	"testing"

	"github.com/KasumiMercury/alchemark/token"
)

func TestCountIndent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  IndentInfo
	}{
		{
			name:  "No indent",
			input: "No indent",
			want: IndentInfo{
				0,
				0,
				0,
			},
		},
		{
			name:  "4 spaces indent",
			input: "    4 spaces indent",
			want: IndentInfo{
				1,
				4,
				0,
			},
		},
		{
			name:  "8 spaces indent",
			input: "        8 spaces indent",
			want: IndentInfo{
				2,
				8,
				0,
			},
		},
		{
			name:  "1 tab indent",
			input: "\t1 tab indent",
			want: IndentInfo{
				1,
				1,
				0,
			},
		},
		{
			name:  "2 tabs indent",
			input: "\t\t2 tabs indent",
			want: IndentInfo{
				2,
				2,
				0,
			},
		},
		{
			name:  "1 tab 4 spaces indent",
			input: "\t    1 tab 4 spaces indent",
			want: IndentInfo{
				2,
				5,
				0,
			},
		},
		{
			name:  "4 spaces 1 tab indent",
			input: "    \t4 spaces 1 tab indent",
			want: IndentInfo{
				2,
				5,
				0,
			},
		},
		{
			name:  "1 tab 3 spaces indent",
			input: "\t   1 tab 3 spaces indent",
			want: IndentInfo{
				1,
				4,
				3,
			},
		},
		{
			name:  "3 spaces 1 tab indent",
			input: "   \t3 spaces 1 tab indent",
			want: IndentInfo{
				1,
				4,
				3,
			},
		},
		{
			name:  "indent sandwiched by spaces",
			input: "  \t  mix 4 spaces and 1 tab",
			want: IndentInfo{
				2,
				5,
				0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := countIndent([]rune(tt.input)); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("countIndent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHeadingDetector(t *testing.T) {
	t.Parallel()

	type args struct {
		input string
	}

	type want struct {
		token  token.BlockToken
		detect bool
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Heading",
			args: args{
				input: "# Heading",
			},
			want: want{
				token.NewHeadingBlock("Heading", 1),
				true,
			},
		},
		{
			name: "Heading2",
			args: args{
				input: "## Heading",
			},
			want: want{
				token.NewHeadingBlock("Heading", 2),
				true,
			},
		},
		{
			name: "Heading6",
			args: args{
				input: "###### Heading",
			},
			want: want{
				token.NewHeadingBlock("Heading", 6),
				true,
			},
		},
		{
			name: "First char is not # will be not Heading",
			args: args{
				input: "Heading",
			},
			want: want{
				nil,
				false,
			},
		},
		{
			name: "Heading7 will be not Heading",
			args: args{
				input: "####### Heading",
			},
			want: want{
				nil,
				false,
			},
		},
		{
			name: "Heading no space will be not Heading",
			args: args{
				input: "#Heading",
			},
			want: want{
				nil,
				false,
			},
		},
		{
			name: "escaped # will be not Heading",
			args: args{
				input: "\\# Heading",
			},
			want: want{
				nil,
				false,
			},
		},
		{
			name: "leading or trailing space will be removed",
			args: args{
				input: "#   Heading    ",
			},
			want: want{
				token.NewHeadingBlock("Heading", 1),
				true,
			},
		},
		{
			name: "Heading with closing hashes",
			args: args{
				input: "### Heading ###   ",
			},
			want: want{
				token.NewHeadingBlock("Heading", 3),
				true,
			},
		},
		{
			name: "Heading keeps inline hashes without separator",
			args: args{
				input: "### Heading###",
			},
			want: want{
				token.NewHeadingBlock("Heading###", 3),
				true,
			},
		},
		{
			name: "Opening sequence will be prioritized over longer closing sequence",
			args: args{
				input: "### Heading ###########",
			},
			want: want{
				token.NewHeadingBlock("Heading", 3),
				true,
			},
		},
		{
			name: "Opening sequence will be prioritized over shorter closing sequence",
			args: args{
				input: "### Heading #",
			},
			want: want{
				token.NewHeadingBlock("Heading", 3),
				true,
			},
		},
		{
			name: "Heading retains hashes when followed by text",
			args: args{
				input: "# Heading ##text",
			},
			want: want{
				token.NewHeadingBlock("Heading ##text", 1),
				true,
			},
		},
		{
			name: "Heading allows tab after marker",
			args: args{
				input: "#\tHeading",
			},
			want: want{
				token.NewHeadingBlock("Heading", 1),
				true,
			},
		},
		{
			name: "inline can be empty",
			args: args{
				input: "#",
			},
			want: want{
				token.NewHeadingBlock("", 1),
				true,
			},
		},
		{
			name: "inline can be empty",
			args: args{
				input: "## ",
			},
			want: want{
				token.NewHeadingBlock("", 2),
				true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got, detect := HeadingDetector([]rune(tt.args.input)); !reflect.DeepEqual(got, tt.want.token) || detect != tt.want.detect {
				t.Errorf("HeadingDetector() = {%v}, %v / want {%v}, %v", got, detect, tt.want.token, tt.want.detect)
			}
		})
	}
}

func TestCodeBlockDetector(t *testing.T) {
	t.Parallel()

	type args struct {
		input string
	}

	type want struct {
		token  token.BlockToken
		detect bool
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "CodeBlock by ```",
			args: args{
				input: "```",
			},
			want: want{
				token.NewCodeBlockFence('`', ""),
				true,
			},
		},
		{
			name: "CodeBlock by ``` with infoString",
			args: args{
				input: "```go",
			},
			want: want{
				token.NewCodeBlockFence('`', "go"),
				true,
			},
		},
		{
			name: "CodeBlock by ~~~",
			args: args{
				input: "~~~",
			},
			want: want{
				token.NewCodeBlockFence('~', ""),
				true,
			},
		},
		{
			name: "CodeBlock by ~~~ with infoString",
			args: args{
				input: "~~~ruby",
			},
			want: want{
				token.NewCodeBlockFence('~', "ruby"),
				true,
			},
		},
		{
			name: "long CodeBlock will be allowed",
			args: args{
				input: "````",
			},
			want: want{
				token.NewCodeBlockFence('`', ""),
				true,
			},
		},
		{
			name: "long CodeBlock with infoString will be allowed",
			args: args{
				input: "````python",
			},
			want: want{
				token.NewCodeBlockFence('`', "python"),
				true,
			},
		},
		{
			name: "shortage will be not CodeBlock",
			args: args{
				input: "``",
			},
			want: want{
				nil,
				false,
			},
		},
		{
			name: "shortage will be not CodeBlock",
			args: args{
				input: "~~",
			},
			want: want{
				nil,
				false,
			},
		},
		{
			name: "not collected char at 2nd will be not CodeBlock",
			args: args{
				input: "`*`",
			},
			want: want{
				nil,
				false,
			},
		},
		{
			name: "not collected char at 3rd will be not CodeBlock",
			args: args{
				input: "``*",
			},
			want: want{
				nil,
				false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got, detect := CodeBlockDetector([]rune(tt.args.input)); !reflect.DeepEqual(got, tt.want.token) || detect != tt.want.detect {
				t.Errorf("CodeBlockDetector() = {%v}, %v / want {%v}, %v", got, detect, tt.want.token, tt.want.detect)
			}
		})
	}
}

func TestHorizontalDetector(t *testing.T) {
	t.Parallel()

	type args struct {
		input string
	}

	type want struct {
		token  token.BlockToken
		detect bool
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Horizontal by ***",
			args: args{
				input: "***",
			},
			want: want{
				token.NewHorizontal(),
				true,
			},
		},
		{
			name: "Horizontal by ___",
			args: args{
				input: "___",
			},
			want: want{
				token.NewHorizontal(),
				true,
			},
		},
		{
			name: "Horizontal by ---",
			args: args{
				input: "---",
			},
			want: want{
				token.NewHorizontal(),
				true,
			},
		},
		{
			name: "shortage will be not Horizontal",
			args: args{
				input: "**",
			},
			want: want{
				nil,
				false,
			},
		},
		{
			name: "shortage will be not Horizontal",
			args: args{
				input: "__",
			},
			want: want{
				nil,
				false,
			},
		},
		{
			name: "shortage will be not Horizontal",
			args: args{
				input: "--",
			},
			want: want{
				nil,
				false,
			},
		},
		{
			name: "more than 3 will be used as Horizontal",
			args: args{
				input: "****",
			},
			want: want{
				token.NewHorizontal(),
				true,
			},
		},
		{
			name: "more than 3 will be used as Horizontal",
			args: args{
				input: "____",
			},
			want: want{
				token.NewHorizontal(),
				true,
			},
		},
		{
			name: "more than 3 will be used as Horizontal",
			args: args{
				input: "----",
			},
			want: want{
				token.NewHorizontal(),
				true,
			},
		},
		{
			name: "space between will allowed",
			args: args{
				input: "- - -",
			},
			want: want{
				token.NewHorizontal(),
				true,
			},
		},
		{
			name: "space between will allowed",
			args: args{
				input: "**  * ** * ** * **",
			},
			want: want{
				token.NewHorizontal(),
				true,
			},
		},
		{
			name: "start with space will allowed",
			args: args{
				input: "-    --",
			},
			want: want{
				token.NewHorizontal(),
				true,
			},
		},
		{
			name: "end with space will allowed",
			args: args{
				input: "---    ",
			},
			want: want{
				token.NewHorizontal(),
				true,
			},
		},
		{
			name: "other character included will be not Horizontal",
			args: args{
				input: "---a",
			},
			want: want{
				nil,
				false,
			},
		},
		{
			name: "other character included will be not Horizontal",
			args: args{
				input: "-*-***",
			},
			want: want{
				nil,
				false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got, detect := HorizontalDetector([]rune(tt.args.input)); !reflect.DeepEqual(got, tt.want.token) || detect != tt.want.detect {
				t.Errorf("HorizontalDetector() = {%v}, %v / want {%v}, %v", got, detect, tt.want.token, tt.want.detect)
			}
		})
	}
}

func TestBlockQuoteDetector(t *testing.T) {
	t.Parallel()

	type args struct {
		input string
	}

	type want struct {
		token  token.BlockToken
		detect bool
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Blockquote with paragraph",
			args: args{
				input: "> Blockquote",
			},
			want: want{
				token.NewBlockQuote(
					1,
					token.NewParagraphBlock("Blockquote", 0),
				),
				true,
			},
		},
		{
			name: "Spaces after > can be omitted",
			args: args{
				input: ">Blockquote",
			},
			want: want{
				token.NewBlockQuote(
					1,
					token.NewParagraphBlock("Blockquote", 0),
				),
				true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got, detect := BlockQuoteDetector([]rune(tt.args.input)); !reflect.DeepEqual(got, tt.want.token) || detect != tt.want.detect {
				t.Errorf("BlockQuoteDetector() = {%v}, %v / want {%v}, %v", got, detect, tt.want.token, tt.want.detect)
			}
		})
	}
}

func TestListItemDetector(t *testing.T) {
	t.Parallel()

	type args struct {
		input string
	}

	type want struct {
		token  token.BlockToken
		detect bool
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "List item",
			args: args{
				input: "- List item",
			},
			want: want{
				token.NewListItem(
					'-',
					0,
					token.NewParagraphBlock("List item", 0),
				),
				true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got, detect := ListItemDetector([]rune(tt.args.input)); !reflect.DeepEqual(got, tt.want.token) || detect != tt.want.detect {
				t.Errorf("ListItemDetector() = {%v}, %v / want {%v}, %v", got, detect, tt.want.token, tt.want.detect)
			}
		})
	}
}

func TestHyphenDetector(t *testing.T) {
	t.Parallel()

	type args struct {
		input []rune
	}

	type want struct {
		token  token.BlockToken
		detect bool
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Horizontal by ---",
			args: args{
				input: []rune{'-', '-', '-'},
			},
			want: want{
				token.NewHyphen(true, []rune{'-', '-', '-'}),
				true,
			},
		},
		{
			name: "not Horizontal",
			args: args{
				input: []rune{'-', '-', 'a'},
			},
			want: want{
				token.NewHyphen(false, []rune{'-', '-', 'a'}),
				true,
			},
		},
		{
			name: "Horizontal by --- with space",
			args: args{
				input: []rune{'-', ' ', '-', ' ', '-'},
			},
			want: want{
				token.NewHyphen(true, []rune{'-', ' ', '-', ' ', '-'}),
				true,
			},
		},
		{
			name: "List item",
			args: args{
				input: []rune{'-', ' ', 'L', 'i', 's', 't', ' ', 'i', 't', 'e', 'm'},
			},
			want: want{
				token.NewListItem(
					'-',
					0,
					token.NewParagraphBlock("List item", 0),
				),
				true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got, detect := HyphenDetector(tt.args.input); !reflect.DeepEqual(got, tt.want.token) || detect != tt.want.detect {
				t.Errorf("HyphenDetector() = {%v}, %v / want {%v}, %v", got, detect, tt.want.token, tt.want.detect)
			}
		})
	}
}

func TestAsteriskDetector(t *testing.T) {
	t.Parallel()

	type args struct {
		input []rune
	}

	type want struct {
		token  token.BlockToken
		detect bool
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Horizontal by ***",
			args: args{
				input: []rune{'*', '*', '*'},
			},
			want: want{
				token.NewAsterisk(true, []rune{'*', '*', '*'}),
				true,
			},
		},
		{
			name: "not Horizontal",
			args: args{
				input: []rune{'*', '*', 'a'},
			},
			want: want{
				token.NewAsterisk(false, []rune{'*', '*', 'a'}),
				true,
			},
		},
		{
			name: "Horizontal by *** with space",
			args: args{
				input: []rune{'*', ' ', '*', ' ', '*'},
			},
			want: want{
				token.NewAsterisk(true, []rune{'*', ' ', '*', ' ', '*'}),
				true,
			},
		},
		{
			name: "List item",
			args: args{
				input: []rune{'*', ' ', 'L', 'i', 's', 't', ' ', 'i', 't', 'e', 'm'},
			},
			want: want{
				token.NewListItem(
					'*',
					0,
					token.NewParagraphBlock("List item", 0),
				),
				true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got, detect := AsteriskDetector(tt.args.input); !reflect.DeepEqual(got, tt.want.token) || detect != tt.want.detect {
				t.Errorf("HyphenDetector() = {%v}, %v / want {%v}, %v", got, detect, tt.want.token, tt.want.detect)
			}
		})
	}
}

// TODO: Add test for EqualDetector

func TestDetectBlockTypeSuccess(t *testing.T) {
	t.Parallel()

	type args struct {
		input string
	}

	tests := []struct {
		name string
		args args
		want token.BlockToken
	}{
		{
			name: "Heading with 1 space",
			args: args{input: " # Heading"},
			want: token.NewHeadingBlock("Heading", 1),
		},
		{
			name: "Heading with 2 spaces",
			args: args{input: "  ## Heading"},
			want: token.NewHeadingBlock("Heading", 2),
		},
		{
			name: "Heading with 3 spaces",
			args: args{input: "   ### Heading"},
			want: token.NewHeadingBlock("Heading", 3),
		},
		{
			name: "4 spaces before # should be IndentedBlock",
			args: args{input: "    # Heading"},
			want: token.NewIndentedBlock(1, []rune("# Heading")),
		},
		{
			name: "Indented List item",
			args: args{input: "    - List item"},
			want: token.NewListItem('-', 1, token.NewParagraphBlock("List item", 0)),
		},
		{
			name: "IndentedBlock",
			args: args{input: "    IndentedBlock"},
			want: token.NewIndentedBlock(1, []rune("IndentedBlock")),
		},
		{
			name: "IndentedBlock with additional space",
			args: args{input: "      IndentedBlock"},
			want: token.NewIndentedBlock(1, []rune("  IndentedBlock")),
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
