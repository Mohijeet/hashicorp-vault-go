package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/mapstructure"
	//"github.com/Shopify/sarama"
)

// type VaultSecrets struct {
// 	KafkaBootstrap string `json:"vault.kafka.bootstrap.server"`
// 	JwtSecret      string `json:"auth.jwt.secret"`
// }

type VaultSecrets struct {
	KafkaBootstrap string `mapstructure:"vault.kafka.bootstrap.server"`
	JwtSecret      string `mapstructure:"auth.jwt.secret"`
}

func (v *VaultSecrets) Init() (*VaultSecrets, error) {
	client, err := getVaultClient()
	if err != nil {
		return nil, err
	}
	secret, err := client.KVv2("api").Get(context.Background(), "staging")

	//fmt.Print(secret)
	if err != nil {
		return nil, err
	}

	//jsonString, _ := json.Marshal(secret.Data)
	var secrets VaultSecrets
	err = mapstructure.Decode(secret.Data, &secrets)
	if err != nil {
		fmt.Print(err)
	}

	//json.Unmarshal(jsonString, &secrets)
	//err = secret.Data.Decode(&secrets)

	if err != nil {
		return nil, err
	}
	return &secrets, nil
}

func getVaultClient() (*api.Client, error) {
	address := os.Getenv("VAULT_ADDR")
	if address == "" {
		address = "http://localhost:8200"
	}
	token := os.Getenv("root")
	if token == "" {
		token = "root"
	}

	config := api.DefaultConfig()
	config.Address = address
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
	client.SetToken(token)
	return client, nil
}

func main() {
	vaultSecrets, err := new(VaultSecrets).Init()
	if err != nil {
		log.Fatalf("Failed to init Vault secrets: %v", err)
	}
	fmt.Print("\n\n\n\n")
	fmt.Printf("JwtSecret: %v, KafkaBootstrap: %v\n", vaultSecrets.JwtSecret, vaultSecrets.KafkaBootstrap)

}
