package checkers

import (
	"os/exec"
	"strings"
)

type prettierChecker struct {
	name           string
	packageManager string
	lockFile       string
	executable     string
	command        []string
	dependency     string
}

func (c *prettierChecker) Name() string {
	return c.name
}

func (c *prettierChecker) Runnable(workDir string) (bool, error) {
	return requires(
		jsProjectInDir(workDir),
		fileInDir(workDir, c.lockFile),
		packageHasDependency(workDir, c.dependency),
		executableOnPath(c.executable),
	)
}

func (c *prettierChecker) Check(workDir string) ([]Issue, error) {
	cmd := exec.Command(c.executable, c.command...)
	cmd.Dir = workDir
	output, err := cmd.CombinedOutput()

	// If prettier --check succeeds (exit code 0), all files are formatted correctly
	if err == nil {
		return nil, nil
	}

	var issues []Issue
	if len(output) > 0 {
		lines := strings.Split(strings.TrimSpace(string(output)), "\n")
		for _, line := range lines {
			if line = strings.TrimSpace(line); line != "" {
				// Skip informational messages
				if strings.Contains(line, "Checking formatting...") ||
					strings.Contains(line, "All matched files use Prettier code style!") {
					continue
				}
				issues = append(issues, Issue{
					Message: line,
				})
			}
		}
	}

	return issues, nil
}

func init() {
	Register(&prettierChecker{
		name:           "npm-prettier",
		packageManager: "npm",
		lockFile:       "package-lock.json",
		executable:     "npm",
		command:        []string{"run", "prettier", "--", "--check", "."},
		dependency:     "prettier",
	})
	Register(&prettierChecker{
		name:           "bun-prettier",
		packageManager: "bun",
		lockFile:       "bun.lock",
		executable:     "bun",
		command:        []string{"run", "prettier", "--", "--check", "."},
		dependency:     "prettier",
	})
}
