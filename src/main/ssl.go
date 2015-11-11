package main

// https://golang.org/src/crypto/tls/generate_cert.go

/*
I am using P224 Elliptic Curve Diffie Hellman keys, 224 bits long, 
Equivalent to RSA 2048 bits. This reduces cycles while ensuring high level encryption
*/

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"crypto/md5"
)

const (
	CERT_ORG  string = "Performance Labs"
	CERT_FILE string = MASTER_DIR + "perflabs.cert.pem"
	CERT_KEY  string = MASTER_DIR + "perflabs.key.pem"
)

func pemBlockForKey(priv ecdsa.PrivateKey) *pem.Block {
	b, err := x509.MarshalECPrivateKey(&priv)
	if err != nil {
		panic("Unable to marshal ECDSA private key: " + err.Error())
	}
	return &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}
}

func certSignature() (sig []byte) {
	pemfile, err := ioutil.ReadFile(CERT_FILE)
	if err != nil {
		panic("Error loading certificate: " + err.Error())
	}
	
	crt, rest := pem.Decode(pemfile)
	if crt == nil {
		panic("unable to decode pem certificate: "+string(rest))
	}
	
	var cert *x509.Certificate
	cert, err = x509.ParseCertificate(crt.Bytes)
	if err != nil {
		panic("Error parsing certificate: " + err.Error())
	}
	sig = cert.Signature
	return sig
}

func hostid() (id string){
	sig := certSignature()
	id = fmt.Sprintf("%x", md5.Sum(sig) )
	return id
}