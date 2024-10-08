package cmd

import (
	"strings"

	"github.com/sommerfeld-io/source2adoc/internal/codefiles"
	"github.com/spf13/cobra"
)

const rootDescShort = "Streamline the process of generating documentation from inline comments within source code files."
const rootDescLong = `
Facilitate the creation of comprehensive and well-structured documentation
directly from code comments. The app supports multiple source code languages.
The common ground is, that these languages mark their comments through the
hash-symbol (#).

For more information, visit the project's documentation:
  https://source2adoc.sommerfeld.io

Quick Start:
  The root command source2adoc [flags] scans the --source-dir for code files
  and starts the conversion process. The output is written to --output-dir.

Example:
  source2adoc --source-dir ./src --output-dir ./docs

Example (Docker):
  docker run -v "$(pwd):$(pwd)" -w "$(pwd)" sommerfeldio/source2adoc:latest -s ./src -o ./docs

Supported Languages:
  `

func supportedLanguagesDesc() string {
	supportedLanguages := codefiles.SupportedCodeFilenames
	keys := make([]string, 0, len(supportedLanguages))
	for k := range supportedLanguages {
		keys = append(keys, k)
	}
	return strings.Join(keys, ", ")
}

// Values from the CLI flags
var (
	sourceDir string
	outputDir string
	exclude   []string
)

var rootCmd = &cobra.Command{
	Use:   "source2adoc",
	Short: rootDescShort,
	Long:  rootDescLong + supportedLanguagesDesc(),

	Args: cobra.ExactArgs(0),

	Run: func(cmd *cobra.Command, args []string) {
		excludes, err := getExcludes(cmd)
		handleError(err)

		sourceCodeFiles := findCodeFiles(excludes)
		sourceCodeFiles = readCodeFiles(sourceCodeFiles)
		sourceCodeFiles = parseFileContent(sourceCodeFiles)
		writeDocsFiles(sourceCodeFiles)
	},
}

func getExcludes(cmd *cobra.Command) ([]string, error) {
	exclude, err := cmd.Flags().GetStringSlice("exclude")
	return exclude, err
}

func findCodeFiles(exclude []string) []*codefiles.CodeFile {
	finder := codefiles.NewFinder(sourceDir)
	finder.SetExcludes(exclude)
	sourceCodeFiles, err := finder.FindSourceCodeFiles()

	if err != nil {
		handleError(err)
	}
	return sourceCodeFiles
}

// readCodeFiles reads the code files from the source directory.
func readCodeFiles(files []*codefiles.CodeFile) []*codefiles.CodeFile {
	for _, file := range files {
		err := file.ReadFileContent()
		handleError(err)
	}
	return files
}

// parseFileContent parses the content of the code files for comments.
func parseFileContent(files []*codefiles.CodeFile) []*codefiles.CodeFile {
	for _, file := range files {
		err := file.Parse()
		handleError(err)
	}
	return files
}

// writeDocsFiles writes the documentation files to the output directory.
func writeDocsFiles(files []*codefiles.CodeFile) {
	for _, file := range files {
		err := file.WriteDocumentationFile(outputDir)
		handleError(err)
	}
}

func init() {
	initMandatoryFlags()
	initMultipleValuesFlags()
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
}

func initMandatoryFlags() {
	var params = []struct {
		name     string
		short    string
		variable *string
		desc     string
	}{
		{name: "source-dir", short: "s", variable: &sourceDir, desc: "Directory containing the source code files"},
		{name: "output-dir", short: "o", variable: &outputDir, desc: "Directory to write the generated documentation to"},
	}

	for _, param := range params {
		rootCmd.Flags().StringVarP(param.variable, param.name, param.short, "", param.desc)
		err := rootCmd.MarkFlagRequired(param.name)
		handleError(err)
	}
}

func initMultipleValuesFlags() {

	var params = []struct {
		name      string
		short     string
		variable  *[]string
		desc      string
		mandatory bool
	}{
		{name: "exclude", short: "x", variable: &exclude, desc: "Exclude files and/or folders when generating documentation"},
	}

	for _, param := range params {
		rootCmd.Flags().StringSliceVarP(param.variable, param.name, param.short, []string{}, param.desc)
	}
}

// Execute acts as the entrypoint for the CLI app.
func Execute() {
	err := rootCmd.Execute()
	handleError(err)
}

// RegisterSubCommand adds a subcommand to the root command.
func RegisterSubCommand(cmd *cobra.Command) {
	rootCmd.AddCommand(cmd)
}
