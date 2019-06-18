package bufpipe_test

import (
	"io"
	"os"

	"github.com/acomagu/bufpipe"
)

func Example() {
	r, w := bufpipe.New(nil)

	done := make(chan struct{})
	go func() {
		io.Copy(os.Stdout, r)
		done <- struct{}{}
	}()

	io.WriteString(w, "abc")
	io.WriteString(w, "def")
	w.Close()
	<-done
	// Output: abcdef
}
