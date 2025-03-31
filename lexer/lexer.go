package lexer

import (
	"bufio"
	"log"
	"os"
)

func Lex(path string) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal("Error opening file", path)
	}
	defer f.Close()
	rdr := bufio.NewReader(f)
	for {
		r, _, err := rdr.ReadRune()
		if err != nil {
			break
		}
		
	}
}