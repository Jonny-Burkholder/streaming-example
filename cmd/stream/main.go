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

	v1 := r.NewGroup("v1")
	v1.Get("/audio", handler.FileHandlerV1(audio)) // should really be adding the headers somewhere else
	v1.Get("/video", handler.FileHandlerV1(video))

	v2 := r.NewGroup("v2")
	// ideally this would include either a path variable or a query param to select a
	// specific song, but this is just an example and I'm too lazy lol
	v2.Get("/audio", http.HandlerFunc(handler.AudioHandlerV2))

	log.Println("Now serving on port 8080")

	log.Panic(http.ListenAndServe(":8080", r))

}
