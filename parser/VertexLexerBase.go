package parser

import "github.com/antlr4-go/antlr/v4"

// VertexLexerBase is the superclass for the ANTLR4-generated VertexLexer.
// The grammar declares `superClass = VertexLexerBase`.
// Add any custom lexer helper methods here as the language grows.
type VertexLexerBase struct {
	*antlr.BaseLexer
}