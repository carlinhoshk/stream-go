package handlers

import (
    "fmt"
    "log"
    "os"
    "os/exec"
    "net/http"
)

func TranscodeFFmpeg(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPut {
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }
    // opens the file from the request
    file, header, err := r.FormFile("file")
    if err != nil {
        fmt.Println("Error opening file:", err)
        return
    }
    defer file.Close()

    // creates the transcoded file
    transcodedFile, err := os.Create(header.Filename + ".transcoded")
    if err != nil {
        fmt.Println("Error creating transcoded file:", err)
        return
    }
    defer transcodedFile.Close()

    // runs the ffmpeg command to 
