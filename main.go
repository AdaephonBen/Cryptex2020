
// From: https://auth0.com/blog/authentication-in-golang/
//
// go run main.go
// http://localhost:3000 for page that authenticates via auth server to Auth0
// u: amit+testauth0.door2door.io

package main

import (
    "strings"
    "io"
    "encoding/json"
    "errors"
    "github.com/auth0/go-jwt-middleware"
    "github.com/dgrijalva/jwt-go"
    "github.com/gorilla/handlers"
    "github.com/gorilla/mux"
    "github.com/joho/godotenv"
    "log"
    "compress/gzip"
    "net/http"
    "os"
)

type Thing struct {
    Id   int
    Name string
    Slug string
}

var things = []Thing{
    Thing{Id: 1, Name: "Thing 1", Slug: "Slug 1"},
    Thing{Id: 2, Name: "Thing 2", Slug: "Slug 2"},
    Thing{Id: 3, Name: "Thing 3", Slug: "Slug 3"},
    Thing{Id: 4, Name: "Thing 4", Slug: "Slug 4"},
    Thing{Id: 5, Name: "Thing 5", Slug: "Slug 5"},
}

func main() {

    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error reading .env file")
    }

    router := mux.NewRouter()

    router.Handle("/", http.FileServer(http.Dir("./dist/")))
   // router.Handle("/callback", http.ServeFile())
    router.Handle("/status", Status).Methods("GET")
    // Without JWT middleware check
    // router.Handle("/things", ThingsHandler).Methods("GET")
    router.Handle("/things", jwtMiddleware.Handler(ThingsHandler)).Methods("GET")

    router.Handle("/thing/{slug}/foo", jwtMiddleware.Handler(AddToThingHandler)).Methods("POST")

    // Not necessary when wired up to Auth0 to get tokens
    // router.Handle("/get-token", GetTokenHandler).Methods("GET")

    router.PathPrefix("/").Handler(http.FileServer(http.Dir("./dist/")))

    http.ListenAndServe(":8080", gzipHandler(handlers.LoggingHandler(os.Stdout, router)))
}

/******************************************/
/* Handlers for respective HTTP responses */
/******************************************/
type gzipResponseWriter struct {
    io.Writer
    http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
    return w.Writer.Write(b)
}

// EnableGZIP will attempt to compress the response if the client has passed a
// header value for Accept-Encoding which allows gzip
func EnableGZIP(fn http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        gz := gzip.NewWriter(w)
        defer gz.Close()
        gzr := gzipResponseWriter{Writer: gz, ResponseWriter: w}
        fn.ServeHTTP(gzr, r)
    })
}
func gzipHandler(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
            h.ServeHTTP(w, r)
            return
        }
        w.Header().Set("Content-Encoding", "gzip")
        gz := gzip.NewWriter(w)
        defer gz.Close()
        h.ServeHTTP(gzipResponseWriter{Writer: gz, ResponseWriter: w}, r)
    })
}
var NotImplemented = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Not Implemented"))
})

var Status = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("We're OK"))
})

var ThingsHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    payload, _ := json.Marshal(things)
    w.Header().Set("Content-Type", "application/json")
    w.Write([]byte(payload))
})

var AddToThingHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    var thing Thing
    vars := mux.Vars(r)
    slug := vars["slug"]

    for _, t := range things {
        if t.Slug == slug {
            thing = t
        }
    }

    w.Header().Set("Content-Type", "application/json")

    if thing.Slug != "" {
        payload, _ := json.Marshal(thing)
        w.Write([]byte(payload))
    } else {
        w.Write([]byte("Thing not found..."))
    }
})

// Uncomment to not generate tokens within Auth0, but manually instead.

// var signingKey = []byte("secret")
// var GetTokenHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//  token := jwt.New(jwt.SigningMethodHS256)
//  claims := token.Claims.(jwt.MapClaims)
//  claims["admin"] = true
//  claims["name"] = "Me!"
//  claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

//  tokenString, _ := token.SignedString(signingKey)
//  w.Write([]byte(tokenString))
// })

var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
    ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
        signingKey := []byte(os.Getenv("AUTH0_CLIENT_SECRET"))
        if len(signingKey) == 0 {
            return nil, errors.New("No Auth0 Client Secret set.")
        }
        return signingKey, nil
    },
    SigningMethod: jwt.SigningMethodHS256,
    // Extractor:     jwtmiddleware.FromAuthHeader,
    // Debug:         true,
})