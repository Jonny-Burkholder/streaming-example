package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Jonny-Burkholder/streaming-example/pkg/netkit"
)

func main() {

	r := netkit.NewRouter(nil)
	r.Get("/home", handleHome)
	r.Get("/home/info", handleHomeInfo)

	v1 := r.NewGroup("v1")
	v1.Get("/home", handleHomeVersion(1))
	v1.Get("/home/info", handleHomeInfoVersion(1))

	v2 := r.NewGroup("v2")
	v2.Get("/home", handleHomeVersion(2))
	v2.Get("/home/info", handleHomeInfoVersion(2))

	log.Panic(http.ListenAndServe(":3000", r))

}

func handleHome(w http.ResponseWriter, r *http.Request) {
	netkit.WriteRaw(w, r, 200, []byte("<h1>This is my home page<h1>"))
}

func handleHomeInfo(w http.ResponseWriter, r *http.Request) {
	netkit.WriteRaw(w, r, 200, []byte("<h1>This is my home page with info<h1>"))
}

func handleHomeVersion(version int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := fmt.Sprintf("<h1>This is my home page<h1><p>(VERSION %d GROUP)</p>", version)
		netkit.WriteRaw(w, r, 200, []byte(data))
	}
}

func handleHomeInfoVersion(version int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := fmt.Sprintf("<h1>This is my home page <br>with info</b><h1><p>(VERSION %d GROUP)</p>", version)
		netkit.WriteRaw(w, r, 200, []byte(data))
	}
}
