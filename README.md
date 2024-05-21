
[![Project Status: Concept – Minimal or no implementation has been done yet, or the repository is only intended to be a limited example, demo, or proof-of-concept.](https://www.repostatus.org/badges/latest/concept.svg)](https://www.repostatus.org/#concept)

pdtmpl
======

A light weight pre-processor for Pandoc. Target use case is JSON object
documents rendered via Pandoc templates. Target use case is Markdown
document with embedded "form" YAML elements.

This is a proof-of-concept Go package which makes it easy to extend
your Go application to incorporate Pandoc template processing. It includes
the Go package and an example command line tool (cli).

cli
---

The command line implementation is a proof of concept Pandoc
pre-processor in Go. An example usage would be to process
[example.json](example.json) JSON document using a Pandoc template
called [example.tpml](example.tmpl).

~~~shell
    pdtmpl tmpl example.tmpl < example.json > example.html
~~~

Render a Markdown file with an embedded YAML block describing
a webform via Pandoc.

~~~shell
    pdtmpl webform < document.md | pandoc -f Markdown -t html5 -s
~~~

The Markdown file "document.md" would have a simple web form
described via YAML. This oleviates the need to embed HTML. Here's
an example of embedding a "search" form in a "document.md"

~~~markdown

# A search page ...

This is followed by my form.

---
form:
  id: search
  action: /search
  method:  POST
  elements:
    - id: search
      type: search
      placeholdertext: enter a search string here
    - id: submit
      name: submit
      type: submit
      value: Search
---

Other text could follow the YAML block.

~~~


When pltmpl webform processes this turns the form object into HTML markup.
form describes the outer HTML form element who's attributes also map
(e.g. aciton, method). Additionally the form as
an "elements" attribute which is a list of HTML elements
that are children of the form element.  The YAML notations
provides a cleaner expression of the webform by avoiding
the HTML markdown but clear maps from the HTML markup.

Go package
----------

Here's some simple use examples of the three functions supplied
in the pdtmpl package.

Given a JSON Object document  as a slice of bytes render formated
output based on the Pandoc template `example.tmpl`

~~~go
    src, err := ioutil.ReadFile("example.json")
    if err != nil {
        // ... handle error
    }
    // options passed to Pandoc
    opt := []string{}
    src, err = pdtmpl.Apply(src, "example.tmpl", opt)
    if err != nil {
        // ... handle error
    }
    fmt.Fprintf(os.Stdout, "%s", src)
~~~

Using an `io.Reader` to retrieve the JSON content, process with the
`example.tmpl` template and write standard output

~~~go
    f, err := Open("example.json")
    if err != nil {
        // ... handle error
    }
    defer f.Close()
    // options passed to Pandoc
    opt := []string{}
    src, err := pdtmpl.ReadAll(f, "example.tmpl", opt)
    if err != nil {
        // ... handle error
    }
    fmt.Fprintf(os.Stdout, "%s", src)
~~~

Using an `io.Reader` and `io.Writer` to read JSON source from standard
input and write the processed pandoc templated standard output.

~~~go
    // options passed to Pandoc
    opt := []string{}
    err := pdtmpl.ApplyIO(os.Stdin, os.Stdout, "example.tmpl", opt)
    if err != nil {
        // ... handle error
    }
~~~

Requirements
------------

- [Pandoc](https://pandoc.org) 2.18 or better
- [Go](https://golang.org) 1.22 or better to compile from source
- [GNU Make](https://www.gnu.org/software/make/) (optional) to automated compilation
- [Git](https://git-scm.com/) or other Git client to retrieve this repository

Installation
------------

1. Clone https://github.com/rsdoiel/pdtmpl to your local machine
2. Change directory into the git repository (i.e. `pdtmpl`
3. Compile using `make`
4. Install using `make install`

~~~shell
    git clone https://github.com/rsdoiel/pdtmpl
    cd pdtmpl
    make
    make install
~~~

NOTE: This recipe assumes' you are familar with setting up a
Go development environment (e.g. You've set GOPATH environment
appropriately). See the [go website](https://golang.org) for
details about setting up and compiler programs.

License
-------

AGPL version 3 or later

