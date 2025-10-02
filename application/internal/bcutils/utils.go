package bcutils

import (
	"fmt"
	"os"

	"crypto/x509"

	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/thcrull/fabric-interface/application/pkg/config"
)

// NewGrpcConnection creates a new gRPC client connection to a Fabric peer.
// It loads the TLS certificate from the given config, adds it to a certificate pool,
// and returns a secure gRPC connection to the peer specified in cfg.Network.PeerEndpoint.
func NewGrpcConnection(cfg *config.Config) (*grpc.ClientConn, error) {
	tlsCertificatePEM, err := os.ReadFile(cfg.Network.TLSCertPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read TLS certificate: %w", err)
	}

	tlsCertificate, err := identity.CertificateFromPEM(tlsCertificatePEM)
	if err != nil {
		return nil, fmt.Errorf("failed to parse TLS certificate: %w", err)
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(tlsCertificate)
	transportCredentials := credentials.NewClientTLSFromCert(certPool, cfg.Network.TLSHostname)

	return grpc.NewClient(cfg.Network.PeerEndpoint, grpc.WithTransportCredentials(transportCredentials))
}

// NewIdentity creates a Fabric client identity using an X.509 certificate.
// It reads and parses the certificate from cfg.Identity.CertPath and associates it
// with the MSP ID specified in cfg.Identity.MspID.
func NewIdentity(cfg *config.Config) (*identity.X509Identity, error) {
	certificatePEM, err := os.ReadFile(cfg.Identity.CertPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate: %w", err)
	}

	certificate, err := identity.CertificateFromPEM(certificatePEM)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %w", err)
	}

	id, err := identity.NewX509Identity(cfg.Identity.MspID, certificate)
	if err != nil {
		return nil, fmt.Errorf("failed to create identity: %w", err)
	}

	return id, nil
}

// NewSign creates a digital signer using a private key.
// The private key is read from cfg.Identity.KeyPath and is used to sign transaction messages.
func NewSign(cfg *config.Config) (identity.Sign, error) {
	privateKeyPEM, err := os.ReadFile(cfg.Identity.KeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key: %w", err)
	}

	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	sign, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create signer: %w", err)
	}

	return sign, nil
}
