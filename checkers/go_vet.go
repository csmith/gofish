package checkers

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func init() {
	Register(&goVetChecker{})
}

type goVetChecker struct{}

func (g *goVetChecker) Name() string {
	return "go vet"
}

func (g *goVetChecker) Runnable(workDir string) (bool, error) {
	return requires(goProjectInDir(workDir), executableOnPath("go"))
}

func (g *goVetChecker) Check(workDir string) ([]Issue, error) {
	cmd := exec.Command("go", "vet", "./...")
	cmd.Dir = workDir
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	// go vet returns non-zero exit code when it finds issues
	err := cmd.Run()
	if err != nil && stderr.Len() == 0 {
		// If there's an error but no stderr output, it's a real error
		return nil, fmt.Errorf("failed to run go vet: %w", err)
	}

	var issues []Issue
	scanner := bufio.NewScanner(&stderr)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "go: downloading") {
			issues = append(issues, Issue{
				Message: line,
			})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to parse go vet output: %w", err)
	}

	return issues, nil
}
