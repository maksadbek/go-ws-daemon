package main

import (
	"flag"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"

	"github.com/Maksadbek/go-ws-daemon/conf"
	"github.com/Maksadbek/go-ws-daemon/route"

	"net"
	"net/http"
	"strings"
	"syscall"
)

func main() {
	confPath := flag.String("conf", "config.toml", "configuration file")
	viewsPath := flag.String("views", "views", "views folder")
	d, err := ioutil.ReadFile(*confPath)
	if err != nil {
		log.Fatal(err)
	}

	reader := strings.NewReader(string(d))

	config, err := conf.Read(reader)
	if err != nil {
		log.Fatal(err)
	}

	t, err := template.ParseFiles(
		*viewsPath+"/index.html",
		*viewsPath+"/header.html",
		*viewsPath+"/active.html",
	)
	if err != nil {
		log.Fatal(err)
	}

	route.Initialize(config, t)

	server := &http.Server{Handler: webHandler()}
	var listener net.Listener

	// start the server
	var gracefulChild bool
	flag.BoolVar(&gracefulChild, "graceful", false, "listen on fd open 3")
	flag.Parse()

	//if gracefulchild flag is set
	// terminate the parent and start child process
	if gracefulChild {
		//terminating parent
		parent := syscall.Getppid()
		log.Printf("main: killing parent pid : %v", parent)
		syscall.Kill(parent, syscall.SIGTERM)
		listener, err = net.Listen("tcp", config.SRV.IP+config.SRV.Port)

		log.Println("main: listening on existing file descriptor 3")
	} else {
		listener, err = net.Listen("tcp", config.SRV.IP+config.SRV.Port)
		if err != nil {
			panic(err)
		}
		log.Println("main: listening on a new file descriptor")
	}

	go func() {
		server.Serve(listener)
	}()
	log.Println(">>>>>>>>>>>>>>>>>>>>>>>..")

	sig := make(chan os.Signal)
	restart := make(chan bool)
	done := make(chan bool)

	signal.Notify(sig, syscall.SIGTERM, syscall.SIGHUP)

	go func(sign chan os.Signal) {
		for {
			select {
			case s := <-sign:
				if s == syscall.SIGTERM {
					log.Println("SIGTERM")
					done <- true
					break
				} else if s == syscall.SIGHUP {
					log.Println("SIGHUP")
					restart <- true
					break
				}
			}
		}
	}(sig)

	go func() {
		<-restart
		listener.Close()
		log.Println("restarting")
		//geting path of the file
		absPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
		log.Println(absPath)

		if err != nil {
			panic(err)
		}
		args := []string{"-graceful"}

		cmd := exec.Command(absPath+"/go-ws-daemon", args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err = cmd.Start()
		if err != nil {
			log.Println(err)
		}
	}()

	<-done
	os.Exit(1)
}

func webHandler() http.Handler {
	web := http.NewServeMux()
	web.Handle("/favicon.ico", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	web.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	web.HandleFunc("/active-orders", route.GetActiveOrders)
	web.HandleFunc("/orders", route.GetOrders)
	web.HandleFunc("/active-logs", route.GetOrderLogs)
	return web
}
