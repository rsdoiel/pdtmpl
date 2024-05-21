%pdtmpl(1) skimmer user manual | version 0.0.2 2024-05-21
% R. S. Doiel
% c2a99e3

# APP

pdtmpl

# SYSNOPSIS

pdtmpl [OPTIONS] VERB

# DESCRIPTION

pdtmpl is a Pandoc preprocessor intended as a simple platform to test various
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

pdtmpl expect a verb to describe the use case be tested. Currently
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
into the command line tool pdtmpl.

  pdtmpl tmpl example.tmpl < example.json

Render example.json as Markdown document. We need to use
Pandoc's own options of "-s" (stand alone) and "-t" (to
tell Pandoc the output format)

  pdtmpl tmpl "" -s -t markdown < example.json

Process a "codemeta.json" file with "codemeta-md.tmpl" to
produce an about page in Markdown via Pandocs template
processing (the "codemeta-md.tmpl" is a Pandoc template
marked up to produce Markdown output).

  pdtmpl -i codemeta.json -o about.md tmpl \
             codemeta-md.tmpl


In this example we have a markdown file called "guestbook.md"
is processed and then sent through Pandoc to render as HTML.

~~~shell
pdtmpl webform guestbook.md | pandoc -f Markdown -t HTML5 -s \
  >guestbook.html
~~~

