package token

const (
	LITERAL         = "Literal"
	GROUP           = "Group"
	UNCAPTURE_GROUP = "Uncapture_group"
	REPEAT          = "Repeat"
	OR              = "Or"
	BRACKET         = "Bracket"
)

type TokenType string

type Token struct {
	Value interface{}
	Type  TokenType
}
