package grpchttpdialer

import (
	"bytes"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
)

var gUpgrade = `grpc-http-dialer/1.0`
var gHandshake = []byte(`grpc-http-dialer`)

// ProxyPath ...
var ProxyPath = `/grpc-http-dialer`

// Server is an HTTP server that supports upgrading from HTTP connection to TCP connection using non HTTP CONNECT method.
type _Server struct {
}

// Handler ...
func Handler() http.HandlerFunc {
	return (&_Server{}).ServeHTTP
}

// ServeHTTP ...
func (s *_Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Println(`method error`)
		return
	}
	addr := r.URL.Query().Get(`addr`)
	upgrade := r.Header.Get("Upgrade")
	if upgrade != gUpgrade {
		log.Println(`upgrade error`)
		return
	}
	outConn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Println(`dial error`)
		return
	}
	w.WriteHeader(http.StatusSwitchingProtocols)
	conn, bio, err := w.(http.Hijacker).Hijack()
	if err != nil {
		log.Println(`hijack error`)
		outConn.Close()
		return
	}

	b := make([]byte, len(gHandshake))
	_, err = io.ReadFull(io.MultiReader(bio, conn), b)
	if err != nil || !bytes.Equal(gHandshake, b) {
		log.Println(err)
		conn.Close()
		outConn.Close()
		return
	}
	log.Println(`recv handshake ok`)

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		defer conn.Close()
		if err := bio.Writer.Flush(); err != nil {
			return
		}
		io.Copy(conn, outConn)
	}()
	go func() {
		defer wg.Done()
		defer outConn.Close()
		if n := bio.Reader.Buffered(); n > 0 {
			if _, err := io.CopyN(outConn, bio.Reader, int64(n)); err != nil {
				return
			}
		}
		io.Copy(outConn, conn)
	}()
	wg.Wait()
}
