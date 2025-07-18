// Code generated by "stringer -type=TokenType"; DO NOT EDIT.

package token

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[LEFT_PAREN-0]
	_ = x[RIGHT_PAREN-1]
	_ = x[LEFT_BRACE-2]
	_ = x[RIGHT_BRACE-3]
	_ = x[COMMA-4]
	_ = x[DOT-5]
	_ = x[MINUS-6]
	_ = x[PLUS-7]
	_ = x[SEMICOLON-8]
	_ = x[SLASH-9]
	_ = x[STAR-10]
	_ = x[QUESTION-11]
	_ = x[COLON-12]
	_ = x[BANG-13]
	_ = x[BANG_EQUAL-14]
	_ = x[EQUAL-15]
	_ = x[EQUAL_EQUAL-16]
	_ = x[GREATER-17]
	_ = x[GREATER_EQUAL-18]
	_ = x[LESS-19]
	_ = x[LESS_EQUAL-20]
	_ = x[IDENTIFIER-21]
	_ = x[STRING-22]
	_ = x[NUMBER-23]
	_ = x[AND-24]
	_ = x[CLASS-25]
	_ = x[ELSE-26]
	_ = x[FALSE-27]
	_ = x[FUN-28]
	_ = x[FOR-29]
	_ = x[IF-30]
	_ = x[NIL-31]
	_ = x[OR-32]
	_ = x[PRINT-33]
	_ = x[RETURN-34]
	_ = x[SUPER-35]
	_ = x[THIS-36]
	_ = x[TRUE-37]
	_ = x[VAR-38]
	_ = x[WHILE-39]
	_ = x[EOF-40]
}

const _TokenType_name = "LEFT_PARENRIGHT_PARENLEFT_BRACERIGHT_BRACECOMMADOTMINUSPLUSSEMICOLONSLASHSTARQUESTIONCOLONBANGBANG_EQUALEQUALEQUAL_EQUALGREATERGREATER_EQUALLESSLESS_EQUALIDENTIFIERSTRINGNUMBERANDCLASSELSEFALSEFUNFORIFNILORPRINTRETURNSUPERTHISTRUEVARWHILEEOF"

var _TokenType_index = [...]uint8{0, 10, 21, 31, 42, 47, 50, 55, 59, 68, 73, 77, 85, 90, 94, 104, 109, 120, 127, 140, 144, 154, 164, 170, 176, 179, 184, 188, 193, 196, 199, 201, 204, 206, 211, 217, 222, 226, 230, 233, 238, 241}

func (i TokenType) String() string {
	if i < 0 || i >= TokenType(len(_TokenType_index)-1) {
		return "TokenType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _TokenType_name[_TokenType_index[i]:_TokenType_index[i+1]]
}
