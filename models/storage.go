package models

import (
	"context"
	"fmt"
	
	"net/url"
	"os"

	"github.com/Azure/azure-storage-blob-go/azblob"
)

// BlobStorage é uma struct que armazena as informações de conta e chave de armazenamento
type BlobStorage struct {
	AccountName string
	AccountKey  string
}

// NewBlobStorage cria uma nova instância de BlobStorage
func NewBlobStorage(accountName, accountKey string) *BlobStorage {
	return &BlobStorage{
		AccountName: accountName,
		AccountKey:  accountKey,
	}
}

// UploadFile faz upload de um arquivo para o Azure Blob Storage
func (s *BlobStorage) UploadFile(containerName, fileName string) error {
	f, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("Error opening file: %v", err)
	}
	defer f.Close()

	credential, err := azblob.NewSharedKeyCredential(s.AccountName, s.AccountKey)
	if err != nil {
		return fmt.Errorf("Error creating storage credentials: %v", err)
	}

	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})
	url, _ := url.Parse(
		fmt.Sprintf("https://%s.blob.core.windows.net/%s", s.AccountName, containerName))

	containerURL := azblob.NewContainerURL(*url, p)
	blobURL := containerURL.NewBlockBlobURL(fileName)
	_, err = azblob.UploadFileToBlockBlob(context.Background(), f, blobURL, azblob.UploadToBlockBlobOptions{
		BlockSize:   4 * 1024 * 1024,
		Parallelism: 16})
	if err != nil {
		return fmt.Errorf("Error uploading file: %v", err)
	}
	return nil
}
