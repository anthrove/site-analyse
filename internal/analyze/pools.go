package analyze

import (
	"context"
	"github.com/InfluxCommunity/influxdb3-go/v2/influxdb3"
	"github.com/anthrove/site-analyse/pkg/object"
	"github.com/anthrove/site-analyse/pkg/util"
	"github.com/influxdata/line-protocol/v2/lineprotocol"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
	"time"
)

func Pools(ctx context.Context, influxClient *influxdb3.Client, fileName string) error {
	_, rightSide, _ := strings.Cut(fileName, "-")
	leftSide, _, _ := strings.Cut(rightSide, ".")

	date, err := time.Parse("2006-01-02", leftSide)

	if err != nil {
		return err
	}

	file, err := os.Open(fileName)

	if err != nil {
		return err
	}

	defer file.Close()

	poolsChan := util.GetStreamingData[object.Pools](ctx, file)

	poolCount := 0

	categories := map[string]int{}
	isActiveSize := map[string]int{}
	poolPostCountGroup := map[string]int{}

	updatedSize := 0
	createdSize := 0

	for pool := range poolsChan {
		poolCount++
		categories[pool.Category]++

		if pool.IsActive {
			isActiveSize["active"]++
		} else {
			isActiveSize["inactive"]++
		}

		postsStr := pool.PostIds
		postsStr = postsStr[1 : len(postsStr)-1]
		postsSize := len(strings.Split(postsStr, ","))

		if postsSize > 0 && postsSize <= 5 {
			poolPostCountGroup["1_5"]++
		} else if postsSize > 5 && postsSize <= 10 {
			poolPostCountGroup["6_10"]++
		} else if postsSize > 10 && postsSize <= 20 {
			poolPostCountGroup["11_20"]++
		} else if postsSize > 20 && postsSize <= 30 {
			poolPostCountGroup["21_30"]++
		} else if postsSize > 30 && postsSize <= 40 {
			poolPostCountGroup["31_40"]++
		} else if postsSize > 40 && postsSize <= 50 {
			poolPostCountGroup["41_50"]++
		} else if postsSize > 50 {
			poolPostCountGroup["50+"]++
		}

		if len(pool.CreatedAt) > 0 {
			updateTime, err := time.Parse("2006-01-02 15:04:05", pool.CreatedAt)

			if err != nil {
				log.WithField("pool_id", pool.ID).WithField("time", pool.CreatedAt).WithError(err).Error("Failed to parse created_at time")
				continue
			}

			// not the best method but it's an idea
			if updateTime.Add(time.Hour * 24).After(date) {
				createdSize++
			}
		}

		if len(pool.UpdatedAt) > 0 {
			updateTime, err := time.Parse("2006-01-02 15:04:05", pool.UpdatedAt)

			if err != nil {
				log.WithField("pool_id", pool.ID).WithField("time", pool.UpdatedAt).WithError(err).Error("Failed to parse updated time")
				continue
			}

			// not the best method but it's an idea
			if updateTime.Add(time.Hour * 24).After(date) {
				updatedSize++
			}
		}

	}

	totalSizePoint := influxdb3.NewPointWithMeasurement("pools_total").
		SetUIntegerField("total", uint64(poolCount)).
		SetTimestamp(date)

	poolCategoryConverted := make(map[string]interface{})
	for k, v := range categories {
		poolCategoryConverted[k] = v
	}
	poolCategoryPoint := influxdb3.NewPoint("pool_category", map[string]string{"site": "e621.net"}, poolCategoryConverted, date)

	poolStateConverted := make(map[string]interface{})
	for k, v := range isActiveSize {
		poolStateConverted[k] = v
	}
	poolStatePoint := influxdb3.NewPoint("pool_state", map[string]string{"site": "e621.net"}, poolStateConverted, date)

	poolPostCountConverted := make(map[string]interface{})
	for k, v := range poolPostCountGroup {
		poolPostCountConverted[k] = v
	}
	poolPostCountPoint := influxdb3.NewPoint("pool_post_count", map[string]string{"site": "e621.net"}, poolPostCountConverted, date)

	poolStatusPoint := influxdb3.NewPointWithMeasurement("pool_status").
		SetTag("site", "e621.net").
		SetUIntegerField("created", uint64(createdSize)).
		SetUIntegerField("updated", uint64(updatedSize)).
		SetTimestamp(date)

	err = influxClient.WritePoints(
		ctx,
		[]*influxdb3.Point{totalSizePoint, poolCategoryPoint, poolStatePoint, poolPostCountPoint, poolStatusPoint},
		influxdb3.WithPrecision(lineprotocol.Second),
		influxdb3.WithDefaultTags(map[string]string{"site": "e621.net"}))

	return err

}
