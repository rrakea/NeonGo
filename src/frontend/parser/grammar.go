package parser

import "fmt"

type Grammar struct {
	rules []Rule
}

type Rule struct {
	left  string
	right []string
}

func Get_Grammar() *Grammar {
	gr := Grammar{}
	rules := []Rule{}
	g := grammar_string()
	buf := ""
	left := ""
	right := []string{}
	for _, r := range []rune(g) {
		switch r {
		case ' ':
			if left == "" {
				left = buf
				buf = ""
			} else {
				right = append(right, buf)
				buf = ""
			}
		case '\n':
			if left != "" {
				rules = append(rules, Rule{left, append(right, buf)})
			} else {
				rules = append(rules, Rule{buf, right})
			}
			left = ""
			buf = ""
			right = []string{}
		default:
			buf += string(r)
		}
	}
	gr.rules = rules
	return &gr
}

func grammar_string() string {
	// The first word is the left side
	// When a rule only has one word it is an epsilon rule

	// TODO: Interfaces, string[]{}
	s := "" +
		"S PACKAGE CODE\n" +
		"CODE FN CODE\n" +
		"CODE STRUCT CODE\n" +
		"CODE ENUM CODE\n" +
		"CODE\n" +
		//"CODE INTERFACE\n" +

		"PACKAGE package name\n" +

		// Function Declaration
		"FN VIS fn RECIEVER name ( ARGS ) RET BLOCK\n" +
		"VIS pub\n" +
		"VIS\n" +
		"RECIEVER ( name TYPE )\n" +
		"RECIEVER\n" +
		"RET -> TYPE RETLIST QUESTION\n" +
		"RET\n" +
		"RETLIST , TYPE RETLIST\n" +
		"QUESTION ?\n" +
		"QUESTION\n" +
		"ARGS ARG ARGLIST\n" +
		"ARG name TYPE\n" +
		"ARGLIST , ARG ARGLIST\n" +
		"ARGLIST\n" +
		"ARGS\n" +

		// Statements
		"BLOCK { STATEMENTBLOCK }\n" +
		"STATEMENTBLOCK STATEMENT STATEMENTBLOCK\n" +
		"STATEMENTBLOCK\n" +
		"STATEMENT RETURN\n" +
		"STATEMENT ERROR\n" +
		"STATEMENT VARDEC\n" +
		"STATEMENT VARASSIGN\n" +
		"STATEMENT FNCALL\n" +
		"STATEMENT IF\n" +
		"STATEMENT FOR\n" +
		"STATEMENT MATCH\n" +
		"STATEMENT CHECK\n" +
		"STATEMENT JUMP\n" +

		"JUMP break\n" +
		"JUMP continue\n" +

		"FNCALL name ( CALLARG )\n" +
		"CALLARG EXPRESSION CALLARGCONT\n" +
		"CALLARG\n" +
		"CALLARGCONT , EXPRESSION CALLARGCONT\n" +
		"CALLARGCONT\n" +

		"IF if EXPRESSION BLOCK ELSE\n" +
		"ELSE else if EXPRESSION BLOCK ELSE\n" +
		"ELSE else BLOCK\n" +

		"RETURN return EXRESSION\n" +
		"RETURN return\n" +

		"ERROR error EXPRESSION\n" +

		"CHECK check name BLOCK\n" +

		"VARDEC name := EXPRESSION\n" +
		"VARDEC TYPE name\n" +

		"VARASSIGN name EQUALS EXPRESSION\n" +

		"EQUALS =\n" +
		"EQUALS +=\n" +
		"EQUALS -=\n" +
		"EQUALS *=\n" +
		"EQUALS /=\n" +

		"FOR for FORCOND BLOCK\n" +
		"FORCOND\n" + // for {}
		"FORCOND ~ EXPRESSION\n" + // for ~ len(a) {}
		"FORCOND name ~ EXPRESSION\n" + // for i ~ len(a) {}
		"FORCOND name , name ~ EXPRESSION\n" + // for key,val ~ imap {}

		// match s {
		//    1 | 2 | 3: return 1
		//    _: return 0
		// }
		"MATCH match EXPRESSION { CASEBLOCK }\n" +
		"CASEBLOCK CASE CASEBLOCK\n" +
		"CASEBLOCK\n" +
		"CASE EXPRESSION ORCASE : STATEMENTBLOCK\n" +
		"ORCASE | EXPRESSION ORCASE\n" +
		"ORCASE\n" +
		"CASE _ : STATEMENTBLOCK\n" +

		// enum FRUIT { APPLE | ORANGE}
		"ENUM enum name { name ENUMLIST }\n" +
		"ENUMLIST | name ENUMLIST\n" +
		"ENUMLIST\n" +

		"STRUCT struct name { STRUCTLIST }\n" +
		"STRUCTLIST type name STRUCTLIST\n" +
		"STRUCTLIST\n" +

		//"INTERFACE interface name { }\n" +

		"TYPE ARRAYTYPE\n" +
		"TYPE SLICETYPE\n" +
		"TYPE MAPTYPE\n" +
		"TYPE bool\n" +
		"TYPE int\n" +
		"TYPE string\n" +
		"TYPE float\n" +
		"TYPE char\n" +
		"TYPE any\n" +
		"TYPE name\n" +

		"ARRAYTYPE TYPE [ EXPRESSION ]\n" +

		"SLICETYPE TYPE [ ]\n" +

		"MAPTYPE map [ TYPE ] TYPE\n" +

		// Operator Presedence:
		// op1: ||
		// op2: &&
		// op3:  == !=
		// op4: < <= > >=
		// op5: + -
		// op6: * / %
		// op7:
		"EXPRESSION EX1 op1 EX2\n" +
		"EXPRESSION EX2\n" +
		"EX2 EX2 op2 EX3\n" +
		"EX2 EX3\n" +
		"EX3 EX3 op3 EX4\n" +
		"EX3 EX4\n" +
		"EX4 EX4 op4 EX5\n" +
		"EX4 EX5\n" +
		"EX5 EX5 op5 EX6\n" +
		"EX5 EX6\n" +
		"EX6 EX6 op6 EX7\n" +
		"EX6 EX7\n" +
		"EX7 LITERAL\n" +
		"EX7 ( EXPRESSION )\n" +
		"EX7 FNCALL\n" +
		"EX7 SLCACCESS\n" +
		"EX7 name\n" +
		"EX7 . name\n" + // id := s.id

		"LITERAL intlit\n" +
		"LITERAL floatlit\n" +
		"LITERAL stringlit\n" +
		"LITERAL charlit\n" +
		"LITERAL boollit\n" +

		"SLCACCESS name [ EXPRESSION ]"

	return s
}

func (g *Grammar) PrintGrammar() {
	grmap := map[string][][]string{}
	grorder := []string{}
	for _, r := range g.rules {
		if len(grmap[r.left]) == 0 {
			grorder = append(grorder, r.left)
		}
		grmap[r.left] = append(grmap[r.left], r.right)
	}
	fmt.Println("Amount Rules:", len(grorder))
	for _, left := range grorder {
		right := ""
		for _, rs := range grmap[left] {
			if len(rs) == 0 {
				right += " e |"
			} else {
				for _, r := range rs {
					right += " " + r
				}
				right += " |"
			}
		}
		fmt.Println(left, "->", right[:len(right)-1])

	}
}
