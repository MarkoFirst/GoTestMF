package main

// SERVER
import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	listener, _ := net.Listen("tcp", "localhost:3000")
	fmt.Println("Start server")

	defer listener.Close()

	conn, err := listener.Accept()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("new connection")

	listenConnection(conn)

}

func listenConnection(conn net.Conn) {
	for {
		buffer := make([]byte, 1024)

		//go func(conn net.Conn) {
		dataSize, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("connection closed")
			break
		}
		//}(conn)
		data := buffer[:dataSize]

		if string(buffer[:3]) == "get" {
			fmt.Println("__________________________get")
			sendFile(string(buffer[4:dataSize]), conn)
		} else if string(buffer[:6]) == "upload" {
			fmt.Println("__________________________upload")
			arrayText := strings.Split(string(buffer[6:dataSize]), "|")
			os.Create(arrayText[0])
			ioutil.WriteFile(arrayText[0], []byte(arrayText[1]), 0644)
		} else if string(buffer[:6]) == "delete" {
			fmt.Println("__________________________delete")
			err2 := os.Remove(string(buffer[7:dataSize]))
			if err2 != nil {
				fmt.Println(err2)
				break
			}
			conn.Write([]byte("_________OK"))
		} else {
			files, err := ioutil.ReadDir(string(data))
			if err != nil {
				fmt.Println("closed")
				break
			}
			for _, f := range files {
				fmt.Println(f.Name())
				conn.Write([]byte(f.Name()))
			}
			fmt.Println("")
		}
	}

}

func sendFile(fileName string, conn net.Conn) {
	var currentByte int64 = 0
	fileBuffer := make([]byte, 1024)
	file, err := os.Open(strings.TrimSpace(fileName))
	if err != nil {
		log.Fatal(err)
	}
	var err1 error

	for {
		_, err1 = file.ReadAt(fileBuffer, currentByte)
		currentByte += 1024

		fmt.Println(fileBuffer) // для проверки

		conn.Write(fileBuffer)

		if err1 == io.EOF {
			break
		}
	}

	file.Close()
	return

}
