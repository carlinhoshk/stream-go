package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/Azure/azure-storage-blob-go/azblob"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	// Verifica se o método da solicitação é POST
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Lê o arquivo enviado na solicitação
	file, header, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Erro ao ler o arquivo enviado: %v\n", err)
		return
	}
	defer file.Close()

	// Lê o arquivo em binário
	fileData := make([]byte, header.Size)
	_, err = file.Read(fileData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Erro ao ler o arquivo em binário: %v\n", err)
		return
	}

	// Define as credenciais de armazenamento
	accountName := "bloobstream"
	accountKey := os.Getenv("TOKEN_AZURE_BLOB")

	// Define o nome do contêiner e do arquivo de vídeo
	containerName := "cinema-storage"
	videoFileName := header.Filename

	// Cria um ponto de extremidade de blob
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Erro ao criar as credenciais: %v\n", err)
		return
	}
	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})
	url, _ := url.Parse(
		fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName))

	// Envia o arquivo para o blob
	containerURL := azblob.NewContainerURL(*url, p)
	blobURL := containerURL.NewBlockBlobURL(videoFileName)
	_, err = azblob.UploadBufferToBlockBlob(context.Background(), fileData, blobURL, azblob.UploadToBlockBlobOptions{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Erro ao fazer upload do arquivo: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Upload do arquivo %s concluído com sucesso!\n", videoFileName)
}