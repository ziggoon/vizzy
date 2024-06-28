package main

import (
	"log"
	"net/http"
	"time"
  "crypto/rand"
  "encoding/hex"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var (
  JWT_SECRET, err = generateRandomString(16)
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
      return []byte(JWT_SECRET), nil
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
func hashPassword(password string) (string, error) {
  cost := bcrypt.DefaultCost
  hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), cost)
  if err != nil {
    return "", err
  }

  return string(hashedPassword), nil
}

func verifyPassword(hash, password string) error {
  return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func generateRandomString(length int) (string, error) {
  byteLength := length / 2

  randomBytes := make([]byte, byteLength)

  _, err := rand.Read(randomBytes)
  if err != nil {
      return "", err
  }

  randomString := hex.EncodeToString(randomBytes)

  return randomString, nil
}

func createJWT(uid int) (string, error) {
  claims := jwt.MapClaims{
    "sub": uid,
    "exp": time.Now().Add(time.Hour*24).Unix(),
  }

  token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

  tokenString, err := token.SignedString([]byte(JWT_SECRET))
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
    return []byte(JWT_SECRET), nil
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
		log.Print(err)
	}

  user, err := getUserById(dbConnection, int(claims["sub"].(float64)))
  if err != nil {
    log.Print(err)
  }

  if !user.Admin {
    return false
  }

  return true
}
