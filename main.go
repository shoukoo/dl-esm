package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

func main() {
	if err := rootCmd().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func rootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dl-esm PAKCAGE DIR[optional]",
		Short: "Download ESM modules from npm and jsdelivr",
		Long:  "Download ESM modules from npm and jsdelivr",
		Example: `
# download latest version of solid js
dl-easm solid-js

# download a specific version of solid js
dl-easm solid-js@1.7.5

# download to a specific dir
dl-easm solid-js@1.7.5 /tmp
		`,
		Args: cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {

			location := "."
			packageArg := args[0]

			if len(args) == 2 {
				location = args[1]
			}

			root := filepath.Join(location)
			if _, err := os.Stat(root); os.IsNotExist(err) {
				os.MkdirAll(root, os.ModePerm)
			}

			var code string

			if strings.HasPrefix(packageArg, "https://") {

				response, err := http.Get(packageArg)
				if err != nil {
					return err
				}

				defer response.Body.Close()

				bytes, err := io.ReadAll(response.Body)
				if err != nil {
					return err
				}

				code = string(bytes)

			} else {

				code = fetchCode(packageArg)
			}

			pkgVer := extractPkgVer(code, packageArg)
			path := simplifyPath(fmt.Sprintf("/npm/%s/+esm", pkgVer))

			// Rewrite code and save to path
			rewrittenCode, capturedPaths := rewriteCode(code)
			err := os.WriteFile(filepath.Join(root, path), []byte(rewrittenCode), 0644)
			if err != nil {
				return err
			}

			fmt.Fprintln(os.Stderr, path)

			// Do the same thing for all the captured paths, recursively
			toFetch := map[string]string{}
			for k, v := range capturedPaths {
				toFetch[k] = v
			}

			for len(toFetch) > 0 {

				for path, simplifiedPath := range toFetch {

					delete(toFetch, path)

					url := "https://cdn.jsdelivr.net" + path
					response, err := http.Get(url)
					if err != nil {
						return err
					}

					defer response.Body.Close()
					bytes, err := io.ReadAll(response.Body)
					if err != nil {
						return err
					}
					code := string(bytes)

					rewrittenCode, moreCapturedPaths := rewriteCode(code)
					err = os.WriteFile(filepath.Join(root, simplifiedPath), []byte(rewrittenCode), 0644)
					if err != nil {
						return err
					}

					fmt.Fprintln(os.Stderr, simplifiedPath)

					for k, v := range moreCapturedPaths {
						toFetch[k] = v
					}

				}
			}

			return nil

		},
	}

	return cmd
}

func fetchCode(packageArg string) string {
	url := fmt.Sprintf("https://cdn.jsdelivr.net/npm/%s/+esm", packageArg)
	response, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	return string(bytes)

}

func extractPkgVer(content string, pkgName string) string {

	// if pkgName includes version then return
	r := regexp.MustCompile(`([^@]+)@(\d+\.\d+\.\d+)`)
	matches := r.FindStringSubmatch(pkgName)
	if len(matches) > 2 {
		return pkgName
	}

	// extract path from the file
	pattern := `Original file: (\/npm\/[^\s]+)`
	r = regexp.MustCompile(pattern)
	matches = r.FindStringSubmatch(content)

	if len(matches) < 1 {
		panic("Could not find original file")
	}

	// split sub dirs if there is one
	pkgNameSubDir := strings.SplitN(pkgName, "/", 2)
	subDir := ""

	if len(pkgNameSubDir) == 2 {
		subDir = "/" + pkgNameSubDir[1]
		pkgName = pkgNameSubDir[0]
	}

	// extract package and version
	pattern = fmt.Sprintf(`/%s@(\d+\.\d+\.\d+)`, pkgName)
	r = regexp.MustCompile(pattern)
	matches = r.FindStringSubmatch(matches[1])

	if len(matches) < 2 {
		panic("unable to extract package info")
	}

	return pkgName + "@" + matches[1] + subDir

}

func simplifyPath(path string) string {

	split := strings.Split(path, "/npm/")
	split[1] = strings.TrimPrefix(split[1], "@")

	packageInfo := strings.SplitN(split[1], "@", 2)

	packageName := strings.ReplaceAll(packageInfo[0], "/", "-")

	packageVersion := strings.Split(packageInfo[1], "/dist")[0]
	packageVersion = strings.ReplaceAll(packageVersion, "/+esm", "")
	packageVersion = strings.ReplaceAll(packageVersion, ".", "-")
	packageVersion = strings.ReplaceAll(packageVersion, "/", "-")

	simplifiedName := fmt.Sprintf("%s-%s.js", packageName, packageVersion)

	return simplifiedName
}

func rewriteCode(code string) (string, map[string]string) {

	pattern := `(?P<keyword>import|export)\s*(?P<imports>\{?[^}]+?\}?)\s*from\s*"(?P<path>\/npm\/[^"]+)"`
	re := regexp.MustCompile(pattern)
	capturedPaths := map[string]string{}

	replaceImport := func(match string) string {

		submatch := re.FindStringSubmatch(match)

		keyword := submatch[1]
		imports := submatch[2]
		path := submatch[3]

		simplifiedPath := simplifyPath(path)

		capturedPaths[path] = simplifiedPath

		return fmt.Sprintf(`%s %s from "./%s";`, keyword, imports, simplifiedPath)
	}

	rewrittenCode := re.ReplaceAllStringFunc(code, replaceImport)
	rewrittenCode = removeSourceMappingComments(rewrittenCode)

	return rewrittenCode, capturedPaths

}

func removeSourceMappingComments(code string) string {

	pattern := `\/\/#\s*sourceMappingURL=.*?\.map`
	re := regexp.MustCompile(pattern)
	cleanedCode := re.ReplaceAllString(code, "")

	return cleanedCode

}
