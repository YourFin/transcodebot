// Copyright © 2018 Patrick Nuckolls <nuckollsp at gmail>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package build

import (
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"crypto/rand"
	"encoding/pem"
	"time"
	"os"
	"net"

	"github.com/yourfin/transcodebot/common"
)

const (
	rootKeyFileName    string = "root.keyfile"
	rootCertFileName   string = "root.crt"
	clientKeyFileName  string = "client.keyfile"
	clientCertFileName string = "client.crt"
)

//Much here taken from https://ericchiang.github.io/post/go-tls

//Generate server and client certificates and dump them to files
//the destinations should be folders, not full paths
func GenRootCert(serverIPs []net.IP, clientDestination, serverDestination string) {
	common.PrintVerbose("Generating certificates...")
	//_, certPEM, _:= genRootCert(serverIPs)
	//common.Println(string(certPEM[:]))
	_, rootCertPEM, rootKey := genRootCert(serverIPs)

	writeCertFile(rootCertPEM, rootCertFileName)
	writeCertFile(privateKeyPEMify(rootKey), rootKeyFileName)
}

func genRootCert(serverIPs []net.IP) (*x509.Certificate, []byte, *rsa.PrivateKey) {
	rootKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		common.PrintError("certificate key err:", err)
	}

	rootCertTmpl := certTemplate()

	rootCertTmpl.IsCA = true
	rootCertTmpl.KeyUsage = x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature
	rootCertTmpl.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth}
	rootCertTmpl.IPAddresses = serverIPs
	cert, certPEM := createCert(rootCertTmpl, rootCertTmpl, &rootKey.PublicKey, rootKey)
	return cert, certPEM, rootKey
}

func createCert(template, parent *x509.Certificate, pub, parentPriv interface{}) (*x509.Certificate, []byte) {
	certDER, err := x509.CreateCertificate(rand.Reader, template, parent, pub, parentPriv)
	if err != nil {
		common.PrintError("root cert gen err:", err)
	}

	//Parse resulting cert for re-use later
	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		common.PrintError("Parse Certificate err:", err)
	}

	//PEM encode the certificate (adds the --BEGIN CERT stuff)
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	return cert, certPEM
}

//Turns private key into file storeable form
func privateKeyPEMify(privateKey *rsa.PrivateKey) []byte {
	return pem.EncodeToMemory(&pem.Block{
		Type: "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
}

// Procedure:
//	writeCertFile
// Purpose:
//	Helper for writing out cert files in $SettingsDir/cert/
// Parameters:
//	The data to be written to the file: data []byte
//  The file name to write to: fileName
// Produces:
//	Side effects
// Preconditions:
//  fileName is a valid filename on the system
//  fileName doesn't contain path seperator characters
// Postconditions:
//  $SettingsDir/cert/$fileName contains $data, or an error is printed
func writeCertFile(data []byte, fileName string) {
	err := common.SettingsWriteFile(data, "cert", fileName)
	if err != nil {
		common.PrintError("Writing Cert file err:", err)
	}
}

func certTemplate() *x509.Certificate {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)

	if err != nil {
		common.PrintError("certificate secure random err:", err)
	}

	hostname, err := os.Hostname()

	tmpl := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{Organization: []string{"transcodebot-" + hostname}},
		SignatureAlgorithm: x509.SHA256WithRSA,
		NotBefore: time.Now(),
		//*Supposedly* the tls protocol is implemented such that
		//certs can't be valid past 2049
		//see www-01.ibm.com/support/docview.wss?uid=swg21220045
		//I'm not totally sure of this, however, as that is dated from 2012
		NotAfter: time.Date(2049, time.December, 1, 1, 1, 1, 1, time.UTC),
		BasicConstraintsValid: true,
	}
	return &tmpl
}
