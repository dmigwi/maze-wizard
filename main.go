package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"golang.org/x/exp/maps"
	// "golang.org/x/exp/maps"
)

var port uint

const (
	MinPort uint16 = 1024
	MaxPort uint16 = 49151

	invalidSearchReply = "Sorry"
)

func init() {
	flag.UintVar(&port, "port", 8080, "Default port to be used to run the server")
}

func NotFound(w http.ResponseWriter, req *http.Request) {
	notFoundHandler(w)
}

func notFoundHandler(w http.ResponseWriter) {
	io.WriteString(w, invalidSearchReply)
}

func root(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		notFoundHandler(w)
		return
	}

	path, err := findShortestPath(body)
	if err != nil {
		log.Printf("Failed to process the shortest path :Error;  %v \n", err)
	}

	if err != nil || len(path) == 0 {
		notFoundHandler(w)
		return
	}

	log.Printf("Shortest path is %v \n", path)
	fmt.Fprint(w, path)
}

func main() {
	flag.Parse()

	log.Println("Beginning the Maze Wizard Path finder")

	if port < uint(MinPort) || port > uint(MaxPort) {
		log.Fatalf("Invalid port of %d was found. Support range is only from %d to %d", port, MinPort, MaxPort)
	}

	http.HandleFunc("/", root)

	url := fmt.Sprintf(":%d", port)
	if err := http.ListenAndServe(url, nil); err != nil {
		log.Fatalf("Serving on %s failed: Error %v", url, err)
	}
	log.Printf("Serving on port %d \n", port)
}

func isMazeEnd(loc string) bool {
	return strings.ToLower(loc) == "exit"
}

func findShortestPath(body []byte) ([]string, error) {
	var data = make(map[string]interface{})
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	var shortesPath []string
	for k, v := range data {
		info := findPath(k, v, []string{})
		if info != nil {
			info = append([]string{k}, info...)
		}
		if len(shortesPath) == 0 {
			shortesPath = info
		} else if len(shortesPath) == 0 && len(info) < len(shortesPath) {
			shortesPath = info
		}
	}

	if shortesPath == nil {
		shortesPath = []string{}
	}

	return shortesPath, nil
}

func findPath(startPos string, data interface{}, info []string) []string {
	str, ok := data.(string)
	if ok && isMazeEnd(str) {
		return []string{}
	}

	val, ok := data.(map[string]interface{})
	if ok {
		startPos = maps.Keys(val)[0]
		data = maps.Values(val)[0]
		p := findPath(startPos, data, info)
		if p == nil {
			return nil
		}
		info = append(info, startPos)
		return append(info, p...)
	}

	return nil
}
