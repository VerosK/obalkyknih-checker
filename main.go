package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var httpAddress = flag.String(
	"listenAddress", ":8000",
	"The address to listen on for HTTP requests.")
var httpCheckInterval = flag.Duration(
	"checkInterval", 30*time.Second,
	"Check interval of backend servers.")

type WatchedServer struct {
	BaseUrl  string `json:"baseUrl"`
	CheckUrl string `json:"checkUrl"`
	Alive    bool   `json:"alive"`
}

var DEFAULT_SERVERS = []WatchedServer{
	WatchedServer{
		BaseUrl:  "https://cache1.obalkyknih.cz/",
		CheckUrl: "https://cache1.obalkyknih.cz/api/runtime/alive",
		Alive:    true},
	WatchedServer{
		BaseUrl:  "https://cache2.obalkyknih.cz/",
		CheckUrl: "https://cache2.obalkyknih.cz/api/runtime/alive",
		Alive:    true},
}

var watchedServers []WatchedServer = DEFAULT_SERVERS
var serverAlive *WatchedServer

func httpRoot(w http.ResponseWriter, r *http.Request) {
	if serverAlive.Alive {
		out, err := json.Marshal(serverAlive.BaseUrl)
		if err != nil {
			panic(err)
		}
		w.Write(out)
		w.Header().Set("Cache-Control", "max-age=10s")
		return
	}
	http.Error(w, "# 500 error # All backends are currently dead", 500)

}

func status(w http.ResponseWriter, r *http.Request) {
	b, err := json.Marshal(watchedServers)

	if err != nil {
		panic(err)
	}

	w.Write(b)
}

func checkAlive(ticker *time.Ticker) {
	var n int

	serverAlive = &watchedServers[0]

	for range ticker.C {
		checkUrl := watchedServers[n].CheckUrl
		// log.Print("Checking: ", checkUrl)
		_, err := http.Get(checkUrl)
		if err == nil {
			watchedServers[n].Alive = true
			serverAlive = &watchedServers[n]
		} else {
			watchedServers[n].Alive = false
		}
		n = (n + 1)
		if n >= len(watchedServers) {
			n = 0
		}
	}
}

func main() {
	flag.Parse()

	log.Print("checkInterval: ", *httpCheckInterval)
	log.Print("tail: ", flag.Args())

	if flag.NArg() > 0 {
		watchedServers = make([]WatchedServer, flag.NArg())
		for n, arg := range flag.Args() {
			if strings.Index(arg, "=") == -1 {
				watchedServers[n].BaseUrl = arg
				watchedServers[n].CheckUrl = arg
			} else {
				parts := strings.SplitN(arg, "=", 2)
				watchedServers[n].BaseUrl = parts[0]
				watchedServers[n].CheckUrl = parts[1]
			}
			watchedServers[n].Alive = true
			log.Print("Added server: ",
				watchedServers[n].BaseUrl,
				" with checkUrl: ",
				watchedServers[n].CheckUrl,
			)
		}
	}

	ticker := time.NewTicker(*httpCheckInterval)
	go checkAlive(ticker)

	http.HandleFunc("/", httpRoot)
	http.HandleFunc("/status", status)
	http.Handle("/metrics", promhttp.Handler())

	log.Print("Listening on ", *httpAddress)

	log.Fatal(http.ListenAndServe(*httpAddress, nil))

}
