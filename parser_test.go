// Copyright 2020 - 2024 The xgen Authors. All rights reserved. Use of this
// source code is governed by a BSD-style license that can be found in the
// LICENSE file.
//
// Package xgen written in pure Go providing a set of functions that allow you
// to parse XSD (XML schema files). This library needs Go version 1.10 or
// later.

package xgen

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testFixtureDir = "test"
	// externalFixtureDir is where one copy their own XSDs to run validation on them. For a set
	// of XSDs to run tests on, see https://github.com/xuri/xsd. Note that external tests leave the
	// generated output for inspection to support use-cases of manual review of generated code
	externalFixtureDir = "data"
)

func TestParseGo(t *testing.T) {
	testParseForSource(t, "Go", "go", "go", testFixtureDir, false, false)
}

// TestParseGoExternal runs tests on any external XSDs within the externalFixtureDir
func TestParseGoExternal(t *testing.T) {
	testParseForSource(t, "Go", "go", "go", externalFixtureDir, true, false)
}

// TestParseGoInMemory runs tests using an in-memory schema
func TestParseGoInMemory(t *testing.T) {
	testParseForSource(t, "Go", "go", "go", testFixtureDir, false, true)
}

// testParseForSource runs parsing tests for a given language. The sourceDirectory specifies the root of the
// input for the tests. The expected structure of the sourceDirectory is as follows:
//
//	source
//	├── xsd (with the input xsd files to run through the parser)
//	└── <langDirName> (with the expected generated code named <xsd-file>.<fileExt>
//
// The test cleans up files it generates unless leaveOutput is set to true. In which case, the generate file is left
// on disk for manual inspection under <sourceDirectory>/<langDirName>/output.
//
// The inMemorySchema flag is used to test parsing with an in-memory schema.
func testParseForSource(t *testing.T, lang string, fileExt string, langDirName string, sourceDirectory string, leaveOutput bool, inMemorySchema bool) {
	codeDir := filepath.Join(sourceDirectory, langDirName)

	outputDir := filepath.Join(codeDir, "output")
	if leaveOutput {
		err := PrepareOutputDir(outputDir)
		require.NoError(t, err)
	} else {
		tempDir, err := ioutil.TempDir(codeDir, "output-*")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		outputDir = tempDir
	}

	inputDir := filepath.Join(sourceDirectory, "xsd")
	files, err := GetFileList(inputDir)
	// Abort testing if the source directory doesn't include a xsd directory with inputs
	if os.IsNotExist(err) {
		return
	}

	require.NoError(t, err)
	for _, file := range files {
		if filepath.Ext(file) == ".xsd" {
			xsdName, err := filepath.Rel(inputDir, file)
			require.NoError(t, err)

			var schema []byte
			if inMemorySchema {
				t.Logf("Using in-memory schema for %s", file)
				schemaFile, err := ioutil.ReadFile(file)
				require.NoError(t, err)
				schema = schemaFile
			}

			t.Run(xsdName, func(t *testing.T) {
				parser := NewParser(&Options{
					FilePath:            file,
					Schema:              schema,
					InputDir:            inputDir,
					OutputDir:           outputDir,
					Lang:                lang,
					IncludeMap:          make(map[string]bool),
					LocalNameNSMap:      make(map[string]string),
					NSSchemaLocationMap: make(map[string]string),
					ParseFileList:       make(map[string]bool),
					ParseFileMap:        make(map[string][]interface{}),
					ProtoTree:           make([]interface{}, 0),
				})
				err = parser.Parse()
				assert.NoError(t, err, file)
				generatedFileName := strings.TrimPrefix(file, inputDir) + "." + fileExt
				actualFilename := filepath.Join(outputDir, generatedFileName)

				actualGenerated, err := ioutil.ReadFile(actualFilename)
				assert.NoError(t, err)

				expectedFilename := filepath.Join(codeDir, generatedFileName)
				expectedGenerated, err := ioutil.ReadFile(expectedFilename)
				assert.NoError(t, err)

				assert.Equal(t, string(expectedGenerated), string(actualGenerated), fmt.Sprintf("error in generated code for %s", file))
			})
		}
	}
}

func TestParseTypeScript(t *testing.T) {
	testParseForSource(t, "TypeScript", "ts", "ts", testFixtureDir, false, false)
}

func TestParseTypeScriptExternal(t *testing.T) {
	testParseForSource(t, "TypeScript", "ts", "ts", externalFixtureDir, true, false)
}

func TestParseC(t *testing.T) {
	testParseForSource(t, "C", "h", "c", testFixtureDir, false, false)
}

func TestParseCExternal(t *testing.T) {
	testParseForSource(t, "C", "h", "c", externalFixtureDir, true, false)
}

func TestParseJava(t *testing.T) {
	testParseForSource(t, "Java", "java", "java", testFixtureDir, false, false)
}

func TestParseJavaExternal(t *testing.T) {
	testParseForSource(t, "Java", "java", "java", externalFixtureDir, true, false)
}

func TestParseRust(t *testing.T) {
	testParseForSource(t, "Rust", "rs", "rs", testFixtureDir, false, false)
}

func TestParseRustExternal(t *testing.T) {
	testParseForSource(t, "Rust", "rs", "rs", externalFixtureDir, true, false)
}

func Test_validateOptions(t *testing.T) {
	type args struct {
		opt *Options
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "empty file path",
			args: args{
				opt: &Options{
					FilePath:            "",
					OutputDir:           "",
					Extract:             false,
					Lang:                "",
					IncludeMap:          nil,
					LocalNameNSMap:      nil,
					NSSchemaLocationMap: nil,
					ParseFileList:       nil,
					ParseFileMap:        nil,
					ProtoTree:           nil,
				},
			},
			wantErr: assert.Error,
		},
		{
			name: "empty output dir",
			args: args{
				opt: &Options{
					FilePath:            "some file",
					OutputDir:           "",
					Extract:             false,
					Lang:                "",
					IncludeMap:          nil,
					LocalNameNSMap:      nil,
					NSSchemaLocationMap: nil,
					ParseFileList:       nil,
					ParseFileMap:        nil,
					ProtoTree:           nil,
				},
			},
			wantErr: assert.Error,
		},
		{
			name: "empty input dir",
			args: args{
				opt: &Options{
					FilePath:            "some file",
					InputDir:            "",
					Extract:             false,
					Lang:                "go",
					IncludeMap:          nil,
					LocalNameNSMap:      nil,
					NSSchemaLocationMap: nil,
					ParseFileList:       nil,
					ParseFileMap:        nil,
					ProtoTree:           nil,
				},
			},
			wantErr: assert.Error,
		},
		{
			name: "empty lang",
			args: args{
				opt: &Options{
					FilePath:            "some file",
					OutputDir:           "some dir",
					Extract:             false,
					Lang:                "",
					IncludeMap:          nil,
					LocalNameNSMap:      nil,
					NSSchemaLocationMap: nil,
					ParseFileList:       nil,
					ParseFileMap:        nil,
					ProtoTree:           nil,
				},
			},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, validateOptions(tt.args.opt), fmt.Sprintf("validateOptions(%v)", tt.args.opt))
		})
	}
}
