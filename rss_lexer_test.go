package rss

import (
	"log"
	"testing"
)

func Test_RssLexerHookup(t *testing.T) {
	log.Printf("Hookup Succeeded - testing RSS lexer")
}

func Test_SimpleTag(t *testing.T) {
	if true {
		return
	}
	input := "<tag>"
	l := lex("simple tag", input)

	lexeme := l.nextItem()
	testLexeme(lexeme, itemOpenTag, "tag", t)
}

func Test_SelfClosingTag(t *testing.T) {
	if true {
		return
	}
	input := "<tag />"
	l := lex("self closing tag", input)

	lexeme := l.nextItem()
	testLexeme(lexeme, itemOpenTag, "tag", t)

	lexeme = l.nextItem()
	testLexeme(lexeme, itemSelfClosingTag, "", t)
}

func Test_TagWithAttributes(t *testing.T) {
	input := "<tag attr1=\"val1\" attr2=\"val2\">"
	l := lex("self closing tag", input)

	lexeme := l.nextItem()
	testLexeme(lexeme, itemOpenTag, "tag", t)

	lexeme = l.nextItem()
	testLexeme(lexeme, itemAttributeName, "attr1", t)
	lexeme = l.nextItem()
	testLexeme(lexeme, itemAttributeValue, "val1", t)
	lexeme = l.nextItem()
	testLexeme(lexeme, itemAttributeName, "attr2", t)
	lexeme = l.nextItem()
	testLexeme(lexeme, itemAttributeValue, "val2", t)
}

func Test_SelfClosingWithAttributes(t *testing.T) {
	input := "<tag a1=\"v1\" />"
	l := lex("self closing tag with attributes", input)

	lexeme := l.nextItem()
	testLexeme(lexeme, itemOpenTag, "tag", t)

	lexeme = l.nextItem()
	testLexeme(lexeme, itemAttributeName, "a1", t)
	lexeme = l.nextItem()
	testLexeme(lexeme, itemAttributeValue, "v1", t)
	lexeme = l.nextItem()
	testLexeme(lexeme, itemSelfClosingTag, "", t)
}

func Test_ClosingTag(t *testing.T) {
	input := "<tag><child></child></tag>"
	l := lex("closing tag", input)

	lexeme := l.nextItem()
	testLexeme(lexeme, itemOpenTag, "tag", t)
	lexeme = l.nextItem()
	testLexeme(lexeme, itemOpenTag, "child", t)
	lexeme = l.nextItem()
	testLexeme(lexeme, itemCloseTag, "child", t)
	lexeme = l.nextItem()
	testLexeme(lexeme, itemCloseTag, "tag", t)
}

func Test_Text(t *testing.T) {
	input := "<tag>Child text</tag>"
	l := lex("text", input)

	lexeme := l.nextItem()
	testLexeme(lexeme, itemOpenTag, "tag", t)
	lexeme = l.nextItem()
	testLexeme(lexeme, itemText, "Child text", t)
	lexeme = l.nextItem()
	testLexeme(lexeme, itemCloseTag, "tag", t)
}

func Test_CData(t *testing.T) {
	input := "<tag><![CDATA[child & <b>text</b>]]></tag>"
	l := lex("cdata", input)

	lexeme := l.nextItem()
	testLexeme(lexeme, itemOpenTag, "tag", t)
	lexeme = l.nextItem()
	testLexeme(lexeme, itemHtml, "child & <b>text</b>", t)
	lexeme = l.nextItem()
	testLexeme(lexeme, itemCloseTag, "tag", t)
}

func Test_NamespaceDeclaration(t *testing.T) {
	input := "<?xml version=\"1.0\" encoding=\"UTF-8\"?>"
	l := lex("namespace decl", input)

	lexeme := l.nextItem()
	testLexeme(lexeme, itemOpenTag, "xml", t)
	lexeme = l.nextItem()
	testLexeme(lexeme, itemAttributeName, "version", t)
	lexeme = l.nextItem()
	testLexeme(lexeme, itemAttributeValue, "1.0", t)
	lexeme = l.nextItem()
	testLexeme(lexeme, itemAttributeName, "encoding", t)
	lexeme = l.nextItem()
	testLexeme(lexeme, itemAttributeValue, "UTF-8", t)
	lexeme = l.nextItem()
	testLexeme(lexeme, itemNamespaceEnd, "", t)
}

func Test_NamespaceTag(t *testing.T) {
	input := "<rss ns1:a1=\"v1\" ns2:a2=\"v2\"></rss>"
	l := lex("namespaced tags", input)

	lexeme := l.nextItem()
	testLexeme(lexeme, itemOpenTag, "rss", t)
	lexeme = l.nextItem()
	testLexeme(lexeme, itemNamespace, "ns1", t)
	lexeme = l.nextItem()
	testLexeme(lexeme, itemAttributeName, "a1", t)
	lexeme = l.nextItem()
	testLexeme(lexeme, itemAttributeValue, "v1", t)
	lexeme = l.nextItem()
	testLexeme(lexeme, itemNamespace, "ns2", t)
	lexeme = l.nextItem()
	testLexeme(lexeme, itemAttributeName, "a2", t)
	lexeme = l.nextItem()
	testLexeme(lexeme, itemAttributeValue, "v2", t)
	lexeme = l.nextItem()
	testLexeme(lexeme, itemCloseTag, "rss", t)
}

func Test_CDataWithBracket(t *testing.T) {
	input := `<content:encoded>
<![CDATA[
<p><strong>EXCITING NEWS EVERYONE! Like Eric Lippert, Neil Gaiman enjoys soup!</strong></p> <p>That probably didn't make a whole lot of sense without context, so I should start by reposting <a href="http://blogs.msdn.com/b/ericlippert/archive/2011/07/08/my-buddy-neil-totally-agrees-with-me.aspx">M
]]>
<![CDATA[
Eric Lippert, who created one of the first Tolkien Web sites in 1993, sees the anti-Tolkien contingent as little more than literary snobs. [...blah blah blah...] Neil Gaiman, author of the fantasy series "The Sandman," said Tolkien "exists outside the orthodox canon of literature. You can't put him in a box." <strong>Like Lippert, Gaiman believes that Tolkien's commercial success is what drove his critics to jealous fury.
]]>
</content:encoded>`

	l := lex("namespaced tags", input)

	lexeme := l.nextItem()
	testLexeme(lexeme, itemNamespace, "content", t)
	lexeme = l.nextItem()
	testLexeme(lexeme, itemOpenTag, "encoded", t)
	lexeme = l.nextItem()
	log.Printf("text: %q\n", lexeme)
	lexeme = l.nextItem()
	log.Printf("text: %q\n", lexeme)
	lexeme = l.nextItem()
	testLexeme(lexeme, itemNamespace, "content", t)
	lexeme = l.nextItem()
	testLexeme(lexeme, itemCloseTag, "encoded", t)
}

func Test_Comments(t *testing.T) {
	input := "<tag><!-- child></child--></tag>"
	l := lex("comments", input)

	lexeme := l.nextItem()
	testLexeme(lexeme, itemOpenTag, "tag", t)
	lexeme = l.nextItem()
	testLexeme(lexeme, itemCloseTag, "tag", t)
}

func testLexeme(l lexeme, expectedType lexItemType, expectedVal string, t *testing.T) {
	if l.typ != expectedType {
		t.Errorf("lexeme item type (%q) not as expected (%q)", l.typ, expectedType)
	}

	if l.val != expectedVal {
		t.Errorf("lexeme val (%q) not as expected (%q)", l.val, expectedVal)
	}
}
