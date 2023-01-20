package main

import (
	"fmt"
	"net/http"

	"github.com/carlinhoshk/stream-go/handlers"
)

func main() {
	http.HandleFunc("/transcode", handlers.TranscodeFFmpeg)
	// Inicia o servidor na porta 9090
	fmt.Println("Servidor iniciado na porta 9090")
	http.ListenAndServe(":9090", nil)
}
