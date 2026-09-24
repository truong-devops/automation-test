package automationtests

import (
	"context"
	"github.com/cucumber/godog"
)

func InitializeScenario(scenario *godog.ScenarioContext) {
	world := newWorld()

	scenario.Before( func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		world.resetScenarioState()

		err := world.connectMongo(ctx)
		if err != nil {
			return ctx, err
		}

		err = world.resetMongo(ctx)
		if err != nil {
			return ctx, err
		}

		return ctx, nil
	})

	scenario.After(func(ctx context.Context, sc *godog.Scenario, scenarioErr error) (context.Context, error) {
		err := world.closeMongo(ctx)

		return ctx, err
	})

	scenario.Given(`^the following orders exist:$`, world.theFollowingOrdersExist)

	scenario.Given(`^the following order request fixtures:$`, world.theFollowingOrderRequestFixtures)

	scenario.When(`^I send a (POST|PUT) request to "([^"]*)" using fixture "([^"]*)"$`, world.sendRequestUsingFixture)

	scenario.When(`^I send a (GET|DELETE) request to "([^"]*)"$`, world.sendRequestWithoutBody)

	scenario.Then(`^the response status should be (\d+)$`, world.responseStatusShouldBe)

	scenario.Then(`^the response field "([^"]*)" should equal "([^"]*)"$`, world.responseFieldShouldEqual)
}
