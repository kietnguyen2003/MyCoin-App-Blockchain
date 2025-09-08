package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"math/big"

	"golang.org/x/crypto/ripemd160"
)

type Wallet struct {
	Address    string            `json:"address"`
	PrivateKey *ecdsa.PrivateKey `json:"-"`
	PublicKey  []byte            `json:"public_key"`
}

func NewWallet() *Wallet {
	// Tạo cặp khóa công khai và khóa riêng
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	publicKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)
	address := generateAddress(publicKey)

	return &Wallet{
		Address:    address,
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}
}

func generateAddress(publicKey []byte) string {
	sha256Hash := sha256.Sum256(publicKey)

	ripemd160Hasher := ripemd160.New()
	ripemd160Hasher.Write(sha256Hash[:])
	ripemd160Hash := ripemd160Hasher.Sum(nil)

	versionedPayload := append([]byte{0x00}, ripemd160Hash...)

	checksum := checksum(versionedPayload)
	fullPayload := append(versionedPayload, checksum...)

	return hex.EncodeToString(fullPayload)
}

func checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])
	return secondSHA[:4]
}

func (w *Wallet) GetPrivateKeyHex() string {
	return hex.EncodeToString(w.PrivateKey.D.Bytes())
}

func (w *Wallet) GetPublicKeyHex() string {
	return hex.EncodeToString(w.PublicKey)
}

func LoadWalletFromPrivateKey(privateKeyHex string) (*Wallet, error) {
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return nil, err
	}

	privateKey := new(ecdsa.PrivateKey)
	privateKey.PublicKey.Curve = elliptic.P256()
	privateKey.D = new(big.Int).SetBytes(privateKeyBytes)
	privateKey.PublicKey.X, privateKey.PublicKey.Y = privateKey.PublicKey.Curve.ScalarBaseMult(privateKeyBytes)

	publicKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)
	address := generateAddress(publicKey)

	return &Wallet{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		Address:    address,
	}, nil
}

func NewWalletFromPassphrase(passphrase string) *Wallet {
	hash := sha256.Sum256([]byte(passphrase))
	privateKey := new(ecdsa.PrivateKey)
	privateKey.PublicKey.Curve = elliptic.P256()
	privateKey.D = new(big.Int).SetBytes(hash[:])
	privateKey.PublicKey.X, privateKey.PublicKey.Y = privateKey.PublicKey.Curve.ScalarBaseMult(hash[:])

	publicKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)
	address := generateAddress(publicKey)

	return &Wallet{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		Address:    address,
	}
}
