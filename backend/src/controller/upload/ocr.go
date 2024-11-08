package upload

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
)

// sendToPaddleOCR envoie le fichier à PaddleOCR et récupère le résultat
func sendToPaddleOCR(fileHeader *multipart.FileHeader) (map[string]interface{}, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)
	part, err := writer.CreateFormFile("file", fileHeader.Filename)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}
	writer.Close()

	resp, err := http.Post("http://ocr:5000/ocr", writer.FormDataContentType(), &requestBody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, errors.New(string(bodyBytes))
	}

	var ocrResult map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&ocrResult)
	if err != nil {
		return nil, err
	}
	return ocrResult, nil
}
