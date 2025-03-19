package main;

import(
	"os"
	"log"
	"fmt"
	"flag"
	"time"
	apiv1 "pxnMetricsBackend/api/apiv1"
	gnet  "github.com/panjf2000/gnet/v2"
);



const DEFAULT_BIND_HOST = "127.0.0.1";
const DEFAULT_BIND_PORT = 9001;
//TODO
//const DEFAULT_CHECKSUM_INIT = 9000;

var EnableReusePort bool;
var EnableThreading bool;



func main() {
	print("\n");
	var bind     string;
	var bindHost string;
	var bindPort int;
	var disableReusePort bool;
	var disableThreading bool;
	flag.StringVar(&bind, "bind", "",
		"The IP and port the server should listen on. Can also use a unix socket /run/pxnMetrics.socket");
	flag.StringVar(&bindHost, "bind-host", "", "Bind to UDP host");
	flag.IntVar(   &bindPort, "bind-port", 0,  "Bind to UDP port");
	flag.BoolVar(&disableReusePort, "disable-reuse-port", false, "Don't reuse ports");
	flag.BoolVar(&disableThreading, "disable-threading",  false, "Don't use multi-core socket handling");
	flag.Parse();
	EnableReusePort = !disableReusePort;
	EnableThreading = !disableThreading;
	if bind == "" {
		if bindHost == "" { bindHost = DEFAULT_BIND_HOST; }
		if bindPort < 1 { bindPort = DEFAULT_BIND_PORT; }
		bind = fmt.Sprintf("%s:%d", bindHost, bindPort);
	}
	if EnableReusePort { log.Print("Enable Reuse-Ports"); }
	if EnableThreading { log.Print("Enable Threading"  ); }
	// udp listener
	log.Printf("API listening on UDP://%s", bind);
	api := apiv1.New();
	log.Fatal(gnet.Run(api, "udp://"+bind, gnet.WithMulticore(EnableThreading), gnet.WithReusePort(EnableReusePort)));
}
