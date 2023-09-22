package testutil

import (
	"regexp"
	"testing"
)

type Task interface {
	Do(t testing.TB) error
}

func Run(t testing.TB, tasks []Task) {
	for i, task := range tasks {
		if err := task.Do(t); err != nil {
			t.Fatalf("run %d task: %v", i, err)
		}
	}
}

func MatchError(t testing.TB, err error, errstr string) bool {
	switch {
	case err == nil && errstr == "":
		return true
	case err != nil && errstr == "":
		return false
	case err == nil && errstr != "":
		return false
	case err != nil && errstr != "":
		matched, e := regexp.MatchString(errstr, err.Error())
		if e != nil {
			t.Fatalf("match: regexp match string: %v", e)
		}
		return matched
	}
	panic("never reach")
}

func RecoverPanic(f func()) {
	defer func() {
		recover()
	}()
	f()
}
