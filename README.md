pdtmpl
======

A light weight pre-processor Pandoc. Target use is JSON documents 
processed via Pandoc templates.

This is a Go package which makes it easy to extend your Go application
to incorporate Pandoc.

cli
---

The command line implementation is a proof of concept writing a
Pandoc pre-processor in Go. An example usage would be to process
[example.json](example.json) JSON document using a Pandoc template
called [example.tpml](example.tmpl).

```shell
    pdtmpl example.tmpl -t html5 < example.json > example.html
```

Go package
----------

Here's some simple use examples of the three functions supplied
in the pdtmpl package.

Given a JSON Object document  as a slice of bytes render formated
output based on the Pandoc template `example.tmpl`

```go
    src, err := ioutil.ReadFile("example.json")
    if err != nil {
        // ... handle error
    }
    // options passed to Pandoc
    opt := []string{"-t", "html5"}
    src, err = pdtmpl.Format(src, "example.tmpl", opt)
    if err != nil {
        // ... handle error
    }
    fmt.Fprintf(os.Stdout, "%s", src)
```

Using an `io.Reader` to retrieve the JSON content, process with the
`example.tmpl` template and write standard output

```go
    f, err := Open("example.json")
    if err != nil {
        // ... handle error
    }
    defer f.Close()
    // options passed to Pandoc
    opt := []string{"-t", "html5"}
    src, err := pdtmpl.ReadAll(f, "example.tmpl", opt)
    if err != nil {
        // ... handle error
    }
    fmt.Fprintf(os.Stdout, "%s", src)
```

Using an `io.Reader` and `io.Writer` to read JSON source from standard
input and write the processed pandoc templated standard output.

```go
    // options passed to Pandoc
    opt := []string{"-t", "html5"}
    err := pdtmpl.ApplyTemplate(os.Stdin, os.Stdout, "example.tmpl", opt)
    if err != nil {
        // ... handle error
    }
```


