package token

type Token interface {
	Type() string
}

type HeadingBlock struct {
	inlineString string
	level        int
}

func NewHeadingBlock(inlineString string, level int) *HeadingBlock {
	return &HeadingBlock{
		inlineString: inlineString,
		level:        level,
	}
}
func (h HeadingBlock) Type() string {
	return "Heading"
}
func (h HeadingBlock) InlineString() string {
	return h.inlineString
}
func (h HeadingBlock) Level() int {
	return h.level
}

type ParagraphBlock struct {
	inlineString string
	depth        int
}

func NewParagraphBlock(inlineString string, depth int) *ParagraphBlock {
	return &ParagraphBlock{
		inlineString: inlineString,
		depth:        depth,
	}
}
func (p ParagraphBlock) Type() string {
	return "Paragraph"
}
func (p ParagraphBlock) InlineString() string {
	return p.inlineString
}
func (p ParagraphBlock) Depth() int {
	return p.depth
}

type CodeBlock struct {
	lang string
	code string
}

func NewCodeBlock(lang string, code string) *CodeBlock {
	return &CodeBlock{
		lang: lang,
		code: code,
	}
}
func (c CodeBlock) Type() string {
	return "CodeBlock"
}
func (c CodeBlock) Lang() string {
	return c.lang
}
func (c CodeBlock) Code() string {
	return c.code
}

type Horizontal struct{}

func NewHorizontal() Horizontal {
	return Horizontal{}
}
func (h Horizontal) Type() string {
	return "Horizontal"
}
