package main

import (
	"flag"
	"net/http"
	"net/http/httptrace"
	"time"

	"github.com/fortnoxab/ginprometheus"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var addr = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")

var promMetrics = NewMetrics()

func main() {
	flag.Parse()
	g := gin.Default()

	m := ginprometheus.New("http")
	m.Use(g)

	g.GET("/test", testEndpoint)

	err := http.ListenAndServe(*addr, g)
	if err != nil {
		logrus.Error(err)
	}
}

func testEndpoint(c *gin.Context) {
	u := c.Query("url")

	req, err := http.NewRequestWithContext(c.Request.Context(), "GET", u, nil)

	if err != nil {
		logrus.Error(err)
		return
	}

	t := http.DefaultTransport.(*http.Transport).Clone()
	t.DisableKeepAlives = true
	tt := newTransport(t)
	httpClient := &http.Client{Timeout: 5 * time.Second, Transport: tt}

	trace := &httptrace.ClientTrace{
		DNSStart:             tt.DNSStart,
		DNSDone:              tt.DNSDone,
		ConnectStart:         tt.ConnectStart,
		ConnectDone:          tt.ConnectDone,
		GotConn:              tt.GotConn,
		GotFirstResponseByte: tt.GotFirstResponseByte,
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))

	resp, err := httpClient.Do(req)
	if err != nil {
		logrus.Error(err)
		return
	}
	resp.Body.Close()
	tt.current.end = time.Now()
	tt.current.Observe()
}
