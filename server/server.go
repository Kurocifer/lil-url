package server

import (
	"fmt"
	"kurocfer/lil-url/utils"
	"log"
	"net/http"
)

var (
	Port        = ":8080"
	urlMap      map[string]string
	storageFile string
)

func Start(filename, port string) error {
	Port = port
	storageFile = filename
	http.HandleFunc("/", redirect)

	log.Printf("Server starting on port %s", port)
	log.Printf("Access links with localhost:%s/shortned-link", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		return fmt.Errorf("failed to start server: %v", err)
	}

	return nil
}

func redirect(w http.ResponseWriter, r *http.Request) {
	reqeustURL := "http:/" + r.URL.String()

	urlMap, _ := utils.LoadURLs(storageFile)

	value, exist := urlMap[reqeustURL]
	if !exist {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("Reponse from lilurl -> Invalid URL: " + reqeustURL + "has no match to a known url"))
		if err != nil {
			log.Printf("Error writing reponse: %v", err)
		}
		return
	}

	http.Redirect(w, r, value, http.StatusFound)
}
