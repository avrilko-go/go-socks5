package main

import (
	"fmt"
	"net"
	"socks5/util"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:8102")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go doConn(conn)
	}
}

func doConn(conn net.Conn) {
	defer conn.Close()
	data := make([]byte, 512)
	dataLen, err := conn.Read(data)
	data = data[:dataLen]
	util.Decode(data)
	if err != nil {
		fmt.Println(err)
		return
	}
	disString := string(data[:])
	fmt.Println(disString)
	l, err := net.Dial("tcp", disString)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()
	_, err = conn.Write([]byte{0x01})
	if err != nil {
		fmt.Println(err)
		return
	}

	go ioCopySend(conn, l)
	ioCopyRe(l, conn)
}

func ioCopySend(local, dis net.Conn) (int, error) {
	var data = make([]byte, 1024)
	writeLen := 0
	for {
		dataLen, err := local.Read(data)
		if err != nil {
			return writeLen, err
		}

		if dataLen > 0 {
			data = data[:dataLen]
			util.Decode(data)
			wLen, err := dis.Write(data)
			if err != nil {
				return writeLen, err
			}
			writeLen = writeLen + wLen
		} else {
			return writeLen, nil
		}
	}
}

func ioCopyRe(dis, local net.Conn) (int, error) {
	var data = make([]byte, 1024)
	writeLen := 0
	for {
		dataLen, err := dis.Read(data)
		if err != nil {
			return writeLen, err
		}

		if dataLen > 0 {
			data = data[:dataLen]
			util.Encode(data)
			wLen, err := local.Write(data)
			if err != nil {
				return writeLen, err
			}
			writeLen = writeLen + wLen
		} else {
			return writeLen, nil
		}
	}
}
