package main

import (
  //"fmt"
  "log"
  "net/http"
  
  "github.com/jritsema/gotoolbox/web"
)

func index(r *http.Request) *web.Response {
	return web.HTML(http.StatusOK, html, "index.html", data, nil)	
}

func login(r *http.Request) *web.Response {
  switch r.Method {
  case http.MethodGet:
    return web.HTML(http.StatusOK, html, "login.html", nil, nil)

  case http.MethodPost:
    row := User{}
    r.ParseForm()
    row.Username = r.Form.Get("username")
    row.Password = r.Form.Get("password")
    
    dbConnection, err := createDbConnection()
    if err != nil {
      log.Fatal(err)
    }

    testy := verifyLogin(dbConnection, row.Username, row.Password) 
    
    if testy != true {
      return web.HTML(http.StatusOK, html, "login.html", nil, nil)
    }

    return web.HTML(http.StatusOK, html, "index.html", data, nil)
  }

  return web.Empty(http.StatusNotImplemented)
}

func creds(r *http.Request) *web.Response {
  switch r.Method {
  case http.MethodGet:
    return web.HTML(http.StatusOK, html, "index.html", data, nil)

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

func scans(r *http.Request) *web.Response {
  //id, segments := web.PathLast(r)
  return web.HTML(http.StatusOK, html, "scans.html", nil, nil)
}
