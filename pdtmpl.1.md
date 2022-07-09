---
title: pdtmpl
section: 1
header: User Manual
footer: pdtmpl 1.0.0
date: July 9, 2022
---

# NAME

pdtmpl - a pandoc preprocessor for JSON and YAML content.

# SYNOPSIS

  pdtmpl [OPTIONS] PANDOC_TEMPLATE_NAME [PANDOC_OPTIONS...]

# DESCRIPTION

pdtmpl s is a Pandoc preprocessor. It reads JSON from standard
input and passes that via a temp file into pandoc where it
applies a pandoc template to render to standard out.

NOTE: If PANDOC_TEMPLATE_NAME is an empty string than the JSON or 
YAML read will be processed without using Pandoc's "--template"
option. You will need to provide some of Pandoc's options after
the empty string placeholder. See example below.

# OPTIONS

**-help**
: display usage

**-license**
: display license

**-version**
: display version

**-i FILENAME**
: read JSON or TOML file

**-o FILENAME**
: write Pandoc output to file

# EXAMPLE

In this example we have a JSON object document called
"example.json" and a Pandoc template called "example.tmpl".
A redirect `"<"` is used to pipe the content of "example.json"
into the command line tool pdtmpl.

```shell
  pdtmpl example.tmpl < example.json
```

Render example.json as Markdown document. We need to use
Pandoc's own options of "-s" (stand alone) and "-t" (to
tell Pandoc the output format)

```shell
  pdtmpl "" -s -t markdown < example.json
```

Process a "codemeta.json" file with "codemeta-md.tmpl" to
produce an about page in Markdown via Pandocs template
processing (the "codemeta-md.tmpl" is a Pandoc template
marked up to produce Markdown output).

```shell
  pdtmpl -i codemeta.json -o about.md \
             codemeta-md.tmpl
```

