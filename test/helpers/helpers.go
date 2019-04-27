package helpers

import (
	cryptorand "crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math/big"
	"math/rand"
	"net"
	"net/http"
	"time"
)

const randomStringLength = 10
const randomStringChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var HTTPClient = &http.Client{
	Timeout: 5 * time.Second,
	Transport: &http.Transport{
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		Dial:                (&net.Dialer{Timeout: 5 * time.Second}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
	},
}

func GetFreeAddr() string {
	conn, err := net.Listen("tcp", "127.0.0.1:0")
	panicOnError(err)
	defer conn.Close()
	return conn.Addr().String()
}

func ServerIsAvailable(address string) bool {
	conn, err := net.Dial("tcp", address)
	if err == nil {
		tls.Client(conn, &tls.Config{InsecureSkipVerify: true}).Handshake()
		conn.Close()
		return true
	}
	return false
}

func WaitForServerToBeAvailable(address string, timeout time.Duration) error {
	timeoutChan := time.After(timeout)
	for {
		select {
		case <-timeoutChan:
			return fmt.Errorf("failed to connect to %s within %s", address, timeout)
		default:
			if ServerIsAvailable(address) {
				return nil
			}
		}
	}
}

func RandomString() string {
	b := make([]byte, randomStringLength)
	for i := range b {
		b[i] = randomStringChars[rand.Intn(len(randomStringChars))]
	}
	return string(b)
}

func PrivateKeyToPEM(privateKey *rsa.PrivateKey) string {
	privateKeyPEM := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}
	privateKeyBytes := pem.EncodeToMemory(privateKeyPEM)
	return string(privateKeyBytes)
}

func PublicKeyToPEM(publicKey *rsa.PublicKey) string {
	asn1Bytes, err := x509.MarshalPKIXPublicKey(publicKey)
	panicOnError(err)

	publicKeyPEM := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	}
	publicKeyBytes := pem.EncodeToMemory(publicKeyPEM)
	return string(publicKeyBytes)
}

func GenerateSelfSignedCert(commonName string) (string, string) {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := cryptorand.Int(cryptorand.Reader, serialNumberLimit)
	if err != nil {
		panic(err)
	}

	template := &x509.Certificate{
		SerialNumber: serialNumber,
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		Subject:      pkix.Name{CommonName: commonName},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(time.Hour * 24 * 365),
	}

	key, err := rsa.GenerateKey(cryptorand.Reader, 4096)
	if err != nil {
		panic(err)
	}

	derBytes, err := x509.CreateCertificate(cryptorand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}

	certPath, err := ioutil.TempFile("", "cert")
	if err != nil {
		panic(err)
	}
	if err := pem.Encode(certPath, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		panic(err)
	}

	keyPath, err := ioutil.TempFile("", "key")
	if err != nil {
		panic(err)
	}
	if err := pem.Encode(keyPath, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)}); err != nil {
		panic(err)
	}

	return certPath.Name(), keyPath.Name()
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
