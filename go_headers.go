package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type Headers struct {
	Ipaddress string `json:"ipaddress"`
	Language  string `json:"language"`
	Software  string `json:"software"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", http.StripPrefix("/", fs))

	http.HandleFunc("/api", handler)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	re := regexp.MustCompile(`\((.*?)\)`)
	headers := Headers{}
	headers.Ipaddress = r.RemoteAddr
	headers.Language = strings.Split(r.Header.Get("Accept-Language"), ",")[0]
	headers.Software = re.FindStringSubmatch(r.Header.Get("User-Agent"))[1]

	js, err := json.Marshal(headers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
