package checkers

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

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

func packageHasScript(workDir, scriptName string) Requirement {
	return func() error {
		packageJSONPath := filepath.Join(workDir, "package.json")
		data, err := os.ReadFile(packageJSONPath)
		if err != nil {
			return fmt.Errorf("failed to read package.json: %w", err)
		}

		var packageJSON struct {
			Scripts map[string]string `json:"scripts"`
		}

		if err := json.Unmarshal(data, &packageJSON); err != nil {
			return fmt.Errorf("failed to parse package.json: %w", err)
		}

		if packageJSON.Scripts != nil {
			if _, exists := packageJSON.Scripts[scriptName]; exists {
				return nil
			}
		}

		return fmt.Errorf("script %s not found in package.json", scriptName)
	}
}
