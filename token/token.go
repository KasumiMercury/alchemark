package token

import "fmt"

const (
	HeadingBlockType    = "Heading"
	ParagraphBlockType  = "Paragraph"
	CodeBlockType       = "CodeBlock"
	CodeBlockFenceType  = "CodeBlockFence"
	HorizontalBlockType = "Horizontal"
	EmptyBlockType      = "Empty"
)

type BlockType string

type Token interface {
	Type() BlockType
	String() string
}

type HeadingBlock struct {
	level        int
	inlineString string
}

func NewHeadingBlock(inlineString string, level int) *HeadingBlock {
	return &HeadingBlock{
		level:        level,
		inlineString: inlineString,
	}
}
func (h HeadingBlock) Type() BlockType {
	return HeadingBlockType
}
func (h HeadingBlock) Level() int {
	return h.level
}
func (h HeadingBlock) InlineString() string {
	return h.inlineString
}
func (h HeadingBlock) String() string {
	return fmt.Sprintf("Type: %s, Level: %d, InlineString: %s", HeadingBlockType, h.level, h.inlineString)
}

type ParagraphBlock struct {
	depth        int
	inlineString string
}

func NewParagraphBlock(inlineString string, depth int) *ParagraphBlock {
	return &ParagraphBlock{
		depth:        depth,
		inlineString: inlineString,
	}
}
func (p ParagraphBlock) Type() BlockType {
	return ParagraphBlockType
}
func (p ParagraphBlock) Depth() int {
	return p.depth
}
func (p ParagraphBlock) InlineString() string {
	return p.inlineString
}
func (p ParagraphBlock) String() string {
	return fmt.Sprintf("Type: %s, Depth: %d, InlineString: %s", ParagraphBlockType, p.depth, p.inlineString)
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
func (c CodeBlock) String() string {
	return fmt.Sprintf("Type: %s, InfoString: %s, CodeLines: %v", CodeBlockType, c.infoString, c.codeLines)
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
func (c CodeBlockFence) String() string {
	return fmt.Sprintf("Type: %s, FenceChar: %c, InfoString: %s", CodeBlockFenceType, c.fenceChar, c.infoString)
}

type Horizontal struct{}

func NewHorizontal() Horizontal {
	return Horizontal{}
}
func (h Horizontal) Type() BlockType {
	return HorizontalBlockType
}
func (h Horizontal) String() string {
	return fmt.Sprintf("Type: %s", HorizontalBlockType)
}

type Empty struct{}

func NewEmpty() Empty {
	return Empty{}
}
func (e Empty) Type() BlockType {
	return EmptyBlockType
}
func (e Empty) String() string {
	return fmt.Sprintf("Type: %s", EmptyBlockType)
}
