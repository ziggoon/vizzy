package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jritsema/gotoolbox"
	"github.com/jritsema/gotoolbox/web"
)

var (
	//go:embed all:templates/*
	templateFS embed.FS

	//parsed templates
	html *template.Template
)

func startHTTPServer(dbConnection *Database, creds []Credential, hosts []Host) {
	var err error
	html, err = web.TemplateParseFSRecursive(templateFS, ".html", true, nil)
	if err != nil {
		panic(err)
	}

	router := http.NewServeMux()
	router.Handle("/xml/", http.StripPrefix("/xml/", http.FileServer(http.Dir("./xml"))))

	router.Handle("/", web.Action(indexHandler))

	router.Handle("/hosts", web.Action(hostsHandler))
	router.Handle("/hosts/add", web.Action(addHostsHandler))

	router.Handle("/creds", web.Action(credsHandler))
	router.Handle("/creds/add", web.Action(addCredsHandler))

	router.Handle("/scans", web.Action(scansHandler))
	router.Handle("/scan/", web.Action(scanHandler))

	router.Handle("/upload", web.Action(uploadHandler))

	nextRequestID := func() string {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}

	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	middleware := tracing(nextRequestID)(logging(logger)(router))
	port := gotoolbox.GetEnvWithDefault("PORT", "42069")
	logger.Println("listening on http://localhost:" + port)

	if err := http.ListenAndServe(":"+port, middleware); err != nil {
		logger.Println("http.ListenAndServe():", err)
		os.Exit(1)
	}
}

func handleSigTerms() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("received SIGTERM, exiting")
		os.Exit(1)
	}()
}

func main() {
	handleSigTerms()

	dbConnection, err := createDbConnection()
	if err != nil {
		log.Fatal(err)
	}
	createDbSchema(dbConnection)

	creds, err := getCredentials(dbConnection)
	if err != nil {
		log.Fatal(err)
	}

	hosts, err := getHosts(dbConnection)
	if err != nil {
		log.Fatal(err)
	}

	startHTTPServer(dbConnection, creds, hosts)
}
