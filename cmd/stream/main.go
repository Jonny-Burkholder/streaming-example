package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Jonny-Burkholder/streaming-example/internal/handler"
	"github.com/joho/godotenv"
)

var audio = os.Getenv("BASEPATH") + "/audio"
var video = os.Getenv("BASEPATH") + "/video"

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Panic("failure loading environment variables")
	}

	http.Handle("/v1/audio", addHeaders(http.FileServer(http.Dir(audio))))
	http.Handle("/v1/video", addHeaders(http.FileServer(http.Dir(video))))

	// ideally this would include either a path variable or a query param to select a
	// specific song, but this is just an example and I'm too lazy lol
	http.Handle("/v2/audio", http.HandlerFunc(handler.AudioHandlerV2))

	log.Println("Now serving on port 8080")

	log.Fatal(http.ListenAndServe(":8080", nil))

}

func addHeaders(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		h.ServeHTTP(w, r)
	}
}
