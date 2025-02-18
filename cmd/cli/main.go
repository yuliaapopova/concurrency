package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"concurrency/app/network"
)

func main() {
	address := flag.String("address", "127.0.0.1:9001", "address to listen on")
	flag.Parse()

	reader := bufio.NewReader(os.Stdin)
	client, err := network.NewTcpClient(*address)
	defer client.Close()
	if err != nil {
		log.Fatal(err)
	}
	for {
		request, _ := reader.ReadString('\n')
		response, err := client.Send([]byte(request))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("response: %s \n", string(response))
	}
}
