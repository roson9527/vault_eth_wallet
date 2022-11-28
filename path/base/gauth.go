package base

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"github.com/roson9527/vault_eth_wallet/modules"
	"strings"
	"time"
)

var GAuth = &gAuth{}

type gAuth struct {
	// contains filtered or unexported fields
}

func (g *gAuth) Generate(secret string) (*modules.Authenticator, error) {
	inputNoSpaces := strings.Replace(secret, " ", "", -1)
	inputNoSpacesUpper := strings.ToUpper(inputNoSpaces)
	key, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(inputNoSpacesUpper)
	if err != nil {
		return nil, err
	}

	epochSeconds := time.Now().Unix()
	pwd := g.oneTimePassword(key, g.toBytes(epochSeconds/30))
	secondsRemaining := 30 - (epochSeconds % 30)
	return &modules.Authenticator{
		Secret: inputNoSpacesUpper,
		Expire: int(secondsRemaining),
		Code:   pwd,
	}, nil
}

func (g *gAuth) toBytes(value int64) []byte {
	var result []byte
	mask := int64(0xFF)
	shifts := [8]uint16{56, 48, 40, 32, 24, 16, 8, 0}
	for _, shift := range shifts {
		result = append(result, byte((value>>shift)&mask))
	}
	return result
}

func (g *gAuth) toUint32(bytes []byte) uint32 {
	return (uint32(bytes[0]) << 24) + (uint32(bytes[1]) << 16) +
		(uint32(bytes[2]) << 8) + uint32(bytes[3])
}

func (g *gAuth) oneTimePassword(key []byte, value []byte) uint32 {
	// sign the value using HMAC-SHA1
	hmacSha1 := hmac.New(sha1.New, key)
	hmacSha1.Write(value)
	hash := hmacSha1.Sum(nil)

	// We're going to use a subset of the generated hash.
	// Using the last nibble (half-byte) to choose the index to start from.
	// This number is always appropriate as it's maximum decimal 15, the hash will
	// have the maximum index 19 (20 bytes of SHA1) and we need 4 bytes.
	offset := hash[len(hash)-1] & 0x0F

	// get a 32-bit (4-byte) chunk from the hash starting at offset
	hashParts := hash[offset : offset+4]

	// ignore the most significant bit as per RFC 4226
	hashParts[0] = hashParts[0] & 0x7F

	number := g.toUint32(hashParts)

	// size to 6 digits
	// one million is the first number with 7 digits so the remainder
	// of the division will always return < 7 digits
	pwd := number % 1000000

	return pwd
}
