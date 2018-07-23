package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/eahydra/socks"
	"github.com/elazarl/goproxy"
	"github.com/jamesbcook/print"
)

var (
	appVersion string
	gitCommit  string
)

//UserAgent contains a slice of useragents to be used
type UserAgent struct {
	Names []string `json:"UserAgents"`
}

//SocksProxy contains a slice of proxies to be used
type SocksProxy struct {
	Names []string `json:"Proxies"`
}

//UpstreamDialer contains a slice of socks dialers to be used
type UpstreamDialer struct {
	forwardDialers []socks.Dialer
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

type flags struct {
	userAgentFile string
	socks5File    string
	socksListener string
	httpListener  string
	verbose       bool
}

func flagSetup() *flags {
	uaFile := flag.String("uaFile", "useragents.json", "Json file that contains useragents to use")
	socksFile := flag.String("socksFile", "socks5-proxies.json", "Socks file that contains socks proxies to use")
	socksListen := flag.String("socks", "localhost:9292", "Local socks listener to accept connections")
	httpListen := flag.String("http", "localhost:9293", "HTTP listener to accept connections, this changes the useragent on each request")
	verbose := flag.Bool("verbose", false, "Verbose output from proxy")
	version := flag.Bool("version", false, "Current Version")
	flag.Parse()
	if *version {
		fmt.Printf("proxy-ng v%s %s\n", appVersion, gitCommit)
		os.Exit(0)
	}
	return &flags{userAgentFile: *uaFile, socks5File: *socksFile,
		socksListener: *socksListen, httpListener: *httpListen,
		verbose: *verbose}
}

func main() {
	myFlags := flagSetup()
	f, err := os.Open(myFlags.userAgentFile)
	if err != nil {
		print.Badln(err)
	}
	f2, err := os.Open(myFlags.socks5File)
	if err != nil {
		print.Badln(err)
	}
	var buf bytes.Buffer
	buf.ReadFrom(f)
	ua := &UserAgent{}
	if err := json.Unmarshal(buf.Bytes(), ua); err != nil {
		print.Badln(err)
	}
	buf.Reset()
	proxies := &SocksProxy{}
	buf.ReadFrom(f2)
	if err := json.Unmarshal(buf.Bytes(), proxies); err != nil {
		print.Badln(err)
	}
	var router socks.Dialer
	proxy := goproxy.NewProxyHttpServer()
	if myFlags.verbose {
		router = logBuildUpStream(BuildUpstreamRouter)(proxies.Names)
		proxy.ConnectDial = func(network, address string) (net.Conn, error) {
			return logDialer(router.Dial)(network, address)
		}
		proxy.Tr.Dial = func(network, address string) (net.Conn, error) {
			return logDialer(router.Dial)(network, address)
		}
	} else {
		router = BuildUpstreamRouter(proxies.Names)
		proxy.ConnectDial = func(network, address string) (net.Conn, error) {
			return router.Dial(network, address)
		}
		proxy.Tr.Dial = func(network, address string) (net.Conn, error) {
			return router.Dial(network, address)
		}
	}
	socksListen, err := net.Listen("tcp", myFlags.socksListener)
	if err != nil {
		print.Badln(err)
	}
	if myFlags.verbose {
		print.Goodf("Started socks listener on %s\n", myFlags.socksListener)
	}
	socksvr, err := socks.NewSocks5Server(router)
	if err != nil {
		print.Badln(err)
	}
	httpListen, err := net.Listen("tcp", myFlags.httpListener)
	if err != nil {
		print.Badln(err)
	}
	if myFlags.verbose {
		print.Goodf("Started http listener on %s\n", myFlags.httpListener)
	}
	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)
	proxy.OnRequest().DoFunc(
		func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			r.Header.Set("User-Agent", ua.randomName())
			return r, nil
		})
	go func() {
		http.Serve(httpListen, proxy)
	}()
	print.Goodln("Ready")
	if err := socksvr.Serve(socksListen); err != nil {
		print.Badln(err)
	}
}

func (ua UserAgent) randomName() string {
	max := len(ua.Names)
	if max == 0 {
		return ua.Names[0]
	}
	randomUA := 0 + rand.Intn(max-0)
	return ua.Names[randomUA]
}

//NewUpstreamDialer to be added to the array of dialers
func NewUpstreamDialer(forwardDialers []socks.Dialer) *UpstreamDialer {
	return &UpstreamDialer{
		forwardDialers: forwardDialers,
	}
}

func (u *UpstreamDialer) getRandomDialer() socks.Dialer {
	max := len(u.forwardDialers)
	if max == 0 {
		return u.forwardDialers[0]
	}
	randomDialer := 0 + rand.Intn(max-0)
	return u.forwardDialers[randomDialer]
}

func logDialer(f func(network, address string) (net.Conn, error)) func(network, address string) (net.Conn, error) {
	return func(network, address string) (net.Conn, error) {
		conn, err := f(network, address)
		print.Statusf("Connecting to %v\n", conn.RemoteAddr())
		return conn, err
	}
}

//Dial is a custom dialer that picks a random dialer before it makes it's connection
func (u *UpstreamDialer) Dial(network, address string) (net.Conn, error) {
	router := u.getRandomDialer()
	conn, err := router.Dial(network, address)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func logBuildUpStream(f func(proxies []string) socks.Dialer) func(proxies []string) socks.Dialer {
	return func(proxies []string) socks.Dialer {
		for x := range proxies {
			print.Statusf("Loading %s\n", proxies[x])
		}
		defer print.Goodln("Loading complete")
		return f(proxies)
	}
}

//BuildUpstreamRouter populates the slice of dialers
func BuildUpstreamRouter(proxies []string) socks.Dialer {
	var allForward []socks.Dialer
	for _, proxy := range proxies {
		forward, err := socks.NewSocks5Client("tcp", proxy, "", "", socks.Direct)
		if err != nil {
			log.Fatal(err)
		}
		allForward = append(allForward, forward)
	}
	return NewUpstreamDialer(allForward)
}
