package main

import (
	"os"
	"github.com/mxk/go-sqlite/sqlite3"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"github.com/leekchan/timeutil"
	"math/big"
	"net"
	"time"
	"strings"
)


func init_dir(){
	dir_exists,_ := File_exists(MASTER_DIR)
	if ! dir_exists {
	  	err := os.MkdirAll(MASTER_DIR,0755)
  		if err != nil {
  			panic("Can't create "+MASTER_DIR+", "+err.Error())
  		}
  		Info(MASTER_DIR +" created")
  	}
}

func init_db(){
	file_exists,_ := File_exists(DB_FILE)
	if ! file_exists {
		Info("Local DB doesn't exist, creating new")
	}
	  
	conn, err := sqlite3.Open( DB_FILE )
	if err != nil {
		panic(err)
	}
	  
	create_table(CPUTIMES_TB, CPUTIMES_CREATE, conn)
	create_table(LOADAVG_TB, LOADAVG_CREATE, conn)
	create_table(MEMORY_TB, MEMORY_CREATE, conn)
	  
	err = conn.Close()
	if err != nil {
	  	Err("Closing local DB: "+err.Error())
	}
}

func init_ssl() {
	key_exists, _ := File_exists(CERT_KEY)
	if key_exists {
		return 
	}
	
	host, _ = os.Hostname()
	priv, err := ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
	if err != nil {
		panic("failed to generate private key: " + err.Error())
	}

	var now time.Time
	now = time.Now()
	td1 := timeutil.Timedelta{Days: 7}
	notBefore := now.Add( -td1.Duration() )
	td2 := timeutil.Timedelta{Weeks: 480} // 480 weeks ~= 10 years
	notAfter  := notBefore.Add( td2.Duration() ) 
	
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		Err("failed to generate serial number: "+ err.Error())
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{CERT_ORG},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	hosts := strings.Split(host, ",")
	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
	}

	template.IsCA = true
	template.KeyUsage |= x509.KeyUsageCertSign
	
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		panic("Failed to create certificate: " + err.Error())
	}

	certOut, err := os.Create(CERT_FILE)
	if err != nil {
		panic("failed to open "+CERT_FILE+" for writing: " + err.Error())
	}
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certOut.Close()
	Info("written "+CERT_FILE)

	keyOut, err := os.OpenFile(CERT_KEY, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		panic("failed to open "+CERT_KEY+" for writing: "+err.Error())
	}
	pem.Encode(keyOut, pemBlockForKey(*priv))
	keyOut.Close()
	Info("written "+CERT_KEY)
}


func initialize(){
	init_dir()
	init_db()
	init_ssl()
	host = hostid()
}