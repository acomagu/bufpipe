package bufpipe_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/acomagu/bufpipe"
)

func Example1() {
	r, w := bufpipe.New(nil)
	io.WriteString(w, "abc") // Don't blocks.
	io.WriteString(w, "def") // Don't blocks, too.
	w.Close()
	io.Copy(os.Stdout, r)
	// Output: abcdef
}

func Example2() {
	r, w := bufpipe.New(nil)
	result := make(chan []byte)
	go func() {
		s, _ := ioutil.ReadAll(r) // Blocks until the write half is closed.
		result <- s
	}()

	io.WriteString(w, "abc")
	io.WriteString(w, "def")
	w.Close()

	fmt.Println(string(<-result))
	// Output: abcdef
}
