package util

import (
	"compress/gzip"
	"context"
	"encoding/csv"
	"errors"
	"github.com/anthrove/site-analyse/pkg/e621"
	"github.com/minio/minio-go/v7"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"reflect"
	"strconv"
	"time"
)

func DownloadE6File(ctx context.Context, minioClient *minio.Client, bucketName string, filename string) (string, error) {

	if stat, err := os.Stat(filename); err == nil {
		// path/to/whatever exists
		if stat.IsDir() {
			return "", errors.New("filename is a directory")
		} else {
			return filename, nil
		}
	}

	err := minioClient.FGetObject(ctx, bucketName, filename, filename, minio.GetObjectOptions{})

	if err == nil {
		return "", nil
	}
	// Problem downloading file from s3 storage

	log.WithError(err).WithField("filename", filename).Warn("Failed to download file from s3 storage")

	err = e621.DownloadFile(ctx, filename, filename)
	if err != nil {
		return "", err
	}

	_, err = minioClient.FPutObject(ctx, bucketName, filename, filename, minio.PutObjectOptions{})

	if err != nil {
		return "", err
	}

	return filename, nil
}

func GetStreamingData[T any](ctx context.Context, rc io.Reader) chan T {
	ch := make(chan T)
	go func() {
		inputChan := make(chan []string)
		uncompressedStream, err := gzip.NewReader(rc)
		defer uncompressedStream.Close()

		if err != nil {
			panic(err)
		}

		r := csv.NewReader(uncompressedStream)
		var header []string
		if header, err = r.Read(); err != nil {
			log.Fatal(err)
		}
		defer close(inputChan)
		go func() {
			defer close(ch)
			returnChannel := parseRecord[T](header, inputChan)
			for data := range returnChannel {
				ch <- data
			}
		}()

		for {
			rec, err := r.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Fatal(err)
			}

			if len(rec) == 0 {
				continue
			}

			inputChan <- rec
		}
	}()
	return ch
}

func parseRecord[T any](header []string, input chan []string) chan T {
	channel := make(chan T)
	go func() {
		defer close(channel)
		var e T
		et := reflect.TypeOf(e)
		var headers = make(map[string]int, et.NumField())
		for i := 0; i < et.NumField(); i++ {
			headers[et.Field(i).Name] = func(element string, array []string) int {
				for k, v := range array {
					if v == element {
						return k
					}
				}
				return -1
			}(et.Field(i).Tag.Get("csv"), header)
		}
		for record := range input {
			if len(record) == 0 {
				continue
			}

			for h, i := range headers {
				if i == -1 {
					continue
				}
				elem := reflect.ValueOf(&e).Elem()
				field := elem.FieldByName(h)
				if field.CanSet() {
					switch field.Type().Name() {
					case "bool":
						a, _ := strconv.ParseBool(record[i])
						field.Set(reflect.ValueOf(a))
					case "int":
						a, _ := strconv.Atoi(record[i])
						field.Set(reflect.ValueOf(a))
					case "float64":
						a, _ := strconv.ParseFloat(record[i], 64)
						field.Set(reflect.ValueOf(a))
					case "Time":
						a, _ := time.Parse("2006-01-02T00:00:00Z", record[i])
						field.Set(reflect.ValueOf(a))
					case "string":
						field.Set(reflect.ValueOf(record[i]))
					default:
						log.Printf("Unknown Fieldtype: %s\n", field.Type().Name())
						field.Set(reflect.ValueOf(record[i]))
					}
				}
			}
			channel <- e
		}
		log.Info("parsing ended")
	}()
	return channel
}
