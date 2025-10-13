package main

import (
	"fmt"
	"ibrahemassa/http_bootdev/internal/request"
	"io"
	"log"
	"net"
	"strings"
)

const PORT = ":42069"

func getLinesChannel(f io.ReadCloser) <-chan string {
	str_chan := make(chan string)
	go func() {
		cur := ""
		defer f.Close()
		defer close(str_chan)

		for {
			bytes := make([]byte, 8, 8)
			n, err := f.Read(bytes)
			if err != nil {
				if cur != "" {
					str_chan <- cur
				}
				if err != io.EOF {
					break
				}
				return
			}

			str := string(bytes[:n])
			parts := strings.Split(str, "\n")
			for i := 0; i < len(parts)-1; i++ {
				str_chan <- fmt.Sprintf("%s%s", cur, parts[i])
				cur = ""
			}
			cur += parts[len(parts)-1]
		}
	}()
	return str_chan
}

func main() {
	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		// req, err := request.RequestFromReader(conn)
		request.RequestFromReader(conn)
		// lines := getLinesChannel(conn)
		//
		// for line := range lines{
		// 	fmt.Println(line)
		// }

		// fmt.Printf("Request Line:\n- Method: %s\n- target: %s\n- Version: %s", req.RequestLine.Method, req.RequestLine.RequestTarget, req.RequestLine.HttpVersion)
		// fmt.Println("Request line:")
		// fmt.Printf("- Method: %s\n", req.RequestLine.Method)
		// fmt.Printf("- Target: %s\n", req.RequestLine.RequestTarget)
		// fmt.Printf("- Version: %s\n", req.RequestLine.HttpVersion)

	}
}
