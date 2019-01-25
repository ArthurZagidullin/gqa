package lib

import (
	"fmt"
	"os"
	"sync"
)

type Message struct {
	Id    string
	Value int
}

type SimpleStorage struct {
}

func (s *SimpleStorage) Save(id string, avg int) error {
	fmt.Printf("SimpleStorage: ID: %s AVG: %d \n", id, avg)
	return nil
}

type FileStorage struct {
	f  *os.File
	mu *sync.Mutex
}

func NewFileStorage(filename string) *FileStorage {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		panic(err)
	}
	mu := &sync.Mutex{}
	return &FileStorage{f: f, mu: mu}
}

func (s *FileStorage) Save(id string, avg int) error {
	s.mu.Lock()
	_, err := s.f.Write([]byte(fmt.Sprintf("FileStorage: ID: %s AVG: %d \n", id, avg)))
	s.mu.Unlock()
	fmt.Printf("FileStorage: ID: %s AVG: %d \n", id, avg)
	return err
}
