package input

import (
	"SDR-Labo4/src/utils/log"
	"bufio"
	"fmt"
	"os"
	"strings"
)

type checkFunc[T any] struct {
	check func(T) bool
	error string
}

type Input[T any] struct {
	label        string
	errorMessage string
	checks       []checkFunc[T]
	read         func() (T, error)
}

func (i Input[T]) AddCheck(check func(T) bool, errorMessage string) Input[T] {
	i.checks = append(i.checks, checkFunc[T]{check, errorMessage})
	return i
}

func (i Input[T]) isValid(value T) (bool, error) {
	for _, check := range i.checks {
		if !check.check(value) {
			return false, fmt.Errorf(check.error)
		}
	}
	return true, nil
}

func (i Input[T]) Read() T {
	for {
		fmt.Print(i.label)
		value, err := i.read()
		if err != nil {
			log.Logf(log.Error, "Error reading input: %s", err)
			continue
		}
		if valid, err := i.isValid(value); !valid {
			log.Logf(log.Error, "Error validating input: %s", err)
			continue
		}
		return value
	}
}

func StringInput(label string, args ...interface{}) Input[string] {
	return Input[string]{
		label:        fmt.Sprintf(label, args...),
		errorMessage: "Please enter a valid string",
		checks:       []checkFunc[string]{},
		read: func() (string, error) {
			reader := bufio.NewReader(os.Stdin)
			text, err := reader.ReadString('\n')
			return strings.TrimSpace(text), err
		},
	}
}

func BasicInput[T any](label string, args ...interface{}) Input[T] {
	return Input[T]{
		label:        fmt.Sprintf(label, args...),
		errorMessage: "Please enter a valid input",
		checks:       []checkFunc[T]{},
		read: func() (T, error) {
			var value T
			reader := bufio.NewReader(os.Stdin)
			_, err := fmt.Fscanln(reader, &value)
			return value, err
		},
	}
}
