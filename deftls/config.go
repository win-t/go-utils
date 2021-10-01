// Package deftls.
//
// this package contain context some useful function for
// configuring tls.Config.
package deftls

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"math/big"
	"net"
	"time"

	"github.com/payfazz/go-errors/v2"
)

type Option func(*tls.Config) error

// Config return default tls config with some options.
func Config(opts ...Option) (*tls.Config, error) {
	config := &tls.Config{
		PreferServerCipherSuites: true,

		MinVersion: tls.VersionTLS12,

		CipherSuites: []uint16{
			// TLS 1.3
			tls.TLS_CHACHA20_POLY1305_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_AES_128_GCM_SHA256,

			// TLS 1.2
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},

		CurvePreferences: []tls.CurveID{
			tls.X25519,
			tls.CurveP256,
			tls.CurveP384,
			tls.CurveP521,
		},
	}

	for _, o := range opts {
		if o != nil {
			if err := o(config); err != nil {
				return config, errors.Trace(err)
			}
		}
	}

	return config, nil
}

// UseCertFile option.
func UseCertFile(certfile, keyfile string) Option {
	return func(config *tls.Config) error {
		cert, err := tls.LoadX509KeyPair(certfile, keyfile)
		if err != nil {
			return errors.Trace(err)
		}
		config.Certificates = []tls.Certificate{cert}
		return nil
	}
}

// UseCertPem option.
func UseCertPem(certpem, keypem string) Option {
	return func(config *tls.Config) error {
		cert, err := tls.X509KeyPair([]byte(certpem), []byte(keypem))
		if err != nil {
			return errors.Trace(err)
		}
		config.Certificates = []tls.Certificate{cert}
		return nil
	}
}

// UseCertSelfSigned option.
func UseCertSelfSigned() Option {
	return func(config *tls.Config) error {
		privkey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			return errors.Trace(err)
		}

		var data [16]byte
		if _, err := rand.Read(data[:]); err != nil {
			return errors.Trace(err)
		}

		subject := "self-signed-" + hex.EncodeToString(data[:])
		notBefore := time.Now()
		notAfter := notBefore.Add(168 * time.Hour) // 1 week

		template := &x509.Certificate{
			SerialNumber: big.NewInt(1),
			Subject:      pkix.Name{CommonName: subject},
			NotBefore:    notBefore,
			NotAfter:     notAfter,
			KeyUsage:     x509.KeyUsageDigitalSignature,
			ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
			DNSNames:     []string{"localhost"},
			IPAddresses: []net.IP{
				net.IPv6loopback,
				net.IPv4(127, 0, 0, 1),
			},
		}

		derCert, err := x509.CreateCertificate(rand.Reader, template, template, privkey.Public(), privkey)
		if err != nil {
			return errors.Trace(err)
		}

		cert := tls.Certificate{
			Certificate: [][]byte{derCert},
			PrivateKey:  privkey,
		}

		config.Certificates = []tls.Certificate{cert}
		return nil
	}
}
