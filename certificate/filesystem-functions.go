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

package certificate

import (
	"encoding/pem"
	"io/ioutil"
	"errors"
	"crypto/x509"
	"path/filepath"
	"crypto/rsa"

	"github.com/yourfin/transcodebot/common"
)

// Procedure:
//  ReadRsaPrivKey
// Purpose:
//  To decode RSA private keys in $SettingsDir()/cert
// Parameters:
//  The name (sans .keyflie) of the key file: name string
// Produces:
//  key *rsa.PrivateKey
// Preconditions:
//  $SettingsDir() has been set
//  $SettingsDir()/cert/$name.keyfile exists and is a readable PEM encoded RSA prvate key
// Postconditions:
//  Errors are handled
//  key is the private key PEM encoded in $SettingsDir()/cert/$name.keyfile
func ReadRsaKey(name string) *rsa.PrivateKey {
	path := filepath.Join(common.SettingsDir(), "cert", name + ".crt")
	data, err := DecodePEMFile(path)
	if err != nil {
		common.PrintError("Read private key err:", err)
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(data)
	if err != nil {
		common.PrintError("Private key data err:", err)
	}
	return privateKey
}

// Procedure:
//  ReadCert
// Purpose:
//  To decode certificates in $SettingsDir()/cert
// Parameters:
//  The name (sans .crt) of the cert file: name string
// Produces:
//  cert *x509.Certificate
// Preconditions:
//  $SettingsDir()/cert exists
//  SettingsDir() is set
// Postconditions:
//  Any errors are handled
//  cert is the certificate PEM encoded in $SettingsDir()/cert/$name.crt
func ReadCert(name string) *x509.Certificate {
	path := filepath.Join(common.SettingsDir(), "cert", name + ".crt")
	data, err := DecodePEMFile(path)
	if err != nil {
		common.PrintError("Read certificate err:", err)
	}
	cert, err := x509.ParseCertificate(data)
	if err != nil {
		common.PrintError("Invalid certificatee data err:", err)
	}
	return cert
}

// Procedure:
//  DecodePEMFile
// Purpose:
//  Convert PEM encoded files into a byte array
// Parameters:
//  Absolute path to the PEM file: path string
// Produces:
//  output []byte
//  err error
// Preconditions:
//  There is a PEM encoded file at $path
//  The process has the read rights to the file at $path
// Postconditions:
//  Will error if there is no PEM data in the file and it is not empty
//  If any errors are generated, they are passed up through err and output will be empty
//  Output contains the contents of the file at $path decoded from PEM
func DecodePEMFile(path string) ([]byte, error) {
	rawData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, remaining := pem.Decode(rawData)
	if len(block.Bytes) == 0 && len(remaining) != 0 {
		return nil, errors.New("PEM Decode: no pem data found")
	}
	return block.Bytes, nil
}

// Procedure:
//	writeCertFile
// Purpose:
//	Helper for writing out cert files in $SettingsDir/cert/
// Parameters:
//	The data to be written to the file: data []byte
//  The file name to write to: fileName
// Produces:
//	File system side effects
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
