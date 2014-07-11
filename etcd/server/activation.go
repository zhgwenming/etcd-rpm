package server

import (
	"crypto/tls"
	"github.com/coreos/go-systemd/activation"
	"log"
	"net"
	"net/http"
)

var activatedSockets []net.Listener

const expectedSockets = 2

const (
	EtcdSock int = iota
	RaftSock
)

func init() {
	files := activation.Files(false)
	if files == nil || len(files) == 0 {
		// no socket activation attempted
		activatedSockets = nil
	} else if len(files) == expectedSockets {
		// socket activation
		activatedSockets = make([]net.Listener, len(files))
		for i, f := range files {
			var err error
			activatedSockets[i], err = net.FileListener(f)
			if err != nil {
				log.Fatal("socket activation failure: ", err)
			}
		}
	} else {
		// socket activation attempted with incorrect number of sockets
		activatedSockets = nil
		log.Fatalf("socket activation failure: %d sockets received, %d expected.", len(files), expectedSockets)
	}
}

func SocketActivated() bool {
	return activatedSockets != nil
}

func ActivateListenAndServe(srv *http.Server, sockno int) error {
	if !SocketActivated() {
		return srv.ListenAndServe()
	} else {
		return srv.Serve(activatedSockets[sockno])
	}
}

func ActivateListenAndServeTLS(srv *http.Server, sockno int, certFile, keyFile string) error {
	if !SocketActivated() {
		return srv.ListenAndServeTLS(certFile, keyFile)
	} else {
		config := &tls.Config{}
		if srv.TLSConfig != nil {
			*config = *srv.TLSConfig
		}
		if config.NextProtos == nil {
			config.NextProtos = []string{"http/1.1"}
		}

		var err error
		config.Certificates = make([]tls.Certificate, 1)
		config.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return err
		}

		tlsListener := tls.NewListener(activatedSockets[sockno], config)
		return srv.Serve(tlsListener)
	}
}

func GetActivatedPort(sockno int) string {
	activatedAddr := activatedSockets[sockno].Addr().String()
	_, port, err := net.SplitHostPort(activatedAddr)
	if err != nil {
		log.Fatal(err)
	}
	return port
}

func UseActivatedPort(hostport string, sockno int) string {
	port := GetActivatedPort(sockno)

	host, _, err := net.SplitHostPort(hostport)
	if err != nil {
		log.Fatal(err)
	}

	return host + ":" + port
}
