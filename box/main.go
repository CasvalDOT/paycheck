package box

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

// Some constants
const baseURL = "https://upload.box.com/api/2.0"

// IBox is the interface of the module
type IBox interface {
	Upload(*os.File, string) error
}

// Box main structure
type Box struct {
	token string
}

// Upload a file to specific
// directory
func (b *Box) Upload(file *os.File, targetID string) error {
	fileInfo, err := os.Stat(file.Name())
	if err != nil {
		return err
	}

	client := http.Client{}

	body := &bytes.Buffer{}

	multipartWriter := multipart.NewWriter(body)
	defer multipartWriter.Close()

	fieldWriter, err := multipartWriter.CreateFormField("attributes")
	if err != nil {
		return err
	}

	attributes, err := json.Marshal(map[string]interface{}{
		"name": fileInfo.Name(),
		"parent": map[string]interface{}{
			"id": targetID,
		},
	})
	if err != nil {
		return err
	}

	_, err = io.Copy(fieldWriter, strings.NewReader(string(attributes)))
	if err != nil {
		return err
	}

	fileWriter, err := multipartWriter.CreateFormFile("file", fileInfo.Name())
	if err != nil {
		return err
	}

	_, err = io.Copy(fileWriter, file)
	if err != nil {
		return err
	}

	multipartWriter.Close()

	request, err := http.NewRequest("POST", baseURL+"/files/content", bytes.NewReader(body.Bytes()))
	if err != nil {
		return err
	}

	request.Header.Set("Authorization", "Bearer "+b.token)
	request.Header.Set("Content-Type", multipartWriter.FormDataContentType())

	response, err := client.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusConflict {
		return nil
	}

	if response.StatusCode != http.StatusCreated {
		return errors.New(response.Status)
	}

	return nil
}

// New initialize a newer instance of
// Box model
func New(token string) IBox {
	b := Box{token: token}
	return &b
}
