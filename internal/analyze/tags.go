package analyze

import (
	"context"
	"github.com/anthrove/site-analyse/pkg/object"
	"github.com/anthrove/site-analyse/pkg/util"
	log "github.com/sirupsen/logrus"
	"os"
)

func Tags(ctx context.Context, fileName string) error {
	//date := fileName[5 : len(fileName)-4]

	file, err := os.Open(fileName)

	if err != nil {
		return err
	}

	defer file.Close()

	tagsChan := util.GetStreamingData[object.Tag](ctx, file)

	size := 0

	generalSize := 0
	artistSize := 0
	contributorSize := 0
	copyrightSize := 0
	characterSize := 0
	speciesSize := 0
	invalidSize := 0
	metaSize := 0
	loreSize := 0
	unknownSize := 0

	tagSizeEmpty := 0
	tagSize0_9 := 0
	tagSize10_99 := 0
	tagSize100_999 := 0
	tagSize1000_9999 := 0
	tagSize10000_99999 := 0
	tagSize100000_999999 := 0
	tagSizeXXL := 0

	for tag := range tagsChan {
		size++

		switch tag.Category {
		case 0:
			generalSize++
		case 1:
			artistSize++
		case 2:
			contributorSize++
		case 3:
			copyrightSize++
		case 4:
			characterSize++
		case 5:
			speciesSize++
		case 6:
			invalidSize++
		case 7:
			metaSize++
		case 8:
			loreSize++
		default:
			unknownSize++
		}

		if tag.PostCount == 0 {
			tagSizeEmpty++
		} else if tag.PostCount > 0 && tag.PostCount < 10 {
			tagSize0_9++
		} else if tag.PostCount > 10 && tag.PostCount < 100 {
			tagSize10_99++
		} else if tag.PostCount > 100 && tag.PostCount < 1000 {
			tagSize100_999++
		} else if tag.PostCount > 1000 && tag.PostCount < 10000 {
			tagSize1000_9999++
		} else if tag.PostCount > 10000 && tag.PostCount < 100000 {
			tagSize10000_99999++
		} else if tag.PostCount > 100000 && tag.PostCount < 10000000 {
			tagSize100000_999999++
		} else {
			tagSizeXXL++
		}
	}

	log.WithField("total_tags", size).Info("Total tags processed")
	log.WithField("general_tags", generalSize).Info("General category tags count")
	log.WithField("artist_tags", artistSize).Info("Artist category tags count")
	log.WithField("contributor_tags", contributorSize).Info("Contributor category tags count")
	log.WithField("copyright_tags", copyrightSize).Info("Copyright category tags count")
	log.WithField("character_tags", characterSize).Info("Character category tags count")
	log.WithField("species_tags", speciesSize).Info("Species category tags count")
	log.WithField("invalid_tags", invalidSize).Info("Invalid category tags count")
	log.WithField("meta_tags", metaSize).Info("Meta category tags count")
	log.WithField("lore_tags", loreSize).Info("Lore category tags count")
	log.WithField("unknown_tags", unknownSize).Info("Unknown category tags count")

	log.WithField("postcount_0", tagSizeEmpty).Info("Tags with PostCount == 0")
	log.WithField("postcount_1_9", tagSize0_9).Info("Tags with PostCount between 1 and 9")
	log.WithField("postcount_11_99", tagSize10_99).Info("Tags with PostCount between 11 and 99")
	log.WithField("postcount_101_999", tagSize100_999).Info("Tags with PostCount between 101 and 999")
	log.WithField("postcount_1001_9999", tagSize1000_9999).Info("Tags with PostCount between 1001 and 9999")
	log.WithField("postcount_10001_99999", tagSize10000_99999).Info("Tags with PostCount between 10001 and 99999")
	log.WithField("postcount_100001_999999", tagSize100000_999999).Info("Tags with PostCount between 100001 and 999999")
	log.WithField("postcount_xxl", tagSizeXXL).Info("Tags with PostCount outside defined ranges (XXL)")

	return nil
}
