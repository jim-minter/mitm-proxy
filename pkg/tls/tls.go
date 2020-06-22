package tls

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"io/ioutil"
	"math/big"
	"sync"
	"time"
)

// DB represents a TLS database
type DB struct {
	mu sync.Mutex

	cakey  *rsa.PrivateKey
	CACert *x509.Certificate
}

// NewDB instantiates a new TLS database
func NewDB() (*DB, error) {
	cakey, cacert, err := getKeyAndCertificate("ca", nil, nil, true)
	if err != nil {
		return nil, err
	}

	return &DB{
		cakey:  cakey,
		CACert: cacert,
	}, nil
}

// GetKeyAndCertificate returns a cached key and certificate, or creates and
// caches new ones
func (db *DB) GetKeyAndCertificate(commonName string) (*rsa.PrivateKey, *x509.Certificate, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	return getKeyAndCertificate(commonName, db.cakey, db.CACert, false)
}

func getKeyAndCertificate(commonName string, parentKey *rsa.PrivateKey, parentCert *x509.Certificate, isCA bool) (*rsa.PrivateKey, *x509.Certificate, error) {
	key, cert, err := readKeyAndCertificate(commonName)
	if err == nil {
		return key, cert, nil
	}

	key, cert, err = generateKeyAndCertificate(commonName, parentKey, parentCert, isCA)
	if err != nil {
		return nil, nil, err
	}

	err = writeKeyAndCertificate(key, cert)
	if err != nil {
		return nil, nil, err
	}

	return key, cert, nil
}

func readKeyAndCertificate(commonName string) (*rsa.PrivateKey, *x509.Certificate, error) {
	keydata, err := ioutil.ReadFile(commonName + "-key.der")
	if err != nil {
		return nil, nil, err
	}

	key, err := x509.ParsePKCS1PrivateKey(keydata)
	if err != nil {
		return nil, nil, err
	}

	certdata, err := ioutil.ReadFile(commonName + "-cert.der")
	if err != nil {
		return nil, nil, err
	}

	cert, err := x509.ParseCertificate(certdata)
	if err != nil {
		return nil, nil, err
	}

	return key, cert, nil
}

func generateKeyAndCertificate(commonName string, parentKey *rsa.PrivateKey, parentCert *x509.Certificate, isCA bool) (*rsa.PrivateKey, *x509.Certificate, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, nil, err
	}

	now := time.Now().UTC().Add(-time.Hour)
	notAfter := now.AddDate(10, 0, 0)

	if parentCert != nil && parentCert.NotAfter.Before(notAfter) {
		notAfter = parentCert.NotAfter
	}

	template := &x509.Certificate{
		SerialNumber:          serialNumber,
		NotBefore:             now,
		NotAfter:              notAfter,
		Subject:               pkix.Name{CommonName: commonName},
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		IsCA:                  isCA,
	}

	if isCA {
		template.KeyUsage |= x509.KeyUsageCertSign
	} else {
		template.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
	}

	if parentCert == nil && parentKey == nil {
		parentCert = template
		parentKey = key
	}

	b, err := x509.CreateCertificate(rand.Reader, template, parentCert, &key.PublicKey, parentKey)
	if err != nil {
		return nil, nil, err
	}

	cert, err := x509.ParseCertificate(b)
	if err != nil {
		return nil, nil, err
	}

	return key, cert, nil
}

func writeKeyAndCertificate(key *rsa.PrivateKey, cert *x509.Certificate) error {
	err := ioutil.WriteFile(cert.Subject.CommonName+"-key.der", x509.MarshalPKCS1PrivateKey(key), 0600)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(cert.Subject.CommonName+"-cert.der", cert.Raw, 0666)
}
