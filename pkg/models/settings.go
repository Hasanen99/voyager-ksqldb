package models

import (
	"encoding/json"
	"fmt"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

type PluginSettings struct {
	// These fields map directly to your plugin.json's jsonData properties
	Ksql     string `json:"ksqlserver"`
	Http     bool   `json:"http"`
	Username string `json:"username"`
	// Pass     string                `json:"password"`
	Secrets *SecretPluginSettings `json:"-"`
}

type SecretPluginSettings struct {
	Pass string `json:"password"`
}

func LoadPluginSettings(source backend.DataSourceInstanceSettings) (*PluginSettings, error) {
	settings := PluginSettings{}
	err := json.Unmarshal(source.JSONData, &settings)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal PluginSettings json: %w", err)
	}

	settings.Secrets = loadSecretPluginSettings(source.DecryptedSecureJSONData)

	return &settings, nil
}

func loadSecretPluginSettings(source map[string]string) *SecretPluginSettings {
	return &SecretPluginSettings{
		Pass: source["password"],
	}
}
