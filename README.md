# bufpipe: Buffered Pipe

The buffered version of io.Pipe. It's safe for concurrent use.

## How different from io.Pipe?

Writes never block because the pipe has variable-sized buffer.

```Go
r, w := bufpipe.New(nil)
io.WriteString(w, "abc") // No blocking.
io.WriteString(w, "def") // No blocking, too.
w.Close()
io.Copy(os.Stdout, r)
// Output: abcdef
```

## How different from bytes.Buffer?

Reads block if the internal buffer is empty until the writer is closed.

```Go
r, w := bufpipe.New(nil)
go io.Copy(os.Stdout, r) // The reads block until the writer is closed.

io.WriteString(w, "abc")
io.WriteString(w, "def")
w.Close()
// Output: abcdef
```
