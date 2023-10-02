package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/jritsema/gotoolbox/web"
)

var (
	uploadDir = "./xml"
)

func indexHandler(r *http.Request) *web.Response {
	return web.HTML(http.StatusOK, html, "index.html", data, nil)
}

func credsHandler(r *http.Request) *web.Response {
	switch r.Method {
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

	default:
		return web.Empty(http.StatusMethodNotAllowed)
	}

}

func addCreds(r *http.Request) *web.Response {
	return web.HTML(http.StatusOK, html, "creds-add.html", data, nil)
}

func getScanFiles() ([]string, error) {
	err := os.MkdirAll(uploadDir, 0755)
	if err != nil {
		log.Fatal(err)
	}

	files, err := os.ReadDir(uploadDir)
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
	scanFiles, err := getScanFiles()
	if err != nil {
		log.Fatal(err)
	}

	data := struct {
		ScanFiles []string
	}{
		ScanFiles: scanFiles,
	}

	return web.HTML(http.StatusOK, html, "scans.html", data, nil)
}

func scanHandler(r *http.Request) *web.Response {
	id, segments := web.PathLast(r)

	data := struct {
		ID string
	}{
		ID: id,
	}
	if segments < 2 {
		log.Fatal("meow")
	}

	return web.HTML(http.StatusOK, html, "scan.html", data, nil)
}

func uploadHandler(r *http.Request) *web.Response {
	uploadDir := "./xml"

	switch r.Method {
	case http.MethodGet:
		return web.HTML(http.StatusOK, html, "upload.html", nil, nil)

	case http.MethodPost:
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			return web.HTML(http.StatusBadRequest, html, "upload.html", nil, nil)
		}

		file, _, err := r.FormFile("file")
		if err != nil {
			return web.HTML(http.StatusBadRequest, html, "upload.html", nil, nil)
		}
		defer file.Close()

		err = os.MkdirAll(uploadDir, 0755)
		if err != nil {
			return web.HTML(http.StatusInternalServerError, html, "upload.html", nil, nil)
		}

		fileName := uuid.New().String() + ".xml"

		filePath := filepath.Join(uploadDir, fileName)
		newFile, err := os.Create(filePath)
		if err != nil {
			return web.HTML(http.StatusInternalServerError, html, "upload.html", nil, nil)
		}
		defer newFile.Close()

		_, err = io.Copy(newFile, file)
		if err != nil {
			return web.HTML(http.StatusInternalServerError, html, "upload.html", nil, nil)
		}

		return web.HTML(http.StatusOK, html, "scans.html", nil, nil)

	default:
		return web.HTML(http.StatusBadRequest, html, "upload.html", nil, nil)
	}
}
