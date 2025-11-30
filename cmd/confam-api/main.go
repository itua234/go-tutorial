package main

import (
	app "confam-api/internal/app"
	"context"
	"log"
)

func main() {
	ctx := context.Background()

	application, err := app.New(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}
	//Ensures resources (like Redis) are closed when main exits
	defer application.Close()

	//Run the Server
	if err := application.Run(); err != nil {
		log.Fatalf("Server crashed: %v", err)
	}

}
