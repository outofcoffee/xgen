package xgen

import (
	"encoding/xml"
	"io/ioutil"
	"path/filepath"
	"testing"

	schema "github.com/outofcoffee/xgen/test/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGeneratedGo runs through test cases to validate Go generated structs. Each test case
// requires a xml fixture file to unmarshal into the receiving struct. Validate first validates
// that the file can be unmarshaled as the receiving struct and then remarshals the content
// to make sure the marshaling is symmetrical
func TestGeneratedGo(t *testing.T) {
	testCases := []struct {
		// xmlFileName is the path to the xml fixture file to unmarshal into the receiving struct
		xmlFileName string
		// receivingStruct is a pointer to the struct to unmarshal the xml file content into. It should match
		// the type of the top level element present in that file
		receivingStruct interface{}
	}{
		{
			xmlFileName:     "base64.xml",
			receivingStruct: &schema.TopLevel{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.xmlFileName, func(t *testing.T) {
			fullPath := filepath.Join("xmlFixtures", tc.xmlFileName)

			input, err := ioutil.ReadFile(fullPath)
			require.NoError(t, err)

			err = xml.Unmarshal(input, tc.receivingStruct)
			require.NoError(t, err)

			// Validate that decoding resulted in a non-zero value
			assert.NotEmpty(t, tc.receivingStruct)

			// Remarshal the parsed content to compare it with the original and make sure that the parsing/encoding
			// is symmetrical
			remarshaled, err := xml.MarshalIndent(tc.receivingStruct, "", "    ")
			require.NoError(t, err)

			assert.Equal(t, string(input), string(remarshaled))
		})
	}
}

func TestToTitle(t *testing.T) {
	test := func(expected, actual string) {
		assert.Equal(t, expected, ToTitle(actual))
	}

	test("", "")
	test("A", "a")
	test("Ab", "ab")
	test("A b", "a b")
	test("Ab cd", "ab cd")

	// Test Сyrillic (`привет мир` → `hello world`)
	test("Привет", "привет")
	test("Привет мир", "привет мир")
}

func TestCodeGeneratorFileWithExtension(t *testing.T) {
	testCases := []struct {
		description string
		filename    string
		extension   string
		expected    string
	}{
		{
			description: "filename without extension and extension without period should add extension",
			filename:    "foo",
			extension:   "java",
			expected:    "foo.java",
		},
		{
			description: "filename without extension and extension with period should add extension",
			filename:    "foo",
			extension:   ".java",
			expected:    "foo.java",
		},
		{
			description: "filename with extension already should not add extension",
			filename:    "foo.java",
			extension:   ".java",
			expected:    "foo.java",
		},
		{
			description: "filename with different extension should add extension",
			filename:    "foo.bar",
			extension:   ".java",
			expected:    "foo.bar.java",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			gen := CodeGenerator{
				File: tc.filename,
			}
			actual := gen.FileWithExtension(tc.extension)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
