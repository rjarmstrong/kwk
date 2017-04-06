package main

import (
	"fmt"
	"os"
	"strings"
	"bytes"
)

var list = []string{
	"apple", "banana", "orange", "pumpkin", "dog", "cat",
}

const (
	esc        = "\033["
	save_pos   = esc + "s"
	return_pos = esc + "u"
)

func main() {
	term := ""
	fmt.Fprint(os.Stdout, search(term))
	moveUp(len(list) + 1)
	//var in string
	for {
		var next string
		l, err := fmt.Scan(&next)
		fmt.Println(l, next)

		if err != nil {
			panic(err)
		}
		//res := search(in)
		//fmt.Print(res)
		//moveUp(len(strings.Split(res, "\n")) - 1)
	}
}

func moveUp(q int) {
	fmt.Fprintf(os.Stdout, "%s%dA", esc, q)
}

func search(term string) string {
	fmt.Println()
	buf := &bytes.Buffer{}
	for _, v := range list {
		if strings.HasPrefix(v, term) || term == "" {
			buf.WriteString(fmt.Sprintf("-> %s\n", v))
		}
	}
	return buf.String()
}
