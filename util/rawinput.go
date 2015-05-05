package util

import (
	"bufio"
	"os"
	"strings"
)

func RawInput() string {
	in := bufio.NewReader(os.Stdin)
	line, _ := in.ReadString('\n')
	line = strings.Trim(line, "\n\r")
	return line
}
