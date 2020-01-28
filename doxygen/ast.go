package doxygen

import (
	"github.com/alecthomas/participle"
	"github.com/alecthomas/participle/lexer"
	"github.com/alecthomas/participle/lexer/regex"
)

var (
	Lexer = lexer.Must(regex.New(`
	        Func    = (@|\\)[a-z][a-z]*
		Word    = [^\s][^\s]*
		Newline = \n

		whitespace = \s+
	`))

	Parser = participle.MustBuild(
		&CommentBlock{},
		participle.Lexer(Lexer),
	)
)

type CommentBlock struct {
	Pos      lexer.Position
	Comments []Comment `( @@ )*`
}

type Comment struct {
	Pos     lexer.Position
	Command *Command `( @@`
	Doc     *Doc     `| @@ )`
	Newline string   `@Newline`
}

type Command struct {
	Pos   lexer.Position
	Func  *Func `@@`
	Words []Word `( @@ )*`
}

type Doc struct {
	Pos   lexer.Position
	Words []Word `( @@ )*`
}

type Func struct {
	Pos  lexer.Position
	Name string `@Func`
}

type Word struct {
	Pos  lexer.Position
	Text string `@Word`
}
