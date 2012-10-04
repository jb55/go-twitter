//
// Copyright 2009 Bill Casarin <billcasarin@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// From gotweet - a command line twitter client by Dmitry Chestnykh
// modified by Bill Casarin
//

package twitter

import (
	"net/http"
	"encoding/base64"
	"io"
	"strings"
	"net"
	"bufio"
	"fmt"
	"bytes"
	"net/url"
)

type readClose struct {
	io.Reader
	io.Closer
}

type badStringError struct {
	what string
	str  string
}

func (e *badStringError) Error() string { return fmt.Sprintf("%s %q", e.what, e.str) }

// Given a string of the form "host", "host:port", or "[ipv6::address]:port",
// return true if the string includes a port.
func hasPort(s string) bool { return strings.LastIndex(s, ":") > strings.LastIndex(s, "]") }

func send(req *http.Request) (resp *http.Response, err error) {
	addr := req.URL.Host
	if !hasPort(addr) {
		addr += ":http"
	}
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	err = req.Write(conn)
	if err != nil {
		conn.Close()
		return nil, err
	}

	reader := bufio.NewReader(conn)
	resp, err = http.ReadResponse(reader, req)
	if err != nil {
		conn.Close()
		return nil, err
	}

	r := io.Reader(reader)
	if n := resp.ContentLength; n != -1 {
		r = io.LimitReader(r, n)
	}
	resp.Body = readClose{r, conn}

	return
}

func encodedUsernameAndPassword(user, pwd string) string {
	bb := &bytes.Buffer{}
	encoder := base64.NewEncoder(base64.StdEncoding, bb)
	encoder.Write([]byte(user + ":" + pwd))
	encoder.Close()
	return bb.String()
}

func authGet(url_, user, pwd string) (r *http.Response, err error) {
	var req http.Request
	h := make(http.Header)
	h.Add("Authorization", "Basic "+
		encodedUsernameAndPassword(user, pwd))
	req.Header = h
	if req.URL, err = url.Parse(url_); err != nil {
		return
	}
	if r, err = send(&req); err != nil {
		return
	}
	return
}

// Post issues a POST to the specified URL.
//
// Caller should close r.Body when done reading it.
func authPost(url_, user, pwd, client, clientURL, version, agent, bodyType string,
	body io.Reader) (r *http.Response, err error) {
	var req http.Request
	req.Method = "POST"
	req.Body = body.(io.ReadCloser)

	h := make(http.Header)
	h.Add("Content-Type", bodyType)
	h.Add("Transfer-Encoding", "chunked")
	h.Add("User-Agent", agent)
	h.Add("X-Twitter-Client", client)
	h.Add("X-Twitter-Client-URL", clientURL)
	h.Add("X-Twitter-Version", version)
	h.Add("Authorization", "Basic "+encodedUsernameAndPassword(user, pwd))
	req.Header = h

	req.URL, err = url.Parse(url_)
	if err != nil {
		return nil, err
	}

	return send(&req)
}

// Do an authenticated Get if we've called Authenticated, otherwise
// just Get it without authentication
func httpGet(url_, user, pass string) (*http.Response, error) {
	var r *http.Response
	var err error

	if user != "" && pass != "" {
		r, err = authGet(url_, user, pass)
	} else {
		r, err = http.Get(url_)
	}

	return r, err
}

// Do an authenticated Post if we've called Authenticated, otherwise
// just Post it without authentication
func httpPost(url_, user, pass, client, clientURL, version, agent,
	data string) (*http.Response, error) {
	var r *http.Response
	var err error

	body := bytes.NewBufferString(data)
	bodyType := "application/x-www-form-urlencoded"

	if user != "" && pass != "" {
		r, err = authPost(url_, user, pass, client, clientURL,
			version, agent, bodyType, body)
	} else {
		r, err = http.Post(url_, bodyType, body)
	}

	return r, err
}
