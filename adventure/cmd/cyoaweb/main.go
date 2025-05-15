package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Georgeygigz/go-quizes/adventure"
)

func main() {
	port := flag.Int("port", 3000, "the port to start the CYOA web application on")
	filename := flag.String("file", "gopher.json", "the JSON file with the CYOA story")
	flag.Parse()

	file, err := os.Open(*filename)
	checkError(err, fmt.Sprintf("reading %s", *filename))

	story, err := adventure.JsonStory(file)
	checkError(err, "decoding JSON")

	h := adventure.NewHandler(story)
	fmt.Printf("Starting the server at:%d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), h))

}

func checkError(err error, action string) {
	if err != nil {
		exit(fmt.Sprintf("Encountered an error when: %s Error:%s", action, err))
	}
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
