// pidigitgen project main.go
package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
)

func main() {

	if os.Args[1] == "client" {
		if os.Args[2] == "TCP" {
			tcpClient()
		} else {
			udpClient()
		}
		return
	}

	piData, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	go tcpPiGenListen(piData)
	udpPiGenListen(piData)
}

func tcpPiGenListen(readFile *os.File) {
	service := ":31415"
	addr, err := net.ResolveTCPAddr("tcp", service)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	listner, err := net.ListenTCP("tcp", addr)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for {
		conn, err := listner.Accept()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		go tcpPiGen(conn, readFile)
	}
}

func tcpPiGen(conn net.Conn, readFile *os.File) {
	var index int64 = 0
	buf := make([]byte, 1)
	for {
		_, err := readFile.ReadAt(buf, index)
		if err != nil {
			conn.Close()
			return
		}
		conn.Write(buf)
		index++
	}
}

func udpPiGenListen(readFile *os.File) {
	service := ":31415"
	addr, err := net.ResolveUDPAddr("udp", service)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	listner, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	buf := make([]byte, 1024)
	readbuf := make([]byte, 1)
	for {
		num, addr, err := listner.ReadFromUDP(buf)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		s := string(buf[0:num])
		ind, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		ind = ind - 1
		if ind < 0 {
			ind = 0
		}

		_, err = readFile.ReadAt(readbuf, ind)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		s = s + ":" + string(readbuf)
		listner.WriteToUDP([]byte(s), addr)
	}
}

func tcpClient() {
	service := os.Args[3] + ":31415"
	addr, err := net.ResolveTCPAddr("tcp", service)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for {
		buf := make([]byte, 1024)
		num, err := conn.Read(buf)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Print(string(buf[0:num]))
	}
}

func udpClient() {
	service := os.Args[3] + ":31415"
	conn, err := net.Dial("udp", service)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	conn.Write([]byte(os.Args[4]))
	buf := make([]byte, 1024)
	num, err := conn.Read(buf)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	} else {
		fmt.Println(string(buf[0:num]))
	}

}
