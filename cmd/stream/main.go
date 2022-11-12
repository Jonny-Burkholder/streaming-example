package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var audio = os.Getenv("BASEPATH") + "/audio"
var video = os.Getenv("BASEPATH") + "/video"

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Panic("failure loading environment variables")
	}

	http.Handle("/audio", addHeaders(http.FileServer(http.Dir(audio))))
	http.Handle("/video", addHeaders(http.FileServer(http.Dir(video))))

	log.Println("Now serving on port 8080")

	log.Fatal(http.ListenAndServe(":8080", nil))

}

func addHeaders(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		h.ServeHTTP(w, r)
	}
}
