package lexer

import (
	"bufio"
	"log"
	"os"
)

type token struct {
	literal     string
	name        string
	line_number int
}

func Lex(path string) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal("Error opening file", path)
	}
	defer f.Close()
	rdr := bufio.NewReader(f)

	line_number := 1
	tokens := []*token{}
	current_token := token{line_number: 1}

	buffer := ""
	is_symbolstring := false
	is_string := false
	is_comment := false
	is_multiline := false
	multiline_buffer := ' '
	for {
		r, _, err := rdr.ReadRune()
		if err != nil {
			break
		}
		switch {
		case is_string:
			if r == '"' {
				is_string = false
				current_token.literal = buffer + string(r)
				current_token.name = "string_lit"
				buffer = ""
				continue
			} else {
				buffer += string(r)
			}
		case is_comment:
			if r == '\n' {
				is_comment = true
			}
			continue
		case is_multiline:
			if r == '*' {
				continue
			}
			if r == '/' && multiline_buffer == '*' {
				is_multiline = false
				continue
			}
		case is_symbolstring:
			if is_symbol(r) {
				if can_conc_symbols(r, buffer) {
					current_token.literal = string(r) + buffer
					current_token.name = string(r) + buffer
					tokens = append(tokens, &current_token)
					buffer = ""
				} else {
					current_token.literal = buffer
					current_token.name = buffer
					tokens = append(tokens, &current_token)
					buffer = string(r)
				}
				continue
			}
			fallthrough
		case is_symbol(r):
			is_symbolstring = true
			buffer += string(r)
		case r == '\r':
			// Do nothing
		case r == '\n':
			line_number++
			current_token.line_number = line_number
		case r == '"':
			is_string = true
		case r == ' ':
			current_token.literal = buffer
			buffer = ""
			tokens = append(tokens, &current_token)
			go determine_token(&current_token)
		default:
			buffer += string(r)
		}
	}
}

func determine_token(*token) {
	if token.literal[0] == "'" {
		token.name = "char_lit"
		return
	}
	name := is_number(token.literal)
	if name != "" {
		token.name = name
		return
	}
	switch token.literal {
	case "fn", "int", "bool", "char", "string", "return", "pub", "float", "map", "const", "if", "else", "for", "match", "case", "break", "continue", "defer", "go":
		token.name = token.literal
	case "true", "false":
		token.name = "bool_lit"
	default:
		token.name = "name"
	}
}

func can_conc_symbols(s1 rune, s2 string) bool {
	switch string(s1) + s2 {
	case "==", ">=", "<=", "||", "&&", "!=":
		return true
	default:
		return false
	}
}

func is_number(num string) string {
	used_period := false
	for _, c := range num {
		switch c {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			//
		case '.':
			if used_period == true {
				log.Fatal("Double period used in double literal: ", num)
			}
			used_period = true
		default:
			return ""
		}
	}
	if used_period {
		return "double_lit"
	} else {
		return "int_lit"
	}
}

func is_symbol(symbol rune) bool {
	switch symbol {
	case '/', '+', '-', '*', '.', '>', '<', '=', '(', ')', '!', '%', '&', '|', '[', ']', '{', '}', '~', ':':
		return true
	default:
		return false
	}
}

func is_digit(digit rune) bool {
	switch digit {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return true
	default:
		return false
	}
}
