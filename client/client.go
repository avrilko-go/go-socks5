package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"socks5/util"
	"strconv"
)

func main() {
	addr, _ := net.ResolveTCPAddr("tcp4", "0.0.0.0:8081")
	l, err := net.ListenTCP("tcp4", addr)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := l.AcceptTCP()
		if err != nil {
			continue
		}
		go doConn(conn)
	}
}

func doConn(conn *net.TCPConn) {
	defer conn.Close()
	var data = make([]byte, 512)
	_, err := conn.Read(data)
	if err != nil {
		fmt.Println(err)
		return
	}
	if data[0] != 0x05 { // 不是socks5协议
		fmt.Println(data)
		fmt.Println("不是socks5协议")
		return
	}
	_, err = conn.Write([]byte{0x05, 0x00})
	if err != nil {
		fmt.Println(err)
		return
	}
	lenHost, err := conn.Read(data)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 解析端口和host
	if data[0] != 0x05 {
		fmt.Println("不是socks5协议")
		return
	}
	netType := data[3]
	dIp := ""
	if netType == 0x01 { // ipv4

		dIp = net.IPv4(data[4], data[5], data[6], data[7]).String()
	} else if netType == 0x04 { // ipv6
		dIp = net.IP{data[4], data[5], data[6], data[7], data[8], data[9], data[10], data[11], data[12], data[13], data[14], data[15], data[16], data[17], data[18], data[19]}.String()
	} else if netType == 0x03 { // host
		hostLen := data[4]
		host := data[5 : 5+hostLen]
		dIp = string(host)
	}

	port := int(binary.BigEndian.Uint16(data[lenHost-2:]))

	addrByte := []byte(net.JoinHostPort(dIp, strconv.Itoa(port)))
	fmt.Println(string(addrByte))
	l, err := net.Dial("tcp", "136.244.114.167:8102")
	//l, err := net.Dial("tcp", "111.231.252.16:8102")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()
	_, err = encodeSend(l, addrByte)

	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = l.Read(data)
	if err != nil {
		fmt.Println(err)
		return
	}
	if data[0] == 0x01 {
		conn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
		go ioCopySend(conn, l)
		ioCopyRe(l, conn)
	}
}

func encodeSend(conn net.Conn, data []byte) (int, error) {
	util.Encode(data)
	return conn.Write(data)
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
			util.Encode(data)
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
			util.Decode(data)
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
