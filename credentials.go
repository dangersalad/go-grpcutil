package grpcutil

import (
	env "github.com/dangersalad/go-environment"
	"github.com/pkg/errors"
	"google.golang.org/grpc/credentials"
)

const (
	rootCert = "/etc/ssl/certs/ca-certificates.crt"
	// EnvKeyCrtFile is the environment variable that contains the
	// location of the SSL crt file. Defaults to "/ssl/tls.crt".
	EnvKeyCrtFile = "SSL_CRT_FILE"
	// EnvKeyKeyFile is the environment variable that contains the
	// location of the SSL key file. Defaults to "/ssl/tls.key".
	EnvKeyKeyFile = "SSL_KEY_FILE"
)

func getCertFiles() (crt, key string, err error) {
	crt = "/ssl/tls.crt"
	key = "/ssl/tls.key"
	conf, err := env.ReadOptions(env.Options{
		EnvKeyCrtFile: crt,
		EnvKeyKeyFile: key,
	})
	if err != nil {
		return "", "", errors.Wrap(err, "reading environment")
	}
	crt = conf[EnvKeyCrtFile]
	key = conf[EnvKeyKeyFile]
	return
}

// GetClientCredentials returns client credentials using the root CAs and a server name override
func GetClientCredentials(servername string) (credentials.TransportCredentials, error) {
	creds, err := credentials.NewClientTLSFromFile(rootCert, servername)
	if err != nil {
		return nil, errors.Wrap(err, "creating grpc client credentials")
	}
	return creds, nil
}

// GetServerCredentials returns server credentials obtained from the secret injected SSL cert
func GetServerCredentials() (credentials.TransportCredentials, error) {
	crt, key, err := getCertFiles()
	if err != nil {
		return nil, errors.Wrap(err, "getting tls filenames")
	}
	creds, err := credentials.NewServerTLSFromFile(crt, key)
	if err != nil {
		return nil, errors.Wrap(err, "creating grpc server credentials")
	}
	return creds, nil
}
