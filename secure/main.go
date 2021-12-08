package secure

import (
	"io"
	"os"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/armor"
	"github.com/ProtonMail/go-crypto/openpgp/packet"
)

// The ISecure main interface
type ISecure interface {
	Encrypt(string, *os.File, string) (*os.File, error)
}

// The Secure structure
type secure struct{}

func (c *secure) readPubKey(keyPath string) (*openpgp.Entity, error) {
	file, err := os.Open(keyPath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	block, err := armor.Decode(file)
	if err != nil {
		return nil, err
	}

	return openpgp.ReadEntity(packet.NewReader(block.Body))
}

// Encrypt a file
func (c *secure) Encrypt(targetPath string, file *os.File, pubkey string) (*os.File, error) {
	pubKey, err := c.readPubKey(pubkey)
	if err != nil {
		return nil, err
	}

	gpgFile, err := os.Create(targetPath)
	if err != nil {
		return nil, err
	}

	defer gpgFile.Close()

	gpgWriter, err := openpgp.Encrypt(
		gpgFile,
		[]*openpgp.Entity{pubKey},
		nil,
		&openpgp.FileHints{IsBinary: true},
		nil,
	)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(gpgWriter, file)
	if err != nil {
		return nil, err
	}

	defer gpgWriter.Close()

	return os.Open(targetPath)
}

// New instance of secure
func New() ISecure {
	return &secure{}
}
