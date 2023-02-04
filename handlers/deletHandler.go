package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/Azure/azure-storage-blob-go/azblob"
)

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	// Verifica se o método da solicitação é DELETE
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Obtém o nome do arquivo a ser excluído
	fileName := r.URL.Query().Get("fileName")
	if fileName == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Nome do arquivo não informado\n")
		return
	}

	// Define as credenciais de armazenamento
	accountName := "bloobstream"
	accountKey := os.Getenv("TOKEN_AZURE_BLOB")

	// Define o nome do contêiner e do arquivo de vídeo
	containerName := "cinema-storage"

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

	// Exclui o arquivo do blob
	containerURL := azblob.NewContainerURL(*url, p)
	blobURL := containerURL.NewBlockBlobURL(fileName)
	_, err = blobURL.Delete(context.Background(), azblob.DeleteSnapshotsOptionInclude, azblob.BlobAccessConditions{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Erro ao excluir o arquivo: %v\n", err)
		return
	}
	
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Exclusão do arquivo %s concluída com sucesso!\n", fileName)
}
