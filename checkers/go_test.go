package checkers

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func init() {
	Register(&goTestChecker{})
}

type goTestChecker struct{}

func (g *goTestChecker) Name() string {
	return "go test"
}

func (g *goTestChecker) Runnable(workDir string) (bool, error) {
	return requires(goProjectInDir(workDir), executableOnPath("go"))
}

func (g *goTestChecker) Check(workDir string) ([]Issue, error) {
	cmd := exec.Command("go", "test", "./...")
	cmd.Dir = workDir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// go test returns non-zero exit code when tests fail
	err := cmd.Run()
	if err != nil && stdout.Len() == 0 && stderr.Len() == 0 {
		// If there's an error but no output, it's a real error
		return nil, fmt.Errorf("failed to run go test: %w", err)
	}

	var issues []Issue

	// Parse stderr for build errors (skip harmless download messages)
	if stderr.Len() > 0 {
		scanner := bufio.NewScanner(&stderr)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line != "" && !strings.HasPrefix(line, "go: downloading") {
				issues = append(issues, Issue{
					Message: line,
				})
			}
		}
	}

	// Parse stdout for test failures
	scanner := bufio.NewScanner(&stdout)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Look for FAIL lines which indicate test failures
		if strings.HasPrefix(line, "FAIL") || strings.Contains(line, "--- FAIL:") {
			issues = append(issues, Issue{
				Message: line,
			})
		}
		// Also capture panic and error messages
		if strings.Contains(line, "panic:") || strings.Contains(line, "Error:") {
			issues = append(issues, Issue{
				Message: line,
			})
		}
	}

	if scanErr := scanner.Err(); scanErr != nil {
		return nil, fmt.Errorf("failed to parse go test output: %w", scanErr)
	}

	return issues, nil
}
