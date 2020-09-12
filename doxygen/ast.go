package doxygen

import (
	"github.com/alecthomas/participle"
	"github.com/alecthomas/participle/lexer"
	"github.com/alecthomas/participle/lexer/stateful"
)

var (
	Lexer = lexer.Must(stateful.New(stateful.Rules{
		"Root": {
			{"Func", `(@|\\)\w+`, nil},
			{"Text", `[^\s]+`, nil},
			{"Newline", `\n`, nil},
			{"whitespace", `\s+`, nil},
		},
	}))
	Parser = participle.MustBuild(
		&CommentBlock{},
		participle.Lexer(Lexer),
	)
)

type CommentBlock struct {
	Pos      lexer.Position
	Comments []Comment `parser:"@@*"`
}

type Comment struct {
	Pos     lexer.Position
	Doc     *Doc     `parser:"( @@"`
	Command *Command `parser:"| @@"`
	Newline *string   `parser:"| @Newline )"`
}

type Command struct {
	Pos   lexer.Position
	Func  *Func  `parser:"@@"`
	Words []Word `parser:"@@*"`
}

type Func struct {
	Pos  lexer.Position
	Name string `parser:"@Func"`
}

type Doc struct {
	Pos   lexer.Position
	Words []Word `parser:"@@+"`
}

type Word struct {
	Pos  lexer.Position
	Text string `parser:"@Text"`
}
