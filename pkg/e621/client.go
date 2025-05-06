package e621

import (
	"fmt"
	"net/http"
)

var httpClient http.Client

func init() {
	httpClient = http.Client{}
}

func buildE6Request(url string) (*http.Request, error) {
	request, err := http.NewRequest("GET", fmt.Sprintf("%s%s", "https://e621.net", url), nil)
	if err != nil {
		return nil, err
	}

	//TODO thinking about adding user agent

	return request, nil
}
