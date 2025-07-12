package checkers

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func init() {
	Register(&gofmtChecker{})
}

type gofmtChecker struct{}

func (g *gofmtChecker) Name() string {
	return "gofmt"
}

func (g *gofmtChecker) Runnable(workDir string) (bool, error) {
	return requires(goProjectInDir(workDir), executableOnPath("gofmt"))
}

func (g *gofmtChecker) Check(workDir string) ([]Issue, error) {
	cmd := exec.Command("gofmt", "-l", ".")
	cmd.Dir = workDir
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run gofmt: %w", err)
	}

	var issues []Issue
	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		file := strings.TrimSpace(scanner.Text())
		if file != "" {
			issues = append(issues, Issue{
				File:    file,
				Message: fmt.Sprintf("File needs formatting: %s", file),
			})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to parse gofmt output: %w", err)
	}

	return issues, nil
}
