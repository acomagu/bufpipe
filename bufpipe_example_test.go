package bufpipe_test

import (
	"io"
	"os"

	"github.com/acomagu/bufpipe"
)

func Example() {
	r, w := bufpipe.New(nil)
	go io.Copy(os.Stdout, r) // => abcdef.

	io.WriteString(w, "abc")
	io.WriteString(w, "def")
	w.Close()
}
