package handlers

import (
    "github.com/Azure/azure-storage-blob-go/azblob"
    "log"
    "os/exec"
    "context"
    "net/http"
    "fmt"
    "os"
    "net/url"
)

func transcode(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPut {
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }

    // Define as credenciais de armazenamento
    accountName := "bloobstream"
    accountKey := os.Getenv("TOKEN_AZURE_BLOB")

    // Define o nome do contêiner e do arquivo de vídeo
    containerName := "cinema-storage"
    videoFileName := "teste.mp4"

    // Cria um ponto de extremidade de blob
    credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
    if err != nil {
        fmt.Println("Erro ao criar as credenciais:", err)
        return
    }

    p := azblob.NewPipeline(credential, azblob.PipelineOptions{})
    url, _ := url.Parse(
        fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName))

    // Abre o contêiner
    containerURL := azblob.NewContainerURL(*url, p)
    blobURL := containerURL.NewBlobURL(videoFileName)
    // opens the file from the request
    file, header, err := r.FormFile("file")
    if err != nil {
        fmt.Println("Error opening file:", err)
        return
    }
    defer file.Close()

    // opens a new file to write the transcoded file
    transcodedFile, err := os.Create(header.Filename)
    if err != nil {
        fmt.Println("Error creating transcoded file:", err)
        return
    }
    defer transcodedFile.Close()

    // runs the ffmpeg command to transcode the video
    cmd := exec.Command("ffmpeg", "-i", header.Filename, "-c:v", "libx264", "-b:v", "3000k", transcodedFile.Name())
    err = cmd.Run()
    if err != nil {
        log.Fatal(err)
    }

    // uploads the transcoded file to the container
    blockBlobURL := blobURL.ToBlockBlobURL()
    _, err = azblob.UploadFileToBlockBlob(context.Background(), transcodedFile, blockBlobURL, azblob.UploadToBlockBlobOptions{
        BlockSize:   4 * 1024 * 1024,
        Parallelism: 16})
    if err != nil {
        fmt.Println("Error uploading file to blob storage:", err)
        return
    }
    fmt.Fprintf(w, "File %s transcoded and uploaded successfully", header.Filename)
}

