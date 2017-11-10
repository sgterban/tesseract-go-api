package main

import (
	"bytes"
	"fmt"
	"github.com/otiai10/gosseract"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

//generate random seed for random string function
var randomSeed *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

//character set used for random string character selection
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func main() {
	http.HandleFunc("/", serveIndex)
	http.HandleFunc("/image", processHTTPImage)

	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}

func serveIndex(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(GetTemplate()))
}

func processHTTPImage(w http.ResponseWriter, req *http.Request) {
	//current directory where images will be temporarily stored
	currentDirectory, err := filepath.Abs(filepath.Dir(os.Args[0]))
	//read ArrayBuffer request of image
	requestBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "Error processing image", 500)
		return
	}
	//setup Reader for []byte array
	reader := bytes.NewReader(requestBytes)
	//decode image from binary data
	img, format, err := image.Decode(reader)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Error decoding image", 500)
		return
	}

	fmt.Println("Image Format: " + format)
	height := img.Bounds().Max.X - img.Bounds().Min.X
	width := img.Bounds().Max.Y - img.Bounds().Min.Y
	fmt.Println("Height: " + strconv.Itoa(height) + ", Width: " + strconv.Itoa(width))

	img = optimizeImage(img)

	//generate random file for temporary processing
	imgPath := currentDirectory + "/" + RandomString(20) + "." + format
	fmt.Println(imgPath)

	//create file, defer close and removal
	toimg, _ := os.Create(imgPath)
	defer toimg.Close()
	defer os.Remove(imgPath)

	//each image format uses a different encoding format
	if format == "png" {
		png.Encode(toimg, img)
	} else if format == "jpeg" {
		jpeg.Encode(toimg, img, &jpeg.Options{jpeg.DefaultQuality})
	}

	//setup tesseract client to process image, defer close of client
	client := gosseract.NewClient()
	defer client.Close()

	//set image to temp file created above
	client.SetImage(imgPath)
	text, err := client.Text()
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Error processing image text", 500)
		return
	}

	fmt.Println("------ Text Found -------")
	fmt.Println(text)
	fmt.Println("-------------------------\n")

	//replace /n characters with html line breaks
	response := strings.Replace(text, "\n", "<br/>", -1)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

func optimizeImage(img image.Image) image.Image {
	return img
}

//creates a random string of passed length, used for temp file creation
func RandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[randomSeed.Intn(len(charset))]
	}
	return string(b)
}

func GetTemplate() string {
	return `
	<html>
		<head>
			<title>OCR TEST</title>
		</head>
		<body>
			<h1>OCR Golang Test</h1>
			<p>Upload an image and hit submit to see what text is detected</p>
			<input name='file' type='file' id='upload_file'>
			<br/><br/>
			<button type='button' onclick='upload_file()' >Submit</button>
			<p id='detected_text'></p>

			<script type='text/javascript'>
				var send_file = function(e) {
					var data = e.target.result;
					var req = new XMLHttpRequest();
					req.open('POST', 'http://localhost:8080/image', true);
					req.onload = function() {
						text_update = document.getElementById('detected_text');
						text_update.innerHTML = req.responseText;
					};
					req.send(data);
				};

				var upload_file = function() {
					var file = document.getElementById('upload_file').files[0];
					var reader = new FileReader();
					reader.onload = send_file;
	    			reader.readAsArrayBuffer(file);
				};
			</script>
		</body>
	</html>
	`
}
