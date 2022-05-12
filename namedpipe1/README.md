# namedpipe blocking test

Namedpipe는 프로세스간 통신(IPC) 방법 중 하나이다. 
하나의 프로세스가 namedpipe 파일에 쓰면, 다른 프로세스에서 읽어간다. 

gist: https://gist.github.com/a00a3a81fb155928a65bb9dd56f99833

## blocking test - readonly

코드로 확인하려는 것은 읽는 프로세스가 블로킹 되는 시점이다. 

첫 번째 코드를 보자 

```go
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
```

생성되어 있는 namedpipe 파일을 `readonly`로 열어서 읽는다. 
그런데 중요한 것은 이때 `os.OpenFile` 함수에서 블로킹이 일어난다는 것이다. 

정확히는 namedpipe 파일을 쓸 수 있게 여는 시점에서 블로킹이 해제된다. 
즉, namedpipe 파일을 `readonly`로 열려면, 이 파일을 어디에선가 `writable` 하게 
열어야 한다는 것이다. 

```bash
test open readonly for read
2022/05/12 13:53:33 r-- opening a namedpipe READONLY
2022/05/12 13:53:33 w-- opening a namedpipe WRITEONLY for write after 3 seconds
2022/05/12 13:53:36 w-- opened a namedpipe WRITEONLY for write. it will write after 3 seconds
2022/05/12 13:53:36 r-- opened a namedpipe READONLY
2022/05/12 13:53:39 w-- wrote string to named pipe file.
2022/05/12 13:53:39 r-- read from namedpipe: test write
2022/05/12 13:53:39 
```

위 터미널 출력을 보면 쓰는 쪽의 `os.OpenFile` 함수가 실행 완료한 다음에, 
읽는 쪽의 `os.OpenFile` 함수가 실행 완료하는 것을 볼 수 있다. 


## blocking test - readwrite

이번에는 읽는 쪽에서 `readwrite`로 읽어보겠다. 

```go
func readwriteRead(wg *sync.WaitGroup) {

	// 생략
	readfile, err := os.OpenFile(pipeFile, os.O_RDWR, os.ModeNamedPipe)
	// 생략
}
```

앞서의 코드와 다른 점은 오직 `os.OpenFile` 시에 `os.O_RDWR` 옵션으로 연 것 뿐이다. 
하지만, 아래 터미널 출력에서 보듯이, 이 함수는 쓰는 쪽에서 파일을 여는 것과 상관없이, 
블로킹없이, 바로 실행 완료 되는 것을 알 수 있다. 

```bash
test open readwrite for read
2022/05/12 13:53:39 r-- opening a namedpipe READWRITE
2022/05/12 13:53:39 r-- opened a namedpipe READWRITE
2022/05/12 13:53:39 w-- opening a namedpipe WRITEONLY for write after 3 seconds
2022/05/12 13:53:42 w-- opened a namedpipe WRITEONLY for write. it will write after 3 seconds
2022/05/12 13:53:45 w-- wrote string to named pipe file.
2022/05/12 13:53:45 r-- read from namedpipe: test write
```

## Code Reference
- https://gist.github.com/matishsiao/fc1601a3a3f37c70d91ab3b1ed8485c4

