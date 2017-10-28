package main

// CLIENT
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
	conn, err := net.Dial("tcp", "Localhost:3000")
	if err != nil {
		log.Fatalln(err)
	}
	go func(conn net.Conn) {
		buffer := make([]byte, 1400)
		for { //Приём
			dataSize, err := conn.Read(buffer)
			data := buffer[:dataSize]
			fmt.Println(string(data))
			if err != nil {
				fmt.Println("connection closed")
				break
			}
		}
	}(conn)
	pac := "C:/Games"
	forWrite(pac, err, conn)
}

func forWrite(pac string, err error, conn net.Conn) {

	_, err = conn.Write([]byte(pac)) //Передача
	if err != nil {
		log.Fatalln(err)
	}
	forNext(pac, err, conn)

}

func forNext(pac string, err error, conn net.Conn) {
	text := ""
	fmt.Println("")
	fmt.Scanln(&text)
	fmt.Println("")
	arrayCommands := strings.Split(text, "|")

	switch arrayCommands[0] {
	case "goto":
		pac = pac + arrayCommands[1]
		forWrite(pac, err, conn)
	case "upload":
		upload(arrayCommands[1], arrayCommands[2], conn)
	case "delete":
		conn.Write([]byte(text))
	}

}

func upload(fileName string, directory string, conn net.Conn) {
	var currentByte int64 = 0
	fileBuffer := make([]byte, 1024)

	files, err := ioutil.ReadDir(directory)
	if err != nil {
		fmt.Println("closed")
		return
	}
	for _, f := range files {
		if f.Name() == fileName {
			file, err := os.Open(strings.TrimSpace(fileName))
			if err != nil {
				log.Fatal(err)
			}
			var err1 error

			for {
				_, err1 = file.ReadAt(fileBuffer, currentByte)
				currentByte += 1024
				if err1 == io.EOF {
					break
				}

				forSend := []byte("upload" + fileName + "|" + string(fileBuffer))
				conn.Write(forSend)
			}
			file.Close()
			return
		}
	}

}
