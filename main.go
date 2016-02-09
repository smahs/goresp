package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
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
			fmt.Printf("%v <-> %v\n", client.LocalAddr(),
				client.RemoteAddr())
			ch <- client
		}
	}()
	return ch
}

func handleConn(client net.Conn) {
	defer client.Close()
	for {
		request, err := handleRequest(client)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		fmt.Println(request)
		reply := genResponse(request)
		reply += DELIMITER
		client.Write([]byte(reply))
		/*if err != nil { // int casting failed, ignore message
			input, err := handleRequest(reader, count)
			if err != nil { // EOF probably
				fmt.Println(err.Error())
				break
			}
			output := new(bytes.Buffer)
			for _, str := range input {
				output.WriteString(str + "\n")
			}
			fmt.Println(output.String())
			client.Write(output.Bytes())
		}*/
		/*line = bytes.TrimRight(line, DELIMITER)
		fmt.Println("IN: " + string(line))
		resp := append(line, DELIMITER...)
		fmt.Println("OUT: " + string(resp))
		client.Write(resp)*/
	}
}

func handleRequest(client net.Conn) ([]string, error) {
	fmt.Println("hr called")
	var args []string
	var i int
	reader := bufio.NewReader(client)
	request, err := reader.ReadString('\n')
	if err != nil { // EOF probably
		fmt.Println(err.Error())
		return nil, err
	}
	argsCount, err := strconv.Atoi(string(request[1]))
	fmt.Println(argsCount)
	if err != nil {
		return nil, errors.New("invalid command")
	}
	for i = 0; i < argsCount; i++ {
		str, err := readString(reader)
		if err != nil {
			return nil, err
		}
		args = append(args, str)
	}
	if err != nil {
		return nil, err
	}
	return args, nil
}

func readString(reader *bufio.Reader) (string, error) {
	length, err := reader.ReadString('\n')
	size, err := strconv.Atoi(string(length[1]))
	str, err := reader.ReadString('\n')
	str = strings.TrimSpace(str)
	if len(str) != size {
		err = errors.New("Arg length mismatched")
	}
	if err == nil {
		fmt.Println(str)
		return str, nil
	} else {
		return "", err
	}
}

func genResponse(input []string) string {
	output := "+ "
	for i, str := range input {
		output += str
		if i != len(input) {
			output += " "
		}
	}
	return output
}
