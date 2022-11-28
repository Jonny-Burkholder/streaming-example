package handler

import (
	"io"
	"net/http"
	"os"
)

func ImageHandlerV2(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content Type", "image/jpg")
	var image = os.Getenv("BASEPATH") + "/image/"
	filePath := image + "Scifi Adventure.mp3"
	f, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "error reading from image file", http.StatusInternalServerError)
		return
	}

	defer f.Close()

	_, err = io.Copy(w, f)
	if err != nil {
		http.Error(w, "error reading song to response", http.StatusInternalServerError)
		return
	}

}
