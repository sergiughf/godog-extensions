package extension

import (
	"fmt"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/gherkin"
	"github.com/onsi/gomega"
)

// NewGomegaFailHandler registers gomega fail handler
func NewGomegaFailHandler(s *godog.Suite) {
	var failures []string

	gomega.RegisterFailHandler(func(message string, _ ...int) {
		failures = append(failures, message)
	})

	s.AfterStep(func(step *gherkin.Step, err error) {
		if err == nil {
			return
		}

		for _, failure := range failures {
			fmt.Printf("%s\n", failure)
		}

		failures = []string{}
	})
}
