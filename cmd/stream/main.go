package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Jonny-Burkholder/streaming-example/internal/handler"
	"github.com/Jonny-Burkholder/streaming-example/pkg/netkit"
	"github.com/joho/godotenv"
)

var audio = os.Getenv("BASEPATH") + "/audio"
var video = os.Getenv("BASEPATH") + "/video"

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Panic("failure loading environment variables")
	}

	r := netkit.NewRouter(nil)

	v1 := r.NewGroup("v1")
	v1.Get("/audio", addHeaders(http.FileServer(http.Dir(audio)))) // should really be adding the headers somewhere else
	v1.Get("/video", addHeaders(http.FileServer(http.Dir(video))))

	v2 := r.NewGroup("v2")
	// ideally this would include either a path variable or a query param to select a
	// specific song, but this is just an example and I'm too lazy lol
	v2.Get("/audio", http.HandlerFunc(handler.AudioHandlerV2))

	log.Println("Now serving on port 8080")

	log.Fatal(http.ListenAndServe(":8080", r))

}

func addHeaders(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		h.ServeHTTP(w, r)
	}
}
