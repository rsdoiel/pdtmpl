package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/rsdoiel/pdtmpl"
)

var (
	helpText = `%{app_name}(1) skimmer user manual | version {version} {release_hash}
% R. S. Doiel
% {release_date}

# APP

{app_name}

# SYSNOPSIS

{app_name} [OPTIONS] VERB

# DESCRIPTION

{app_name} is a Pandoc preprocessor intended as a simple platform to test various
ideas about how to better support Pandoc generated static websites.

The first use case is to take a JSON or YAML file
turning the data into a YAML block that then is sent through Pandoc using
the a Pandoc template and any additional Pandoc processing options. 

The second use vase is to take a Markdown document, scan it for YAML blogs
containing a "form" object then replacing those blocks with an HTML block
holding a web form. The goal of the usecase is to see if that improves
readability of the original Markdown file also making the aggregated 
HTML generated from Pandoc easier to manage.

# VERB

{app_name} expect a verb to describe the use case be tested. Currently
support verbs are

help
: Display this help page.

tmpl
: Apply the template preprosor for turning raw JSON and YAML into
a Markdown stream sent to Pandoc over standard io.

webform
: This reads and writes to standard io replace any embedded YAML blocks
with a form object with HTML blocks containing a webform defined by the
form object.

# OPTIONS

-help
: display usage

-license
: display license

-version
: display version

-i INPUT
: Is a text file with embedded YAML blocks to be transformed
into HTML blocks (e.g. a Markdown file)

-o OUTPUT
: write Pandoc output to file

# EXAMPLES

In this example we have a JSON object document called
"example.json" and a Pandoc template called "example.tmpl".
A redirect "<" is used to pipe the content of "example.json"
into the command line tool {app_name}.

  {app_name} tmpl example.tmpl < example.json

Render example.json as Markdown document. We need to use
Pandoc's own options of "-s" (stand alone) and "-t" (to
tell Pandoc the output format)

  {app_name} tmpl "" -s -t markdown < example.json

Process a "codemeta.json" file with "codemeta-md.tmpl" to
produce an about page in Markdown via Pandocs template
processing (the "codemeta-md.tmpl" is a Pandoc template
marked up to produce Markdown output).

  {app_name} -i codemeta.json -o about.md tmpl \
             codemeta-md.tmpl


In this example we have a markdown file called "guestbook.md"
is processed and then sent through Pandoc to render as HTML.

~~~shell
{app_name} webform guestbook.md | pandoc -f Markdown -t HTML5 -s \
  >guestbook.html
~~~

`

)

// handleError prints the error message and exists with status code 1
func handleError(eout io.Writer, err error) {
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		os.Exit(1)
	}
}

func main() {
	var (
		showHelp    bool
		showLicense bool
		showVersion bool
		verbose     bool
		input       string
		output      string
		err         error
	)

	appName := path.Base(os.Args[0])
	licenseText := pdtmpl.LicenseText
	version, releaseHash, releaseDate := pdtmpl.Version, pdtmpl.ReleaseHash, pdtmpl.ReleaseDate
	verb, verbs := "help", []string{ "help", "tmpl", "webform" }
	fmtHelp := pdtmpl.FmtHelp
	
	flag.BoolVar(&showHelp, "help", false, "display usage")
	flag.BoolVar(&showVersion, "version", false, "display version")
	flag.BoolVar(&showLicense, "license", false, "display license")
	flag.BoolVar(&verbose, "verbose", false, "show Pandoc envocation")
	flag.StringVar(&input, "i", "", "read JSON or YAML from file")
	flag.StringVar(&output, "o", "", "write Pandoc output to file")
	flag.Parse()

	in := os.Stdin
	out := os.Stdout
	eout := os.Stderr

	if showHelp {
		fmt.Fprintf(out, "%s", fmtHelp(helpText, appName, version, releaseHash, releaseDate))
		os.Exit(0)
	}
	if showVersion {
		fmt.Fprintf(out, "%s %s %s\n", appName, version, releaseHash)
		os.Exit(0)
	}
	if showLicense {
		fmt.Fprintf(out, "%s\n", licenseText)
		os.Exit(0)
	}

	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintf(eout, "error, expected %s, see %s help for details", strings.Join(verbs, ", "), appName)
		os.Exit(1)

	}
	verb, args = args[0], args[1:]
	pdtmpl.SetVerbose(verbose)

	if input != "" && input != "-" {
		in, err = os.Open(input)
		handleError(eout, err)
		defer in.Close()
	}
	if output != "" && output != "-" {
		out, err = os.Create(output)
		handleError(eout, err)
		defer out.Close()
	}

	switch verb {
	case "help":
		fmt.Fprintf(out, "%s", fmtHelp(helpText, appName, version, releaseHash, releaseDate))
		os.Exit(0)
	case "tmpl":
		if len(args) == 0 {
			handleError(eout, fmt.Errorf("missing template name"))
		}
		err = pdtmpl.ApplyIOTemplate(in, out, args[0], args[1:])
		handleError(eout, err)
	case "webform":
		err := pdtmpl.ApplyWebForm(in, out, eout, args)
		handleError(eout, err)
	default:
		fmt.Fprintf(eout, "error, expected %s, see %s help for details", strings.Join(verbs, ", "), appName)
		os.Exit(1)
	}
}
