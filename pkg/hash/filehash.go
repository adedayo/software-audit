package hash

import (
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

func FileSha(path string) (sha string, err error) {
	f, err := os.Open(path)
	if err != nil {
		return sha, err
	}
	defer f.Close()

	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
		return sha, err
	}
	hashSum := fmt.Sprintf("%x", h.Sum(nil))
	return hashSum, nil

}

func FileSha256(path string) (sha2 string, err error) {
	f, err := os.Open(path)
	if err != nil {
		return sha2, err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return sha2, err
	}
	hashSum := fmt.Sprintf("%x", h.Sum(nil))
	return hashSum, nil

}
