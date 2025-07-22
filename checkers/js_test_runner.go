package checkers

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type jsTestChecker struct {
	name           string
	packageManager string
	lockFile       string
	executable     string
	command        []string
}

func (c *jsTestChecker) Name() string {
	return c.name
}

func (c *jsTestChecker) Runnable(workDir string) (bool, error) {
	return requires(
		jsProjectInDir(workDir),
		fileInDir(workDir, c.lockFile),
		packageHasScript(workDir, "test"),
		executableOnPath(c.executable),
	)
}

func (c *jsTestChecker) Check(workDir string) ([]Issue, error) {
	cmd := exec.Command(c.executable, c.command...)
	cmd.Dir = workDir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil && stdout.Len() == 0 && stderr.Len() == 0 {
		return nil, fmt.Errorf("failed to run %s: %w", strings.Join(c.command, " "), err)
	}

	var issues []Issue

	if stderr.Len() > 0 {
		scanner := bufio.NewScanner(&stderr)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line != "" && !strings.HasPrefix(line, "npm WARN") {
				issues = append(issues, Issue{
					Message: line,
				})
			}
		}
	}

	scanner := bufio.NewScanner(&stdout)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.Contains(line, "FAIL") ||
			strings.Contains(line, "✗") ||
			strings.Contains(line, "×") ||
			strings.Contains(line, "failed") ||
			strings.Contains(line, "Error:") ||
			strings.Contains(line, "AssertionError") {
			issues = append(issues, Issue{
				Message: line,
			})
		}
	}

	if scanErr := scanner.Err(); scanErr != nil {
		return nil, fmt.Errorf("failed to parse test output: %w", scanErr)
	}

	return issues, nil
}

func init() {
	Register(&jsTestChecker{
		name:           "npm-test",
		packageManager: "npm",
		lockFile:       "package-lock.json",
		executable:     "npm",
		command:        []string{"test"},
	})
	Register(&jsTestChecker{
		name:           "bun-test",
		packageManager: "bun",
		lockFile:       "bun.lock",
		executable:     "bun",
		command:        []string{"test"},
	})
}
