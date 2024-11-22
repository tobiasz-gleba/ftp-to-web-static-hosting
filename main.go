package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/jlaffaye/ftp"
)

var (
	ftpServer   = getEnv("FTP_HOSTNAME", "localhost") + ":" + getEnv("FTP_PORT", "21")
	ftpUser     = getEnv("FTP_USERNAME", "admin")
	ftpPassword = getEnv("FTP_PASSWORD", "admin")
	ftpBaseDir  = getEnv("FTP_BASEDIR", "/public")
)

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

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

type cacheEntry struct {
	data      []byte
	timestamp time.Time
}

var (
	cache      = make(map[string]cacheEntry)
	cacheMutex sync.Mutex
	cacheSize  int64
	maxCacheSize int64 = int64(getEnvInt("CACHE_SIZE_MB", 300)) * 1024 * 1024 // Default 300MB
	cacheTTL   time.Duration = time.Duration(getEnvInt("CACHE_TTL_MINUTES", 30)) * time.Minute // Default 30 minutes
)

// serveFTPFile fetches a file from the FTP server and writes it to the HTTP response
func serveFTPFile(w http.ResponseWriter, r *http.Request) {
	filePath := ftpBaseDir + r.URL.Path
	log.Printf("Fetching file: %s", filePath)

	cacheMutex.Lock()
	entry, found := cache[filePath]
	if found && time.Since(entry.timestamp) < cacheTTL {
		log.Printf("Serving %s from cache", filePath)
		w.Header().Set("Content-Type", http.DetectContentType(entry.data))
		w.WriteHeader(http.StatusOK)
		w.Write(entry.data)
		cacheMutex.Unlock()
		return
	} else if found {
		// Delete expired cache entry
		delete(cache, filePath)
		cacheSize -= int64(len(entry.data))
	}
	cacheMutex.Unlock()

	// Fetch the file from the FTP server
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

	data := buf.Bytes()

	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	// Update the cache
	if found {
		cacheSize -= int64(len(entry.data))
	}
	cache[filePath] = cacheEntry{data: data, timestamp: time.Now()}
	cacheSize += int64(len(data))

	// Evict old entries if cache size exceeds the limit
	for cacheSize > maxCacheSize {
		oldestKey := ""
		oldestTime := time.Now()
		for key, entry := range cache {
			if entry.timestamp.Before(oldestTime) {
				oldestTime = entry.timestamp
				oldestKey = key
			}
		}
		if oldestKey != "" {
			cacheSize -= int64(len(cache[oldestKey].data))
			delete(cache, oldestKey)
		}
	}

	w.Header().Set("Content-Type", http.DetectContentType(data))
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func main() {
	initFTP()
	defer ftpClient.Quit()

	http.HandleFunc("/", serveFTPFile)

	port := 80
	log.Printf("Serving files on port %d...", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
