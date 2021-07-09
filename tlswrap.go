package tlswrap

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/big"
	"net"
	"net/http"
	"os"
	"time"

	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
)

var (
	ACMELetsEncryptStagingUrl = "https://acme-staging-v02.api.letsencrypt.org/directory"
	DevCertName               = "localhost"
)

// Config holds
// certCache, the path to the certs
// domains the whiteliste of domains allowd for tls. An empty list is a wildcard
// the directoryURL for requesting certs.
type Config struct {
	certCache    string
	domains      []string
	directoryURL string
}

// NewConfig returns a default config
func NewConfig(certCache string, domains []string) Config {
	return Config{
		certCache:    certCache,
		domains:      domains,
		directoryURL: autocert.DefaultACMEDirectory,
	}
}

// NewStageConfig returns a stage/test config, use this when testing tls provider
func NewStageConfig(certCache string, domains []string) Config {
	return Config{
		certCache:    certCache,
		domains:      domains,
		directoryURL: ACMELetsEncryptStagingUrl,
	}
}

// NewStageConfigFromConfig returns a stage/test config from a config
func NewStageConfigFromConfig(config Config) Config {
	return Config{
		certCache:    config.certCache,
		domains:      config.domains,
		directoryURL: ACMELetsEncryptStagingUrl,
	}
}

// NewServer returns a server with TLS self provision and a handler for port 80
// the handler takes care of acme challanges and redirects to https
// the handler blocks, so start in a go rutine
func NewServer(address string, handler http.Handler, config Config) (http.Server, func()) {
	server := http.Server{
		Addr:    address,
		Handler: handler,
	}
	manager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(config.domains...),
		Cache:      autocert.DirCache(config.certCache),
		Client:     &acme.Client{DirectoryURL: config.directoryURL},
	}

	httpHandlerFunc := func() {
		err := http.ListenAndServe("0.0.0.0:80", manager.HTTPHandler(nil))
		if err != nil {
			log.Fatal("Critical. Could not start port 80 handler due to", err)
		}
	}

	// TLS Ciphers
	server.TLSConfig = manager.TLSConfig()
	server.TLSConfig.MinVersion = tls.VersionTLS12
	server.TLSConfig.CurvePreferences = []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256}
	server.TLSConfig.CipherSuites = []uint16{
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256, // needed for http/2
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
	}
	// server.TLSConfig.PreferServerCipherSuites.TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0), disable http/s (uses 128 bit)
	server.TLSConfig.PreferServerCipherSuites = true

	return server, httpHandlerFunc
}

// NewDevServer returns a dev server that uses a self-signed cert
func NewDevServer(address string, handler http.Handler, config Config) (http.Server, error) {
	// TODO check if certs exists and not expired
	// get the private key and cert, if it does not exists, creates it
	keyPath := fmt.Sprintf("%s/%s.key", config.certCache, DevCertName)
	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return http.Server{}, fmt.Errorf("failed to read cert: %s", err)
		}
		if err = CreateSelfSignedCert(config.certCache); err != nil {
			return http.Server{}, err
		}
		key, err = ioutil.ReadFile(keyPath)
		if err != nil {
			return http.Server{}, fmt.Errorf("failed to read cert after creation %s", err)
		}
	}
	certPath := fmt.Sprintf("%s/%s.pem", config.certCache, DevCertName)
	pem, err := ioutil.ReadFile(certPath)
	if err != nil {
		return http.Server{}, fmt.Errorf("failed to read private key: %s", err)
	}
	// configure a server instance with mostly safe ciphers and tls
	cert, err := tls.X509KeyPair(pem, key)
	if err != nil {
		return http.Server{}, err
	}
	server := http.Server{
		Addr:    address,
		Handler: handler,
		TLSConfig: &tls.Config{
			Certificates:     []tls.Certificate{cert},
			MinVersion:       tls.VersionTLS12,
			CurvePreferences: []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256, // needed for http/2
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
			},
			// 		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0), disable http/s (uses 128 bit)
			PreferServerCipherSuites: true,
		},
	}

	return server, err
}

// StartServer starts a production server
func StartServer(address string, config Config) error {
	return StartServerWithHandler(address, config, nil)
}

// StartServerHandler starts a production server
func StartServerWithHandler(address string, config Config, handler http.Handler) error {
	server, httpHandler := NewServer(address, handler, config)
	go httpHandler()

	return server.ListenAndServeTLS("", "")
}

// StartStageServer starts a server with stage acme endpoint
func StartStageServer(address string, config Config) error {
	return StartStageServerWithHandler(address, config, nil)
}

// StartStageServerWithHandler starts a server with stage acme endpoint
func StartStageServerWithHandler(address string, config Config, handler http.Handler) error {
	config = NewStageConfigFromConfig(config)
	server, httpHandler := NewServer(address, handler, config)
	go httpHandler()

	return server.ListenAndServeTLS("", "")
}

// StartDevServer starts a development server with self signed TLS.
func StartDevServer(address string, config Config) error {
	return StartDevServerWithHandler(address, config, nil)
}

// StartDevServerWithHandler starts a development server with self signed TLS
func StartDevServerWithHandler(address string, config Config, handler http.Handler) error {
	server, err := NewDevServer(address, handler, config)
	if err != nil {
		return fmt.Errorf("failed to start development server: %s", err)
	}
	return server.ListenAndServeTLS("", "")
}

// createSelfSignedCert creates a snakeoil certificate and saves it at supplied path
func CreateSelfSignedCert(path string) error {
	certPath := fmt.Sprintf("%s/%s.pem", path, DevCertName)
	privKeyPath := fmt.Sprintf("%s/%s.key", path, DevCertName)
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			Organization:  []string{"Dev"},
			Country:       []string{"Dev"},
			Province:      []string{"Dev"},
			Locality:      []string{"Dev"},
			StreetAddress: []string{"Dev"},
			PostalCode:    []string{"0"},
		},
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}
	certPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}
	certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &certPrivKey.PublicKey, certPrivKey)
	if err != nil {
		return err
	}

	certPEM := new(bytes.Buffer)
	pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})
	certPrivKeyPEM := new(bytes.Buffer)
	pem.Encode(certPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certPrivKey),
	})

	certPrivKeyPEMBytes, err := io.ReadAll(certPrivKeyPEM)
	if err != nil {
		return err
	}
	err = os.WriteFile(privKeyPath, certPrivKeyPEMBytes, os.ModePerm)
	if err != nil {
		return err
	}

	certPEMBytes, err := io.ReadAll(certPEM)
	if err != nil {
		return err
	}
	err = os.WriteFile(certPath, certPEMBytes, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
