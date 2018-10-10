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
	readIo, err := client.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		return err
	}

	if readIo != nil {
		// Get output
		buf := new(bytes.Buffer)
		buf.ReadFrom(readIo)
		fmt.Printf("%s update response: %s", image, buf.String())
	}

	return nil
}

func removeImage(ctx context.Context, client *client.Client, image types.ImageSummary) (error) {
	delResp, err := client.ImageRemove(ctx, image.ID, types.ImageRemoveOptions{Force:true, PruneChildren:true})
	if err != nil {
		return err
	}
	for _, d := range delResp {
		fmt.Printf("Delete response %s untagged response %s", d.Deleted, d.Untagged)
	}
	return nil
}


func hostImages(ctx context.Context, client *client.Client) ([]types.ImageSummary, error ){
	images, err := client.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		return nil, err
	}
	return images, nil
}


func run() {
	var updateImages arrayFlags
	var interval int
	var prune bool
	var keepUntagged bool

	flag.Var(&updateImages, "update", "A list of images that should be monitored for update pulls")
	flag.IntVar(&interval, "interval", 10,
		"How often should the service check for image updates in minutes")
	flag.BoolVar(&prune, "prune", false,
		"Whether non listed images should be pruned/removed from the host")
	flag.BoolVar(&keepUntagged,"keep_untagged", true,
		"A flag on whether untagged/nontagged images should be kept on the host")
	flag.Parse()

	if updateImages.String() == "[]" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	for {
		fmt.Println("Monitor stage start ", time.Now().Format(time.RFC3339))
		fmt.Printf("Prune nonlisted images %v \n", prune)

		if prune {
			images, err := hostImages(ctx, cli)
			if err != nil {
				fmt.Printf("Failed to retrieve host images, (err): %v \n", err)
			}

			for _, i := range images {



				if err := removeImage(ctx, cli, i); err != nil {
					fmt.Printf("Failed to remove %v, (err): %v \n", i.ID, err)
				}
			}
		}

		fmt.Println("Checking for new images", time.Now().Format(time.RFC3339))
		for _, i := range updateImages {
			fmt.Printf("Checking %v for a new version \n", i)
			if err := updateImage(ctx, cli, i); err != nil {
				fmt.Printf("Failed to check %v for updates, (err): %v \n", i, err)
			}
		}

		fmt.Println("Monitor stage finished ", time.Now().Format(time.RFC3339))
		time.Sleep(time.Duration(interval) * time.Minute)
	}
}
