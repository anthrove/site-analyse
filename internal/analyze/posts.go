package analyze

import (
	"context"
	"github.com/anthrove/site-analyse/pkg/object"
	"github.com/anthrove/site-analyse/pkg/util"
	"github.com/prometheus/client_golang/prometheus/push"
	log "github.com/sirupsen/logrus"
	"os"
)

func Posts(ctx context.Context, promPusher *push.Pusher, fileName string) error {
	file, err := os.Open(fileName)

	if err != nil {
		return err
	}

	defer file.Close()

	postsChan := util.GetStreamingData[object.Post](ctx, file)

	size := 0

	safeSize := 0
	questionableSize := 0
	explicitSize := 0

	deletedSize := 0

	fileExts := map[string]int{}

	for post := range postsChan {
		size++

		switch post.Rating {
		case "e":
			explicitSize++
		case "q":
			questionableSize++
		case "s":
			safeSize++
		}

		if post.IsDeleted {
			deletedSize++
		}

		fileExts[post.FileExt]++
	}

	log.Infof("Size: %d", size)
	log.Infof("Deleted: %d", deletedSize)

	return nil
}
