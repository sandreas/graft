package sftpd

import (
	"golang.org/x/crypto/ssh"
	"log"
	"net"
	"fmt"
	"os"
	"crypto/subtle"
	"io/ioutil"
	"crypto/rsa"
	"encoding/pem"
	"crypto/x509"
	"io"
	"crypto/rand"
	"github.com/pkg/sftp"
	"strconv"
)

func NewGraftServer(graftHomePath, listenAddress string, listenPort int, username, password string, pathMapper *PathMapper) {
	 //An SSH server is represented by a ServerConfig, which holds
	 //certificate details and handles authentication of ServerConns.



	config := &ssh.ServerConfig{
		PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			// Should use constant-time compare (or better, salt+hash) in
			// a production setting.

			log.Printf("Login: %s\n", c.User())

			if subtle.ConstantTimeCompare([]byte(username), []byte(c.User())) == 1 && subtle.ConstantTimeCompare(pass, []byte(password)) == 1 {
				return nil, nil
			}
			return nil, fmt.Errorf("password rejected for %q", c.User())
		},
	}

	generateKeysIfNotExist(graftHomePath)

	// graftHomePath = "/Users/andreas/.ssh"
	privateBytes, err := ioutil.ReadFile(graftHomePath + "/id_rsa")
	if err != nil {
		log.Fatal("Failed to load private key", err)
	}
	// println(privateBytes)
	//
	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		log.Fatal("Failed to parse private key", err)
	}
	//
	config.AddHostKey(private)
	// println("Server key generation worked")

	// Once a ServerConfig has been configured, connections can be
	// accepted.
	listener, err := net.Listen("tcp", listenAddress+":"+strconv.Itoa(listenPort))
	// listener, err := net.Listen("tcp", "0.0.0.0:2023")
	if err != nil {
		log.Fatal("failed to listen for connection", err)
	}
	fmt.Printf("Listening on %v\n", listener.Addr())

	for {
		conn, e := listener.Accept()
		if e != nil {
			os.Exit(2)
		}
		go HandleConn(conn, config, pathMapper)
	}

}

func HandleConn(conn net.Conn, config *ssh.ServerConfig,  pathMapper *PathMapper) {
	defer conn.Close()
	e := handleConn(conn, config, pathMapper)
	if e != nil {
		log.Println("sftpd connection errored:", e)
	}
}
func handleConn(conn net.Conn, config *ssh.ServerConfig,  pathMapper *PathMapper) error {
	sconn, chans, reqs, e := ssh.NewServerConn(conn, config)
	if e != nil {
		return e
	}
	defer sconn.Close()

	// The incoming Request channel must be serviced.
	log.Println( "login detected:", sconn.User())

	// The incoming Request channel must be serviced.
	go ssh.DiscardRequests(reqs)

	// Service the incoming Channel channel.
	for newChannel := range chans {
		if newChannel.ChannelType() != "session" {
			newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
			continue
		}
		channel, requests, err := newChannel.Accept()
		if err != nil {
			return err
		}

		go func(in <-chan *ssh.Request) {
			for req := range in {
				log.Printf( "Request: %v\n", req.Type)
				ok := false
				switch req.Type {
				case "subsystem":
					log.Printf( "Subsystem: %s\n", req.Payload[4:])
					if string(req.Payload[4:]) == "sftp" {
						ok = true
					}
				}
				log.Printf( " - accepted: %v\n", ok)
				req.Reply(ok, nil)
			}
		}(requests)


		root := VfsHandler(pathMapper)
		server := sftp.NewRequestServer(channel, root)
		if err := server.Serve(); err == io.EOF {
			server.Close()
			log.Print("sftp client exited session.")
		} else if err != nil {
			log.Fatal("sftp server completed with error:", err)
		}

	}
	return nil
}




func generateKeysIfNotExist(homeDir string) {

	privateKeyFile := homeDir + "/id_rsa"
	publicKeyFile := homeDir + "/id_rsa.pub"

	if _, err := os.Stat(privateKeyFile); os.IsNotExist(err) {
		makeSSHKeyPair(publicKeyFile, privateKeyFile)
	}
}

func makeSSHKeyPair(pubKeyPath, privateKeyPath string) error {

	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return err
	}

	// generate and write private key as PEM
	privateKeyFile, err := os.Create(privateKeyPath)
	defer privateKeyFile.Close()
	if err != nil {
		return err
	}
	privateKeyPEM := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}
	if err := pem.Encode(privateKeyFile, privateKeyPEM); err != nil {
		return err
	}

	// generate and write public key
	pub, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(pubKeyPath, ssh.MarshalAuthorizedKey(pub), 0655)
}