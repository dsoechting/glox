package scanner

import (
	"dsoechting/glox/error"
	"dsoechting/glox/token"
	"fmt"
	"strconv"
)

type Scanner struct {
	source  string
	tokens  []token.Token
	start   int
	current int
	line    int
}

func Create(source string) *Scanner {
	return &Scanner{
		source: source,
		//TODO maybe smarter init size?
		tokens:  make([]token.Token, 0),
		start:   0,
		current: 0,
		line:    1,
	}
}

func (s *Scanner) ScanTokens() ([]token.Token, error) {

	for !s.isAtEnd() {
		s.start = s.current
		err := s.scanToken()
		if err != nil {
			return nil, err
		}

	}
	newToken := token.Create(token.EOF, "", nil, s.line)
	s.tokens = append(s.tokens, *newToken)
	return s.tokens, nil
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) scanToken() error {
	currentRune := s.advance()
	switch currentRune {
	case '(':
		s.addTokenSimple(token.LEFT_PAREN)
		break
	case ')':
		s.addTokenSimple(token.RIGHT_PAREN)
		break
	case '{':
		s.addTokenSimple(token.LEFT_BRACE)
		break
	case '}':
		s.addTokenSimple(token.RIGHT_BRACE)
		break
	case ',':
		s.addTokenSimple(token.COMMA)
		break
	case '.':
		s.addTokenSimple(token.DOT)
		break
	case '-':
		s.addTokenSimple(token.MINUS)
		break
	case '+':
		s.addTokenSimple(token.PLUS)
		break
	case ';':
		s.addTokenSimple(token.SEMICOLON)
		break
	case '*':
		s.addTokenSimple(token.STAR)
		break
	case '!':
		if s.match('=') {
			s.addTokenSimple(token.BANG_EQUAL)
		} else {
			s.addTokenSimple(token.BANG)
		}
		break
	case '=':
		if s.match('=') {
			s.addTokenSimple(token.EQUAL_EQUAL)
		} else {
			s.addTokenSimple(token.EQUAL)
		}
		break
	case '<':
		if s.match('=') {
			s.addTokenSimple(token.LESS_EQUAL)
		} else {
			s.addTokenSimple(token.LESS)
		}
		break
	case '>':
		if s.match('=') {
			s.addTokenSimple(token.GREATER_EQUAL)
		} else {
			s.addTokenSimple(token.GREATER)
		}
		break
	case '?':
		s.addTokenSimple(token.QUESTION)
		break
	case ':':
		s.addTokenSimple(token.COLON)
		break
	case '/':
		// Single line comments
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
			// Block comments
		} else if s.match('*') {
			for !(s.peek() == '*' && s.peekNext() == '/') {
				s.advance()
			}
			s.match('*')
			s.match('/')
		} else {
			s.addTokenSimple(token.SLASH)
		}
	case ' ':
	case '\r':
	case '\t':
		// Ignore whitespace
		break
	case '\n':
		s.line++
		break
	case '"':
		s.scanString()
		break
	default:
		if isDigit(currentRune) {
			s.number()
		} else if isAlphaNumeric(currentRune) {
			s.identifier()
		} else {
			errorString := fmt.Sprintf("Unexpected character: %c", currentRune)
			return glox_error.Create(s.line, "", errorString)
		}
	}
	return nil
}

func isAlphaNumeric(r rune) bool {
	return isAlpha(r) || isDigit(r)
}

func isAlpha(r rune) bool {
	return (r >= 'a' && r <= 'z') ||
		(r >= 'A' && r <= 'Z') ||
		r == '_'
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func (s *Scanner) identifier() error {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}
	text := s.source[s.start:s.current]
	tokenType, ok := Keywords[text]
	if !ok {
		tokenType = token.IDENTIFIER
	}
	s.addTokenSimple(tokenType)
	return nil
}

func (s *Scanner) number() error {
	for isDigit(s.peek()) {
		s.advance()
	}
	if s.peek() == '.' && isDigit(s.peekNext()) {
		s.advance()
		for isDigit(s.peek()) {
			s.advance()
		}
	}
	value, err := strconv.ParseFloat(s.source[s.start:s.current], 64)
	if err != nil {
		return glox_error.Create(s.line, "", "Could not parse number")
	}
	s.addToken(token.NUMBER, value)
	return nil
}

func (s *Scanner) scanString() error {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}
	if s.isAtEnd() {
		return glox_error.Create(s.line, "", "Unterminated string literal")
	}
	s.advance()

	value := s.source[s.start+1 : s.current-1]
	s.addToken(token.STRING, value)
	return nil
}

// Only advance if it's the rune that we want
func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}
	currentRune := rune(s.source[s.current])
	if currentRune != expected {
		return false
	}
	s.current += 1
	return true
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return rune(0)
	}
	return rune(s.source[s.current])
}

func (s *Scanner) peekNext() rune {
	nextPos := s.current + 1
	if s.isAtEnd() || (nextPos >= len(s.source)) {
		return rune(0)
	}
	return rune(s.source[nextPos])
}

func (s *Scanner) advance() rune {
	currentRune := s.source[s.current]
	s.current += 1
	return rune(currentRune)
}

func (s *Scanner) addTokenSimple(tokenType token.TokenType) {
	s.addToken(tokenType, nil)
}

func (s *Scanner) addToken(tokenType token.TokenType, literal any) {
	text := s.source[s.start:s.current]
	newToken := token.Create(tokenType, text, literal, s.line)
	s.tokens = append(s.tokens, *newToken)
}
