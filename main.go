package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)

const PORT = "6379"
const DELIMITER = "\r\n"

func main() {
	server, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		panic("Couldn't start listening: " + err.Error())
	} else {
		fmt.Println("Listening on port: " + PORT)
	}

	conns := clientConns(server)
	for {
		go handleConn(<-conns)
	}
}

func clientConns(listener net.Listener) chan net.Conn {
	ch := make(chan net.Conn)
	go func() {
		for {
			client, err := listener.Accept()
			if err != nil {
				fmt.Printf("Couldn't accept: " + err.Error())
				continue
			}
			fmt.Printf("%v <-> %v\n", client.LocalAddr(), client.RemoteAddr())
			ch <- client
		}
	}()
	return ch
}

func handleConn(client net.Conn) {
	defer client.Close()
	b := bufio.NewReader(client)
	for {
		line, err := b.ReadBytes('\n')
		if err != nil { // EOF, or worse
			break
		}
		fmt.Println(string(line))
		count := binary.LittleEndian.Uint16(line[1:])
		fmt.Println(count)
		if err != nil { // int casting failed, ignore message
			input, err := handleRequest(b, count)
			if err != nil { // EOF probably
				fmt.Println(err.Error())
				break
			}
			output := bytes.Join(input, []byte(" "))
			fmt.Println(output)
			client.Write(output)
		}
		/*line = bytes.TrimRight(line, DELIMITER)
		fmt.Println("IN: " + string(line))
		resp := append(line, DELIMITER...)
		fmt.Println("OUT: " + string(resp))
		client.Write(resp)*/
	}
}

func handleRequest(reader *bufio.Reader, count uint16) ([][]byte, error) {
	fmt.Println("hr called")
	var arr [][]byte
	var i uint16
	for i = 0; i < count*2; i++ {
		if i%2 != 0 {
			raw, err := reader.ReadBytes('\n')
			if err == nil {
				fmt.Println(string(raw))
				arr = append(arr, raw)
			} else {
				return nil, err
			}
		}
	}
	return arr, nil
}

func sendResponse(client net.Conn) {
}
