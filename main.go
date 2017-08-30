package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"strings"
	"time"

)

var httpAddress = flag.String(
	"listenAddress", ":8000",
	"The address to listen on for HTTP requests.")
var httpCheckInterval = flag.Duration(
	"checkInterval", 15*time.Second,
	"Check interval of backend servers.")
var aliveFile = flag.String(
	"aliveFile", "/dev/stdout",
	"File to store alive serves")
var aliveTemplate = flag.String(
	"aliveTemplate", "$obalkyknih=\"{{ server }}\";",
	"Template string to write.")

type WatchedServer struct {
	BaseUrl  string `json:"baseUrl"`
	CheckUrl string `json:"checkUrl"`
	Alive    bool   `json:"alive"`
}

var DEFAULT_SERVERS = []WatchedServer{
	WatchedServer{
		BaseUrl:  "https://cache1.obalkyknih.cz/",
		CheckUrl: "https://cache1.obalkyknih.cz/api/runtime/alive",
		Alive:    false},
	WatchedServer{
		BaseUrl:  "https://cache2.obalkyknih.cz/",
		CheckUrl: "https://cache2.obalkyknih.cz/api/runtime/alive",
		Alive:    false},
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

func updateStatusFile() {
	log.Print("updating output file to ", serverAlive.BaseUrl)
}

func getWorkingServer() bool {
	var checkUrl string
        var server *WatchedServer
	for n := range watchedServers {
                server = &watchedServers[n]
		checkUrl = server.CheckUrl
		_, err := http.Get(checkUrl)
		if err == nil {
			// no change
			if serverAlive.BaseUrl == server.BaseUrl {
				return true
			}
			server.Alive = true
			serverAlive = server
			updateStatusFile()
			return true
		} else {
			server.Alive = false
		}
	}
	log.Print("all servers are dead")
	return false
}

func checkAlive(ticker *time.Ticker) {

	serverAlive = &watchedServers[0]
	updateStatusFile()
	getWorkingServer()

	for range ticker.C {
		getWorkingServer()
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
			watchedServers[n].Alive = false
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

	log.Print("Listening on ", *httpAddress)

	log.Fatal(http.ListenAndServe(*httpAddress, nil))
}
