package providers

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/fatih/color"
)

func TestStdout(t *testing.T) {
	var w io.Writer
	w = os.Stdout

	msg := color.GreenString("hello\n")
	fmt.Printf("%x\n", []byte(msg))
	if w == os.Stdout {
		w.Write([]byte(msg))
	}
	fmt.Printf("%v", msg)
}
