package handlers

import (
	"io"
	"net/http"

	"github.com/carlinhoshk/stream-go/models"
)

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	// Obtém o nome do contêiner e do arquivo a ser baixado
	containerName := r.URL.Query().Get("container")
	fileName := r.URL.Query().Get("file")

	// Inicializa o BlobStorage com as credenciais de armazenamento
	storage := models.NewBlobStorage("myAccountName", "myAccountKey")

	// Faz o download do arquivo
	file, err := storage.DownloadFile(containerName, fileName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Envia o arquivo para o cliente
	w.Header().Set("Content-Type", "video/mp4")
	io.Copy(w, file)
}
