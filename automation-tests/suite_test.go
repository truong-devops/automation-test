package automationtests

import (
	"github.com/cucumber/godog"
	"testing"
)

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		Name:                "orders-api",
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:      "pretty",
			Paths:       []string{"features"},
			Strict:      true,
			Concurrency: 1,
			TestingT:    t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("Godog feature tests failed")
	}

}
