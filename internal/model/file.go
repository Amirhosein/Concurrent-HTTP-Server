package model

import (
	"encoding/json"

	"anbar.bale.ai/a.iravanimanesh/concurrent-http-server/internal/pkg"
)

type File struct {
	Data []byte
}

func (s File) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

func (s *File) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, s)
}

type FileRepo struct {
	Files map[uint64]File
}

type FileDB interface {
	Get(key string) (File, bool)
	Set(key string, value File)
}

func (s FileRepo) Get(key uint64) (File, bool) {
	value, ok := s.Files[key]
	return value, ok
}

func (s *FileRepo) Set(name string, key uint64, value File) {
	s.Files[key] = value
	pkg.ConcurrentFileWrite(name, value.Data)
}
