package analyze

import (
	"context"
	"github.com/anthrove/site-analyse/pkg/object"
	"github.com/anthrove/site-analyse/pkg/util"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"os"
)

var (
	TagsTotal = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "tags_total",
		Help: "Total number of tags processed",
	})

	TagsCategory = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "tags_category_count",
		Help: "Number of tags per category",
	}, []string{"category"})

	TagsPostCountBuckets = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "tags_postcount_bucket",
		Help: "Number of tags per post count bucket",
	}, []string{"bucket"})
)

func Tags(ctx context.Context, promPusher *push.Pusher, fileName string) error {
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

	TagsTotal.Set(float64(size))

	TagsCategory.WithLabelValues("general").Set(float64(generalSize))
	TagsCategory.WithLabelValues("artist").Set(float64(artistSize))
	TagsCategory.WithLabelValues("contributor").Set(float64(contributorSize))
	TagsCategory.WithLabelValues("copyright").Set(float64(copyrightSize))
	TagsCategory.WithLabelValues("character").Set(float64(characterSize))
	TagsCategory.WithLabelValues("species").Set(float64(speciesSize))
	TagsCategory.WithLabelValues("invalid").Set(float64(invalidSize))
	TagsCategory.WithLabelValues("meta").Set(float64(metaSize))
	TagsCategory.WithLabelValues("lore").Set(float64(loreSize))
	TagsCategory.WithLabelValues("unknown").Set(float64(unknownSize))

	TagsPostCountBuckets.WithLabelValues("0").Set(float64(tagSizeEmpty))
	TagsPostCountBuckets.WithLabelValues("1_9").Set(float64(tagSize0_9))
	TagsPostCountBuckets.WithLabelValues("11_99").Set(float64(tagSize10_99))
	TagsPostCountBuckets.WithLabelValues("101_999").Set(float64(tagSize100_999))
	TagsPostCountBuckets.WithLabelValues("1001_9999").Set(float64(tagSize1000_9999))
	TagsPostCountBuckets.WithLabelValues("10001_99999").Set(float64(tagSize10000_99999))
	TagsPostCountBuckets.WithLabelValues("100001_999999").Set(float64(tagSize100000_999999))
	TagsPostCountBuckets.WithLabelValues("xxl").Set(float64(tagSizeXXL))

	promPusher.Collector(TagsTotal)
	promPusher.Collector(TagsCategory)
	promPusher.Collector(TagsPostCountBuckets)
	return nil
}
