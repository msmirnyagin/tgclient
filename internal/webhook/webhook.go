package webhook

import (
	"net/http"
	"strings"
)

func GetPost(url string, message string) error {
	payload := strings.NewReader(message)
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}
