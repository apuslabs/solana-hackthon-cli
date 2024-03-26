package ca

import (
	"encoding/json"
	"fmt"
	"github.com/blocto/solana-go-sdk/types"
	"io"
	"os"
)

var pubkey_file = "/etc/apus-miner-pubkey-id.json"
var pubkey_id string

type KeyPaire struct {
	Pubkey string `json:"pubkey"`
}

func Init() {
	// 判断是否需要初始化 pubkeyid, 不存在则创建pubkeyid，存在则读取放到模块变量中
	_, err := os.Stat(pubkey_file)
	if err == nil {
		pubkey_id = ReadLocalKey()
	} else {
		pubkey_id = GenerateKey()
	}
}

func GenerateKey() string {
	wallet := types.NewAccount()
	pubkeyid := wallet.PublicKey.ToBase58()
	fmt.Printf("pubkey-id: %s\n", pubkeyid)

	keypaire := KeyPaire{Pubkey: pubkeyid}
	content, err := json.Marshal(keypaire)
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
	return pubkeyid
}

func ReadLocalKey() string {
	f, err := os.Open(pubkey_file)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	contentByte, err := io.ReadAll(f)
	fmt.Printf("%s", string(contentByte))
	var keyPaire KeyPaire
	err = json.Unmarshal(contentByte, &keyPaire)
	if err != nil {
		panic(err)
	}
	return keyPaire.Pubkey
}

func GetPubkey() string {
	return pubkey_id
}
