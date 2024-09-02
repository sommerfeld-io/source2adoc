package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cucumber/godog"
	"github.com/sommerfeld-io/source2adoc-acceptance-tests/testhelper"
)

type Paths struct {
	codeFile string
	adocFile string
}

type TestState struct {
	sourceDir   string
	outputDir   string
	cmdWithArgs []string
	exitCode    int
	paths       []Paths
}

func (ts *TestState) reset() {
	ts.cmdWithArgs = []string{}
	ts.sourceDir = ""
	ts.outputDir = ""
	ts.exitCode = 0
	ts.paths = []Paths{}
}

// appendCommand appends to the command which will be passed to the app (with arguments and flags).
func (ts *TestState) appendCommand(cmd ...string) {
	ts.cmdWithArgs = append(ts.cmdWithArgs, cmd...)
}

// cleanup removes the output directory if it was created during the test.
func (ts *TestState) cleanup() error {
	if ts.outputDir == "" {
		return nil
	}

	err := os.RemoveAll(ts.outputDir)
	if err != nil {
		return fmt.Errorf("error cleaning up target directory: %v", err)
	}
	return nil
}

func Test_BasicFeatures(t *testing.T) {
	featureFile := "basic.feature"
	opts := testhelper.Options(t, featureFile)

	suite := godog.TestSuite{
		Name:                 featureFile,
		ScenarioInitializer:  initializeBasicScenario,
		TestSuiteInitializer: initializeTestSuite,
		Options:              opts,
	}

	exitcode := suite.Run()
	if exitcode != 0 {
		t.Fatal(suite.Name, "|", "non-zero status returned.", "failed to run tests.", "finished with exit code", exitcode)

	}
}

func initializeTestSuite(sc *godog.TestSuiteContext) {
}

func initializeBasicScenario(sc *godog.ScenarioContext) {
	ts := &TestState{}

	sc.Step(`^I use the root command of the source2adoc CLI tool$`, ts.iUseTheRootCommand)
	sc.Step(`^I specify the "([^"]*)" flag$`, ts.iSpecifyTheFlag)
	sc.Step(`^I specify the "([^"]*)" flag with value "([^"]*)"$`, ts.iSpecifyTheFlagWithValue)
	sc.Step(`^I run the app$`, ts.iRunTheApp)
	sc.Step(`^exit code should be (\d+)$`, ts.exitCodeShouldBe)
	sc.Step(`^no AsciiDoc files should be generated$`, ts.noAsciiDocFilesShouldBeGenerated)
	sc.Step(`^AsciiDoc files should be generated for all source code files$`, ts.asciiDocFilesShouldBeGeneratedForAllSourceCodeFiles)
	sc.Step(`^the path of the generated docs in the --output-dir directory should mimic the source code file\'s path$`, ts.theAdocPathShouldMimicCodeFilesPath)

	sc.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		ts.reset()
		return ctx, nil
	})

	sc.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		return ctx, ts.cleanup()
	})
}

func (ts *TestState) iUseTheRootCommand() error {
	info, err := os.Stat(testhelper.BinaryPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("source2adoc binary does not exist: %v", err)
		} else {
			return fmt.Errorf("failed to check source2adoc binary existence: %v", err)
		}
	}

	if info.Mode()&os.ModeType != 0 {
		return fmt.Errorf("source2adoc binary is not executable")
	}

	// The root cmd does not require a dedicated command name
	return nil
}

func (ts *TestState) iSpecifyTheFlag(flag string) error {
	ts.appendCommand(flag)
	return nil
}

func (ts *TestState) iSpecifyTheFlagWithValue(flag, value string) error {
	if flag == "--source-dir" {
		ts.sourceDir = value
	}
	if flag == "--output-dir" {
		ts.outputDir = value
	}

	ts.appendCommand(flag, value)
	return nil
}

func (ts *TestState) iRunTheApp() error {
	cmd := exec.Command(testhelper.BinaryPath, ts.cmdWithArgs...)
	err := cmd.Run()
	if err != nil {
		exitErr, ok := err.(*exec.ExitError)
		ts.exitCode = exitErr.ExitCode()
		if !ok || ts.exitCode != 1 {
			return fmt.Errorf("failed to run the app: %v", err)
		}
	}
	return nil
}

func (ts *TestState) exitCodeShouldBe(expected int) error {
	if ts.exitCode != expected {
		return fmt.Errorf("expected exit code %d, got %d", expected, ts.exitCode)
	}
	return nil
}

func (ts *TestState) noAsciiDocFilesShouldBeGenerated() error {
	if ts.outputDir == "" {
		return fmt.Errorf("output directory not set")
	}

	if _, err := os.Stat(ts.outputDir); !os.IsNotExist(err) {
		return fmt.Errorf(ts.outputDir, "directory should not exist")
	}
	return nil
}

func (ts *TestState) asciiDocFilesShouldBeGeneratedForAllSourceCodeFiles() error {
	if ts.sourceDir == "" {
		return fmt.Errorf("source directory not set")
	}
	if ts.outputDir == "" {
		return fmt.Errorf("output directory not set")
	}

	paths, err := testhelper.FindSourceCodeFiles(ts.sourceDir)
	for _, path := range paths {
		Paths := Paths{
			codeFile: path,
		}
		ts.paths = append(ts.paths, Paths)
	}

	if err != nil {
		return fmt.Errorf("failed to find source code files: %v", err)
	}
	if len(ts.paths) == 0 {
		return fmt.Errorf("no source code files found")
	}

	tmp := []Paths{}
	for _, path := range ts.paths {
		adoc := testhelper.TranslateFilename(path.codeFile)
		path.adocFile = adoc
		tmp = append(tmp, path)
	}
	ts.paths = []Paths{}
	ts.paths = append(ts.paths, tmp...)

	for _, paths := range ts.paths {
		fullPath := filepath.Join(ts.outputDir, paths.adocFile)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			return fmt.Errorf("AsciiDoc file %s does not exist in output directory", paths.adocFile)
		}
	}

	return nil
}

func (ts *TestState) theAdocPathShouldMimicCodeFilesPath() error {
	if ts.paths == nil {
		return fmt.Errorf("no paths found")
	}

	for _, paths := range ts.paths {
		codeFile := paths.codeFile
		codeDir := filepath.Dir(codeFile)

		if !strings.Contains(paths.adocFile, codeDir) {
			return fmt.Errorf("AsciiDoc file path %s does not mimic the source code file path %s", paths.adocFile, codeDir)
		}
	}

	return nil
}
