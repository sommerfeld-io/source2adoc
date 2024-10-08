package codefiles

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ShouldSplitPathAndFilename(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		path             string
		expectedPath     string
		expectedFilename string
	}{
		{path: "/path/to/source.sh", expectedPath: "/path/to", expectedFilename: "source.sh"},
		{path: "path/to/Dockerfile", expectedPath: "path/to", expectedFilename: "Dockerfile"},
		{path: "source.sh", expectedPath: "", expectedFilename: "source.sh"},
	}

	for _, test := range tests {
		path, filename := splitPathAndFilename(test.path)
		assert.Equal(test.expectedPath, path, "Incorrect path")
		assert.Equal(test.expectedFilename, filename, "Incorrect filename")
	}
}

func Test_ShouldIdentifyLanguage(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		filename  string
		expected  string
		supported bool
	}{
		{filename: "config.yml", expected: LanguageYml, supported: true},
		{filename: "config.yaml", expected: LanguageYml, supported: true},
		{filename: "Dockerfile", expected: LanguageDockerfile, supported: true},
		{filename: "Dockerfile", expected: LanguageDockerfile, supported: true},
		{filename: "Dockerfile.docs", expected: LanguageDockerfile, supported: true},
		{filename: "Vagrantfile.prod", expected: LanguageVagrant, supported: true},
		{filename: "Makefile", expected: LanguageMake, supported: true},
		{filename: "script.sh", expected: LanguageBash, supported: true},
		{filename: "script.go", expected: LanguageNotSupported, supported: false},
	}

	for _, test := range tests {
		lang, supported := identifyLanguage(test.filename)
		assert.Equal(test.expected, lang, "Incorrect language identification for: "+test.filename)
		assert.Equal(test.supported, supported, "Invalid supported status for: "+test.filename)
	}
}
func Test_ShouldGetDataFromGetterFunctions(t *testing.T) {
	assert := assert.New(t)

	codeFile := &CodeFile{
		path:          "/path/to",
		name:          "source.sh",
		lang:          LanguageBash,
		supportedLang: true,
	}

	expectedPath := "/path/to"
	actualPath := codeFile.Path()
	assert.Equal(expectedPath, actualPath, "Incorrect path")

	expectedName := "source.sh"
	actualName := codeFile.Filename()
	assert.Equal(expectedName, actualName, "Incorrect filename")

	expectedLang := LanguageBash
	actualLang := codeFile.Language()
	assert.Equal(expectedLang, actualLang, "Incorrect path language")

	expectedSupported := true
	actualSupported := codeFile.IsSupportedLanguage()
	assert.Equal(expectedSupported, actualSupported, "Incorrect path supported status")
}

func Test_ShouldReadFileContent(t *testing.T) {
	assert := assert.New(t)

	codeFile := &CodeFile{
		path:          filepath.Join(TestSourceDir, "good"),
		name:          "small-comment.sh",
		lang:          LanguageBash,
		supportedLang: true,
	}

	err := codeFile.ReadFileContent()
	assert.Nil(err, "Error reading file content")

	expectedContent := `#!/bin/bash
## Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt
## ut labore et dolore magna aliquyam erat, sed diam voluptua.

## Not part of the header comment
`
	actualContent := codeFile.FileContent()
	assert.Equal(expectedContent, actualContent, "Incorrect file content")
}

func Test_ShouldParseDocumentation(t *testing.T) {
	assert := assert.New(t)

	codeFile := &CodeFile{
		path:          filepath.Join(TestSourceDir, "good"),
		name:          "small-comment.sh",
		lang:          LanguageBash,
		supportedLang: true,
		fileContent: `#!/bin/bash
## Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt
# ignore me because I do not follow the comment convention. Maybe I am a typo.
## ut labore et dolore magna aliquyam erat, sed diam voluptua.

## Not part of the header comment
`,
		documentationParts: []DocumentationPart{},
	}

	expectedDocs := `= small-comment.sh

[cols="1,5"]
|===
|Language |` + LanguageBash + `
|Path |` + TestSourceDir + `/good/small-comment.sh
|===

Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt
ut labore et dolore magna aliquyam erat, sed diam voluptua.
`

	err := codeFile.Parse()
	assert.Nil(err, "Error parsing documentation")

	docs := codeFile.parsedDocumentation()
	assert.Equal(expectedDocs, docs, "Incorrect parsed documentation")
}
func Test_ShouldTranslateDocumentationFileName(t *testing.T) {
	codeFile := &CodeFile{
		path: filepath.Join(TestSourceDir, "good"),
		name: "small-comment.sh",
	}

	expectedFileName := "small-comment-sh.adoc"
	actualFileName := codeFile.documentationFileName()
	assert.Equal(t, expectedFileName, actualFileName, "Incorrect documentation file name")
}

func Test_ShouldWriteDocumentationFile(t *testing.T) {
	assert := assert.New(t)

	codeFile := &CodeFile{
		path:          "some/path",
		name:          "unittest.sh",
		lang:          LanguageBash,
		supportedLang: true,
		documentationParts: []DocumentationPart{
			{
				sectionType:    DocumentationPartHeader,
				sectionContent: "Lorem ipsum dolor sit amet, consetetur sadipscing elitr",
			},
		},
	}

	expectedAdocFile := TestOutputDir + "/" + codeFile.Path() + "/unittest-sh.adoc"

	err := codeFile.WriteDocumentationFile(TestOutputDir)
	assert.Nil(err, "Error writing documentation file")

	_, err = os.Stat(expectedAdocFile)
	assert.False(os.IsNotExist(err), "Documentation file does not exist")

	os.Remove(expectedAdocFile)
}
