package main

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/h2non/gock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestSolidJSEsm(t *testing.T) {
	defer gock.Disable()
	htmlResponse := `
/**
 * Bundled by jsDelivr using Rollup v2.79.1 and Terser v5.17.1.
 * Original file: /npm/solid-js@1.7.5/html/dist/html.js
 *
 * Do NOT use SRI with dynamically generated files! More information: https://www.jsdelivr.com/using-sri-with-dynamic-files
 */
import{ascending as t,descending as n,least as $r}from"/npm/solid-js@1.7.5/web/+esm"
//# sourceMappingURL=/sm/01413fe1f7a8c2da69e83bcc6c3f16a63658e3f27f5a8ffeda7da895d71e4aa2.map
`

	webResponse := `
/**
 * Bundled by jsDelivr using Rollup v2.79.1 and Terser v5.17.1.
 * Original file: /npm/solid-js@1.7.5/web/dist/web.js
 *
 * Do NOT use SRI with dynamically generated files! More information: https://www.jsdelivr.com/using-sri-with-dynamic-files
 */
import{ascending as t,descending as n,least as $r}from"/npm/solid-js@1.7.5/+esm"
//# sourceMappingURL=/sm/01413fe1f7a8c2da69e83bcc6c3f16a63658e3f27f5a8ffeda7da895d71e4aa2.map
`

	solidjsResponse := `
/**
 * Bundled by jsDelivr using Rollup v2.79.1 and Terser v5.17.1.
 * Original file: /npm/solid-js@1.7.5/dist/solid.js
 *
 * Do NOT use SRI with dynamically generated files! More information: https://www.jsdelivr.com/using-sri-with-dynamic-files
 */
export default null;
//# sourceMappingURL=/sm/01413fe1f7a8c2da69e83bcc6c3f16a63658e3f27f5a8ffeda7da895d71e4aa2.map
`
	tmpDir, err := os.MkdirTemp("", "*")
	assert.NoError(t, err)
	gock.New("https://cdn.jsdelivr.net").
		Get("/npm/solid-js@1.7.5/html/+esm").
		Reply(200).
		BodyString(htmlResponse)

	gock.New("https://cdn.jsdelivr.net").
		Get("/npm/solid-js@1.7.5/web/+esm").
		Reply(200).
		BodyString(webResponse)

	gock.New("https://cdn.jsdelivr.net").
		Get("/npm/solid-js@1.7.5/+esm").
		Reply(200).
		BodyString(solidjsResponse)

	_, _, err = executeCmd(rootCmd(), "solid-js@1.7.5/html", tmpDir)
	assert.NoError(t, err)
	assert.FileExists(t, fmt.Sprintf("%s/solid-js-1-7-5-html.js", tmpDir))
	assert.FileExists(t, fmt.Sprintf("%s/solid-js-1-7-5-web.js", tmpDir))
	assert.FileExists(t, fmt.Sprintf("%s/solid-js-1-7-5.js", tmpDir))

	// check the code was rewritten

	b, _ := os.ReadFile(fmt.Sprintf("%s/solid-js-1-7-5-html.js", tmpDir))
	assert.Contains(t, string(b), "import {ascending as t,descending as n,least as $r} from \"./solid-js-1-7-5-web.js\"")

	b, _ = os.ReadFile(fmt.Sprintf("%s/solid-js-1-7-5-web.js", tmpDir))
	assert.Contains(t, string(b), "import {ascending as t,descending as n,least as $r} from \"./solid-js-1-7-5.js\"")

	b, _ = os.ReadFile(fmt.Sprintf("%s/solid-js-1-7-5.js", tmpDir))
	assert.Contains(t, string(b), "export default null")
}

