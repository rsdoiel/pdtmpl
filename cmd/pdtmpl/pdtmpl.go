package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/rsdoiel/pdtmpl"
)

func version(appName string) string {
	return fmt.Sprintf("%s %s\n", path.Base(appName), pdtmpl.Version)
}

func license(appName string) string {
	appName = path.Base(appName)
	src := `
BSD 3-Clause License

Copyright (c) 2022, R. S. Doiel
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this
   list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its
   contributors may be used to endorse or promote products derived from
   this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
`
	return strings.ReplaceAll(src, "{app_name}", appName)
}

func usage(appName string) string {
	appName = path.Base(appName)
	return strings.ReplaceAll(`
USAGE:

    {app_name} [OPTIONS] PANDOC_TEMPLATE_NAME [PANDOC_OPTIONS...]

{app_name} s is a Pandoc preprocessor. It reads JSON from standard
input and passes that via a temp file into pandoc where it
applies a pandoc template to render to standard out.

OPTIONS

    -help      display usage
	-license   display license

EXAMPLE

In this example we have a JSON object document called
"example.json" and a Pandoc template called "example.tmpl".
A redirect "<" is used to pipe the content of "example.json"
into the command line tool {app_name}.

    {app_name} example.tmpl < example.json

`, "{app_name}", appName)
}

func handleArgs(appName string, args []string) (string, []string, error) {
	var (
		template string
		params   []string
	)
	if len(args) == 1 {
		template, params = args[0], []string{}
	} else if len(args) > 1 {
		template, params = args[0], args[1:]
	} else {
		return "", []string{}, fmt.Errorf(`%s`, usage(appName))
	}
	return template, params, nil
}

func handleError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func main() {
	var (
		showHelp    bool
		showLicense bool
		showVersion bool
	)
	flag.BoolVar(&showHelp, "help", false, "display usage")
	flag.BoolVar(&showVersion, "version", false, "display version")
	flag.BoolVar(&showLicense, "license", false, "display license")
	flag.Parse()
	if showHelp {
		fmt.Print(usage(os.Args[0]))
		os.Exit(0)
	}
	if showVersion {
		fmt.Print(version(os.Args[0]))
		os.Exit(0)
	}
	if showLicense {
		fmt.Print(license(os.Args[0]))
		os.Exit(0)
	}

	args := flag.Args()
	template, args, err := handleArgs(os.Args[0], args)
	handleError(err)
	err = pdtmpl.ApplyTemplate(os.Stdin, os.Stdout, template, args)
	handleError(err)
}
