package main


// ServeAFileHandler takes a pointer to a file, and returns a handler func from it. Make
// sure you don't forget to close the file (outside outside of the handler) when finished.
// Also, I think it works, but I really havent tested it fully.
func ServeAFileHandler(file *os.File) http.HandlerFunc {
	// Initialize a small buffer.
	buf := make([]byte, 512)
    // Read the fist few bytes of the file so we can attempt
    // to auto detect the Content-Type. 
	_, err := file.Read(buf)
	if err != nil {
		panic(err)
	}
	// Attempt to detect the Content-Type
	var contentType string
	contentType = http.DetectContentType(buf)
	// Finally, rewind our file reader back to beginning.
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		panic(err)
	}
	// Now we can actually handle our web request...
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", contentType)
		_, err := io.CopyBuffer(w, file, buf)
		if err != nil {
			log.Print(err)
			http.Error(w, http.StatusText(http.StatusExpectationFailed), http.StatusExpectationFailed)
			return
		}
	}
}
