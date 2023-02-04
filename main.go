package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/carlinhoshk/stream-go/handlers"
	
	
)

func main() {
	
	http.HandleFunc("/upload", handlers.UploadHandler)

	http.HandleFunc("/download", handlers.DownloadHandler)

	http.HandleFunc("/delete", handlers.DeleteHandler)
	
	//http.HandleFunc("/transcode", handlers.TranscodeFFmpeg)
	// Inicia o servidor na porta 9090
	fmt.Println("Servidor iniciado na porta 9090")
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		println(err)
		os.Exit(1)
	}
}

