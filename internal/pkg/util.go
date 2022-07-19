package pkg

import (
	"bytes"
	"io"
	"log"
	"math"
	"os"
	"sync"
)

const (
	Concurrecy = 10
)

func Read(file *os.File, d chan []byte, size, ofset int64) {
	data := make([]byte, size)
	_, err := file.ReadAt(data, ofset)
	if err == io.EOF {
		log.Fatalf("failed creating file: %s", err)
	}
	d <- data
}
func ConcurrentFileRead(filename string) []byte {
	var chans [Concurrecy]chan []byte
	for i := range chans {
		chans[i] = make(chan []byte)
	}
	info, err := os.Stat(filename)
	if err != nil {
		log.Fatalf("cant open")
	}
	file, _ := os.Open(filename)
	data := make([]byte, 0)
	size := int64(math.Ceil(float64(info.Size()) / float64(Concurrecy)))
	for i := 0; i < Concurrecy-1; i++ {
		go Read(file, chans[i], size, size*int64(i))
	}
	go Read(file, chans[Concurrecy-1], int64(info.Size()-size*(Concurrecy-1)), size*(Concurrecy-1))
	for i := 0; i < Concurrecy; i++ {
		d := <-chans[i]
		data = append(data, d...)
	}
	log.Printf("%s", data)
	return data
}

func Write(filename string, b []byte, ofset int64, wg *sync.WaitGroup) {
	defer wg.Done()
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	defer file.Close()

	b = bytes.TrimRight(b, "\x00")

	file.WriteAt(b, ofset)
}

func ConcurrentFileWrite(filename string, data []byte) {
	var wg sync.WaitGroup
	size := int64(math.Ceil(float64(len(data)) / float64(Concurrecy)))
	for i := 0; i < Concurrecy; i++ {
		wg.Add(1)
		go Write(filename, data[int64(i)*size:int64(i+1)*size], int64(i)*size, &wg)
	}
	wg.Wait()
}
