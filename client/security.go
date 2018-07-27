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

package main

import (
	"crypto/x509"
	"crypto/rsa"
	"encoding/base64"

	"github.com/yourfin/transcodebot/common"
)

//Set at compile time
var (
	serverCert *x509.Certificate
	clientKey *rsa.PrivateKey
	clientCert *x509.Certificate
	b64serverCert string
	b64clientPrivateKey string
	b64clientCert string
)

// Procedure:
//  unmarshalStaticVars
// Purpose:
//  To transform the base64 encoded static variables into what they represent
// Parameters:
//  None
// Produces:
//  Side effects:
//    serverCert, clientKey, and clientCert all set
// Preconditions:
//  This binary was built with the build flags as seen in:
//    github.com/yourfin/transcodebot/build.handleBuildCerts
// Postconditions:
//  all mentioned variables are unmarshaled into the variables they represent
func unmarshalStaticVars() {

}