func TestObservablehqEsm(t *testing.T) {
	defer gock.Disable()
	plotResponse := `
/**
 * Bundled by jsDelivr using Rollup v2.79.1 and Terser v5.17.1.
 * Original file: /npm/@observablehq/plot@0.6.6/src/index.js
 *
 * Do NOT use SRI with dynamically generated files! More information: https://www.jsdelivr.com/using-sri-with-dynamic-files
 */
import{ascending as t,descending as n,least as $r}from"/npm/d3@7.8.4/+esm";import{parse as Mr,format as Lr}from"/npm/isoformat@0.2.1/+esm";function Er(t){return null}
//# sourceMappingURL=/sm/01413fe1f7a8c2da69e83bcc6c3f16a63658e3f27f5a8ffeda7da895d71e4aa2.map
`

	d3Response := `
/**
 * Bundled by jsDelivr using Rollup v2.74.1 and Terser v5.15.1.
 * Original file: /npm/d3@7.8.4/src/index.js
 *
 * Do NOT use SRI with dynamically generated files! More information: https://www.jsdelivr.com/using-sri-with-dynamic-files
 */
export*from"/npm/d3-array@3.2.3/+esm";
`

	isoFormatResponse := `
/**
 * Bundled by jsDelivr using Rollup v2.74.1 and Terser v5.15.1.
 * Original file: /npm/isoformat@0.2.1/src/index.js
 *
 * Do NOT use SRI with dynamically generated files! More information: https://www.jsdelivr.com/using-sri-with-dynamic-files
 */
function t(t,n){};
//# sourceMappingURL=/sm/83fe3f74d02cac187ee0f8c70305b5a8a44813bc43c57abfb7582eb11b5b40df.map
`

	d3ArrayResponse := `
/**
 * Bundled by jsDelivr using Rollup v2.74.1 and Terser v5.15.1.
 * Original file: /npm/d3-array@3.2.3/src/index.js
 *
 * Do NOT use SRI with dynamically generated files! More information: https://www.jsdelivr.com/using-sri-with-dynamic-files
 */
alert('hi')
//# sourceMappingURL=/sm/ed2296044476ebbad3a766409a273f11d2a4aa63582d856d2768d3187e250781.map
`
	tmpDir, err := os.MkdirTemp("", "*")
	assert.NoError(t, err)
	gock.New("https://cdn.jsdelivr.net").
		Get("/npm/@observablehq/plot/+esm").
		Reply(200).
		BodyString(plotResponse)

	gock.New("https://cdn.jsdelivr.net").
		Get("/npm/d3@7.8.4/+esm").
		Reply(200).
		BodyString(d3Response)

	gock.New("https://cdn.jsdelivr.net").
		Get("/npm/isoformat@0.2.1/+esm").
		Reply(200).
		BodyString(isoFormatResponse)

	gock.New("https://cdn.jsdelivr.net").
		Get("/npm/d3-array@3.2.3/+esm").
		Reply(200).
		BodyString(d3ArrayResponse)

	_, _, err = executeCmd(rootCmd(), "@observablehq/plot", tmpDir)
	assert.NoError(t, err)
	fmt.Printf("%s\n", tmpDir)
	assert.FileExists(t, fmt.Sprintf("%s/observablehq-plot-0-6-6.js", tmpDir))
	assert.FileExists(t, fmt.Sprintf("%s/d3-array-3-2-3.js", tmpDir))
	assert.FileExists(t, fmt.Sprintf("%s/isoformat-0-2-1.js", tmpDir))
	assert.FileExists(t, fmt.Sprintf("%s/d3-7-8-4.js", tmpDir))

	// check the code was rewritten
	b, _ := os.ReadFile(fmt.Sprintf("%s/observablehq-plot-0-6-6.js", tmpDir))
	assert.Contains(t, string(b), "import {ascending as t,descending as n,least as $r} from \"./d3-7-8-4.js\"")
	assert.Contains(t, string(b), "import {parse as Mr,format as Lr} from \"./isoformat-0-2-1.js\"")

	b, _ = os.ReadFile(fmt.Sprintf("%s/d3-7-8-4.js", tmpDir))
	assert.Contains(t, string(b), "export * from \"./d3-array-3-2-3.js\"")

}

func executeCmd(root *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	c, err = root.ExecuteC()

	return c, buf.String(), err
}
