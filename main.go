package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"time"

	"os"
	"strconv"

	"github.com/jlaffaye/ftp"
)

var (
    ftpServer   = os.Getenv("FTP_HOSTNAME") + ":" + os.Getenv("FTP_PORT")
    ftpUser     = os.Getenv("FTP_USERNAME")
    ftpPassword = os.Getenv("FTP_PASSWORD")
    ftpBaseDir  = "/public" // Replace with your FTP base directory
    ftpTLS, _   = strconv.ParseBool(os.Getenv("FTP_TLS"))
)

var ftpClient *ftp.ServerConn

// initFTP initializes the FTP connection
func initFTP() {
	var err error
	ftpClient, err = ftp.Dial(ftpServer, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		log.Fatalf("Failed to connect to FTP server: %v", err)
	}

	err = ftpClient.Login(ftpUser, ftpPassword)
	if err != nil {
		log.Fatalf("Failed to log in to FTP server: %v", err)
	}

	log.Println("Connected to FTP server.")
}

// serveFTPFile fetches a file from the FTP server and writes it to the HTTP response
func serveFTPFile(w http.ResponseWriter, r *http.Request) {
	filePath := ftpBaseDir + r.URL.Path
	log.Printf("Fetching file: %s", filePath)

	response, err := ftpClient.Retr(filePath)
	if err != nil {
		http.Error(w, "File not found on FTP server", http.StatusNotFound)
		return
	}
	defer response.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(response)
	if err != nil {
		http.Error(w, "Error reading file from FTP server", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
}

func main() {
	initFTP()
	defer ftpClient.Quit()

	http.HandleFunc("/", serveFTPFile)

	port := 80
	log.Printf("Serving files on port %d...", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
