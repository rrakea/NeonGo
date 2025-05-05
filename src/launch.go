package src

import (
	"log"
	"os"

	flag "github.com/spf13/pflag"

	"neon/src/frontend/config"
	"neon/src/frontend/lexer"
	"neon/src/frontend/parser"
)

func Launch() {
	output := flag.StringP("output", "o", "out", "Specify the name of the output file")
	release := flag.BoolP("release", "r", false, "Release build")
	file := flag.StringP("file", "f", "", "Specify file to compile")
	lex := flag.Bool("lex", false, "Print lex output")
	parse := flag.Bool(" parse", false, "Print parse output")
	typecheck := flag.Bool("typecheck", false, "Print typecheck output")
	rules := flag.Bool("rules", false, "Print grammar rules")

	flag.Parse()
	config.File_name = *output
	config.Print_lex = *lex
	config.Print_parse = *parse
	config.Print_typecheck = *typecheck
	config.Release_build = *release
	/*
		currentDir, err := os.Getwd()
		if err != nil {
			log.Fatal("Could not get working directory", err.Msg(f))
		}
	*/
	if *file == "" {
		log.Fatal("Not implemented yet")
	}

	files := []string{}
	files = append(files, *file)
	token_chan := make(chan []*lexer.Token)
	tokens := [][]*lexer.Token{}
	for _, f := range files {
		open_file, err := os.Open(f)
		if err != nil {
			log.Fatal("Couldnt open file\n", err)
		}
		go lexer.Lex(open_file, token_chan)
	}

	// Lex the files
	for i := range len(files) {
		tokens = append(tokens, []*lexer.Token{})
		pull := <-token_chan
		tokens[i] = append(tokens[i], pull...)
	}

	if config.Print_lex {
		lexer.Print_tokens(tokens[0])
	}

	if *rules {
		parser.Get_Grammar().PrintGrammar()
	}
}
