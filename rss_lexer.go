package rss

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// Thanks again, Rob Pike
// http://blog.golang.org/2011/09/two-go-talks-lexical-scanning-in-go-and.html
// http://golang.org/src/pkg/text/template/parse/lex.go

// Lexeme is a parsed element
type lexeme struct {
	typ lexItemType
	val string
}

// String() for item
func (l lexeme) String() string {
	if len(l.val) > 10 {
		return fmt.Sprintf("{%s, %.10q...}", l.typ, l.val)
	}
	return fmt.Sprintf("{%s, %q}", l.typ, l.val)
}

// types of lex items
type lexItemType int

const (
	itemEOF lexItemType = iota
	itemError
	itemNewline
	itemNamespace
	itemNamespaceEnd
	itemOpenTag
	itemCloseTag
	itemSelfClosingTag
	itemAttributeName
	itemAttributeValue
	itemText
)

// for pretty printing
var lexItemNames map[lexItemType]string

var itemName = map[lexItemType]string{
	itemEOF:            "End of File",
	itemError:          "error",
	itemNewline:        "newline",
	itemNamespace:      "namespace",
	itemNamespaceEnd:   "end namespace",
	itemOpenTag:        "open tag",
	itemCloseTag:       "close tag",
	itemSelfClosingTag: "self-closing tag",
	itemAttributeName:  "attribute name",
	itemAttributeValue: "attribute value",
	itemText:           "text",
}

func (item lexItemType) String() string {
	s := itemName[item]
	if s == "" {
		return fmt.Sprintf("item %d", int(item))
	}
	return s
}

// state function represents 'current state' combined with 'next action'
type stateFn func(*lexer) stateFn

// lexer holds the state of the scanner
type lexer struct {
	name    string      // name of the input (for error reporting)
	input   string      // the string being scanned
	state   stateFn     // next lexing function
	pos     int         // current position in the input string
	start   int         // start position of this item
	width   int         // length of the last input rune
	lexemes chan lexeme // channel of scanned lexemes
}

// create a new lexer
func lex(name, input string) *lexer {
	l := &lexer{
		name:    name,
		input:   input,
		state:   lexContentStart,
		lexemes: make(chan lexeme, 3),
	}
	return l
}

// represent EOF when we're parsing an input string
const eof = -1

// next returns the next rune in the string
func (l *lexer) next() (r rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return r
}

// backup steps back one rune. Should only be called once per next()
func (l *lexer) backup() {
	l.pos -= l.width
}

// peek returns but does not consume the next rune
func (l *lexer) peek() (r rune) {
	r = l.next()
	l.backup()
	return r
}

func (l *lexer) peekForward(amnt int) (r rune) {
	for i := 0; i < amnt; i++ {
		r = l.next()
		if r == eof {
			l.errorf("peeked past end of content")
			return eof
		}
	}
	for i := 0; i < amnt; i++ {
		l.backup()
	}
	return r
}

// emit an item back to the client
func (l *lexer) emit(t lexItemType) {
	l.lexemes <- lexeme{t, l.previewCurrent()}
	l.start = l.pos
}

func (l *lexer) previewCurrent() string {
	return l.input[l.start:l.pos]
}

// skips the pending input
func (l *lexer) ignore() {
	l.start = l.pos
}

