package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

func updateImage(ctx context.Context, client *client.Client, image string) error {
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

func removeImage(ctx context.Context, client *client.Client, image types.ImageSummary) error {
	delResp, err := client.ImageRemove(ctx, image.ID, types.ImageRemoveOptions{Force: true, PruneChildren: true})
	if err != nil {
		return err
	}
	for _, d := range delResp {
		fmt.Printf("Delete response %s untagged response %s \n", d.Deleted, d.Untagged)
	}
	return nil
}

func hostImages(ctx context.Context, client *client.Client) ([]types.ImageSummary, error) {
	images, err := client.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		return nil, err
	}
	return images, nil
}

func run() {
	var updateImages arrayFlags
	var protectImages arrayFlags
	var interval int
	var prune bool
	var pruneUntagged bool

	flag.Var(&updateImages, "update",
		"A list of images that should be monitored for update pulls")
	flag.Var(&protectImages, "protect",
		"A list of images that should not be pruned on the hosts")
	flag.IntVar(&interval, "interval", 10,
		"How often should the service check for image updates in minutes")
	flag.BoolVar(&prune, "prune", false,
		"Whether non update images should be pruned/removed from the host")
	flag.BoolVar(&pruneUntagged, "prune-untagged", false,
		"Requires prune, Whether untagged/nontagged images should be kept on the host")

	flag.Parse()

	fmt.Printf(updateImages.String())
	if updateImages.String() == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Checking update for: %v\n", updateImages)
	fmt.Printf("Pruning images not in: %v or %s \n", updateImages, protectImages)
	for {
		fmt.Println("Monitor stage start ", time.Now().Format(time.RFC3339))
		// Prune non-monitored images
		if prune {
			images, err := hostImages(ctx, cli)
			if err != nil {
				fmt.Printf("Failed to retrieve host images, (err): %v \n", err)
			}

			for _, i := range images {
				for _, tag := range i.RepoTags {
					beingUpdated := false
					protected := false
					if _, ok := updateImages[tag]; ok {
						beingUpdated = true
						// Check weather it matches without image sha
					} else if _, ok := updateImages[tag[0:strings.IndexByte(tag, ':')]]; ok {
						beingUpdated = true
					}

					if _, ok := protectImages[tag]; ok {
						protected = true
						// Check weather it matches without image sha
					} else if _, ok := protectImages[tag[0:strings.IndexByte(tag, ':')]]; ok {
						protected = true
					}

					if !beingUpdated && !protected {
						fmt.Printf("Pruning %v \n", i.ID)
						if err := removeImage(ctx, cli, i); err != nil {
							fmt.Printf("Failed to remove %v, (err): %v \n", i.ID, err)
						}
					}
				}
				// Remove untagged
				if pruneUntagged {
					if len(i.RepoTags) == 0 {
						if err := removeImage(ctx, cli, i); err != nil {
							fmt.Printf("Failed to remove %v, (err): %v \n", i.ID, err)
						}
					}
				}
			}
		}

		fmt.Println("Checking for new images", time.Now().Format(time.RFC3339))
		for k := range updateImages {
			fmt.Printf("Checking %v for a new version \n", k)
			if err := updateImage(ctx, cli, k); err != nil {
				fmt.Printf("Failed to check %v for updates, (err): %v \n", k, err)
			}
		}

		fmt.Println("Monitor stage finished ", time.Now().Format(time.RFC3339))
		time.Sleep(time.Duration(interval) * time.Minute)
	}
}
