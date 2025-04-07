package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"crypto/ed25519"

	"github.com/rs/cors"
	"golang.org/x/crypto/curve25519"
)

type KeyPair struct {
	Private string `json:"private"`
	Public  string `json:"public"`
}

/*
@description:

	IdentityKey is used for DH encrypting.
	SignedPreKey is used for DH encrypting. TODO: Refreshed daily/weekly/etc...
	OTSK is used for DH encrypting. TODO: Refreshes every time when message is sent.

	IK_Sign is used for signing SignedPreKey.
	SPK_signature is signature of SignedPreKey by IK_Sign private key.
*/
type UserKeys struct {
	IdentityKey  KeyPair   `json:"identity_key"`
	SignedPreKey KeyPair   `json:"signed_pre_key"`
	OneTimeKeys  []KeyPair `json:"one_time_keys"`

	IdentityKeySigner KeyPair `json:"identity_key_signer"`
	SPKSignatureB64   string  `json:"spk_signature"`
}

func generateCurve25519KeyPair() (KeyPair, error) {
	var priv [32]byte
	var pub [32]byte
	_, err := rand.Read(priv[:])
	if err != nil {
		return KeyPair{}, err
	}
	// Create ECDH key pair
	curve25519.ScalarBaseMult(&pub, &priv)

	return KeyPair{
		Private: base64.StdEncoding.EncodeToString(priv[:]),
		Public:  base64.StdEncoding.EncodeToString(pub[:]),
	}, nil
}

func generateUserKeys() UserKeys {
	ik, _ := generateCurve25519KeyPair()
	spk, _ := generateCurve25519KeyPair()

	var otpk []KeyPair
	for i := 0; i < 5; i++ {
		key, _ := generateCurve25519KeyPair()
		otpk = append(otpk, key)
	}

	ikSignerPub, ikSignerPriv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}
	spkPubBytes, _ := base64.StdEncoding.DecodeString(spk.Public)
	signature := ed25519.Sign(ikSignerPriv, spkPubBytes)

	keys := UserKeys{
		IdentityKey:  ik,
		SignedPreKey: spk,
		OneTimeKeys:  otpk,
		IdentityKeySigner: KeyPair{
			Private: base64.StdEncoding.EncodeToString(ikSignerPriv),
			Public:  base64.StdEncoding.EncodeToString(ikSignerPub),
		},
		SPKSignatureB64: base64.StdEncoding.EncodeToString(signature),
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

	// CORS
	handler := cors.Default().Handler(mux)

	fmt.Println(">>> Users service is listening to 8002")
	if err := http.ListenAndServe(":8002", handler); err != nil {
		log.Fatal(">>> Error starting users service:", err)
	}
}
