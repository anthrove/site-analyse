package analyze

import (
	"context"
	"github.com/InfluxCommunity/influxdb3-go/v2/influxdb3"
	"github.com/anthrove/site-analyse/pkg/object"
	"github.com/anthrove/site-analyse/pkg/util"
	"os"
	"strings"
	"time"
)

func Tags(ctx context.Context, influxClient *influxdb3.Client, fileName string) error {
	rightSide := strings.SplitN(fileName, "-", 2)[1]
	leftSide := strings.SplitN(rightSide, ".", 2)[0]

	date, err := time.Parse("2006-01-02", leftSide)

	if err != nil {
		return err
	}

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

	tagsTotalPoint := influxdb3.NewPointWithMeasurement("tags_total").
		SetTag("site", "e621.net").
		SetIntegerField("total", int64(size)).
		SetTimestamp(date)
	tagsCategoryPoint := influxdb3.NewPointWithMeasurement("tags_category_count").
		SetTag("site", "e621.net").
		SetIntegerField("general", int64(generalSize)).
		SetIntegerField("artist", int64(artistSize)).
		SetIntegerField("contributor", int64(contributorSize)).
		SetIntegerField("copyright", int64(copyrightSize)).
		SetIntegerField("character", int64(characterSize)).
		SetIntegerField("species", int64(speciesSize)).
		SetIntegerField("invalid", int64(invalidSize)).
		SetIntegerField("meta", int64(metaSize)).
		SetIntegerField("lore", int64(loreSize)).
		SetIntegerField("unknown", int64(unknownSize)).
		SetTimestamp(date)
	tagsPostCountPoint := influxdb3.NewPointWithMeasurement("tags_post_count").
		SetTag("site", "e621.net").
		SetUIntegerField("0", uint64(tagSizeEmpty)).
		SetUIntegerField("1_9", uint64(tagSize0_9)).
		SetUIntegerField("11_99", uint64(tagSize10_99)).
		SetUIntegerField("101_999", uint64(tagSize100_999)).
		SetUIntegerField("1001_9999", uint64(tagSize1000_9999)).
		SetUIntegerField("10001_99999", uint64(tagSize10000_99999)).
		SetUIntegerField("100001_999999", uint64(tagSize100000_999999)).
		SetUIntegerField("xxl", uint64(tagSizeXXL)).
		SetTimestamp(date)

	err = influxClient.WritePoints(ctx, []*influxdb3.Point{
		tagsTotalPoint,
		tagsCategoryPoint,
		tagsPostCountPoint,
	})

	if err != nil {
		return err
	}

	return nil
}
