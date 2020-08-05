package extension

import (
	"fmt"

	"github.com/cucumber/godog"
	"github.com/onsi/gomega"
)

// NewGomegaFailHandler registers gomega fail handler
func NewGomegaFailHandler(ctx *godog.ScenarioContext) {
	var failures []string

	gomega.RegisterFailHandler(func(message string, _ ...int) {
		failures = append(failures, message)
	})

	ctx.AfterStep(func(step *godog.Step, err error) {
		if err == nil {
			return
		}

		for _, failure := range failures {
			fmt.Printf("%s\n", failure)
		}

		failures = []string{}
	})
}
