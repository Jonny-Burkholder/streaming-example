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
var image = os.Getenv("BASEPATH") + "/image"

func main() {

	l := netkit.NewLogger(netkit.LevelInfo)

	defer func() {
		if r := recover(); r != nil {
			l.Error("Recovered from panic", r)
		}
	}()

	err := godotenv.Load(".env")
	if err != nil {
		log.Panic("failure loading environment variables")
	}

	r := netkit.NewRouter(nil)
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/ping" && r.Method == http.MethodGet {
			netkit.WriteRaw(w, r, http.StatusOK, []byte("PONG"))
		}
	})

	cors := netkit.CORSHandler(nil) // nil = default config
	r.Get("/", cors.ServeHTTP)      // this should handle cors site wide, I *THINK*
	// but it is just an example, so you don't have to use it, I just put it here in
	// case you wanted to play around with it.

	v1 := r.NewGroup("v1")
	v1.Get("/audio", handler.FileHandlerV1(audio)) // should really be adding the headers somewhere else
	v1.Get("/video", handler.FileHandlerV1(video))
	v1.Get("/image", handler.FileHandlerV1(image))

	v2 := r.NewGroup("v2")
	// ideally this would include either a path variable or a query param to select a
	// specific song, but this is just an example and I'm too lazy lol
	v2.Get("/audio", http.HandlerFunc(handler.AudioHandlerV2))
	v2.Get("/video", http.HandlerFunc(handler.ImageHandlerV2))

	log.Println("Now serving on port 8080")

	log.Panic(http.ListenAndServe(":8080", r))

}
