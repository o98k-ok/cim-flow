package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
)

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func ExtractURI(stuct map[string]interface{})  string {
	ims, ok := stuct["images"].([]interface{})
	if !ok || len(ims) == 0 {
		return ""
	}

	im, ok := ims[0].(map[string]interface{})
	if !ok {
		return ""
	}

	return im["url"].(string)
}

func main()  {
	if len(os.Args) <= 1 {
		return
	}
	endPoint := "http://www.bing.com"
	url := fmt.Sprintf("%s/HPImageArchive.aspx?format=js&idx=%s&n=1&mkt=en-US", endPoint, os.Args[1])

	imagePath := os.Getenv("image_base_path")
	if imagePath == "" {
		_, _ = io.WriteString(os.Stderr, "please config local image path\n")
		return
	}


	_, err := os.Stat(imagePath)
	if err == os.ErrNotExist {
		_ = os.MkdirAll(imagePath, 0755)
	}

	resp, err := http.Get(url)
	if err != nil {
		_, _ = io.WriteString(os.Stderr, err.Error()+"\n")
		return
	}

	dat, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		_, _ = io.WriteString(os.Stderr, err.Error()+"\n")
		return
	}

	var stuct map[string]interface{}
	err = json.Unmarshal(dat, &stuct)
	if err != nil {
		_, _ = io.WriteString(os.Stderr, err.Error()+"\n")
		return
	}

	uri := ExtractURI(stuct)
	if uri == "" {
		_, _ = io.WriteString(os.Stderr, "get image url failed\n")
		return
	}

	imageUrl := fmt.Sprintf("%s%s", endPoint, uri)
	filepath := path.Join(imagePath, GetMD5Hash(imageUrl)+".jpg")

	resp, err = http.Get(imageUrl)
	if err != nil {
		_, _ = io.WriteString(os.Stderr, err.Error()+"\n")
		return
	}

	dat, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		_, _ = io.WriteString(os.Stderr, err.Error()+"\n")
		return
	}

	_ = os.WriteFile(filepath, dat, 0744)
	fmt.Print(filepath)
}
