package model

import "encoding/json"

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
	Files map[string]File
}

type FileDB interface {
	Get(key string) (File, bool)
	Set(key string, value File)
}

func (s FileRepo) Get(key string) (File, bool) {
	value, ok := s.Files[key]
	return value, ok
}

func (s *FileRepo) Set(key string, value File) {
	s.Files[key] = value
}
