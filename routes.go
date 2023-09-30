package main

import (
	"log"
	"net/http"
	"os"

	"github.com/jritsema/gotoolbox/web"
)

func indexHandler(r *http.Request) *web.Response {
	return web.HTML(http.StatusOK, html, "index.html", data, nil)
}

func credsHandler(r *http.Request) *web.Response {
	switch r.Method {
	/*case http.MethodGet:
	return web.HTML(http.StatusOK, html, "creds.html", data, nil)
	*/
	case http.MethodPost:
		row := Credential{}
		r.ParseForm()
		row.Username = r.Form.Get("username")
		row.Password = r.Form.Get("password")
		row.Host = r.Form.Get("host")
		row.Information = r.Form.Get("information")

		dbConnection, err := createDbConnection()
		if err != nil {
			log.Fatal(err)
		}

		insertCredential(dbConnection, row)
		data = append(data, row)

		return web.HTML(http.StatusOK, html, "creds.html", data, nil)
	}
	return web.Empty(http.StatusNotImplemented)
}

func addCreds(r *http.Request) *web.Response {
	return web.HTML(http.StatusOK, html, "creds-add.html", data, nil)
}

func getScanFiles() ([]string, error) {
	files, err := os.ReadDir("./scans")
	if err != nil {
		log.Fatal(err)
	}

	var fileNames []string
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}

	return fileNames, nil
}

func scansHandler(r *http.Request) *web.Response {
	id, _ := web.PathLast(r)
	scanFiles, err := getScanFiles()
	if err != nil {
		log.Fatal(err)
	}

	data := struct {
		ID        string
		ScanFiles []string
	}{
		ID:        id,
		ScanFiles: scanFiles,
	}

	//fmt.Println("scans.html returned")
	return web.HTML(http.StatusOK, html, "scans.html", data, nil)
}
