package crypto

// nolint:gosec // md5/sha1 hash functions used by some exchanges
import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/quantstop/quantstopterminal/internal/log"
	"github.com/quantstop/quantstopterminal/pkg/system/file"
	"hash"
	"io"
	"io/ioutil"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"
)

// const declarations for common crypto operations
const (
	HashSHA1 = iota
	HashSHA256
	HashSHA512
	HashSHA512_384
	HashMD5
)

var (
	errCertExpired     = errors.New("gRPC TLS certificate has expired")
	errCertDataIsNil   = errors.New("gRPC TLS certificate PEM data is nil")
	errCertTypeInvalid = errors.New("gRPC TLS certificate type is invalid")
)

// HexEncodeToString takes in a hexadecimal byte array and returns a string
func HexEncodeToString(input []byte) string {
	return hex.EncodeToString(input)
}

// Base64Decode takes in a Base64 string and returns a byte array and an error
func Base64Decode(input string) ([]byte, error) {
	result, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Base64Encode takes in a byte array then returns an encoded base64 string
func Base64Encode(input []byte) string {
	return base64.StdEncoding.EncodeToString(input)
}

// GetRandomSalt returns a random salt
func GetRandomSalt(input []byte, saltLen int) ([]byte, error) {
	if saltLen <= 0 {
		return nil, errors.New("salt length is too small")
	}
	salt := make([]byte, saltLen)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}

	var result []byte
	if input != nil {
		result = input
	}
	result = append(result, salt...)
	return result, nil
}

// GetMD5 returns a MD5 hash of a byte array
func GetMD5(input []byte) ([]byte, error) {
	m := md5.New() // nolint:gosec // hash function used by some exchanges
	_, err := m.Write(input)
	return m.Sum(nil), err
}

// GetSHA512 returns a SHA512 hash of a byte array
func GetSHA512(input []byte) ([]byte, error) {
	sha := sha512.New()
	_, err := sha.Write(input)
	return sha.Sum(nil), err
}

// GetSHA256 returns a SHA256 hash of a byte array
func GetSHA256(input []byte) ([]byte, error) {
	sha := sha256.New()
	_, err := sha.Write(input)
	return sha.Sum(nil), err
}

// GetHMAC returns a keyed-hash message authentication code using the desired hash type
func GetHMAC(hashType int, input, key []byte) ([]byte, error) {
	var hasher func() hash.Hash

	switch hashType {
	case HashSHA1:
		hasher = sha1.New
	case HashSHA256:
		hasher = sha256.New
	case HashSHA512:
		hasher = sha512.New
	case HashSHA512_384:
		hasher = sha512.New384
	case HashMD5:
		hasher = md5.New
	}

	h := hmac.New(hasher, key)
	_, err := h.Write(input)
	return h.Sum(nil), err
}

// Sha1ToHex takes a string, sha1 hashes it and return a hex string of the result
func Sha1ToHex(data string) (string, error) {
	h := sha1.New() // nolint:gosec // hash function used by some exchanges
	_, err := h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil)), err
}

// GetTLSDir returns the default TLS dir
func GetTLSDir(dir string) string {
	return filepath.Join(dir, "tls")
}

// CheckCerts will check to see if certificates exist in the supplied directory,
// and then verify that the certificate is not expired. If no certificates exist,
// or they are expired, it will create new self-signed certificates.
func CheckCerts(certDir string) error {
	certFile := filepath.Join(certDir, "cert.pem")
	keyFile := filepath.Join(certDir, "key.pem")

	if !file.Exists(certFile) || !file.Exists(keyFile) {
		log.Warnln(log.Global, "TLS certificate/key file missing, recreating...")
		return genCert(certDir)
	}

	pemData, err := ioutil.ReadFile(certFile)
	if err != nil {
		return fmt.Errorf("unable to open TLS certificate file: %s", err)
	}

	if err = verifyCert(pemData); err != nil {
		if err != errCertExpired {
			return err
		}
		log.Warnln(log.Global, "TLS certificate has expired, regenerating...")
		return genCert(certDir)
	}

	log.Infoln(log.Global, "TLS certificate and key files exist, will use them.")
	return nil
}

// verifyCert will check to make sure a certificate is not expired.
// If the certificate is expired, will return error.
func verifyCert(pemData []byte) error {
	var pemBlock *pem.Block
	pemBlock, _ = pem.Decode(pemData)
	if pemBlock == nil {
		return errCertDataIsNil
	}

	if pemBlock.Type != "CERTIFICATE" {
		return errCertTypeInvalid
	}

	cert, err := x509.ParseCertificate(pemBlock.Bytes)
	if err != nil {
		return err
	}

	if time.Now().After(cert.NotAfter) {
		return errCertExpired
	}
	return nil
}

// genCert will create a new certificate pair, in the supplied directory.
func genCert(targetDir string) error {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return fmt.Errorf("failed to generate ecdsa private key: %s", err)
	}

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return fmt.Errorf("failed to generate serial number: %s", err)
	}

	host, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("failed to get hostname: %s", err)
	}

	dnsNames := []string{host}
	if host != "localhost" {
		dnsNames = append(dnsNames, "localhost")
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"quantstop.com"},
			CommonName:   host,
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour * 24 * 365),
		IsCA:                  true,
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IPAddresses: []net.IP{
			net.ParseIP("127.0.0.1"),
			net.ParseIP("::1"),
		},
		DNSNames: dnsNames,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privKey.PublicKey, privKey)
	if err != nil {
		return fmt.Errorf("failed to create certificate: %s", err)
	}

	certData := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	if certData == nil {
		return fmt.Errorf("cert data is nil")
	}

	b, err := x509.MarshalECPrivateKey(privKey)
	if err != nil {
		return fmt.Errorf("failed to marshal ECDSA private key: %s", err)
	}

	keyData := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: b})
	if keyData == nil {
		return fmt.Errorf("key pem data is nil")
	}

	err = file.Write(filepath.Join(targetDir, "key.pem"), keyData)
	if err != nil {
		return fmt.Errorf("failed to write key.pem file %s", err)
	}

	err = file.Write(filepath.Join(targetDir, "cert.pem"), certData)
	if err != nil {
		return fmt.Errorf("failed to write cert.pem file %s", err)
	}

	log.Infof(log.Global, "gRPC TLS key.pem and cert.pem files written to %s\n", targetDir)
	return nil
}
