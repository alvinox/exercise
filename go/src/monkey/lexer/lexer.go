package lexer

import (
    "monkey/token"
)

type Lexer struct {
    input        string
    position     int  // current position in input (points to current char)
    readPosition int  // current reading position in input (after current char)
    ch           byte // current char under examination
}

func New(input string) *Lexer {
    var l *Lexer = &Lexer{input : input}
    l.readChar()
    return l
}

func (l *Lexer) NextToken() token.Token {
    var tok token.Token

    l.skipWhiteSpace()

    switch l.ch {
    case '=':
        if l.peekChar() == '=' {
            // ch := l.ch
            l.readChar()
            tok = token.Token{Type:token.EQ, Literal:"=="}
        } else {
            tok = newToken(token.ASSIGN, l.ch)
        }
    case ';':
        tok = newToken(token.SEMICOLON, l.ch)
    case ':':
        tok = newToken(token.COLON, l.ch)
    case '(':
        tok = newToken(token.LPAREN, l.ch)
    case ')':
        tok = newToken(token.RPAREN, l.ch)
    case ',':
        tok = newToken(token.COMMA, l.ch)
    case '+':
        tok = newToken(token.PLUS, l.ch)
    case '-':
        tok = newToken(token.MINUS, l.ch)
    case '*':
        tok = newToken(token.ASTERISK, l.ch)
    case '/':
        tok = newToken(token.SLASH, l.ch)
    case '!':
        if l.peekChar() == '=' {
            l.readChar()
            tok = token.Token{Type:token.NOT_EQ, Literal:"!="}
        } else {
            tok = newToken(token.BANG, l.ch)
        }
    case '<':
        tok = newToken(token.LT, l.ch)
    case '>':
        tok = newToken(token.GT, l.ch)
    case '{':
        tok = newToken(token.LBRACE, l.ch)
    case '}':
        tok = newToken(token.RBRACE, l.ch)
    case '[':
        tok = newToken(token.LBRACKET, l.ch)
    case ']':
        tok = newToken(token.RBRACKET, l.ch)
    case '"':
        tok.Type = token.STRING
        tok.Literal = l.readString()
    case 0:
        tok.Type = token.EOF
        tok.Literal = ""
    default:
        if isLetter(l.ch) {
            tok.Literal = l.readIdentifier()
            tok.Type = token.LookupIdent(tok.Literal)
            return tok
        } else if isDigit(l.ch) {
            tok.Type = token.INT
            tok.Literal = l.readNumber()
            return tok
        } else {
            return newToken(token.ILLEGAL, l.ch)
        }
    }
    l.readChar()
    return tok
}

func (l *Lexer) readChar() {
    if l.readPosition >= len(l.input) {
        l.ch = 0
    } else {
        l.ch = l.input[l.readPosition]
    }

    l.position = l.readPosition
    l.readPosition += 1
}

func (l *Lexer) peekChar() byte {
    if l.readPosition > len(l.input) {
        return 0
    } else {
        return l.input[l.readPosition]
    }
}

func (l *Lexer) readIdentifier() string {
    var position int = l.position
    for isLetter(l.ch) {
        l.readChar()
    }
    return l.input[position : l.position]
}

func (l *Lexer) readNumber() string {
    var position int = l.position
    for isDigit(l.ch) {
        l.readChar()
    }
    return l.input[position : l.position]
}

func (l *Lexer) readString() string {
    var position int = l.position + 1
    for {
        l.readChar()
        if l.ch == '"' {
            break
        }
    }
    return l.input[position : l.position]
}

func (l *Lexer) skipWhiteSpace() {
    for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
        l.readChar()
    }
}

func isLetter(ch byte) bool {
    return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
    return '0' <= ch && ch <= '9'
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
    return token.Token{Type:tokenType, Literal:string(ch)}
}