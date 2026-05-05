package file_utils

import (
	"abstractions/logger"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func DownloadFileFromURL(ctx context.Context, URL string) ([]byte, error) {
	logger.Info(ctx, "Downloading file from url: "+URL)
	//Get the response bytes from the url
	response, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		logger.Error(ctx, err)
		return nil, errors.New("Received non 200 response code")
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.Error(ctx, err)
		return nil, err
	}

	logger.Info(ctx, "Downloaded file from URL - "+URL+" - successfully. Content Length = "+fmt.Sprintf("%d", len(responseBody)))

	return responseBody, nil
}
