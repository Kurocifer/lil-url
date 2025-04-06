package utils

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"
)

func LoadURLs(storageFile string) (map[string]string, error) {
	urlMap := make(map[string]string)

	file, err := os.Open(storageFile)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return urlMap, nil
		} else {
			return nil, err
		}
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, " ", 2)
		if len(parts) == 2 {
			urlMap[parts[0]] = parts[1]
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return urlMap, nil
}

func SaveURLs(urlMap map[string]string, storageFile string) error {
	file, err := os.OpenFile(storageFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for key, value := range urlMap {
		value = AppendProtocol(value)
		_, err := fmt.Fprintf(writer, "%s %s\n", key, value)
		if err != nil {
			return err
		}
	}

	return writer.Flush()
}

func AppendProtocol(url string) string {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}

	return url
}
