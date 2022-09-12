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

	var info [][]string
	path := processShortestPath(body, info)
	// if err != nil {
	// 	log.Printf("Failed to process the shortest path :Error;  %v \n", err)
	// }

	if len(path) == 0 {
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

// var stack = make([]map[string]string, 0)

func processShortestPath(mazeMap []byte, info [][]string) [][]string {
	var data = make(map[string]string)
	if err := json.Unmarshal(mazeMap, &data); err != nil {
		return nil
	}

	for k, v := range data {
		if isMazeEnd(k) {
			return [][]string{[]string{k}}
		}

		if isMazeEnd(v) {
			return [][]string{[]string{k}, []string{v}}
		}

		if strings.ContainsAny(v, "{}") {
			return append(info, processShortestPath([]byte(v), info)...)
		}

	}
	return info
}
