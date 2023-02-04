package handlers

import (
	"io"
	"net/http"
	"strconv"

	"github.com/Azure/azure-sdk-for-go/storage"
)

func StreamVideo(w http.ResponseWriter, r *http.Request) {
	// Definir as credenciais de armazenamento do Azure
	accountName := "bloobstream"
	accountKey := "iDdGLcoy7DFiIEMJcaFHlFioFYHjPvdRjN8pdsPnujYGqz/QKzKBWcTen7jGAvBSgbn3eO37zpT6+AStFzGXKw=="

	// Conectar ao Azure Blob
	client, err := storage.NewBasicClient(accountName, accountKey)
	if err != nil {
		println(err)
		return
	}
	blobClient := client.GetBlobService()

	// Obter o nome do arquivo
	fileName := "fibonacciDisney.mp4"

	// Obter o contêiner
	container := blobClient.GetContainerReference("cinema-storage")

	// Obter o blob
	blob := container.GetBlobReference(fileName)

	// Definir o tipo de conteúdo como MP4
	w.Header().Set("Content-Type", "video/mp4")

	// Definir o tamanho do conteúdo
	w.Header().Set("Content-Length", strconv.FormatInt(blob.Properties.ContentLength, 10))

	// Obter o fluxo de dados do blob
	reader, err := blob.Get(nil)
	if err != nil {
		println(err)
		return
	}

	// Copiar o conteúdo do blob para a resposta HTTP
	_, err = io.Copy(w, reader)
	if err != nil {
		http.Error(w, "Error while copying content to response", http.StatusInternalServerError)
		return
	}
}	