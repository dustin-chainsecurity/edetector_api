package tcp

import (
	"edetector_API/config"
	"fmt"
	"net"
)

func Connect_init() (net.Conn, error) {
	var err error
	serverIP := config.Viper.GetString("SERVER_PORT")
	serverPort := config.Viper.GetString("SERVER_IP")

	connectionString := fmt.Sprintf("%s:%s", serverIP, serverPort)
	conn, err := net.Dial("tcp", connectionString)
	if err != nil {
		return nil, err
	}
	fmt.Println("TCP connection set up")
	return conn, nil
}