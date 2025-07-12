package checkers

import (
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func isInformationalLine(line string) bool {
	return strings.Contains(line, "====================================") ||
		strings.Contains(line, "Getting Svelte diagnostics") ||
		strings.Contains(line, "svelte-check found 0 errors") ||
		strings.Contains(line, "svelte-check found 0 warnings")
}

type svelteCheckChecker struct {
	name           string
	packageManager string
	lockFile       string
	executable     string
	command        []string
	dependency     string
}

func (c *svelteCheckChecker) Name() string {
	return c.name
}

func (c *svelteCheckChecker) Runnable(workDir string) (bool, error) {
	return requires(
		jsProjectInDir(workDir),
		fileInDir(workDir, c.lockFile),
		packageHasDependency(workDir, c.dependency),
		executableOnPath(c.executable),
	)
}

func (c *svelteCheckChecker) Check(workDir string) ([]Issue, error) {
	cmd := exec.Command(c.executable, c.command...)
	cmd.Dir = workDir
	output, err := cmd.CombinedOutput()

	if err != nil && len(output) == 0 {
		return nil, err
	}

	var issues []Issue
	if len(output) > 0 {
		lines := strings.Split(strings.TrimSpace(string(output)), "\n")

		// Check summary line first
		summaryRegex := regexp.MustCompile(`svelte-check found (\d+) errors? and (\d+) warnings?`)
		for _, line := range lines {
			if matches := summaryRegex.FindStringSubmatch(line); matches != nil {
				errors, _ := strconv.Atoi(matches[1])
				warnings, _ := strconv.Atoi(matches[2])
				if errors == 0 && warnings == 0 {
					return nil, nil
				}
				break
			}
		}

		// Process error/warning lines
		for _, line := range lines {
			if line = strings.TrimSpace(line); line != "" && !isInformationalLine(line) {
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

func init() {
	Register(&svelteCheckChecker{
		name:           "npm-svelte-check",
		packageManager: "npm",
		lockFile:       "package-lock.json",
		executable:     "npm",
		command:        []string{"run", "svelte-check"},
		dependency:     "svelte-check",
	})
	Register(&svelteCheckChecker{
		name:           "bun-svelte-check",
		packageManager: "bun",
		lockFile:       "bun.lock",
		executable:     "bun",
		command:        []string{"run", "svelte-check"},
		dependency:     "svelte-check",
	})
}
