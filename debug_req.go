package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func main() {
	url := "http://localhost:8080/api/checkout"
	jsonBody := []byte(`{"items":[{"produk_id":37,"quantity":1}]}`)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Response Status:", resp.Status)
	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Response Body:", string(body))
}
