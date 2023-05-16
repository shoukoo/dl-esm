# dl-esm
Download ESM modules from npm and jsdelivr

This is the Go version of Simon Willison's [download-esm](https://github.com/simonw/download-esm) work

## Installation

Install this tool using `go`:
```bash
    go install github.com/shoukoo/dl-esm@latest
```

## Usage

To download Solid JS and all of its dependencies as ECMAScript modules:

    dl-esm solid-js/html

This will download around 3 `.js` files to the current directory.

To put them in another directory, add that as an argument:

    dl-esm solid-js/html /tmp

To download a specific version of solid js

    dl-esm solid-js@1.7.5/html

Each file will have any `import` and `export` statements rewritten as relative paths.

You can then use the library in your own HTML and JavaScript something like this:

```html
...
<script type="importmap">
{
  "imports": {
    "solid-js": "/static/solid-js-1-7-5.js",
    "solid-js/html": "/static/solid-js-1-7-5-html.js",
    "solid-js/web": "/static/solid-js-1-7-5-web.js"
  }
}
</script>
<script type="module">
import html from "solid-js/html";
import { render, hydrate } from "solid-js/web";
import { onCleanup, createEffect, createSignal, onMount, For, Show } from "solid-js";
...
```

## Development

To run dl-esm locally run:
```bash
    go run main.go solid-js@1.7.5
```

To run the tests:
```bash
    go test ./... 
```
