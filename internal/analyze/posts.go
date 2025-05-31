package analyze

import (
	"context"
	"github.com/InfluxCommunity/influxdb3-go/v2/influxdb3"
	"github.com/anthrove/site-analyse/pkg/object"
	"github.com/anthrove/site-analyse/pkg/util"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/publicsuffix"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

func Posts(ctx context.Context, influxClient *influxdb3.Client, fileName string) error {
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

	postsChan := util.GetStreamingData[object.Post](ctx, file)

	size := 0

	safeSize := 0
	questionableSize := 0
	explicitSize := 0

	deletedSize := 0
	pendingSize := 0
	updatedSize := 0

	fileExts := map[string]int{}
	sources := map[string]int{}
	sourceSize := map[int]int{}

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

		if post.IsPending {
			pendingSize++
		}
		fileExts[post.FileExt]++

		if len(post.Source) > 0 {
			sourceSize[len(strings.Split(post.Source, "\n"))]++
			for _, source := range strings.Split(post.Source, "\n") {
				parsedURL, err := url.Parse(source)
				if err != nil {
					// log.WithField("post_id", post.ID).WithField("source", source).Error(err)
					// to many invalid urls.... -.-
					sources["invalid"]++
					continue
				}

				host := parsedURL.Host
				if h, _, err := net.SplitHostPort(host); err == nil {
					host = h
				}

				// Extract the eTLD+1 (effective top-level domain + one label)
				domain, err := publicsuffix.EffectiveTLDPlusOne(host)
				if err != nil {
					domain = "invalid"
				}

				sources[domain]++
			}
		} else {
			sourceSize[0]++
		}

		if len(post.UpdatedAt) > 0 {
			updateTime, err := time.Parse("2006-01-02 15:04:05", post.UpdatedAt)

			if err != nil {
				log.WithField("post_id", post.ID).WithField("time", post.UpdatedAt).WithError(err).Error("Failed to parse updated time")
				continue
			}

			// not the best method but it's an idea
			if updateTime.Add(time.Hour * 24).After(time.Now()) {
				updatedSize++
			}
		}
	}

	totalSizePoint := influxdb3.NewPointWithMeasurement("posts_total").
		SetTag("site", "e621.net").
		SetUIntegerField("total", uint64(size)).
		SetTimestamp(date)

	postStatusPoint := influxdb3.NewPointWithMeasurement("posts_status").
		SetTag("site", "e621.net").
		SetUIntegerField("deleted", uint64(deletedSize)).
		SetIntegerField("pending", int64(pendingSize)).
		SetUIntegerField("updated", uint64(updatedSize)).
		SetTimestamp(date)

	postRatingPoint := influxdb3.NewPointWithMeasurement("post_rating_count").
		SetTag("site", "e621.net").
		SetUIntegerField("safe", uint64(safeSize)).
		SetUIntegerField("questionable", uint64(questionableSize)).
		SetUIntegerField("explicit", uint64(explicitSize)).
		SetTimestamp(date)

	// Convert to map[string]interface{}
	extConverted := make(map[string]interface{})
	for k, v := range fileExts {
		extConverted[k] = v
	}
	postExtsPoint := influxdb3.NewPoint("post_ext_count", map[string]string{"site": "e621.net"}, extConverted, date)

	sourceDomainConverted := make(map[string]interface{})
	for k, v := range sources {
		if v >= 100 {
			sourceDomainConverted[k] = v
		}
	}
	postDomainCountPoint := influxdb3.NewPoint("post_source_domain_count", map[string]string{"site": "e621.net"}, sourceDomainConverted, date)

	sourcePostCountConverted := make(map[string]interface{})
	for k, v := range sourceSize {
		sourcePostCountConverted[strconv.Itoa(k)] = v
	}
	postSourceCountPoint := influxdb3.NewPoint("post_source_count", map[string]string{"site": "e621.net"}, sourcePostCountConverted, date)

	err = influxClient.WritePoints(ctx, []*influxdb3.Point{
		totalSizePoint,
		postStatusPoint,
		postRatingPoint,
		postExtsPoint,
		postDomainCountPoint,
		postSourceCountPoint,
	})

	if err != nil {
		return err
	}

	return nil
}
