package main

import (
	"bufio"
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

	fmt.Println("\ntest open readonly for read")
	wg.Add(2)
	go scheduleWrite(&wg)
	go readonlyRead(&wg)
	wg.Wait()

	log.Println("\n--")

	fmt.Println("test open readwrite for read")
	wg.Add(2)
	go scheduleWrite(&wg)
	go readwriteRead(&wg)
	wg.Wait()

	os.Remove(pipeFile)
}

func readonlyRead(wg *sync.WaitGroup) {

	log.Println("r-- opening a namedpipe READONLY")
	readfile, err := os.OpenFile(pipeFile, os.O_RDONLY, os.ModeNamedPipe)
	if err != nil {
		log.Fatal("open named pipe file error:", err)
	}
	defer readfile.Close()
	log.Println("r-- opened a namedpipe READONLY")

	reader := bufio.NewReader(readfile)
	line, err := reader.ReadBytes('\n')
	if err == nil {
		log.Print("r-- read from namedpipe: " + string(line))
	}
	wg.Done()
}

func readwriteRead(wg *sync.WaitGroup) {

	log.Println("r-- opening a namedpipe READWRITE")
	readfile, err := os.OpenFile(pipeFile, os.O_RDWR, os.ModeNamedPipe)
	if err != nil {
		log.Fatal("open named pipe file error:", err)
	}
	defer readfile.Close()
	log.Println("r-- opened a namedpipe READWRITE")

	reader := bufio.NewReader(readfile)
	line, err := reader.ReadBytes('\n')
	if err == nil {
		log.Print("r-- read from namedpipe: " + string(line))
	}
	wg.Done()
}

func scheduleWrite(wg *sync.WaitGroup) {
	log.Println("w-- opening a namedpipe WRITEONLY for write after 3 seconds")
	time.Sleep(3 * time.Second)
	writefile, err := os.OpenFile(pipeFile, os.O_WRONLY, os.ModeNamedPipe)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer writefile.Close()
	log.Println("w-- opened a namedpipe WRITEONLY for write. it will write after 3 seconds")

	time.Sleep(3 * time.Second)
	writefile.WriteString("test write\n")
	log.Println("w-- wrote string to named pipe file.")
	wg.Done()
}

// test open readonly for read
// 2022/05/12 13:53:33 r-- opening a namedpipe READONLY
// 2022/05/12 13:53:33 w-- opening a namedpipe WRITEONLY for write after 3 seconds
// 2022/05/12 13:53:36 w-- opened a namedpipe WRITEONLY for write. it will write after 3 seconds
// 2022/05/12 13:53:36 r-- opened a namedpipe READONLY
// 2022/05/12 13:53:39 w-- wrote string to named pipe file.
// 2022/05/12 13:53:39 r-- read from namedpipe: test write
// 2022/05/12 13:53:39
// --
// test open readwrite for read
// 2022/05/12 13:53:39 r-- opening a namedpipe READWRITE
// 2022/05/12 13:53:39 r-- opened a namedpipe READWRITE
// 2022/05/12 13:53:39 w-- opening a namedpipe WRITEONLY for write after 3 seconds
// 2022/05/12 13:53:42 w-- opened a namedpipe WRITEONLY for write. it will write after 3 seconds
// 2022/05/12 13:53:45 w-- wrote string to named pipe file.
// 2022/05/12 13:53:45 r-- read from namedpipe: test write
