package shortner

import (
	"bufio"
	"errors"
	"fmt"
	"hash/crc32"
	"kurocfer/lil-url/utils"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type Shortner struct {
	storageFile string
}

var baseURL = "/lilurl"

func NewShortner(storageFile string) *Shortner {
	return &Shortner{
		storageFile: storageFile,
	}
}

func (s *Shortner) Shorten(longURL string) (string, error) {
	urlMap, err := utils.LoadURLs(s.storageFile)
	if err != nil {
		return "", err
	}

	shortURL := compressURL(longURL)
	if _, exists := urlMap[shortURL]; exists {
		return shortURL, nil
	}
	urlMap[shortURL] = longURL
	err = utils.SaveURLs(urlMap, s.storageFile)
	if err != nil {
		return "", fmt.Errorf("encountered an error while saving new link: %w", err)
	}

	return shortURL, nil
}

func (s *Shortner) LookupURL(shortURL string) error {
	urlMap, err := utils.LoadURLs(s.storageFile)
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

func (s *Shortner) List(numLines int) error {
	urlMap, err := utils.LoadURLs(s.storageFile)
	if err != nil {
		return err
	}

	linesPrinted := 0
	for key, value := range urlMap {
		if linesPrinted == numLines {
			break
		}
		fmt.Printf("Shortened -> %s : Original -> %s\n", key, value)
		linesPrinted++
	}
	fmt.Printf("total %d\n", linesPrinted)

	return nil
}

func (s *Shortner) Clear() error {
	file, err := os.OpenFile(s.storageFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	if _, err := file.WriteString(""); err != nil {
		return fmt.Errorf("an error occured while attempting to clear file %w", err)
	}

	return nil
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
