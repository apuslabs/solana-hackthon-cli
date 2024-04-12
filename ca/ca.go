package ca

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/blocto/solana-go-sdk/types"
)

var pubkey_file = "/etc/apus-miner-pubkey-id.json"
var inited = false

type KeyPaire struct {
	Pubkey    string `json:"pubkey"`
	SecretKey []byte `json:"secretKey"`
}

var keyPaire KeyPaire

func Init() {
	// 判断是否需要初始化 pubkeyid, 不存在则创建pubkeyid，存在则读取放到模块变量中
	_, err := os.Stat(pubkey_file)
	if err == nil {
		ReadLocalKey()
		inited = true
	} else {
		// GenerateKey()
	}
}

func GenerateKey() {
	wallet := types.NewAccount()
	pubkeyid := wallet.PublicKey.ToBase58()
	fmt.Printf("pubkey-id: %s\n", pubkeyid)

	keyPaire = KeyPaire{Pubkey: pubkeyid, SecretKey: wallet.PrivateKey.Seed()}
	content, err := json.Marshal(keyPaire)
	if err != nil {
		panic(err)
	}
	f, err := os.Create(pubkey_file)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	_, err = f.Write(content)
	if err != nil {
		panic(err)
	}
}

func SaveLocalKey(pubKey string, secretKey string) {
	keyPaire = KeyPaire{Pubkey: pubKey, SecretKey: []byte(secretKey)}
	content, err := json.Marshal(keyPaire)
	if err != nil {
		panic(err)
	}
	// check if file exists
	_, err = os.Stat(pubkey_file)
	if err != nil {
		if _, ok := err.(*os.PathError); ok {
			// create file
			f, err := os.Create(pubkey_file)
			if err != nil {
				panic(err)
			}
			defer func() {
				if cerr := f.Close(); cerr != nil {
					panic(cerr)
				}
			}()
		} else {
			panic(err)
		}
	}
	// read file
	f, err := os.OpenFile(pubkey_file, os.O_WRONLY, 0666)
	defer func() {
		if cerr := f.Close(); cerr != nil {
			panic(cerr)
		}
	}()
	if err != nil {
		panic(err)
	}
	// write file
	_, err = f.Write(content)
	if err != nil {
		panic(err)
	}
}

func ReadLocalKey() {
	f, err := os.Open(pubkey_file)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	contentByte, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(contentByte, &keyPaire)
	fmt.Printf("pubkey: %s\n", keyPaire.Pubkey)
	if err != nil {
		panic(err)
	}
}

func GetPubkey() KeyPaire {
	return keyPaire
}

func Inited() bool {
	return inited
}
