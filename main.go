package main

import (
	"context"
	"service/container"
)

func main() {
	ctx := context.Background()
	container.StartApp(ctx)
}
