package main

import (
    "crypto/md5"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strings"
    "time"
)

//TIP To run your code, right-click the code and select <b>Run</b>. Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.

type FileMetadata struct {
    Size     int64
    Created  time.Time
    Modified time.Time
    Accessed time.Time
}

func getFileMetadata(path string) (*FileMetadata, error) {
    info, err := os.Stat(path)
    if err != nil {
        return nil, err
    }

    return &FileMetadata{
        Size:     info.Size(),
        Created:  info.ModTime(), // Note: Go doesn't provide creation time in a cross-platform way
        Modified: info.ModTime(),
        Accessed: info.ModTime(), // Note: Go doesn't provide access time in a cross-platform way
    }, nil
}

func calculateMD5(path string) (string, error) {
    file, err := os.Open(path)
    if err != nil {
        return "", err
    }
    defer func(file *os.File) {
        err := file.Close()
        if err != nil {

        }
    }(file)

    hash := md5.New()
    if _, err := io.Copy(hash, file); err != nil {
        return "", err
    }

    return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func scanDirectory(dir string) error {
    return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            fmt.Printf("Error accessing path %s: %v\n", path, err)
            return nil
        }

        if !info.IsDir() {
            metadata, err := getFileMetadata(path)
            if err != nil {
                fmt.Printf("Error getting metadata for %s: %v\n", path, err)
                return nil
            }

            md5Hash, err := calculateMD5(path)
            if err != nil {
                fmt.Printf("Error calculating MD5 for %s: %v\n", path, err)
                return nil
            }

            fmt.Printf("File: %s\n", path)
            fmt.Printf("Size: %d bytes\n", metadata.Size)
            fmt.Printf("Modified: %v\n", metadata.Modified)
            fmt.Printf("MD5 Hash: %s\n", md5Hash)
            fmt.Println(strings.Repeat("-", 50))
        }

        return nil
    })
}

func main() {
    var dir string
    fmt.Print("Enter the directory path to scan: ")
    _, err := fmt.Scanln(&dir)
    if err != nil {
        return
    }

    err = scanDirectory(dir)
    if err != nil {
        fmt.Printf("Error scanning directory: %v\n", err)
    }
}

//TIP See GoLand help at <a href="https://www.jetbrains.com/help/go/">jetbrains.com/help/go/</a>.
// Also, you can try interactive lessons for GoLand by selecting 'Help | Learn IDE Features' from the main menu.
