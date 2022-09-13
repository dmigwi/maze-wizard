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
)

var port uint

const (
	// MINPORT defines the minimum port value allowed to serve.
	MINPORT uint16 = 1024
	// MAXPORT defines the maximum port value allowed to server
	MAXPORT uint16 = 49151

	// invalidSearchReply defines the default search return value if an error occures.
	invalidSearchReply = "Sorry"
)

func init() {
	flag.UintVar(&port, "port", 8080, "Default port to be used to run the server")
}

// NotFound returns a customized response for the non supported routes.
func NotFound(w http.ResponseWriter, req *http.Request) {
	notFoundHandler(w)
}

func notFoundHandler(w http.ResponseWriter) {
	io.WriteString(w, invalidSearchReply)
}

// root handles the default root handler function.
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

	if port < uint(MINPORT) || port > uint(MAXPORT) {
		log.Fatalf("Invalid port of %d was found. Support range is only from %d to %d", port, MINPORT, MAXPORT)
	}

	http.HandleFunc("/", root)

	url := fmt.Sprintf(":%d", port)
	if err := http.ListenAndServe(url, nil); err != nil {
		log.Fatalf("Serving on %s failed: Error %v", url, err)
	}
	log.Printf("Serving on port %d \n", port)
}

// isMazeEnd return true if the maze exit is found in the path
func isMazeEnd(loc string) bool {
	return strings.ToLower(loc) == "exit"
}

// findShortestPath accepts the maze map and returns the shortest path from the
// start position to the exit.
func findShortestPath(body []byte) ([]string, error) {
	var data = make(map[string]interface{})
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	var shortestPath []string
	for k, v := range data {
		info := findPath(k, v, []string{})
		if info != nil {
			info = append([]string{k}, info...)
		} else {
			// ignore empty paths returned
			continue
		}
		if len(shortestPath) == 0 {
			shortestPath = info
		} else if len(info) < len(shortestPath) {
			shortestPath = info
		}
	}

	// initiaze the path array if its still null
	if shortestPath == nil {
		shortestPath = []string{}
	}

	return shortestPath, nil
}

// findPath for the provided starting position and all possible routes to check
// it returns a path the maze exit if exists otherwise it returns null.
// It recursively navigates through all the paths.
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
