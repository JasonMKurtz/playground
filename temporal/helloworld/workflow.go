package app

import (
	"go.temporal.io/sdk/workflow"
)

func GreetWorkflow(ctx workflow.Context, name string) (string, error) {
	var result string
	err := workflow.ExecuteActivity(ctx, Greet, name).Get(ctx, &result)
	return result, err
}
