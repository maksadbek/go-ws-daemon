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
	"strconv"
	"strings"
	"syscall"
)

func main() {
	//write pid into file
	f, err := os.Create("/var/run/go-ws-daemon.pid")
	if err != nil {
		panic(err)
	}
	pid := os.Getpid()
	b := []byte(strconv.Itoa(pid))
	_, err = f.Write([]byte(b))
	if err != nil {
		panic(err)
	}
	f.Close()
	//

	confPath := flag.String("conf", "config.toml", "configuration file")
	viewsPath := flag.String("views", "views", "views folder")
	gracefulChild := flag.Bool("graceful", false, "listen on fd open 3")
	flag.Parse()
	log.Println(*confPath)
	log.Println(*viewsPath)
	log.Println(*gracefulChild)
	//TCP listener
	var listener net.Listener
	//Channels
	sig := make(chan os.Signal)
	restart := make(chan bool)
	done := make(chan bool)

	log.Println(*confPath)
	d, err := ioutil.ReadFile(*confPath)
	if err != nil {
		log.Fatal(err)
	}

	reader := strings.NewReader(string(d))

	config, err := conf.Read(reader)
	if err != nil {
		log.Fatal(err)
	}

	//template files
	t, err := template.ParseFiles(
		*viewsPath+"/index.html",
		*viewsPath+"/header.html",
		*viewsPath+"/active.html",
	)
	if err != nil {
		log.Fatal(err)
	}

	//initialize config into route
	route.Initialize(config, t)

	//if gracefulchild flag is set
	// terminate the parent and start child process
	if *gracefulChild {
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

	//setup server handlers
	server := &http.Server{Handler: webHandler()}
	//start server
	go func() {
		server.Serve(listener)
	}()
	log.Println("server is running on", config.SRV.IP, config.SRV.Port)

	//catch only term and hup signals
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGHUP)

	//catch signals in infinite loop
	go func(sign chan os.Signal) {
		for {
			select {
			case s := <-sign:
				//if signal is term, then terminate the program
				if s == syscall.SIGTERM {
					log.Println("SIGTERM")
					done <- true
					break
					// if signal is hup, restart the program
				} else if s == syscall.SIGHUP {
					log.Println("SIGHUP")
					restart <- true
					break
				}
			}
		}
	}(sig)

	// if restart channel is active, then create fork and terminate parent
	go func() {
		<-restart
		listener.Close()
		log.Println("restarting")
		//get the path of executive file
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
