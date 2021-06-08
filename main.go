package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"text/template"
	"time"
)

var httpAddress = flag.String(
	"listenAddress", ":8000",
	"The address to listen on for HTTP requests.")
var httpCheckInterval = flag.Duration(
	"checkInterval", 90*time.Second,
	"Check interval of backend servers.")
var aliveFile = flag.String(
	"aliveFile", "/tmp/obalkyknih.php",
	"File to store alive serves")
var aliveTemplateStr = flag.String(
	"aliveTemplate", "<?php\n$OBALKYKNIH_BASEURL=\"{{.Server.BaseUrl}}\";\n",
	"Template string to write.")
var aliveTemplate *template.Template

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

func updateStatusFile() {
	log.Print("updating output file to ", serverAlive.BaseUrl)
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)
	data := struct{ Server WatchedServer }{*serverAlive}
	err1 := aliveTemplate.Execute(writer, data)
	if err1 != nil {
		panic(err1)
	}
	writer.Flush()
	err2 := ioutil.WriteFile(*aliveFile, buf.Bytes(), 0644)
	if err2 != nil {
		panic(err2)
	}
}

func getWorkingServer() bool {
	var checkUrl string
	var server *WatchedServer
        var wasAlive bool;
	for n := range watchedServers {
		server = &watchedServers[n]
		checkUrl = server.CheckUrl
		_, err := http.Get(checkUrl)
		wasAlive = server.Alive
                server.Alive = (err == nil)
                if (wasAlive != server.Alive) {
                        log.Print("Server ", server.BaseUrl,
                                        " status is ", server.Alive);
                }

		if err == nil {
			serverAlive = server
			if serverAlive.BaseUrl == server.BaseUrl {
				return true
			}
			updateStatusFile()
			return true
		} else {
                        continue
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
	log.Print("outputFile: ", *aliveFile)
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

	aliveTemplate, _ = template.New("alive").Parse(*aliveTemplateStr)

	ticker := time.NewTicker(*httpCheckInterval)
	go checkAlive(ticker)

	http.HandleFunc("/", httpRoot)
	http.HandleFunc("/status", status)

	log.Print("Listening on ", *httpAddress)

	log.Fatal(http.ListenAndServe(*httpAddress, nil))
}
