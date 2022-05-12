package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"syscall"
	"time"
)

var pipeFile = "pipe.log"

func main() {
	// make namedpipe
	os.Remove(pipeFile)
	if err := syscall.Mkfifo(pipeFile, 0666); err != nil {
		log.Fatal("Make named pipe file error:", err)
	}

	var wg sync.WaitGroup

	fmt.Println("test open readwrite for read")
	wg.Add(2)
	go scheduleWrite(&wg)
	go readwriteRead(&wg)
	wg.Wait()

	os.Remove(pipeFile)
}

func readwriteRead(wg *sync.WaitGroup) {

	log.Println("r-- opening a namedpipe READWRITE")
	readfile, err := os.OpenFile(pipeFile, os.O_RDWR|syscall.O_NONBLOCK, os.ModeNamedPipe)
	if err != nil {
		log.Fatal("open named pipe file error:", err)
	}
	defer readfile.Close()
	log.Println("r-- opened a namedpipe READWRITE")

	b := make([]byte, 8)
	for {
		n, err := readfile.Read(b)
		if err == nil {
			log.Print("r-- read from namedpipe: " + string(b[:n]))
		}
	}
	wg.Done()
}

func scheduleWrite(wg *sync.WaitGroup) {

	for i := 0; i < 5; i++ {
		writefile, err := os.OpenFile(pipeFile, os.O_WRONLY, os.ModeNamedPipe)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		log.Println("w-- opened a namedpipe WRITEONLY for write. it will be closed after 2 seconds")
		// writefile.Write([]byte("hi!"))
		time.Sleep(2 * time.Second)
		writefile.Close()
	}

	wg.Done()
}
