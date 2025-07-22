package checkers

import (
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
