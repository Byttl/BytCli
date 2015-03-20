package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func main() {

	flag.Parse()

	filename := flag.Arg(0)

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(filename))
	if err != nil {
		fmt.Println(err)
		return
	}
	if _, err = io.Copy(part, file); err != nil {
		fmt.Println(err)
		return
	}

	if err = writer.Close(); err != nil {
		fmt.Println(err)
		return
	}

	req, err := http.NewRequest("POST", "http://byt.tl/f/upload", body)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	result := make(map[string]interface{})
	if err := json.Unmarshal(data, &result); err != nil {
		fmt.Println(err)
		return
	}

	if resp.StatusCode != 200 {
		fmt.Println(result["error"].(string))
	} else {
		fmt.Println("wget", result["file"].(string))
	}
}