// accepts a rune if it's in the validRunes string
func (l *lexer) accept(validRunes string) bool {
	if strings.IndexRune(validRunes, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

// accepts a series of runes that match the validRunes characters
func (l *lexer) acceptRun(validRunes string) {
	for strings.IndexRune(validRunes, l.next()) >= 0 {
	}
	l.backup()
}

// accepts a series of runes that are not in the invalidRunes characters
func (l *lexer) acceptRunUntil(invalidRunes string) {
	for {
		rune := l.next()
		if rune == eof {
			break
		}

		isFound := strings.IndexRune(invalidRunes, rune)
		if isFound >= 0 {
			break
		}
	}
	l.backup()
}

// skips all spaces and tabs from the current position
func (l *lexer) skipWhitespace() {
	l.acceptRun(" \t\r\n")
	l.ignore()
}

// which line are we currently on?
func (l *lexer) lineNumber() int {
	return 1 + strings.Count(l.input[:l.pos], "\n")
}

// error returns an error token and terminates the scan by passing
// back a nil pointer that will be the next state, terminating l.nextItem.
func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.lexemes <- lexeme{itemError, fmt.Sprintf(format, args...)}
	return nil
}

// run runs the state machine for the lexer.
func (l *lexer) run() {
	for l.state = lexContentStart; l.state != nil; {
		l.state = l.state(l)
	}
}

// nextItem returns the next item from the input
func (l *lexer) nextItem() lexeme {
	for {
		select {
		case lexeme := <-l.lexemes:
			if lexeme.typ == itemError {
				panic("ERROR: " + lexeme.val + "\n")
			}
			return lexeme
		default:
			l.state = l.state(l)
		}
	}
	panic("Not reached")
}

func lexContentStart(l *lexer) stateFn {
	l.skipWhitespace()
	switch l.peek() {
	case eof:
		l.emit(itemEOF)
		return nil
	case '<':
		return lexTagStart
	default:
		return lexTagContents
	}

	l.errorf("Fell through content start")
	return nil
}

func lexTagStart(l *lexer) stateFn {
	isClosingTag := false

	l.accept("<")
	switch l.peek() {
	case '!':
		l.backup()
		return lexCData
	case '?':
		l.backup()
		return lexNamespace
	case '/':
		l.next()
		isClosingTag = true
	}

	l.ignore()

	// we might come across a namespaced tag, e.g. <ns:tag>
	// so we're going to loop and break out as needed
	for {
		l.acceptRunUntil("/:> \t\r\n")
		switch l.peek() {
		case ':':
			l.emit(itemNamespace) // namespace, just loop around and keep parsing
			l.accept(":")
			l.ignore()
		case ' ':
			l.emit(itemOpenTag)
			return lexAttributes
		case '>':
			if isClosingTag {
				l.emit(itemCloseTag)
				l.accept(">")
				return lexContentStart
			} else {
				l.emit(itemOpenTag)
				l.accept(">")
				return lexTagContents
			}
		default:
			l.errorf("error parsing tag, unexpected symbol: %v", l.peek())
			return nil
		}
	}

	panic("broke out of parsing loop")
	return nil
}

func lexAttributes(l *lexer) stateFn {
	l.skipWhitespace()

	// check our breakout conditions
	switch l.peek() {
	case '/':
		l.accept("/")
		if l.peek() == '>' {
			l.accept(">")
			l.ignore()
			l.emit(itemSelfClosingTag)
			return lexContentStart
		} else {
			l.backup()
		}
	case '?', '!':
		l.acceptRun("!?>")
		l.ignore()
		l.emit(itemNamespaceEnd)
		return lexTagContents
	case '>':
		l.accept(">")
		l.ignore()
		return lexTagContents
	}

Loop:
	// OK, we should have more attributes, so parse them out
	for {
		l.acceptRunUntil(":=")
		switch l.peek() {
		case ':':
			l.emit(itemNamespace)
			l.accept(":")
			l.ignore()
		case '=':
			l.emit(itemAttributeName)
			break Loop
		case eof:
			l.errorf("lex attributes: did not find :=")
			return nil
		default:
			fmt.Printf("Stuck in loop? %v\n", l.previewCurrent())
		}
	}
	l.acceptRun("=\"")
	l.ignore()
	l.acceptRunUntil("\"")
	l.emit(itemAttributeValue)
	l.accept("\"")

	return lexAttributes
}

func lexCData(l *lexer) stateFn {
	if l.peek() == '<' && l.peekForward(2) == '!' && l.peekForward(3) == '-' && l.peekForward(4) == '-' {
		return lexComment
	}
	l.acceptRun("<![")
	// skip CDATA
	l.acceptRunUntil("[")
	l.accept("[")
	l.skipWhitespace()
	l.ignore()

Loop:
	for {
		l.acceptRunUntil("]")
		if l.peek() == ']' && l.peekForward(2) == ']' && l.peekForward(3) == '>' {
			break Loop
		} else if l.peek() == eof {
			l.errorf("CDATA EOF reached")
			return nil
		}
		l.accept("]")
	}

	l.emit(itemText)
	l.acceptRun("]]>")
	l.ignore()
	return lexContentStart
}

func lexComment(l *lexer) stateFn {
	l.acceptRun("<!-")

Loop:
	for {
		l.acceptRunUntil("-")
		if l.peekForward(2) == '-' && l.peekForward(3) == '>' {
			l.accept("->")
			break Loop
		}
		l.accept("-")
	}
	l.acceptRun("->")
	l.ignore()
	return lexContentStart
}

func lexNamespace(l *lexer) stateFn {
	l.acceptRun("<!?")
	l.ignore()
	for {
		l.acceptRunUntil("?!:> ")
		switch l.peek() {
		case ':':
			l.emit(itemNamespace) // namespace, just loop around and keep parsing
		case ' ':
			l.emit(itemOpenTag)
			return lexAttributes
		case '!', '?':
			// I'm not sure we will get here
			l.acceptRun("?!>")
			l.ignore()
			l.emit(itemNamespaceEnd)
			return lexContentStart
		default:
			l.errorf("error parsing tag, unexpected symbol: %v", l.peek())
			return nil
		}
	}

	return nil
}

func lexTagContents(l *lexer) stateFn {
	l.skipWhitespace()
	if l.peek() == '<' {
		return lexContentStart
	}
	l.acceptRunUntil("<")
	l.emit(itemText)
	return lexContentStart
}
