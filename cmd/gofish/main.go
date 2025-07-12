package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/csmith/gofish/checkers"
)

func findProjectDirectories() ([]string, error) {
	var dirs []string

	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			// Skip known-bad directories and any directory starting with "." (except root)
			if (info.Name()[0] == '.' && path != ".") || info.Name() == "node_modules" || info.Name() == "__pycache__" {
				return filepath.SkipDir
			}

			if hasFile(path, "go.mod") || hasFile(path, "package.json") {
				dirs = append(dirs, path)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if len(dirs) == 0 {
		dirs = append(dirs, ".")
	}

	return dirs, nil
}

func hasFile(dir, filename string) bool {
	_, err := os.Stat(filepath.Join(dir, filename))
	return err == nil
}

func main() {
	showChecks := flag.Bool("checks", false, "show available checks and their status")
	flag.Parse()

	dirs, err := findProjectDirectories()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding project directories: %v\n", err)
		os.Exit(1)
	}

	if *showChecks {
		fmt.Println("Available checks:")
		fmt.Println()
		for _, dir := range dirs {
			fmt.Printf("Directory: %s\n", dir)
			for _, checker := range checkers.GetAll() {
				runnable, err := checker.Runnable(dir)
				if runnable && err == nil {
					fmt.Printf("  ✓ %s: ready to run\n", checker.Name())
				} else if err != nil {
					fmt.Printf("  ✗ %s: %v\n", checker.Name(), err)
				} else {
					fmt.Printf("  ✗ %s: not runnable in this directory\n", checker.Name())
				}
			}
			fmt.Println()
		}
		return
	}

	var hasIssues bool
	var allIssues []struct {
		checker string
		dir     string
		issue   checkers.Issue
	}

	for _, dir := range dirs {
		for _, checker := range checkers.GetAll() {
			runnable, err := checker.Runnable(dir)
			if err != nil || !runnable {
				continue
			}

			issues, err := checker.Check(dir)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s in %s: error: %v\n", checker.Name(), dir, err)
				os.Exit(1)
			}

			for _, issue := range issues {
				allIssues = append(allIssues, struct {
					checker string
					dir     string
					issue   checkers.Issue
				}{checker.Name(), dir, issue})
				hasIssues = true
			}
		}
	}

	if hasIssues {
		fmt.Fprintln(os.Stderr, "The following issues were detected by gofish, please address them:")
		fmt.Fprintln(os.Stderr)
		for _, item := range allIssues {
			if item.dir != "." {
				fmt.Fprintf(os.Stderr, "%s in %s: %s\n", item.checker, item.dir, item.issue.Message)
			} else {
				fmt.Fprintf(os.Stderr, "%s: %s\n", item.checker, item.issue.Message)
			}
		}
		os.Exit(2)
	}
}
