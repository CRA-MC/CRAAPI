package cmd

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"

	"github.com/spf13/cobra"
)

var Bits int

// GenRsaKey generates an PKCS#1 RSA keypair of the given bit size in PEM format.
func GenRsaKey(bits int) {
	// Generates private key.
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		panic(err)
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	priv := pem.EncodeToMemory(block)
	os.WriteFile("ID_RSA", priv, 0)

	// Generates public key from private key.
	publicKey := &privateKey.PublicKey
	derPkix := x509.MarshalPKCS1PublicKey(publicKey)
	block = &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: derPkix,
	}
	publ := pem.EncodeToMemory(block)
	os.WriteFile("ID_RSA.pub", publ, 0)
}

var rsagenCmd = &cobra.Command{
	Use:   "rsa",
	Short: "generate rsa file for yggdrasilapi",
	Long:  `...`,
	Run: func(cmd *cobra.Command, args []string) {
		GenRsaKey(Bits)
	},
}

func init() {
	rsagenCmd.Flags().IntVarP(&Bits, "length", "l", 4096, "bits length for RSA")
	rootCmd.AddCommand(rsagenCmd)
}
