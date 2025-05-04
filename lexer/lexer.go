package lexer

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type Token struct {
	literal     string
	name        string
	line_number int
}

func Lex(f *os.File, retchan chan []*Token) {
	defer f.Close()
	rdr := bufio.NewReader(f)

	line_number := 1
	tokens := []*Token{}
	buffer := ""
	is_symbolstring := false
	is_string := false
	is_comment := false
	is_multiline := false
	multiline_escape := false

	for {
		r, _, err := rdr.ReadRune()
		if err != nil {
			break
		}
		switch {
		case is_string:
			if r == '"' {
				is_string = false
				tokens = append(tokens, &Token{buffer, "string_lit", line_number})
				buffer = ""
			} else {
				buffer += string(r)
			}
		case is_comment:
			if r == '\n' {
				is_comment = false
			}
		case is_multiline:
			if r == '*' {
				multiline_escape = true
				continue
			}
			if r == '/' && multiline_escape {
				is_multiline = false
				multiline_escape = false
			}
			multiline_escape = false
		case r == '"':
			is_string = true
		case is_symbolstring:
			if is_symbol(r) {
				if can_conc_symbols(r, buffer) {
					tokens = append(tokens, &Token{buffer + string(r), buffer + string(r), line_number})
					buffer = ""
					is_symbolstring = false
				} else {
					tokens = append(tokens, &Token{buffer, buffer, line_number})
					buffer = string(r)
				}
			} else {
				tokens = append(tokens, &Token{buffer, buffer, line_number})
				buffer = ""
				is_symbolstring = false
				// Need to emulate the rest of the switch:
				switch {
				case r == '\r':
					//
				case r == '\n':
					line_number++
				case r == ' ' || r == '\t':
					//
				default:
					buffer = string(r)
				}
			}
		case is_symbol(r):
			if buffer != "" {
				tokens = append(tokens, &Token{buffer, get_name(buffer), line_number})
			}
			is_symbolstring = true
			buffer = string(r)
		case r == '\r':
			// Do nothing
		case r == '\n':
			if buffer == "" {
				line_number++
				continue
			}
			tokens = append(tokens, &Token{buffer, get_name(buffer), line_number})
			buffer = ""
			line_number++
		case r == ' ' || r == '\t':
			if buffer == "" {
				continue
			}
			tokens = append(tokens, &Token{buffer, get_name(buffer), line_number})
			buffer = ""
		default:
			buffer += string(r)
		}
	}
	retchan <- tokens
}

func get_name(name string) string {
	if name[0] == '\'' {
		return "char_lit"
	}
	is_num, lit_string := is_number(name)
	if is_num {
		return lit_string
	}
	switch name {
	case "fn", "int", "bool", "char", "string", "return", "pub", "float", "map", "const", "if", "else", "for", "match", "case", "break", "continue", "defer", "go", "struct", "package":
		return name
	case "true", "false":
		return "bool_lit"
	default:
		return "name"
	}
}

func can_conc_symbols(s1 rune, s2 string) bool {
	switch string(s1) + s2 {
	case "==", ">=", "<=", "||", "&&", "!=", "->":
		return true
	default:
		return false
	}
}

func is_number(num string) (bool, string) {
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
			return false, ""
		}
	}
	if used_period {
		return true, "double_lit"
	} else {
		return true, "int_lit"
	}
}

func is_symbol(symbol rune) bool {
	switch symbol {
	case '/', ',', '+', '-', '*', '.', '>', '<', '=', '(', ')', '!', '%', '&', '|', '[', ']', '{', '}', '~', ':':
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

func Print_tokens(tokens []*Token) {
	fmt.Println("Amount:", len(tokens))
	for _, pt := range tokens {
		t := *pt
		fmt.Println(t.name, "\""+t.literal+"\"")
	}
	fmt.Println()
}
