package token

const (
	HeadingBlockType    = "Heading"
	ParagraphBlockType  = "Paragraph"
	CodeBlockType       = "CodeBlock"
	CodeBlockFenceType  = "CodeBlockFence"
	HorizontalBlockType = "Horizontal"
)

type BlockType string

type Token interface {
	Type() BlockType
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
func (h HeadingBlock) Type() BlockType {
	return HeadingBlockType
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
func (p ParagraphBlock) Type() BlockType {
	return ParagraphBlockType
}
func (p ParagraphBlock) InlineString() string {
	return p.inlineString
}
func (p ParagraphBlock) Depth() int {
	return p.depth
}

type CodeBlock struct {
	infoString string
	codeLines  []string
}

func NewCodeBlock(infoString string, codeLines []string) *CodeBlock {
	return &CodeBlock{
		infoString: infoString,
		codeLines:  codeLines,
	}
}

func (c CodeBlock) Type() BlockType {
	return CodeBlockType
}
func (c CodeBlock) InfoString() string {
	return c.infoString
}
func (c CodeBlock) CodeLines() []string {
	return c.codeLines
}

type CodeBlockFence struct {
	fenceChar  rune
	infoString string
}

func NewCodeBlockFence(fenceChar rune, infoString string) *CodeBlockFence {
	return &CodeBlockFence{
		fenceChar:  fenceChar,
		infoString: infoString,
	}
}
func (c CodeBlockFence) Type() BlockType {
	return CodeBlockFenceType
}
func (c CodeBlockFence) FenceChar() rune {
	return c.fenceChar
}
func (c CodeBlockFence) InfoString() string {
	return c.infoString
}

type Horizontal struct{}

func NewHorizontal() Horizontal {
	return Horizontal{}
}
func (h Horizontal) Type() BlockType {
	return HorizontalBlockType
}
