package main

import (
	"net/http"
  "log"
	"time"

	"github.com/golang-jwt/jwt"
)

// auth & admin middleware
func authMiddleware(next http.Handler) http.Handler {
  return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
    jwtCookie, err := r.Cookie("jwt")
    if err != nil || jwtCookie == nil {
      http.Redirect(rw, r, "/login", http.StatusSeeOther)
      return
    }

    token, err := jwt.Parse(jwtCookie.Value, func(token *jwt.Token) (interface{}, error) {
      return []byte("poop"), nil
    })

    if err != nil || !token.Valid {
      http.Redirect(rw, r, "/login", http.StatusSeeOther)
      return
    }

    _, ok := token.Claims.(jwt.MapClaims)
    if !ok {
      http.Redirect(rw, r, "/login", http.StatusSeeOther)
      return
    }

    next.ServeHTTP(rw, r)
  })
}

func adminMiddleware(next http.Handler) http.Handler {
  return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
    if !isUserAdmin(r) {
      http.Redirect(rw, r, "/", http.StatusSeeOther)
      return
    }

    next.ServeHTTP(rw, r)
  })
}


// helper funcs
func createJWT(uid int) (string, error) {
  claims := jwt.MapClaims{
    "sub": uid,
    "exp": time.Now().Add(time.Hour*24).Unix(),
  }

  token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

  tokenString, err := token.SignedString([]byte("poop"))
  if err != nil {
    return "", err
  }

  return tokenString, nil
}

func isUserAdmin(r *http.Request) bool {
  jwtCookie, err := r.Cookie("jwt")
  if err != nil || jwtCookie == nil {
    return false
  }

  token, err := jwt.Parse(jwtCookie.Value, func(token *jwt.Token) (interface{}, error) {
    return []byte("poop"), nil
  })
  if err != nil || !token.Valid {
    return false
  }

  claims, ok := token.Claims.(jwt.MapClaims)
  if !ok {
    return false
  }

  dbConnection, err := createDbConnection()
	if err != nil {
		log.Fatal(err)
	}

  user, err := getUserById(dbConnection, int(claims["sub"].(float64)))
  if err != nil {
    log.Fatal(err)
  }

  if !user.Admin {
    return false
  }

  return true
}
