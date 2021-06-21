package main

import (
	"context"
	"fmt"
	"log"

	"go.temporal.io/sdk/client"

	"helloworld/app"
)

func main() {
	c, err := client.NewClient(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create temporal client", err)
	}

	defer c.Close()

	options := client.StartWorkflowOptions{
		ID:        "greeting-workflow",
		TaskQueue: app.GreetTaskQueue,
	}

	name := "Foobar"
	exec, err := c.ExecuteWorkflow(
		context.Background(), options, app.GreetWorkflow, name,
	)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}

	var result string
	if err := exec.Get(context.Background(), &result); err != nil {
		log.Fatalln("Unable to get workflow result", err)
	}

	fmt.Printf("\nWorkflow: %s\nRunID: %s\n%s\n", exec.GetID(), exec.GetRunID(), result)
}
