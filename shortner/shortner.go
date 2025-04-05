package shortner

import (
	"bufio"
	"errors"
	"fmt"
	"hash/crc32"
	"io/fs"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type Shortner struct {
	storageFile string
}

var baseURL = "http://lilurl"

func NewShortner(storageFile string) *Shortner {
	return &Shortner{
		storageFile: storageFile,
	}
}

func (s *Shortner) Shorten(longURL string) (string, error) {
	urlMap, err := s.loadURLs()
	if err != nil {
		return "", err
	}

	shortURL := compressURL(longURL)
	if _, exists := urlMap[shortURL]; exists {
		return shortURL, nil
	}
	urlMap[shortURL] = longURL
	err = s.saveURLs(urlMap)
	if err != nil {
		return "", fmt.Errorf("encountered an error while saving new link: %w", err)
	}

	return shortURL, nil
}

func (s *Shortner) LookupURL(shortURL string) error {
	urlMap, err := s.loadURLs()
	if err != nil {
		return fmt.Errorf("an error occured while loading URLs: %s", err)
	}

	longURL, exists := urlMap[shortURL]
	if !exists {
		return fmt.Errorf("URL %s does not map to any known URL", shortURL)
	}

	fmt.Printf("Original URL: %s ", longURL)

	if prompt("Open URL in broswr ?") {
		err := openInBrowser(longURL)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Shortner) loadURLs() (map[string]string, error) {
	urlMap := make(map[string]string)

	file, err := os.Open(s.storageFile)
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

func (s *Shortner) saveURLs(urlMap map[string]string) error {
	file, err := os.OpenFile(s.storageFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for key, value := range urlMap {
		value = appendProtocol(value)
		_, err := fmt.Fprintf(writer, "%s %s\n", key, value)
		if err != nil {
			return err
		}
	}

	return writer.Flush()
}

func compressURL(value string) string {
	hash := crc32.ChecksumIEEE([]byte(value))
	shortnedValue := fmt.Sprintf("%08X", hash)

	return fmt.Sprintf("%s/%s/", baseURL, shortnedValue)
}

func openInBrowser(longURL string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "linux":
		cmd = "xdg-open"
		args = []string{longURL}

	case "darwin":
		cmd = "open"
		args = []string{longURL}

	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start", longURL}

	default:
		return errors.New("unsupported platform")
	}

	fmt.Printf("opening %s in broswer...\n", longURL)
	return exec.Command(cmd, args...).Start()
}

func prompt(message string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s (y/n)", message)

	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))

	return response == "y"
}

func appendProtocol(url string) string {
	if !strings.HasPrefix(url, "http://") || !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}

	return url
}
