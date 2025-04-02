package external

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"mgarnier11.fr/go/mineager/config"
	"mgarnier11.fr/go/mineager/server/objects/bo"

	"mgarnier11.fr/go/libs/logger"
)

type OfflineStatus struct {
	Motd string `json:"motd"`
}

type CreateProxyRequest struct {
	DomainName        string        `json:"domainName"`
	ListenTo          string        `json:"listenTo"`
	ProxyTo           string        `json:"proxyTo"`
	DisconnectMessage string        `json:"disconnectMessage"`
	OfflineStatus     OfflineStatus `json:"offlineStatus"`
}

func CreateProxy(host *bo.HostBo, serverBo *bo.ServerBo) error {
	createProxyRequest := &CreateProxyRequest{
		DomainName:        serverBo.Url,
		ListenTo:          "0.0.0.0:25565",
		ProxyTo:           fmt.Sprintf("%s:%d", host.ProxyIp, serverBo.Port),
		DisconnectMessage: "Goodbye",
		OfflineStatus: OfflineStatus{
			Motd: "Server is offline",
		},
	}

	// Convert the data to JSON
	jsonData, err := json.Marshal(createProxyRequest)
	if err != nil {

		return fmt.Errorf("error marshalling JSON: %v", err)
	}

	requestUrl := fmt.Sprintf("%s/proxies/%s.json", config.Config.AppConfig.InfraredUrl, serverBo.Name)

	// Create the POST request
	resp, err := http.Post(requestUrl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error sending POST request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %v", err)
	}

	logger.Infof("Status: %s", resp.Status)
	logger.Infof("Response: %s", body)

	return nil
}

func DeleteProxy(host *bo.HostBo, serverBo *bo.ServerBo) error {
	requestUrl := fmt.Sprintf("%s/proxies/%s.json", config.Config.AppConfig.InfraredUrl, serverBo.Name)

	// Create the DELETE request
	req, err := http.NewRequest(http.MethodDelete, requestUrl, nil)
	if err != nil {
		return fmt.Errorf("error creating DELETE request: %v", err)
	}

	// Send the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending DELETE request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %v", err)
	}

	logger.Infof("Status: %s", resp.Status)
	logger.Infof("Response: %s", body)

	return nil

}
