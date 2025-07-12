package checkers

import "sync"

type Issue struct {
	File    string
	Message string
}

type Checker interface {
	Name() string
	Runnable(workDir string) (bool, error)
	Check(workDir string) ([]Issue, error)
}

var (
	mu       sync.RWMutex
	checkers []Checker
)

func Register(c Checker) {
	mu.Lock()
	defer mu.Unlock()
	checkers = append(checkers, c)
}

func GetAll() []Checker {
	mu.RLock()
	defer mu.RUnlock()
	result := make([]Checker, len(checkers))
	copy(result, checkers)
	return result
}
