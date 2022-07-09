// pdtmpl.go is a preprocessor package for Pandoc. The driving use case
// is passing JSON or YAML documents to Pandoc and applying a Pandoc
// template.
//
// @Author R. S. Doiel, <rsdoiel@gmail.com>
//
// copyright 2022 R. S. Doiel
// All rights reserved.
//
// License under the 3-Clause BSD License
// See https://opensource.org/licenses/BSD-3-Clause
package pdtmpl

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

var verbose bool

// SetVerbose when set true will show the Pandoc command
// envocation before running Pandoc to process the JSON document
// and template. Mainly useful for debugging.
func SetVerbose(onoff bool) {
	verbose = onoff
}

// ReadAll reads JSON from as a stream using an io.Reader.
// Buffers it. Then uses Apply and options return
// a slice of bytes and error value.
//
//```shell
//    // Options passed to Pandoc
//    opt := []string{}
//    out, err := pdtmpl.ReadAll(os.Stdin, "page.tmpl", opt)
//    if err != nil {
//       // ... handle error
//    }
//    fmt.Fprintf(os.Stdout, "%s\n", out)
//```
//
func ReadAll(r io.Reader, template string, options []string) ([]byte, error) {
	// Read the JSON input
	src, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return Apply(src, template, options)
}

// ReadFile reads a JSON or YAML document from a file then uses Apply
// and options returning a slice of bytes and error value.
func ReadFile(name string, template string, options []string) ([]byte, error) {
	// Read the JSON or YAML file
	src, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}
	return Apply(src, template, options)
}

// ApplyIO reads in JSON from an io.Reader, applies the template
// and parameters via Format() writing the result to the io.Writer.
// returns an error value.
//
//```
//  // Options passed to Pandoc
//  opt := []string{}
//  err := pdtmpl.ApplyIO(os.Stdin, os.Stdout, "example.tmpl", opt)
//  if err != nil {
//     // ... handle error
//  }
//```
//
func ApplyIO(r io.Reader, w io.Writer, template string, options []string) error {
	src, err := ReadAll(r, template, options)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "%s\n", src)
	return err
}

// Apply takes a byte array (like you could read from os.Stdin
// containing JSON or YAML. It creates a temp file and passes that to
// Pandoc via `--metadata-file` option along with any additional
// pandoc options provided. Pandoc then renders the output either
// using the template name (if non-empty string) and the
// additional options passed to Pandoc.
//
//```
//  src, err := ioutil.ReadFile("example.json")
//  if err != nil {
//     // ... handle error
//  }
//  // Options passed to Pandoc
//  opt := []string{}
//  src, err := pdtmpl.Format(src, "example.tmpl", opt)
//  if err != nil {
//     // ... handle error
//  }
//  fmt.Printf("%s\n", src)
//```
//
// NOTE: If the template name is an empty string then the
// template option of Pandoc will not be automatically generated.
// This can be helpful when turning JSON into non-HTML formats like
// Markdown or using the default Pandoc templates.
//
// Turning "example.json" into a Markdown document. The "body"
// attribute will become the Markdown document body. The whole
// JSON document becomes the YAML frontmatter of the document.
//
//```
//  src, err := ioutil.ReadFile("example.json")
//  if err != nil {
//     // ... handle error
//  }
//  // Options passed to Pandoc
//  opt := []string{"-t", "markdown"}
//  src, err := pdtmpl.Format(src, "", opt)
//  if err != nil {
//     // ... handle error
//  }
//  fmt.Printf("%s\n", src)
//```
//
func Apply(src []byte, template string, options []string) ([]byte, error) {
	pandoc, err := exec.LookPath("pandoc")
	if err != nil {
		return nil, err
	}
	f, err := os.CreateTemp(".", "pandoc.*.json")
	if err != nil {
		return nil, err
	}
	tmpFile := f.Name()
	defer func() {
		f.Close()
		//Remove tmpFile
		os.RemoveAll(tmpFile)
	}()
	// Write our JSON or YAML to the temp file, check for errors.
	_, err = fmt.Fprintf(f, "%s", src)
	if err != nil {
		return nil, err
	}
	// Download on experience from text link for Monday's
	// virtual inspection. Install app and stop at conference code
	// entry.
	vargs := []string{}
	vargs = append(vargs, "--metadata-file", tmpFile)
	if template != "" {
		vargs = append(vargs, "--template", template)
	}
	if (options != nil) && (len(options) > 0) {
		vargs = append(vargs, options...)
	}
	if verbose {
		fmt.Fprintf(os.Stderr, "%s %s\n", pandoc, strings.Join(vargs, " "))
	}
	cmd := exec.Command(pandoc, vargs...)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	errMsg, _ := ioutil.ReadAll(stderr)
	src, _ = ioutil.ReadAll(stdout)
	if err := cmd.Wait(); err != nil {
		if len(errMsg) > 0 {
			return nil, fmt.Errorf("%s, %s\n", errMsg, err)
		}
		return nil, err
	}
	return src, nil
}
