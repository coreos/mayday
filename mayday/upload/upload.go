package upload

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/spf13/viper"
)

type UploadResponse struct {
	Sha    string `json:"sha"`
	Access string `json:"access_key"`
}

var apiVersion = "/v1"

func Upload(reader io.Reader) (*UploadResponse, error) {

	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	// TODO don't just ignore this error
	machineId, _ := ioutil.ReadFile("/etc/machine-id")
	bodyWriter.WriteField("machine_id", string(machineId))

	var res *UploadResponse

	fileWriter, err := bodyWriter.CreateFormFile("targz", "mayday.tar.gz")
	if err != nil {
		fmt.Println("error writing to buffer")
		return res, err
	}

	_, err = io.Copy(fileWriter, reader)
	if err != nil {
		return res, err
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	log.Print(viper.GetString("upload-server") + apiVersion + "/upload")
	resp, err := http.Post(viper.GetString("upload-server")+apiVersion+"/upload", contentType, bodyBuf)
	if err != nil {
		return res, err
	}
	defer resp.Body.Close()
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}

	if resp.StatusCode != 200 {
		return res, errors.New(string(resp_body))
	}
	err = json.Unmarshal(resp_body, &res)
	return res, nil
}
