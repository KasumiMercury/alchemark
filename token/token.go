package token

import "fmt"

const (
	HeadingBlockType   = "Heading"
	ParagraphBlockType = "Paragraph"

	// IndentedBlockType indented block can be used for a nested list, nested paragraph or IndentedCodeBlock
	IndentedBlockType     = "Indented"
	IndentedCodeBlockType = "IndentedCodeBlock"

	CodeBlockType      = "CodeBlock"
	CodeBlockFenceType = "CodeBlockFence"

	// HyphenBlockType hyphen can be used for horizontal line or setext heading
	HyphenBlockType = "HyphenToken"
	// EqualBlockType equal can be used for setext heading
	EqualBlockType = "EqualToken"

	HorizontalBlockType = "Horizontal"
	SetextBlockType     = "SetextHeading"

	EmptyBlockType = "Empty"
)

type BlockType string

type BlockToken interface {
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

type IndentedBlock struct {
	depth int
	self  []rune
}

func NewIndentedBlock(level int, self []rune) *IndentedBlock {
	return &IndentedBlock{
		depth: level,
		self:  self,
	}
}
func (i IndentedBlock) Type() BlockType {
	return IndentedBlockType
}
func (i IndentedBlock) Depth() int {
	return i.depth
}
func (i IndentedBlock) InlineString() string {
	return string(i.self)
}
func (i IndentedBlock) ConvertBlockToIndentedCodeBlock(aboveType BlockType) BlockToken {
	if aboveType == ParagraphBlockType {
		return NewParagraphBlock(string(i.self), i.depth)
	}

	return NewIndentedCodeBlock(i.depth, i.self)
}
func (i IndentedBlock) String() string {
	return fmt.Sprintf("Type: %s, Depth: %d, InlineString: %s", IndentedBlockType, i.depth, i.self)
}

type IndentedCodeBlock struct {
	depth int
	self  []rune
}

func NewIndentedCodeBlock(level int, self []rune) *IndentedCodeBlock {
	return &IndentedCodeBlock{
		depth: level,
		self:  self,
	}
}
func (i IndentedCodeBlock) Type() BlockType {
	return IndentedCodeBlockType
}
func (i IndentedCodeBlock) Depth() int {
	return i.depth
}
func (i IndentedCodeBlock) InlineString() string {
	return string(i.self)
}
func (i IndentedCodeBlock) String() string {
	return fmt.Sprintf("Type: %s, Depth: %d, InlineString: %s", IndentedCodeBlockType, i.depth, i.self)
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

type HyphenToken struct {
	canHorizontal bool
	self          []rune
}

func NewHyphen(canHorizontal bool, self []rune) HyphenToken {
	return HyphenToken{
		canHorizontal: canHorizontal,
		self:          self,
	}
}
func (h HyphenToken) Type() BlockType {
	return HyphenBlockType
}
func (h HyphenToken) CanHorizontal() bool {
	return h.canHorizontal
}
func (h HyphenToken) ConvertBlockToSetextHeading(target BlockToken) (BlockToken, BlockToken) {
	if target.Type() == ParagraphBlockType {
		return NewHeadingBlock(target.(ParagraphBlock).InlineString(), 2), NewSetextHeading()
	}

	if h.canHorizontal {
		return NewEmpty(), NewHorizontal()
	} else {
		return NewEmpty(), NewParagraphBlock(string(h.self), 0)
	}
}
func (h HyphenToken) String() string {
	return fmt.Sprintf("Type: %s, CanHorizontal: %t", HyphenBlockType, h.canHorizontal)
}

type EqualToken struct {
	self []rune
}

func NewEqual(self []rune) EqualToken {
	return EqualToken{
		self: self,
	}
}
func (e EqualToken) Type() BlockType {
	return EqualBlockType
}
func (e EqualToken) ConvertBlockToSetextHeading(target BlockToken) (BlockToken, BlockToken) {
	if target.Type() == ParagraphBlockType {
		return NewHeadingBlock(target.(ParagraphBlock).InlineString(), 1), NewSetextHeading()
	}

	return NewEmpty(), NewParagraphBlock(string(e.self), 0)
}
func (e EqualToken) String() string {
	return fmt.Sprintf("Type: %s", EqualBlockType)
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

type SetextHeading struct{}

func NewSetextHeading() SetextHeading {
	return SetextHeading{}
}
func (s SetextHeading) Type() BlockType {
	return SetextBlockType
}
func (s SetextHeading) String() string {
	return fmt.Sprintf("Type: %s", SetextBlockType)
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
