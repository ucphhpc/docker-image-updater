package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
	"os"
	"time"
)

func updateImage(ctx context.Context, client *client.Client, image string) (error) {
	readio, err := client.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		return err
	}

	if readio != nil {
		// Get output
		buf := new(bytes.Buffer)
		buf.ReadFrom(readio)
		fmt.Printf("%s update response: %s", image, buf.String())
	}

	return nil
}

func run() {
	var imageFlags arrayFlags
	var updateInterval int

	flag.Var(&imageFlags, "update", "A list of images that should be monitored for update pulls")
	flag.IntVar(&updateInterval, "interval", 10, "How often should the service check for image updates in minutes")
	flag.Bool("prune", false, "Whether non listed image should be pruned/removed from the host")
	flag.Bool("keep_untagged", true, "A flag on whether untagged/nontagged images should be kept on the host")
	flag.Parse()

	if imageFlags.String() == "[]" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	for {
		fmt.Println("Checking for new images",
			time.Now().Format(time.RFC3339))
		for _, i := range imageFlags {
			fmt.Printf("Checking %s for a new version \n", i)
			if err := updateImage(ctx, cli, i); err != nil {
				fmt.Printf("Failed to check %s for updates, (err): %s \n", i, err)
			}
		}
		time.Sleep(time.Duration(updateInterval) * time.Minute)
	}
}
