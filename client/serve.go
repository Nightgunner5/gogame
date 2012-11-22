package client

import (
	"io"
	"net/http"
)

func init() {
	http.HandleFunc("/engine.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")

		io.WriteString(w, Script)
	})
}
