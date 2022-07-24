package pkg

import (
	"crypto/sha256"
	"log"
	"math/big"

	uuid "github.com/nu7hatch/gouuid"
)

func sha256Of(input string) []byte {
	algorithm := sha256.New()
	algorithm.Write([]byte(input))

	return algorithm.Sum(nil)
}

func GenerateFileId(filename string) uint64 {
	u, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	urlHashBytes := sha256Of(filename + u.String())
	id := new(big.Int).SetBytes(urlHashBytes).Uint64()

	log.Println(id)

	return id
}
