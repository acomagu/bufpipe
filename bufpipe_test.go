package bufpipe_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"sort"
	"testing"
	"time"

	"github.com/acomagu/bufpipe"
	"github.com/matryer/is"
)

func TestPipeWriter_NoBlocking(t *testing.T) {
	is := is.New(t)

	r, w := bufpipe.New(nil)
	io.WriteString(w, "abc")
	io.WriteString(w, "def")
	w.Close()

	b, err := ioutil.ReadAll(r)
	is.NoErr(err)
	is.Equal(b, []byte("abcdef"))
}

func TestMultiBlocking(t *testing.T) {
	is := is.New(t)

	results := make(chan []byte)
	block := func(r io.Reader) {
		b := make([]byte, 3)
		n, err := r.Read(b)
		is.NoErr(err)
		results <- b[:n]
	}

	r, w := bufpipe.New(nil)
	go block(r)
	go block(r)
	go block(r)

	time.Sleep(time.Millisecond) // Ensure blocking.

	data := []string{"abc", "def", "ghi"}
	for _, s := range data {
		n, err := w.Write([]byte(s))
		is.NoErr(err)
		is.Equal(n, 3)
	}

	var ss []string
	for i := 0; i < 3; i++ {
		ss = append(ss, string(<-results))
	}
	sort.Strings(ss)
	is.Equal(ss, data)
}

func TestPipeWriter_Close(t *testing.T) {
	is := is.New(t)

	r, w := bufpipe.New([]byte("abc"))
	n, err := w.Write([]byte("def"))
	is.NoErr(err)
	is.Equal(n, 3)

	is.NoErr(w.Close())

	buf := make([]byte, 3)
	n, err = r.Read(buf)
	is.NoErr(err)
	is.Equal(buf[:n], []byte("abc"))

	n, err = r.Read(buf)
	is.NoErr(err)
	is.Equal(buf[:n], []byte("def"))

	_, err = r.Read(buf)
	is.Equal(err, io.EOF)
}

func TestPipeWriter_CloseWithError(t *testing.T) {
	is := is.New(t)

	r, w := bufpipe.New([]byte("abc"))
	n, err := w.Write([]byte("def"))
	is.NoErr(err)
	is.Equal(n, 3)

	expect := fmt.Errorf("original error")
	is.NoErr(w.CloseWithError(expect))

	buf := make([]byte, 3)
	n, err = r.Read(buf)
	is.NoErr(err)
	is.Equal(buf[:n], []byte("abc"))

	n, err = r.Read(buf)
	is.NoErr(err)
	is.Equal(buf[:n], []byte("def"))

	_, err = r.Read(buf)
	is.Equal(err, expect)
}

func TestPipeReader_Close(t *testing.T) {
	is := is.New(t)

	r, w := bufpipe.New([]byte("abc"))
	is.NoErr(r.Close())

	n, err := w.Write([]byte("abc"))
	is.Equal(err, bufpipe.ErrClosedPipe)
	is.Equal(n, 0)
}

func TestPipeReader_CloseWithError(t *testing.T) {
	is := is.New(t)

	expect := fmt.Errorf("original error")

	r, w := bufpipe.New([]byte("abc"))
	is.NoErr(r.CloseWithError(expect))

	n, err := w.Write([]byte("abc"))
	is.Equal(err, expect)
	is.Equal(n, 0)
}

func TestPipeReader_WriterCloseNoDeadlock(t *testing.T) {
	r, w := bufpipe.New(nil)

	done := make(chan struct{})
	go func(t *testing.T) {
		buf := make([]byte, 800)
		r.Read(buf)
		done <- struct{}{}
	}(t)

	time.Sleep(300 * time.Millisecond)
	w.Close()

	<-done
}

