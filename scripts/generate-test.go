// +build ignore

// Copyright (c) 2020 Markus Mobius
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	nurl "net/url"
	"os"
	"path/filepath"
	fp "path/filepath"
	"time"

	distiller "github.com/markusmobius/go-domdistiller"
	"github.com/markusmobius/go-domdistiller/data"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var httpClient = &http.Client{Timeout: time.Minute}

func main() {
	rootCmd := &cobra.Command{
		Use:   "go run scripts/generate-test.go",
		Short: "Generate test files for go-domdistiller.",
		Long: "Generate test files for go-domdistiller. This script has several behaviors " +
			"depending on number of arguments:\n" +
			"- no arguments, regenerate all `expected.html` inside test-files directory;\n" +
			"- 1 argument, regenerate `expected.html` for specified test file name;\n" +
			"- 2 arguments, load URL in 2nd args and save it as test file with name in 1st args.",
		Args: cobra.RangeArgs(0, 2),
		Run: func(cmd *cobra.Command, args []string) {
			var testName, sourceURL string
			switch len(args) {
			case 1:
				testName = args[0]
			case 2:
				testName = args[0]
				sourceURL = args[1]
			}

			// If test name is empty, generate test case for all existing test directory
			if testName == "" {
				dirItems, err := ioutil.ReadDir("test-files")
				if err != nil {
					logrus.Fatalf("failed to read test dir: %v\n", err)
				}

				for _, item := range dirItems {
					if !item.IsDir() {
						continue
					}

					sourcePath := fp.Join("test-files", item.Name(), "source.html")
					if !fileExists(sourcePath) {
						continue
					}

					err = generateTest(item.Name(), "")
					if err != nil {
						logrus.Fatalf("failed to generate test for %s: %v\n", item.Name(), err)
					}
				}

				return
			}

			err := generateTest(testName, sourceURL)
			if err != nil {
				logrus.Panicln(err)
			}
		},
	}

	rootCmd.Execute()
}

func generateTest(testName, sourceURL string) error {
	logrus.Println("Generating test for", testName)

	// Check if source file for test exists
	// If source file doesn't exist, download it first.
	// If it exist, but URL is defined as well, redownload it
	sourcePath := filepath.Join("test-files", testName, "source.html")
	if !fileExists(sourcePath) || sourceURL != "" {
		// Download HTML file from URL.
		logrus.Printf("downloading source for %s from %s\n", testName, sourceURL)
		err := downloadWebPage(sourceURL, sourcePath)
		if err != nil {
			return fmt.Errorf("failed to download source: %w", err)
		}
	}

	// If source URL defined, we use it as distiller's URL.
	// If not and expected.json exist, use the page URL inside it
	var url *nurl.URL
	expectedJSON := fp.Join("test-files", testName, "expected.json")

	if sourceURL != "" {
		url, _ = nurl.ParseRequestURI(sourceURL)
	} else if fileExists(expectedJSON) {
		data, err := decodeExpectedJSON(expectedJSON)
		if err != nil {
			return fmt.Errorf("failed to decode JSON: %w", err)
		}

		url, _ = nurl.ParseRequestURI(data.URL)
	}

	// Run distiller for the file
	opts := &distiller.Options{OriginalURL: url, LogFlags: 0}
	result, err := distiller.ApplyForFile(sourcePath, opts)
	if err != nil {
		return fmt.Errorf("failed to distill source: %w", err)
	}

	// Save result to file
	return writeResultToFile(result, testName)
}

func fileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	return !os.IsNotExist(err) && !info.IsDir()
}

func downloadWebPage(srcURL string, dstPath string) error {
	// Verify that URL is valid.
	if _, err := nurl.ParseRequestURI(srcURL); err != nil {
		return fmt.Errorf("failed to parse URL: %w", err)
	}

	// Download HTML file from URL.
	resp, err := httpClient.Get(srcURL)
	if err != nil {
		return fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	// Save to file
	os.MkdirAll(fp.Dir(dstPath), os.ModePerm)
	dst, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	return nil
}

func writeResultToFile(result *distiller.Result, testName string) error {
	// Save HTML file
	expectedHTMLPath := filepath.Join("test-files", testName, "expected.html")
	err := ioutil.WriteFile(expectedHTMLPath, []byte(result.HTML), os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to write HTML: %w", err)
	}

	// Save JSON file
	result.HTML = ""
	result.TimingInfo = data.TimingInfo{}
	bt, err := json.MarshalIndent(&result, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	expectedJSONPath := filepath.Join("test-files", testName, "expected.json")
	err = ioutil.WriteFile(expectedJSONPath, bt, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to write JSON: %w", err)
	}

	return nil
}

func decodeExpectedJSON(path string) (*distiller.Result, error) {
	// Open file
	src, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %w", path, err)
	}
	defer src.Close()

	var result distiller.Result
	err = json.NewDecoder(src).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	return &result, nil
}
