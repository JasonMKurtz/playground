package app

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

func GreetWorkflow(ctx workflow.Context, name string) (string, error) {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 5,
	})

	var result string
	err := workflow.ExecuteActivity(ctx, Greet, name).Get(ctx, &result)
	return result, err
}
