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

type Data struct {
	Data   []byte
	Offset int64
}

func Read(file *os.File, d chan Data, size, ofset int64, Readed []byte) {
	data := make([]byte, size)

	_, err := file.ReadAt(data, ofset)
	if err == io.EOF {
		log.Fatalf("failed creating file: %s", err)
	}

	d <- Data{
		Data:   data,
		Offset: ofset,
	}
}
func ConcurrentFileRead(filename string) []byte {
	var chans chan Data
	chans = make(chan Data)

	info, err := os.Stat(filename)
	if err != nil {
		log.Fatalf("cant open")
	}

	file, _ := os.Open(filename)
	data := make([]byte, info.Size())
	size := int64(math.Ceil(float64(info.Size()) / float64(Concurrecy)))

	for i := 0; i < Concurrecy-1; i++ {
		go Read(file, chans, size, size*int64(i), data)
	}

	go Read(file, chans, int64(info.Size()-size*(Concurrecy-1)), size*(Concurrecy-1), data)

	for i := 0; i < Concurrecy; i++ {
		temp := <-chans
		copy(data[temp.Offset:], temp.Data)
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

	_, err = file.WriteAt(b, ofset)
	if err != nil {
		log.Fatalf("failed writing file: %s", err)
	}
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
