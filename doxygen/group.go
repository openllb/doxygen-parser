package doxygen

import (
	"fmt"
	"io"
	"strings"

	"github.com/alecthomas/participle/lexer"
)

type Group struct {
	Doc    string
	Params []Param
	Return Return
}

type Param struct {
	Name        string
	Description string
}

type Return struct {
	Description string
}

func Parse(r io.Reader) (*Group, error) {
	var cb CommentBlock
	err := Parser.Parse(r, &cb)
	if err != nil {
		return nil, err
	}

	group := &Group{}

	var (
		fun        *Func
		docs       []string
		words      []string
		endCommand bool
	)
	for _, comment := range cb.Comments {
		endCommand = false
		if comment.Command == nil {
			if fun == nil {
				if len(words) > 0 {
					docs = append(docs, "\n")
				}
				words = []string{}
				for _, word := range comment.Doc.Words {
					words = append(words, word.Text)
				}
				docs = append(docs, fmt.Sprintf("%s\n", strings.Join(words, " ")))
				continue
			}

			if len(comment.Doc.Words) == 0 {
				endCommand = true
			} else {
				for _, word := range comment.Doc.Words {
					words = append(words, word.Text)
				}
			}
		} else if fun != nil {
			endCommand = true
		}

		if endCommand {
			err = group.AddCommand(fun, words)
			if err != nil {
				return group, err
			}
			fun = nil
		}

		if comment.Doc != nil {
			continue
		}

		fun = comment.Command.Func
		words = []string{}
		for _, word := range comment.Command.Words {
			words = append(words, word.Text)
		}
	}

	if fun != nil {
		err = group.AddCommand(fun, words)
		if err != nil {
			return group, err
		}
	}

	group.Doc = strings.Join(docs, "")
	return group, nil
}

func (g *Group) AddCommand(fun *Func, words []string) error {
	switch strings.TrimLeft(fun.Name, "@//") {
	case "param":
		if len(words) == 0 {
			return ErrAtToken{fun.Pos, "param must have a name"}
		}

		g.Params = append(g.Params, Param{
			Name:        words[0],
			Description: strings.Join(words[1:], " "),
		})
	case "return", "returns":
		g.Return = Return{
			Description: strings.Join(words, " "),
		}
	}
	return nil
}

type ErrAtToken struct {
	Pos     lexer.Position
	Message string
}

func (e ErrAtToken) Error() string {
	return fmt.Sprintf("%s:%d:%d: %s", e.Pos.Filename, e.Pos.Line, e.Pos.Column, e.Message)
}
