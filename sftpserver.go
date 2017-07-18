// An example SFTP server implementation using the golang SSH package.
// Serves the whole filesystem visible to the user, and has a hard-coded username and password,
// so not for real use!
package main

import (
	"flag"
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
	"encoding/asn1"
	"encoding/gob"
	"github.com/pkg/sftp"
	"github.com/sandreas/graft/sftphandler"
	"io"
	"net"
)

// Based on example server code from golang.org/x/crypto/ssh and server_standalone
func main() {

	var (
		readOnly bool
		debugStderr bool
	)

	flag.BoolVar(&readOnly, "R", false, "read-only server")
	flag.BoolVar(&debugStderr, "e", false, "debug to stderr")
	flag.Parse()

	debugStream := ioutil.Discard
	if debugStderr {
		debugStream = os.Stderr
	}

	// An SSH server is represented by a ServerConfig, which holds
	// certificate details and handles authentication of ServerConns.
	config := &ssh.ServerConfig{
		PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			// Should use constant-time compare (or better, salt+hash) in
			// a production setting.
			fmt.Fprintf(debugStream, "Login: %s\n", c.User())
			if c.User() == "test" && string(pass) == "test" {
				return nil, nil
			}
			return nil, fmt.Errorf("password rejected for %q", c.User())
		},
	}

	usr, err := user.Current()
	graftHomePath := usr.HomeDir + "/.graft";
	mode := int(0755)
	if _, err := os.Stat(graftHomePath); err != nil {
		err := os.Mkdir(graftHomePath, os.FileMode(mode))
		if err != nil {
			println("Could not create home directory " + graftHomePath)
			os.Exit(1)
		}
	}

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
		fmt.Fprintf(debugStream, "SSH server established\n")

		// The incoming Request channel must be serviced.
		go ssh.DiscardRequests(reqs)



		// Service the incoming Channel channel.
		for newChannel := range chans {
			// Channels have a type, depending on the application level
			// protocol intended. In the case of an SFTP session, this is "subsystem"
			// with a payload string of "<length=4>sftp"
			fmt.Fprintf(debugStream, "Incoming channel: %s\n", newChannel.ChannelType())
			if newChannel.ChannelType() != "session" {
				newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
				fmt.Fprintf(debugStream, "Unknown channel type: %s\n", newChannel.ChannelType())
				continue
			}
			channel, requests, err := newChannel.Accept()
			if err != nil {
				log.Fatal("could not accept channel.", err)
			}
			fmt.Fprintf(debugStream, "Channel accepted\n")

			// Sessions have out-of-band requests such as "shell",
			// "pty-req" and "env".  Here we handle only the
			// "subsystem" request.
			go func(in <-chan *ssh.Request) {
				for req := range in {
					fmt.Fprintf(debugStream, "Request: %v\n", req.Type)
					ok := false
					switch req.Type {
					case "subsystem":
						fmt.Fprintf(debugStream, "Subsystem: %s\n", req.Payload[4:])
						if string(req.Payload[4:]) == "sftp" {
							ok = true
						}
					}
					fmt.Fprintf(debugStream, " - accepted: %v\n", ok)
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
}

func generateKeysIfNotExist(homeDir string) {

	privateKeyFile := homeDir + "/id_rsa"
	publicKeyFile := homeDir + "/id_rsa.pub"

	if _, err := os.Stat(privateKeyFile); os.IsNotExist(err) {
		makeSSHKeyPair(publicKeyFile, privateKeyFile)
	}

	//reader := rand.Reader
	//bitSize := 2048
	//key, err := rsa.GenerateKey(reader, bitSize)
	//checkError(err)
	//
	//publicKey := key.PublicKey
	//
	//if _, err := os.Stat(privateKeyFile); os.IsNotExist(err) {
	//	savePEMKey(privateKeyFile, key)
	//}
	//
	//// saveGobKey(homeDir + "/public.key", publicKey)
	//// saveGobKey(homeDir + "/private.key", key)
	//
	//if _, err := os.Stat(publicKeyFile); os.IsNotExist(err) {
	//	savePublicPEMKey(publicKeyFile, publicKey)
	//}
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

func saveGobKey(fileName string, key interface{}) {
	outFile, err := os.Create(fileName)
	checkError(err)
	defer outFile.Close()

	encoder := gob.NewEncoder(outFile)
	err = encoder.Encode(key)
	checkError(err)
}

func savePEMKey(fileName string, key *rsa.PrivateKey) {
	outFile, err := os.Create(fileName)
	checkError(err)
	defer outFile.Close()

	var privateKey = &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}

	err = pem.Encode(outFile, privateKey)
	checkError(err)
}

func savePublicPEMKey(fileName string, pubkey rsa.PublicKey) {
	asn1Bytes, err := asn1.Marshal(pubkey)
	checkError(err)

	var pemkey = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	pemfile, err := os.Create(fileName)
	checkError(err)
	defer pemfile.Close()

	err = pem.Encode(pemfile, pemkey)
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}