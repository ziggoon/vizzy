package main

import (
	"io"
  "time"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/jritsema/gotoolbox/web"
)

var (
	uploadDir = "./xml"
)

// handlers
func indexHandler(r *http.Request) *web.Response {
	return web.HTML(http.StatusOK, html, "index.html", creds, nil)
}

func adminHandler(r *http.Request) *web.Response {
  return web.HTML(http.StatusOK, html, "admin.html", nil, nil)
}

func hostsHandler(r *http.Request) *web.Response {
	switch r.Method {
	case http.MethodPost:
		row := Host{}
		r.ParseForm()
		row.Hostname = r.Form.Get("hostname")
		row.IpAddress = r.Form.Get("ipaddress")
		row.Os = r.Form.Get("os")
		row.Information = r.Form.Get("information")

		dbConnection, err := createDbConnection()
		if err != nil {
			log.Print(err)
		}

		insertHost(dbConnection, row)
		hosts = append(hosts, row)

		return web.HTML(http.StatusOK, html, "hosts.html", hosts, nil)

	case http.MethodGet:
		return web.HTML(http.StatusOK, html, "hosts.html", hosts, nil)

	default:
		return web.Empty(http.StatusMethodNotAllowed)
	}
}

func addHostsHandler(r *http.Request) *web.Response {
	switch r.Method {
	case http.MethodGet:
		return web.HTML(http.StatusOK, html, "hosts-add.html", hosts, nil)

	default:
		return web.HTML(http.StatusMethodNotAllowed, html, "index.html", creds, nil)
	}
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
			log.Print(err)
		}

		insertCredential(dbConnection, row)
		creds = append(creds, row)

		return web.HTML(http.StatusOK, html, "creds.html", creds, nil)

	case http.MethodGet:
		return web.HTML(http.StatusOK, html, "creds.html", creds, nil)

	default:
		return web.Empty(http.StatusMethodNotAllowed)
	}

}

func addCredsHandler(r *http.Request) *web.Response {
	switch r.Method {
	case http.MethodGet:
		return web.HTML(http.StatusOK, html, "creds-add.html", creds, nil)

	default:
		return web.HTML(http.StatusMethodNotAllowed, html, "index.html", creds, nil)
	}
}

func scansHandler(r *http.Request) *web.Response {
	switch r.Method {
	case http.MethodGet:
		scanFiles, err := getScanFiles()
		if err != nil {
			log.Print(err)
		}

		scans := struct {
			ScanFiles []string
		}{
			ScanFiles: scanFiles,
		}

		return web.HTML(http.StatusOK, html, "scans.html", scans, nil)

	default:
		return web.HTML(http.StatusMethodNotAllowed, html, "index.html", nil, nil)
	}
}

func scanHandler(r *http.Request) *web.Response {
	switch r.Method {
	case http.MethodGet:
		id, segments := web.PathLast(r)
		scan := struct {
			ID string
		}{
			ID: id,
		}
		if segments < 2 {
			log.Print("meow")
		}

		return web.HTML(http.StatusOK, html, "scan.html", scan, nil)

	default:
		return web.HTML(http.StatusMethodNotAllowed, html, "index.html", nil, nil)
	}
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

		file, header, err := r.FormFile("file")
		if err != nil {
			return web.HTML(http.StatusBadRequest, html, "upload.html", nil, nil)
		}
		defer file.Close()

		err = os.MkdirAll(uploadDir, 0755)
		if err != nil {
			return web.HTML(http.StatusInternalServerError, html, "upload.html", nil, nil)
		}

		fileName := header.Filename

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

func loginHandler(r *http.Request) *web.Response {
  switch r.Method {
  case http.MethodGet:
    return web.HTML(http.StatusOK, html, "login.html", nil, nil)

  case http.MethodPost:
    row := User{}
		r.ParseForm()
		row.Username = r.Form.Get("username")
		row.PasswordHash = r.Form.Get("password")
    
    // need to implement authentication logic
    dbConnection, err := createDbConnection()
		if err != nil {
			log.Print(err)
	  }

    user, err := getUserByUsername(dbConnection, row.Username)
    if err != nil {
      log.Print(err)
    }

    err = verifyPassword(user.PasswordHash, row.PasswordHash)
    if err != nil {
      return web.HTML(http.StatusBadRequest, html, "login.html", nil, nil)
    }

    jwt, err := createJWT(user.Id)
    if err != nil {
      log.Print(err)
    }
 
    cookieValue := "jwt=" + jwt + "; Path=/; HttpOnly; Expires=" + time.Now().Add(time.Hour*24).Format(time.RFC1123)

    headers := map[string]string{
      "Set-Cookie": cookieValue,
      "Location": "/",
    }

    return web.HTML(http.StatusAccepted, html, "index.html", nil, headers)
  
  default:
    return web.HTML(http.StatusBadRequest, html, "login.html", nil, nil)
  }
}

func logoutHandler(r *http.Request) *web.Response {
  switch r.Method {
  case http.MethodGet:
    cookieValue := "jwt=" + "" + "; Path=/; HttpOnly; Expires=" + time.Now().Add(time.Hour*24).Format(time.RFC1123)
    
    headers := map[string]string{
      "Set-Cookie": cookieValue,
      "Location": "/",
    }
 
    return web.HTML(http.StatusSeeOther, html, "logout.html", nil, headers)

  default:
    return web.HTML(http.StatusOK, html, "index.html", nil, nil)
  }
}

// helpers
func getScanFiles() ([]string, error) {
	err := os.MkdirAll(uploadDir, 0755)
	if err != nil {
		log.Print(err)
	}

	files, err := os.ReadDir(uploadDir)
	if err != nil {
		log.Print(err)
	}

	var fileNames []string
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}

	return fileNames, nil
}
