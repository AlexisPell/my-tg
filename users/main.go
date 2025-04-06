package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/rs/cors"
	"golang.org/x/crypto/curve25519"
)

type KeyPair struct {
	Private string `json:"private"`
	Public  string `json:"public"`
}

type UserKeys struct {
	IdentityKey  KeyPair   `json:"identity_key"`
	SignedPreKey KeyPair   `json:"signed_pre_key"`
	OneTimeKeys  []KeyPair `json:"one_time_keys"`
}

func generateKeyPair() (KeyPair, error) {
	var priv, pub [32]byte
	_, err := rand.Read(priv[:])
	if err != nil {
		return KeyPair{}, err
	}
	curve25519.ScalarBaseMult(&pub, &priv)

	return KeyPair{
		Private: base64.StdEncoding.EncodeToString(priv[:]),
		Public:  base64.StdEncoding.EncodeToString(pub[:]),
	}, nil
}

func generateUserKeys() UserKeys {
	ik, _ := generateKeyPair()
	spk, _ := generateKeyPair()

	var otpk []KeyPair
	for i := 0; i < 5; i++ {
		key, _ := generateKeyPair()
		otpk = append(otpk, key)
	}

	keys := UserKeys{
		IdentityKey:  ik,
		SignedPreKey: spk,
		OneTimeKeys:  otpk,
	}

	jsonData, _ := json.MarshalIndent(keys, "", "  ")
	fmt.Println(string(jsonData))

	return keys
}

func main() {
	// static for now
	uk := generateUserKeys()

	mux := http.NewServeMux()
	mux.HandleFunc("/keys", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(uk)
	})

	// Настройка CORS
	handler := cors.Default().Handler(mux)

	fmt.Println(">>> Users service is listening to 8002")
	if err := http.ListenAndServe(":8002", handler); err != nil {
		log.Fatal(">>> Error starting users service:", err)
	}
}
