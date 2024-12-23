package ntfy

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"mgarnier11/go/utils"
	"net/http"
	"strings"
	"time"
)

func SendNotification(title, message, tags string) error {
	ntfyTopic := utils.GetEnv("NTFY_TOPIC", "")
	ntfyServer := utils.GetEnv("NTFY_SERVER", "")

	if ntfyTopic == "" || ntfyServer == "" {
		return errors.New("invalid ntfy configuration")
	}

	ntfyServer = strings.TrimSuffix(ntfyServer, "/")
	ntfyTopic = strings.TrimPrefix(ntfyTopic, "/")

	ntfyUrl := fmt.Sprintf("http://%s/%s", ntfyServer, ntfyTopic)

	req, err := http.NewRequest("POST", ntfyUrl, bytes.NewBufferString(message))
	if err != nil {
		return err
	}

	currentTime := time.Now().In(time.FixedZone("Europe/Paris", 1*60*60)) // Set to Paris timezone
	titleWithTime := fmt.Sprintf("%s - %s", title, currentTime.Format("02/01/2006 15:04:05"))

	req.Header.Set("Title", titleWithTime)
	req.Header.Set("Tags", tags)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error while sending a NTFY notification, url: %s, error: %w, ", ntfyUrl, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("error while sending a NTFY notification, url: %s, status code: %d, ", ntfyUrl, resp.StatusCode)
	}
	return nil
}
