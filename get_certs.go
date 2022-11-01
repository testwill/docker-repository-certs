package main

import (
	"bytes"
	"crypto/tls"
	"encoding/pem"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
//	"os/exec"

	//"log"
)

const (
	certsPath = "/etc/docker/certs.d/"
)


func main() {
	var (
		addr string
	)
	flag.StringVar(&addr, "addr", "", "仓库地址")
	flag.Parse()
	os.Mkdir(certsPath, 0666)
	paths := certsPath + addr
	if checkFileIsExist(paths) {
		fmt.Println(paths, " exists")
		return
	} else {
		err := os.Mkdir(paths, 0666)
		if err != nil {
			fmt.Println("mkdir ", paths, " : ",err)
		}
		certs2, _:= GetCertificatesPEM(addr)
		fmt.Printf(certs2)
		filename := paths + "/" + addr + ".crt"
		err1 := ioutil.WriteFile(filename, []byte(certs2), 0666)
		if err1 != nil {
			fmt.Println(err1)
		}
	}
}

func GetCertificatesPEM(address string) (string, error) {
	conn, err := tls.Dial("tcp", address, &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		return "", err
	}
	defer conn.Close()
	var b bytes.Buffer
	for _, cert := range conn.ConnectionState().PeerCertificates {
		err := pem.Encode(&b, &pem.Block{
			Type: "CERTIFICATE",
			Bytes: cert.Raw,
		})
		if err != nil {
			return "", err
		}
	}
	return b.String(), nil
}


func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

