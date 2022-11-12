package handler

import (
	"io"
	"net/http"
	"os"
)

func AudioHandlerV2(w http.ResponseWriter, r *http.Request) {
	var audio = os.Getenv("BASEPATH") + "/audio/"
	filePath := audio + "Scifi Adventure.mp3"
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
