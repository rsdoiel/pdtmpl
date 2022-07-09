// pdtmpl.go is a preprocessor package for Pandoc's template engine.
// The general idea is you have some JSON and you want to apply a
// Pandoc template to it. The JSON is written to a temp file and passed
// to Pandoc via `--metadata-file`. The JSON is expected to represent an
// object where attributes map to variables in the template.
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
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

// ReadAll is a proof of concept pre-processor for Pandoc.
// It takes a byte array (like you could read from os.Stdin),
// containing JSON, applies a pandoc template by envoking
// pandoc along with any additional pandoc options you provide.
// The temp file is removed on success.
//
//```shell
//    out, err := pdtmpl.ReadAll(os.Stdin, "page.tmpl", nil)
//    if err != nil {
//       // ... handle error
//    }
//    fmt.Fprintf(os.Stdout, "%s\n", out)
//```
//
func ReadAll(r io.Reader, template string, args []string) ([]byte, error) {
	// Read the JSON input
	jsonObjectSrc, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	obj := map[string]*interface{}{}
	if err := json.Unmarshal(jsonObjectSrc, &obj); err != nil {
		return nil, err
	}
	return Format(jsonObjectSrc, template, args)
}

// ApplyTemplate reads in JSON from an io.Reader, applies the template
// and parameters via Format() writing the result to the io.Writer.
// returns an error value.
//
//```
//  args := os.Args[1:]
//  err := pdtmpl.ApplyTemplate(os.Stdin, os.Stdout, "example.tmpl", args)
//  if err != nil {
//     // ... handle error
//  }
//```
//
func ApplyTemplate(r io.Reader, w io.Writer, template string, args []string) error {
	src, err := ReadAll(r, template, args)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "%s\n", src)
	return err
}

// Format takes a byte array (like you could read from os.Stdin0
// containing JSON. It creates a temp file and passes that to
// Pandoc via `--metadata-file` option along with any additional
// pandoc options provided. Pandoc then renders the template
// and that is returned as a byte array and error.
//
//```
//  src, err := ioutil.ReadFile("example.json")
//  if err != nil {
//     // ... handle error
//  }
//  src, err := pdtmpl.Format(src, "example.tmpl", nil)
//  if err != nil {
//     // ... handle error
//  }
//  fmt.Printf("%s\n", src)
//```
//
func Format(src []byte, template string, args []string) ([]byte, error) {
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
	// Write our JSON to the temp file, check for errors.
	_, err = fmt.Fprintf(f, "%s", src)
	if err != nil {
		return nil, err
	}
	// Download on experience from text link for Monday's
	// virtual inspection. Install app and stop at conference code
	// entry.
	params := []string{}
	params = append(params, "--metadata-file", tmpFile)
	if template != "" {
		params = append(params, "--template", template)
	}
	if (args != nil) && (len(args) > 0) {
		params = append(params, args...)
	}
	cmd := exec.Command(pandoc, params...)
	return cmd.Output()
}
