package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/elves/upgrade-scripts-for-0.7/parse"
)

func main() {
	filenames := os.Args[1:]
	if len(filenames) == 0 {
		fixed, ok := fix("[stdin]", func() ([]byte, error) {
			return ioutil.ReadAll(os.Stdin)
		})
		if ok {
			os.Stdout.Write(fixed)
		}
	} else {
		// Fix files.
		for _, filename := range filenames {
			fixed, ok := fix(filename, func() ([]byte, error) {
				return ioutil.ReadFile(filename)
			})
			if ok {
				ioutil.WriteFile(filename, fixed, 0644)
			}
		}
	}
}

func fix(filename string, readall func() ([]byte, error)) ([]byte, bool) {
	src, err := readall()
	if err != nil {
		log.Printf("cannot read %s: %s; skipping", filename, err)
		return nil, false
	}
	if !utf8.Valid(src) {
		log.Printf("%s not utf8; skipping", filename)
		return nil, false
	}
	// os.Stdout.WriteString(fix("[stdin]", src))
	chunk, err := parse.Parse(filename, string(src))
	if err != nil {
		log.Printf("cannot parse %s: %s; skipping", filename, err)
		return nil, false
	}

	buf := new(bytes.Buffer)
	fixNode(chunk, buf)
	return buf.Bytes(), true
}

var (
	suppress = map[string]bool{
		"\n": true, ";": true, "in": true,
		"do": true, "done": true,
		"then": true, "fi": true,
		"tried": true,
	}
	addSpace = map[string]bool{
		"elif": true, "else": true,
		"except": true, "finally": true,
	}
)

func fixNode(n parse.Node, w io.Writer) {
	if ctrl, ok := n.(*parse.Control); ok {
		for _, child := range n.Children() {
			switch {
			case suppress[child.SourceText()]:
				// suppress it
			case addSpace[child.SourceText()]:
				w.Write([]byte(" " + child.SourceText()))
			case child == ctrl.Array:
				// for i in< lorem ipsum>
				w.Write([]byte("["))
				for i, item := range ctrl.Array.Compounds {
					if i > 0 {
						w.Write([]byte(" "))
					}
					fixNode(item, w)
				}
				w.Write([]byte("]"))
			case isCondition(child, ctrl):
				// if <a | b>
				w.Write([]byte(" ("))
				for _, ch := range trimChunk(child.(*parse.Chunk)) {
					fixNode(ch, w)
				}
				w.Write([]byte(")"))
			case parse.IsChunk(child):
				// for i in lorem ipsum; do <
				//      echo haha >
				// done
				w.Write([]byte(" {"))
				if child == ctrl.ExceptBody {
					w.Write([]byte("\n"))
				}
				fixNode(child, w)
				w.Write([]byte("}"))
			default:
				fixNode(child, w)
			}
		}
	} else if assigns := fixableAssignment(n); assigns != nil {
		for i, assign := range assigns {
			if i > 0 {
				w.Write([]byte("; "))
			}
			if assign.Left.Head.Type == parse.Braced {
				for i, v := range assign.Left.Head.Braced {
					if i > 0 {
						w.Write([]byte(" "))
					}
					fixNode(v, w)
				}
			} else {
				fixNode(assign.Left, w)
			}
			w.Write([]byte(" = "))
			fixNode(assign.Right, w)
		}
	} else if len(n.Children()) == 0 {
		text := n.SourceText()
		if text == "?(" {
			text = "("
		}
		w.Write([]byte(text))
	} else {
		for _, child := range n.Children() {
			fixNode(child, w)
		}
	}
}

func fixableAssignment(n parse.Node) []*parse.Assignment {
	fn, ok := n.(*parse.Form)
	if ok && fn.Vars == nil && fn.Head == nil {
		return fn.Assignments
	}
	return nil
}

func isCondition(ch parse.Node, ctrl *parse.Control) bool {
	if ch == ctrl.Condition {
		return true
	}
	for _, cond := range ctrl.Conditions {
		if ch == cond {
			return true
		}
	}
	return false
}

func allSpaces(s string) bool {
	return strings.TrimFunc(s, parse.IsSpace) == ""
}

func trimChunk(ch *parse.Chunk) []parse.Node {
	var first, last int
	var foundFirst bool
	for i, ch := range ch.Children() {
		if !trunkTrimable(ch.SourceText()) {
			if !foundFirst {
				first = i
				foundFirst = true
			}
			last = i
		}
	}
	return ch.Children()[first : last+1]
}

func trunkTrimable(s string) bool {
	return strings.TrimFunc(s, func(r rune) bool {
		return parse.IsSpaceOrNewline(r) || r == ';'
	}) == ""
}
