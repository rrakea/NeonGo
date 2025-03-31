package main

import (
	"log"
	"neon/config"
	"neon/lexer"
	"os"

	flag "github.com/spf13/pflag"
)

func main() {
	output := flag.StringP("output", "o", "out", "Specify the name of the output file")
	release := flag.BoolP("release", "r", false, "Release build")
	file := flag.StringP("file", "f", "", "Specify file to compile")
	only_lex := flag.Bool("only_lexer", false, "Stop after lexing")
	only_typecheck := flag.Bool("only_typecheck", false, "Stop after typechecking")
	only_parse := flag.Bool("only_parse", false, "Stop after parsing")

	flag.Parse()
	config.File_name = *output
	config.Only_lex = *only_lex
	config.Only_parse = *only_parse
	config.Only_typecheck = *only_typecheck
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
	for i := range len(files) {
		tokens = append(tokens, []*lexer.Token{})
		pull := <-token_chan
		tokens[i] = append(tokens[i], pull...)
	}
	if config.Only_lex {
		for _, t := range tokens {
			lexer.Print_tokens(t)
		}
		return
	}
}
