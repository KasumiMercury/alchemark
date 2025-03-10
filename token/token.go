package token

const (
	Heading    = "HEADING"
	Paragraph  = "PARAGRAPH"
	CodeBlock  = "CODE_BLOCK"
	Horizontal = "HORIZONTAL"
)

type Type string

type Token struct {
	Type Type
}
