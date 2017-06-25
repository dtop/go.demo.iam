package models

import (
	"crypto/rsa"

	"crypto/rand"

	"crypto/x509"
	"time"

	"log"

	"github.com/dtop/go.ginject"
	"github.com/mendsley/gojwk"
	"gopkg.in/redis.v4"
)

const (
	privateKeyName = "__privateKey"
)

type (

	// KeyProvider is the interface for a key provider
	KeyProvider interface {
		GetPrivateKey() *rsa.PrivateKey
		GetPublicKey() *rsa.PublicKey
		GetJWK() *gojwk.Key
	}

	keyProv struct {
		pr *rsa.PrivateKey
		pu *rsa.PublicKey

		Redis *redis.Client `inject:"redis"`
	}
)

// NewKeyProvider creates a new key provider
func NewKeyProvider(dep ginject.Injector) (KeyProvider, error) {

	kp := &keyProv{}
	if err := dep.Apply(kp); err != nil {
		panic(err)
	}

	err := kp.init()
	return kp, err
}

// ################### keyProv

func (kp *keyProv) init() (err error) {

	if val, err := kp.Redis.Get(privateKeyName).Result(); err == nil {

		pkey, err := x509.ParsePKCS1PrivateKey([]byte(val))
		if err == nil {

			kp.pr = pkey
			kp.pu = pkey.Public().(*rsa.PublicKey)
			return nil
		}
	}

	pkey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return
	}

	strRef := string(x509.MarshalPKCS1PrivateKey(pkey))
	kp.Redis.Set(privateKeyName, strRef, 86400*time.Second)

	kp.pr = pkey
	kp.pu = pkey.Public().(*rsa.PublicKey)

	return
}

// GetPrivateKey returns the private key
func (kp *keyProv) GetPrivateKey() *rsa.PrivateKey {
	return kp.pr
}

// GetPublicKey returns the public key
func (kp *keyProv) GetPublicKey() *rsa.PublicKey {
	return kp.pu
}

// GetJWK returns the JWK representation of the public key
func (kp *keyProv) GetJWK() *gojwk.Key {

	key, err := gojwk.PublicKey(kp.GetPublicKey())
	if err != nil {
		log.Println(err)
		return nil
	}

	key.Use = "sign"
	return key
}
