package lexer

type GoLine []byte

func (*GoLine) AppendSemicolon(line []byte) (newLine byte) {
	newLine = append(line, byte(';'))
	return newLine
}

func main() {
}
