package main

import (
	"./silk"
	"flag"
	"log"
	"net"
)

var listen_port int
var sender bool

func init() {
	flag.IntVar(&listen_port, "port", 17159, "broadcast listen port")
	flag.BoolVar(&sender, "sender", false, "debug")
}

func getLocalIp() ([]string, error) {
	ret := make([]string, 0)
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ret, err
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if r := ipnet.IP.To4(); r == nil {
				continue
			}
			ret = append(ret, ipnet.IP.String())
		}
	}
	return ret, nil
}

func recvBroadcast(port int) error {
	socket, err := net.ListenUDP("udp4", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: port,
	})
	for {
		data := make([]byte, 4096)
		read, _, err := socket.ReadFromUDP(data)
		if err != nil {
			return err
		}
		s := string(data[:read])
		log.Println(s)
	}
	return err
}

func broadcast(port int) error {
	broadcast_addr := net.IPv4(255, 255, 255, 255)
	socket, err := net.DialUDP("udp4", nil, &net.UDPAddr{
		IP:   broadcast_addr,
		Port: port,
	})
	if err != nil {
		log.Println(err)
		return err
	}
	ips, err := getLocalIp()
	if err != nil {
		log.Println(err)
		return err
	}
	for _, v := range ips {
		log.Println(v)
	}
	socket.Write([]byte("hello world!"))
	return nil
}

func main() {
	flag.Parse()
	if sender {
		broadcast(listen_port)
	} else {
		go recvBroadcast(listen_port)
		silk.ListenAndServe()
	}
}
