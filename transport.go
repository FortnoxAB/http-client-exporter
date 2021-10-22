package main

import (
	"net/http"
	"net/http/httptrace"
	"sync"
	"time"
)

type roundTripTrace struct {
	start         time.Time
	dnsDone       time.Time
	connectDone   time.Time
	gotConn       time.Time
	responseStart time.Time
	end           time.Time
}

func (rtt roundTripTrace) Total() float64 {
	dur := rtt.end.Sub(rtt.start)
	return dur.Seconds()
}
func (rtt roundTripTrace) TotalDNS() float64 {
	dur := rtt.dnsDone.Sub(rtt.start)
	return dur.Seconds()
}
func (rtt roundTripTrace) TotalTransfer() float64 {
	dur := rtt.end.Sub(rtt.responseStart)
	return dur.Seconds()
}
func (rtt roundTripTrace) TotalProcessing() float64 {
	dur := rtt.responseStart.Sub(rtt.gotConn)
	return dur.Seconds()
}

func (rtt roundTripTrace) TotalConnect() float64 {
	dur := rtt.gotConn.Sub(rtt.dnsDone)
	return dur.Seconds()
}
func (rtt roundTripTrace) Observe() {
	promMetrics.ScanDur.WithLabelValues("resolve").Observe(rtt.TotalDNS())
	if rtt.gotConn.IsZero() {
		return
	}

	promMetrics.ScanDur.WithLabelValues("connect").Observe(rtt.TotalConnect())

	if rtt.responseStart.IsZero() {
		return
	}

	promMetrics.ScanDur.WithLabelValues("processing").Observe(rtt.TotalProcessing())

	if rtt.end.IsZero() {
		return
	}

	promMetrics.ScanDur.WithLabelValues("total").Observe(rtt.Total())
	promMetrics.ScanDur.WithLabelValues("transfer").Observe(rtt.TotalTransfer())
}

type transport struct {
	Transport http.RoundTripper

	mu      sync.Mutex
	traces  []*roundTripTrace
	current *roundTripTrace
}

func newTransport(t http.RoundTripper) *transport {
	return &transport{
		Transport: t,
		traces:    []*roundTripTrace{},
	}
}

// RoundTrip switches to a new trace, then runs embedded RoundTripper.
func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	trace := &roundTripTrace{}
	t.current = trace
	t.traces = append(t.traces, trace)

	return t.Transport.RoundTrip(req)
}
func (t *transport) DNSStart(_ httptrace.DNSStartInfo) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.current.start = time.Now()
}
func (t *transport) DNSDone(_ httptrace.DNSDoneInfo) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.current.dnsDone = time.Now()
}
func (ts *transport) ConnectStart(_, _ string) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	t := ts.current
	// No DNS resolution because we connected to IP directly.
	if t.dnsDone.IsZero() {
		t.start = time.Now()
		t.dnsDone = t.start
	}
}
func (t *transport) ConnectDone(net, addr string, err error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.current.connectDone = time.Now()
}
func (t *transport) GotConn(_ httptrace.GotConnInfo) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.current.gotConn = time.Now()
}
func (t *transport) GotFirstResponseByte() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.current.responseStart = time.Now()
}

/*
func (t *transport) TLSHandshakeStart() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.current.tlsStart = time.Now()
}
func (t *transport) TLSHandshakeDone(_ tls.ConnectionState, _ error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.current.tlsDone = time.Now()
}
*/
