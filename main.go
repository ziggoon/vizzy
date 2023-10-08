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

  router.Handle("/admin", authMiddleware(adminMiddleware(web.Action(adminHandler))))
  
	router.Handle("/xml/", authMiddleware(http.StripPrefix("/xml/", http.FileServer(http.Dir("./xml")))))

  router.Handle("/login", web.Action(loginHandler))
  router.Handle("/logout", authMiddleware(web.Action(logoutHandler)))

	router.Handle("/hosts", authMiddleware(web.Action(hostsHandler)))
	router.Handle("/hosts/add", authMiddleware(web.Action(addHostsHandler)))

	router.Handle("/creds", authMiddleware(web.Action(credsHandler)))
	router.Handle("/creds/add", authMiddleware(web.Action(addCredsHandler)))

	router.Handle("/scans", authMiddleware(web.Action(scansHandler)))
	router.Handle("/scan/", authMiddleware(web.Action(scanHandler)))

	router.Handle("/upload", authMiddleware(web.Action(uploadHandler)))

	router.Handle("/", authMiddleware(web.Action(indexHandler)))

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
		log.Print(err)
	}

	hosts, err := getHosts(dbConnection)
	if err != nil {
		log.Print(err)
	}

  password, err := createAdmin(dbConnection)
  if err != nil {
    log.Print(err)
  }

  fmt.Println("Please login with the following credentials:", "admin", password)
	startHTTPServer(dbConnection, creds, hosts)
}
