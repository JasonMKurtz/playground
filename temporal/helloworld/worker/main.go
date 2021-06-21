package main

import (
	"log"

	"helloworld/app"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	client, err := client.NewClient(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create temporal client", err)
	}

	defer client.Close()

	w := worker.New(client, app.GreetTaskQueue, worker.Options{})
	w.RegisterWorkflow(app.GreetWorkflow)
	w.RegisterActivity(app.Greet)

	if err = w.Run(worker.InterruptCh()); err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
