package checkers

import (
	"os/exec"
	"strings"
)

type staticcheckChecker struct{}

func init() {
	Register(&staticcheckChecker{})
}

func (c *staticcheckChecker) Name() string {
	return "staticcheck"
}

func (c *staticcheckChecker) Runnable(workDir string) (bool, error) {
	return requires(goProjectInDir(workDir), executableOnPath("staticcheck"), fileInDir(workDir, "staticcheck.conf"))
}

func (c *staticcheckChecker) Check(workDir string) ([]Issue, error) {
	cmd := exec.Command("staticcheck", "./...")
	cmd.Dir = workDir
	output, err := cmd.CombinedOutput()

	if err != nil && len(output) == 0 {
		return nil, err
	}

	var issues []Issue
	if len(output) > 0 {
		lines := strings.Split(strings.TrimSpace(string(output)), "\n")
		for _, line := range lines {
			if line = strings.TrimSpace(line); line != "" {
				parts := strings.SplitN(line, ":", 2)
				if len(parts) >= 2 {
					issues = append(issues, Issue{
						File:    parts[0],
						Message: strings.TrimSpace(parts[1]),
					})
				} else {
					issues = append(issues, Issue{
						Message: line,
					})
				}
			}
		}
	}

	return issues, nil
}
