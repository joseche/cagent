package main

// https://golang.org/src/crypto/tls/generate_cert.go

import (
	"crypto/ecdsa"
	"encoding/pem"
	"crypto/x509"	
	"io/ioutil"
)

const (
	CERT_ORG string  = "Performance Labs"
	CERT_FILE string = MASTER_DIR + "perflabs.cert.pem"
	CERT_KEY string  = MASTER_DIR + "perflabs.key.pem"
)

var (
	host string
	ecdsaCurve = "P224"
)

func pemBlockForKey(priv ecdsa.PrivateKey) *pem.Block {
	b, err := x509.MarshalECPrivateKey(&priv)
	if err != nil {
		panic("Unable to marshal ECDSA private key: "+err.Error())
	}
	return &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}
}

func certSignature() (sig []byte) {
	pem, err := ioutil.ReadFile(CERT_FILE)
	if err != nil {
		panic("Error loading certificate: "+err.Error())
	}
	var cert* x509.Certificate
    cert, err = x509.ParseCertificate(pem)
    if err != nil {
		panic("Error parsing certificate: "+err.Error())
	}
    sig = cert.Signature
    return sig
}


