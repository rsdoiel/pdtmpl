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
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	// 3rd Party libraries
	"gopkg.in/yaml.v3"
)

var verbose bool

// SetVerbose when set true will show the Pandoc command
// envocation before running Pandoc to process the JSON document
// and template. Mainly useful for debugging.
func SetVerbose(onoff bool) {
	verbose = onoff
}

//
// Pandoc preprosor for JSON and YAML experiment
//

// ReadAllTemplate reads JSON from as a stream using an io.Reader.
// Buffers it. Then uses Apply and options return
// a slice of bytes and error value.
//
//```shell
//    // Options passed to Pandoc
//    opt := []string{}
//    out, err := pdtmpl.ReadAllTemplate(os.Stdin, "page.tmpl", opt)
//    if err != nil {
//       // ... handle error
//    }
//    fmt.Fprintf(os.Stdout, "%s\n", out)
//```
//
func ReadAllTemplate(r io.Reader, template string, options []string) ([]byte, error) {
	// Read the JSON input
	src, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return ApplyTemplate(src, template, options)
}

// ReadFileTemplate reads a JSON or YAML document from a file then uses Apply
// and options returning a slice of bytes and error value.
func ReadFileTemplate(name string, template string, options []string) ([]byte, error) {
	// Read the JSON or YAML file
	src, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}
	return ApplyTemplate(src, template, options)
}

// ApplyIOTemplate reads in JSON from an io.Reader, applies the template
// and parameters via ApplyTemplate() writing the result to the io.Writer.
// returns an error value.
//
//```
//  // Options passed to Pandoc
//  opt := []string{}
//  err := pdtmpl.ApplyIOTemplate(os.Stdin, os.Stdout, "example.tmpl", opt)
//  if err != nil {
//     // ... handle error
//  }
//```
//
func ApplyIOTemplate(r io.Reader, w io.Writer, template string, options []string) error {
	src, err := ReadAllTemplate(r, template, options)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "%s\n", src)
	return err
}

// ApplyTemplate takes a byte array (like you could read from os.Stdin
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
//  src, err := pdtmpl.ApplyTemplate(src, "example.tmpl", opt)
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
//  // Options passed to Pandoc, use "-s" to trigger writing
//  // the document via Pandoc's default templates.
//  opt := []string{"-s", "-t", "markdown"}
//  src, err := pdtmpl.ApplyTemplate(src, "", opt)
//  if err != nil {
//     // ... handle error
//  }
//  fmt.Printf("%s\n", src)
//```
//
func ApplyTemplate(src []byte, template string, options []string) ([]byte, error) {
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

//
// WebForm experiment
//

// MkWebForm takes a map[string]interface{} and translates the form structure into HTML
func MkWebForm(out io.Writer, eout io.Writer, m map[string]interface{}) error {
	var (
		eType string
		eId string
	)
	fmt.Fprintf(out, "<form")
	for _, k := range []string{ "id", "name", "action", "method", "encoding" } {
		if val, ok := m[k]; ok {
			fmt.Fprintf(out, " %s=%q", k, val)
		}
	}
	fmt.Fprintf(out, " >\n")
	if l, ok := m["elements"].([]interface{}); ok {
		for _, item := range l {
			elem := item.(map[string]interface{})
			if eId, ok = elem["id"].(string); ! ok {
				eId = ""
			}
			if eType, ok = elem["type"].(string); !ok {
				eType = "text"
			}
			// Map the button alias back to the text input type.
			if eType == "button" {
				eType = "text"
			}
			prefix := "input"
			if eType == "select" || eType == "textarea" {
				prefix = eType
			}
			if label, ok := elem["label"].(string); ok {
				label = ""
				fmt.Fprintf(out, "\t<label for=%q>%s</label><%s type=%q", eId, label, prefix, eType)
			} else {
				fmt.Fprintf(out, "\t<%s type=%q", prefix, eType)
			}
			if eId != "" {
				fmt.Fprintf(out, " id=%q", eId)
			}
			for _, k := range []string{"class", "name", "value", "required", "placeholdertext", "title", "required", "pattern" } {
				if val, ok := elem[k].(string); ok {
					fmt.Fprintf(out, " %s=%q", k, val)
				}
			}
			// Handle special case of select element type, handle the submap of options' value and labels
			// Handle case of select and textarea otherwise end form element
			switch eType {
			case "select":
				fmt.Fprintf(out, " >\n")
				if l, ok := elem["options"]; ok {
					option := l.(map[string]interface{})
					for val, label := range option {
						fmt.Fprintf(out, "\t\t<option value=%q>%s</option\n", val, label)
					}
				}
				fmt.Fprintf(out, "\t</select>\n")
			case "textarea":
				fmt.Fprintf(out, " ></textarea>\n")
			default:
				fmt.Fprintf(out, " >\n")
			}
		}
	}
   	fmt.Fprintf(out, "</form>\n")	
	return nil
}

// ApplyWebForm reads Markdown present as input and converts the YAML blogs
// with a form object into HTML blocks with a webform in them.
//
//```shell
//    // Data is read from standard input and written to standard out.
//    opt := []string{}
//    if err := pdtmpl.ApplyIOWebForm(os.Stdin, os.Stdout, os.Stderr, opt); err != nil {
//       // ... handle error
//    }
//```
//
func ApplyWebForm(in io.Reader, out io.Writer, eout io.Writer, options []string) error {
	scanner := bufio.NewScanner(in)
	inYaml, inCodeBlock := false, false
	ymlText := []string{}
	lineNo := 0
	eCnt := 0
	for scanner.Scan() {
		lineNo++
		line := scanner.Text()
		if strings.HasPrefix(line, "~~~") {
			inCodeBlock = !inCodeBlock
		}
		if inCodeBlock {
			fmt.Fprintf(out, "%s\n", line)
		} else {
			if line == "---" {
				if inYaml {
					// Gather Yaml, decode it and see if we have a "form" object
					// If not write it out and continue
					m := map[string]interface{}{}
					txt := strings.Join(ymlText, "\n")
					if err := yaml.Unmarshal([]byte(txt), &m); err != nil {
						fmt.Fprintf(eout, "line %d: %s\n", lineNo, err)
						eCnt++
					}
					form, ok := m["form"].(map[string]interface{});
					if ! ok {
						// This isn't a form block.
						fmt.Fprintf(out, "---\n%s\n---\n", txt)
					} else {
						if err := MkWebForm(out, eout, form); err != nil {
							fmt.Fprintf(eout, "line %d: %s\n", lineNo, err)
							eCnt++
						}
					}
					ymlText = []string{}
				}
				inYaml = !inYaml
			} else if inYaml {
				ymlText = append(ymlText, line)
			} else {
				fmt.Fprintf(out, "%s\n", line)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	if eCnt > 0 {
		return fmt.Errorf("%d errors encountered rending output\n", eCnt)
	}
	return nil
}

