package grpchttpdialer

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
)

// _Client ...
type _Client struct {
	serverAddr string
}

// Dialer ...
func Dialer(serverHTTPAddr string) func(context.Context, string) (net.Conn, error) {
	return (&_Client{
		serverAddr: serverHTTPAddr,
	}).Dial
}

// Dial ...
func (c *_Client) Dial(ctx context.Context, addr string) (net.Conn, error) {
	conn, err := net.Dial("tcp", c.serverAddr)
	if err != nil {
		return nil, err
	}

	closeConn := conn
	defer func() {
		if closeConn != nil {
			closeConn.Close()
		}
	}()

	req, err := http.NewRequestWithContext(ctx,
		http.MethodGet,
		ProxyPath+`?addr=`+addr,
		nil)
	if err != nil {
		log.Println(`new request error`, err)
		return nil, err
	}

	req.URL.Host = c.serverAddr
	req.Header.Set(`Connection`, `Upgrade`)
	req.Header.Set(`Upgrade`, gUpgrade)
	if err := req.Write(conn); err != nil {
		log.Println(`write upgrade failed`, err)
		return nil, err
	}

	br := bufio.NewReader(conn)
	resp, err := http.ReadResponse(br, req)
	if err != nil {
		log.Println(`read response error`, err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 101 {
		return nil, fmt.Errorf("statusCode != 101")
	}

	if br.Buffered() > 0 {
		n := br.Buffered()
		b := make([]byte, n)
		br.Read(b)
		log.Println("br.Buffered > 0")
		return nil, fmt.Errorf("br.Buffered() > 0: %s", string(b))
	}

	if _, err := conn.Write(gHandshake); err != nil {
		log.Println(`write handshake failed`, err)
		return nil, err
	}

	closeConn = nil
	return conn, nil
}
