package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

// var video = os.Getenv("BASEPATH") + "/video"

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Panic("failure loading environment variables")
	}

	// ideally this would include either a path variable or a query param to select a
	// specific song, but this is just an example and I'm too lazy lol
	http.Handle("/audio", http.HandlerFunc(audioHandler))

	// http.Handle("/video", videoHandler)

	log.Println("Now serving on port 8080")

	log.Fatal(http.ListenAndServe(":8080", nil))

}

func audioHandler(w http.ResponseWriter, r *http.Request) {
	var audio = os.Getenv("BASEPATH") + "/audio/"
	filePath := audio + "Scifi Adventure.mp3"
	fmt.Println(filePath)
	f, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "error reading from song file", http.StatusInternalServerError)
		return
	}

	defer f.Close()

	_, err = io.Copy(w, f)
	if err != nil {
		http.Error(w, "error reading song to response", http.StatusInternalServerError)
		return
	}

}
