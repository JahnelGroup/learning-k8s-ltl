package main

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/gomodule/redigo/redis"
)

var redisPath string

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/images/", handleImageRequest)

	http.HandleFunc("/upload", uploadHandler)
	redisPath = os.Getenv("REDIS_URL")
	fmt.Println("Service is ready")
	http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("SERVICE_PORT")), nil)

}

func getImageFromRedis(key string) ([]byte, error) {
	// Establish a connection to Redis
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// Load the image from Redis using the given key
	imageData, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}

	return imageData, nil
}

func handleImageRequest(w http.ResponseWriter, r *http.Request) {
	// Extract the image key from the URL path
	key := r.URL.Path[len("/image/"):]
	fmt.Println(key)
	// Load the image from Redis
	imageData, err := getImageFromRedis(strings.ReplaceAll(key, "/", ""))
	if err != nil {
		// Handle the error
		http.Error(w, "Error loading image from Redis", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "image/gif")

	// Write the image data to the response writer
	_, err = w.Write(imageData)
	if err != nil {
		// Handle the error
		http.Error(w, "Error writing image data to response", http.StatusInternalServerError)
		return
	}
}

func storeImage(w http.ResponseWriter, r *http.Request, fileName string, filePath string) {
	// Read the image file into memory
	imageBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	conn, err := redis.Dial("tcp", redisPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Store the image data in Redis
	_, err = conn.Do("SET", fileName, imageBytes)
	fmt.Println("set", fileName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Redirect the user back to the index page
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func gifHandler(w http.ResponseWriter, r *http.Request) {
	// Get the image file path from the URL parameter
	filePath := r.URL.Query().Get("file_path")
	if filePath == "" {
		http.Error(w, "Image file path not provided", http.StatusBadRequest)
		return
	}

	// Connect to Redis
	conn, err := redis.Dial("tcp", redisPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	// Read the image file into memory
	imageBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Store the image data in Redis
	_, err = conn.Do("SET", "image_key", imageBytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Retrieve the image data from Redis
	redisImageBytes, err := redis.Bytes(conn.Do("GET", "image_key"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write the image data to the response writer
	w.Header().Set("Content-Type", "image/gif")
	_, err = w.Write(redisImageBytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getImages() ([]string, error) {
	conn, err := redis.Dial("tcp", redisPath)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer conn.Close()
	images, err := redis.Strings(conn.Do("KEYS", "*.gif"))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return images, nil
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// Define a struct to hold the list of GIF filenames
	type PageData struct {
		GIFs []string
	}

	// Read the names of all GIF files in the "gifs" directory
	gifFiles, err := getImages() //filepath.Glob("gifs/*.gif")
	if err != nil {
		fmt.Fprintf(w, "Error reading GIF files: %v", err)
		return
	}

	// Create a new PageData struct and fill it with the list of GIF filenames
	data := PageData{GIFs: gifFiles}

	// Parse the "index.html" template file and execute it with the PageData struct
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		fmt.Fprintf(w, "Error parsing template: %v", err)
		return
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		fmt.Fprintf(w, "Error executing template: %v", err)
		return
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is POST
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the multipart form data
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	// Get the "gif" file from the form data
	file, handler, err := r.FormFile("gif")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error retrieving file from form data", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create a new file with the same name as the uploaded file
	f, err := os.OpenFile("gifs/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error creating file", http.StatusInternalServerError)
		return
	}

	// Copy the uploaded file to the new file
	io.Copy(f, file)
	f.Close()

	storeImage(w, r, handler.Filename, "gifs/"+handler.Filename)
}
