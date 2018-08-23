package main


import (
	"context"
	client2 "github.com/docker/docker/client"
	"github.com/docker/docker-ce/components/engine/client"
)


func main() {
	ctx := context.Background()
	client := client.NewClientWithOpts()

	client.ImagePull()

}