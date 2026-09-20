package automationtests

import (
	"context"
	"github.com/cucumber/godog"
)

func InitializeScenario(scenario *godog.ScenarioContext) {
	world := newWorld()

	scenario.Given(
		`^the following orders exist:$`, world.theFollowingOrdersExist,
	)

	scenario.When(
		`^I send a (GET|DELETE) request to "([^"]*)"$`, world.sendRequestWithoutBody,
	)

	scenario.Then(
		`^the response status should be (\d+)$`, world.responseStatusShouldBe,
	)

	scenario.Then(
		`^the response field "([^"]*)" should equal "([^"]*)"$`, world.responseFieldShouldEqual,
	)
}

func (w *World) theFollowingOrdersExist(
	table *godog.Table,
) error {
	return godog.ErrPending
}

func (w *World) sendRequestWithoutBody(
	ctx context.Context,
	method string,
	path string,
) error {
	return godog.ErrPending
}

func (w *World) responseStatusShouldBe(
	expected int,
) error {
	return godog.ErrPending
}

func (w *World) responseFieldShouldEqual(
	fieldPath string,
	expected string,
) error {
	return godog.ErrPending
}
