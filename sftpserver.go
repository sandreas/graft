// An example SFTP server implementation using the golang SSH package.
// Serves the whole filesystem visible to the user, and has a hard-coded username and password,
// so not for real use!
package main

import (
	"fmt"
	//"io"
	"io/ioutil"
	"log"
	//"net"
	"os"

	//"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"os/user"
	"crypto/rsa"
	"crypto/rand"
	"encoding/pem"
	"crypto/x509"
	"github.com/pkg/sftp"
	"github.com/sandreas/graft/sftphandler"
	"io"
	"net"
)

var debugStream os.File;

// Based on example server code from golang.org/x/crypto/ssh and server_standalone
func main() {
	// debugStream := ioutil.Discard

	//debugStream := os.Stderr

	// An SSH server is represented by a ServerConfig, which holds
	// certificate details and handles authentication of ServerConns.
	config := &ssh.ServerConfig{
		PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			// Should use constant-time compare (or better, salt+hash) in
			// a production setting.
			println( "Login: %s\n", c.User())
			if c.User() == "test" && string(pass) == "test" {
				return nil, nil
			}
			return nil, fmt.Errorf("password rejected for %q", c.User())
		},
	}


	graftHomePath := createGraftHomePath();
	generateKeysIfNotExist(graftHomePath)

	// graftHomePath = "/Users/andreas/.ssh"
	privateBytes, err := ioutil.ReadFile(graftHomePath + "/id_rsa")
	if err != nil {
		log.Fatal("Failed to load private key", err)
	}
	println(privateBytes)
	//
	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		log.Fatal("Failed to parse private key", err)
	}
	//
	config.AddHostKey(private)
	println("Server key generation worked")


	// Once a ServerConfig has been configured, connections can be
	// accepted.
	listener, err := net.Listen("tcp", "0.0.0.0:2022")
	if err != nil {
		log.Fatal("failed to listen for connection", err)
	}
	fmt.Printf("Listening on %v\n", listener.Addr())


	for {
		conn, e := listener.Accept()
		if e != nil {
			os.Exit(2)
		}
		go HandleConn(conn, config)
	}

	/*
	for {
		nConn, err := listener.Accept()
		if err != nil {
			log.Fatal("failed to accept incoming connection", err)
		}

		// Before use, a handshake must be performed on the incoming net.Conn.
		sconn, chans, reqs, err := ssh.NewServerConn(nConn, config)
		if err != nil {
			log.Fatal("failed to handshake", err)
		}
		log.Println("login detected:", sconn.User())
		println( "SSH server established\n")

		// The incoming Request channel must be serviced.
		go ssh.DiscardRequests(reqs)



		// Service the incoming Channel channel.
		for newChannel := range chans {
			// Channels have a type, depending on the application level
			// protocol intended. In the case of an SFTP session, this is "subsystem"
			// with a payload string of "<length=4>sftp"
			println( "Incoming channel: %s\n", newChannel.ChannelType())
			if newChannel.ChannelType() != "session" {
				newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
				println( "Unknown channel type: %s\n", newChannel.ChannelType())
				continue
			}
			channel, requests, err := newChannel.Accept()
			if err != nil {
				log.Fatal("could not accept channel.", err)
			}
			println( "Channel accepted\n")

			// Sessions have out-of-band requests such as "shell",
			// "pty-req" and "env".  Here we handle only the
			// "subsystem" request.
			go func(in <-chan *ssh.Request) {
				for req := range in {
					println( "Request: %v\n", req.Type)
					ok := false
					switch req.Type {
					case "subsystem":
						println( "Subsystem: %s\n", req.Payload[4:])
						if string(req.Payload[4:]) == "sftp" {
							ok = true
						}
					}
					println( " - accepted: %v\n", ok)
					req.Reply(ok, nil)
				}
			}(requests)

			var matchingPaths []string
			matchingPaths = append(matchingPaths, "graft.go")
			matchingPaths = append(matchingPaths, "LICENSE")
			matchingPaths = append(matchingPaths, "README.md")

			root := sftphandler.CustomHandler(matchingPaths)
			server := sftp.NewRequestServer(channel, root)
			if err := server.Serve(); err == io.EOF {
				server.Close()
				log.Print("sftp client exited session.")
			} else if err != nil {
				log.Fatal("sftp server completed with error:", err)
			}
		}
	}
	*/
}

func HandleConn(conn net.Conn, config *ssh.ServerConfig) {
	defer conn.Close()
	e := handleConn(conn, config)
	if e != nil {
		log.Println("sftpd connection errored:", e)
	}
}
func handleConn(conn net.Conn, config *ssh.ServerConfig) error {
	sconn, chans, reqs, e := ssh.NewServerConn(conn, config)
	if e != nil {
		return e
	}
	defer sconn.Close()

	// The incoming Request channel must be serviced.
	println( "login detected:", sconn.User())
	println( "SSH server established\n")

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
				println( "Request: %v\n", req.Type)
				ok := false
				switch req.Type {
				case "subsystem":
					println( "Subsystem: %s\n", req.Payload[4:])
					if string(req.Payload[4:]) == "sftp" {
						ok = true
					}
				}
				println( " - accepted: %v\n", ok)
				req.Reply(ok, nil)
			}
		}(requests)

		var matchingPaths []string
		matchingPaths = append(matchingPaths, "graft.go")
		matchingPaths = append(matchingPaths, "LICENSE")
		matchingPaths = append(matchingPaths, "README.md")
		matchingPaths = append(matchingPaths, "data/fixtures/global/file.txt")
		matchingPaths = append(matchingPaths, "data/fixtures")
		matchingPaths = append(matchingPaths, "data")

		root := sftphandler.TestHandler(matchingPaths)
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

//func TestingHandler() {
//	return sftp.Handlers{
//		//root,
//		//root,
//		//root,
//		//root
//
//		//FileGet  FileReader
//		//FilePut  FileWriter
//		//FileCmd  FileCmder
//		//FileInfo FileInfoer
//	}
//}

func PrintDiscardRequests(in <-chan *ssh.Request) {
	for req := range in {
		log.Println("Discarding ssh request", req.Type, *req)
		if req.WantReply {
			req.Reply(false, nil)
		}
	}
}

func createGraftHomePath() string {
	usr, err := user.Current()

	if err != nil {
		println("Could not determine current user ", err)
		os.Exit(1)
	}
	graftHomePath := usr.HomeDir + "/.graft";
	mode := int(0755)
	if _, err := os.Stat(graftHomePath); err != nil {
		err := os.Mkdir(graftHomePath, os.FileMode(mode))
		if err != nil {
			println("Could not create home directory " + graftHomePath)
			os.Exit(1)
		}
	}
	return graftHomePath

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
//
//func checkError(err error) {
//	if err != nil {
//		fmt.Println("Fatal error ", err.Error())
//		os.Exit(1)
//	}
//}