package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/rsdoiel/pdtmpl"
)

func usage(appName string) string {
	appName = path.Base(appName)
	return fmt.Sprintf(`
USAGE:

    %s [OPTIONS] PANDOC_TEMPLATE_NAME [PANDOC_OPTIONS...]

%s is a Pandoc preprocessor. It reads JSON from standard
input and passes that via a temp file into pandoc where it
applies a pandoc template to render to standard out.

EXAMPLE

    cat example.json | %s example.tmpl

`, appName, appName, appName)
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
	showHelp := false
	flag.BoolVar(&showHelp, "help", false, "display usage")
	flag.Parse()
	if showHelp {
		fmt.Print(usage(os.Args[0]))
		os.Exit(0)
	}

	args := flag.Args()
	template, args, err := handleArgs(os.Args[0], args)
	handleError(err)
	err = pdtmpl.ApplyTemplate(os.Stdin, os.Stdout, template, args)
	handleError(err)
}
