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

func updateImage(ctx context.Context, client *client.Client, image string, debug bool) error {
	readIo, err := client.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		return err
	}

	if readIo != nil {
		// Get output
		buf := new(bytes.Buffer)
		buf.ReadFrom(readIo)
		if debug {
			fmt.Printf("%s update response: %s \n", image, buf.String())
		}
	}

	return nil
}

func removeImage(ctx context.Context, client *client.Client, image types.ImageSummary, debug bool) error {
	delResp, err := client.ImageRemove(ctx, image.ID, types.ImageRemoveOptions{Force: true, PruneChildren: true})
	if err != nil {
		return err
	}
	for _, d := range delResp {
		if debug {
			fmt.Printf("Delete response %s untagged response %s \n", d.Deleted, d.Untagged)
		}
	}
	return nil
}

func usedImage(ctx context.Context, client *client.Client, image types.ImageSummary, debug bool) (bool, error) {
	containers, err := client.ContainerList(ctx, types.ContainerListOptions{All: true});
	if err != nil {
		return false, err
	}

	for _, container := range containers {
		if container.ImageID == image.ID {
			if debug {
				fmt.Printf("Image: %v is being used by container: %v \n", image.ID, container.ID)
			}
			return true, nil
		}
	}
	return false, nil
}


func hostImages(ctx context.Context, client *client.Client) ([]types.ImageSummary, error) {
	images, err := client.ImageList(ctx, types.ImageListOptions{All: true})
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
	// debug defines whether the 
	var debug bool

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
	flag.BoolVar(&debug, "debug", false,
		"Set the debug flag to run the updater in debug mode")

	flag.Parse()

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
		if debug {
			fmt.Println("Monitor stage start ", time.Now().Format(time.RFC3339))
		}
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
					beingUsed := false
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
					
					if used, err := usedImage(ctx, cli, i, debug); err != nil {
						fmt.Printf("%v failed to check if an image is used, (err): %v \n", i.ID, err)
					} else {
						beingUsed = used
					}

					if protected && debug {
						fmt.Printf("%v wont be pruned since it is protected \n", tag)
					}

					if !beingUpdated && !protected && !beingUsed {
						if debug {
							fmt.Printf("Pruning %v \n", i.ID)
						}
						if err := removeImage(ctx, cli, i, debug); err != nil {
							fmt.Printf("Failed to remove %v, (err): %v \n", i.ID, err)
						}
					}
				}
				// Remove untagged
				if pruneUntagged {
					if len(i.RepoTags) == 0 {
						if err := removeImage(ctx, cli, i, debug); err != nil {
							fmt.Printf("Failed to remove %v, (err): %v \n", i.ID, err)
						}
					}
				}
			}
		}
		
		fmt.Println("Checking for new images", time.Now().Format(time.RFC3339))
		for k := range updateImages {
			if debug {
				fmt.Printf("Checking %v for a new version \n", k)
			}
			if err := updateImage(ctx, cli, k, debug); err != nil {
				fmt.Printf("Failed to check %v for updates, (err): %v \n", k, err)
			}
		}

		fmt.Println("Update stage finished ", time.Now().Format(time.RFC3339))
		time.Sleep(time.Duration(interval) * time.Minute)
	}
}
