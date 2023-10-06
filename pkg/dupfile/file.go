package dupfile

import (
	"crypto/sha256"
	"io"
	"os"
)

type File struct {
	Path string `json:"path"`
	Name string `json:"name"`
	Size int64  `json:"size"`
	Hash []byte `json:"hash"`
}

func (f File) AbsPath() string {
	return f.Path
}

func (fl *File) CryptoHash(cache *FileCache) ([]byte, error) {
	if len(fl.Hash) == 0 {
		// Check if MD5 is in cache
		if cache != nil {
			hash := cache.Get(fl.Path, fl.Size)
			if hash != nil {
				fl.Hash = hash
				return fl.Hash, nil
			}
		}

		// Calculate and store cryptographic hash
		f, err := os.Open(fl.AbsPath())
		if err != nil {
			return []byte{}, err
		}
		defer f.Close()

		h := sha256.New()
		if _, err := io.Copy(h, f); err != nil {
			return []byte{}, err
		}
		fl.Hash = h.Sum(nil)
		//Add to cache
		if cache != nil {
			cache.Put(fl)
		}
	}

	return fl.Hash, nil
}
