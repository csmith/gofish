package checkers

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type Requirement func() error

func requires(reqs ...Requirement) (bool, error) {
	for _, req := range reqs {
		if err := req(); err != nil {
			return false, err
		}
	}
	return true, nil
}

func executableOnPath(name string) Requirement {
	return func() error {
		if _, err := exec.LookPath(name); err != nil {
			return fmt.Errorf("%s not found in PATH", name)
		}
		return nil
	}
}

func fileInDir(workDir, filename string) Requirement {
	return func() error {
		if _, err := os.Stat(filepath.Join(workDir, filename)); os.IsNotExist(err) {
			return fmt.Errorf("%s not found in %s", filename, workDir)
		}
		return nil
	}
}

func goProjectInDir(workDir string) Requirement {
	return fileInDir(workDir, "go.mod")
}

func jsProjectInDir(workDir string) Requirement {
	return fileInDir(workDir, "package.json")
}

func packageHasDependency(workDir, packageName string) Requirement {
	return func() error {
		packageJSONPath := filepath.Join(workDir, "package.json")
		data, err := os.ReadFile(packageJSONPath)
		if err != nil {
			return fmt.Errorf("failed to read package.json: %w", err)
		}

		var packageJSON struct {
			Dependencies    map[string]string `json:"dependencies"`
			DevDependencies map[string]string `json:"devDependencies"`
		}

		if err := json.Unmarshal(data, &packageJSON); err != nil {
			return fmt.Errorf("failed to parse package.json: %w", err)
		}

		if packageJSON.Dependencies != nil {
			if _, exists := packageJSON.Dependencies[packageName]; exists {
				return nil
			}
		}

		if packageJSON.DevDependencies != nil {
			if _, exists := packageJSON.DevDependencies[packageName]; exists {
				return nil
			}
		}

		return fmt.Errorf("package %s not found in dependencies or devDependencies", packageName)
	}
}
