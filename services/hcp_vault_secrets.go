package services

import (
	"fmt"

	secret "github.com/hashicorp/hcp-sdk-go/clients/cloud-vault-secrets/preview/2023-06-13/client/secret_service"
	"github.com/hashicorp/hcp-sdk-go/httpclient"
)

var secretClient secret.ClientService

func initSecretClient() error {
	if secretClient == nil {
		cl, err := httpclient.New(httpclient.Config{})
		if err != nil {
			return err
		}
		secretClient = secret.New(cl, nil)
	}
	return nil
}

func GetHcpVaultSecret(appName, secretName string) ([]byte, error) {
	err := initSecretClient()
	if err != nil {
		return nil, err
	}
	params := secret.NewOpenAppSecretParams()
	params.AppName = appName
	params.SecretName = secretName
	resp, err := secretClient.OpenAppSecret(params, nil)
	if err != nil {
		return nil, err
	}
	if resp.Payload == nil || resp.Payload.Secret == nil || resp.Payload.Secret.Version == nil {
		return nil, fmt.Errorf("invalid response")
	}
	return []byte(resp.Payload.Secret.Version.Value), nil
}
