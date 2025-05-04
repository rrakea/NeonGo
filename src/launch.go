package main

import (
	"log"
	"os"

	"neon/config"
	"neon/lexer"

	flag "github.com/spf13/pflag"
)

func main() {
	output := flag.StringP("output", "o", "out", "Specify the name of the output file")
	release := flag.BoolP("release", "r", false, "Release build")
	file := flag.StringP("file", "f", "", "Specify file to compile")
	lex := flag.Bool("lexer", false, "Print lex output")
	parse := flag.Bool(" parse", false, "Print parse output")
	typecheck := flag.Bool("typecheck", false, "Print typecheck output")

	flag.Parse()
	config.File_name = *output
	config.print_lex = *lex
	config.print_parse = *parse
	config.print_typecheck = *typecheck
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

	if config.Only_lex {
		lexer.Print_tokens(tokens[0])
	}
}
