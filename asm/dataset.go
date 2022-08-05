package asm

import (
	"bufio"
	"os"

	"github.com/alexbezu/gobol/pl"
)

type DCB struct {
	MACRF   string
	DDNAME  string
	DSORG   string
	f       *os.File
	scanner *bufio.Scanner
	writer  *bufio.Writer
}

func OPEN(file *DCB, mode string) {
	switch mode {
	case "INPUT":
		f, err := os.Open(file.DDNAME)
		if err != nil {
			return
		}
		file.f = f
		file.scanner = bufio.NewScanner(file.f)
	case "OUTPUT":
		f, err := os.Create(file.DDNAME)
		if err != nil {
			return
		}
		file.f = f
		file.writer = bufio.NewWriter(file.f)
	}
}

func GET(file *DCB, D1 pl.Objer) bool {
	notend := file.scanner.Scan()
	D1.Set(file.scanner.Text())
	return !notend
}

func PUT(file *DCB, D1 pl.Objer) {
	_, err := file.writer.WriteString(D1.String() + "\n")
	if err != nil {
		return
	}
}

func CLOSE(file *DCB) {
	if file.writer != nil {
		file.writer.Flush()
	}
	file.f.Close()
}
