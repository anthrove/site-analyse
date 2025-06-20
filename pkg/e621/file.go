package e621

import (
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"os"
)

func DownloadFile(ctx context.Context, filename string, targetPath string) error {
	req, err := buildE6Request(fmt.Sprintf("/db_export/%s", filename))

	if err != nil {
		return err
	}

	req = req.WithContext(ctx)
	resp, err := httpClient.Do(req)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(targetPath)

	if err != nil {
		return err
	}

	defer out.Close()
	_, err = io.Copy(out, resp.Body)

	return err
}

func DownloadData(ctx context.Context, filename string, targetPath string) error {
	req, err := buildE6Request(fmt.Sprintf("/db_export/%s", filename))

	if err != nil {
		return err
	}

	req = req.WithContext(ctx)
	resp, err := httpClient.Do(req)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	uncompressedStream, err := gzip.NewReader(resp.Body)
	defer uncompressedStream.Close()

	if err != nil {
		return err
	}

	out, err := os.Create(targetPath)

	if err != nil {
		return err
	}

	defer out.Close()
	_, err = io.Copy(out, uncompressedStream)

	return err
}
