package handlers

import (
	"context"
	"fmt"
	"io"
	"os"
	"net/http"
	"net/url"

	"github.com/Azure/azure-storage-blob-go/azblob"
)

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	// Define as credenciais de armazenamento
	accountName := "bloobstream"
	accountKey := os.Getenv("TOKEN_AZURE_BLOB")

	// Define o nome do contêiner e do arquivo de vídeo
	containerName := "cinema-storage"
	videoFileName := r.URL.Query().Get("file")

	// Cria um ponto de extremidade de blob
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})
	url, _ := url.Parse(
		fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName))

	// Abre o contêiner
	containerURL := azblob.NewContainerURL(*url, p)
	blobURL := containerURL.NewBlobURL(videoFileName)

	// Cria uma nova requisição para baixar o arquivo
	download, err := blobURL.Download(context.Background(), 0, azblob.CountToEnd, azblob.BlobAccessConditions{}, false, azblob.ClientProvidedKeyOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if download.Response().StatusCode != http.StatusOK {
		http.Error(w, "Erro ao baixar o arquivo de vídeo", http.StatusInternalServerError)
		return
	}

	// Envia o arquivo para o cliente
	w.Header().Set("Content-Type", "video/mp4")
	io.Copy(w, download.Body(azblob.RetryReaderOptions{}))
}
