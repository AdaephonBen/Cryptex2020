
// From: https://auth0.com/blog/authentication-in-golang/
//
// go run main.go
// http://localhost:3000 for page that authenticates via auth server to Auth0
// u: amit+testauth0.door2door.io

package main

import (
    "context"
    "time"
    jose "gopkg.in/square/go-jose.v2"
    "github.com/auth0-community/auth0"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "strings"
    "io"
    "fmt"
    "github.com/auth0/go-jwt-middleware"
    "github.com/dgrijalva/jwt-go"
    "github.com/gorilla/handlers"
    "github.com/gorilla/mux"
    "compress/gzip"
    "net/http"
    "os"
)

var ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
var client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
var collection = client.Database("cryptex").Collection("users")

// Define globals here
var mySigningKey = []byte("secret")


func main() {

    fmt.Println("Server started... ")
    router := mux.NewRouter()

    router.Handle("/", http.FileServer(http.Dir("./dist/")))
   // router.Handle("/callback", http.ServeFile())
    // Without JWT middleware check
    // router.Handle("/things", ThingsHandler).Methods("GET")



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




var newUser = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    secret := vars["secret"]
    username := vars["username"]
    fmt.Println(secret + " : " + username)
    ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
    res, err := collection.InsertOne(ctx, bson.M{"secret": secret, "username": username, "level": -1})
    id := res.InsertedID
    fmt.Println(id)
    fmt.Println(err)
})

// Uncomment to not generate tokens within Auth0, but manually instead.

var signingKey = []byte("secret")
var GetTokenHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
 token := jwt.New(jwt.SigningMethodHS256)
 claims := token.Claims.(jwt.MapClaims)
 claims["admin"] = true
 claims["name"] = "Me!"
 claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

 tokenString, _ := token.SignedString(signingKey)
 w.Write([]byte(tokenString))
})

func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        secret := []byte("yx1C48TyLz4XzZzdeXiAOrRz5enXnyjK")
        secretProvider := auth0.NewKeyProvider(secret)
        audience := []string{"AUDI"}

        configuration := auth0.NewConfiguration(secretProvider, audience, "https://cryptex.auth0.com/", jose.HS256)
        validator := auth0.NewValidator(configuration,nil)

        token, err := validator.ValidateRequest(r)

        if err != nil {
            fmt.Println(err)
            fmt.Println("Token is not valid:", token)
            w.WriteHeader(http.StatusUnauthorized)
            w.Write([]byte("Unauthorized"))
        } else {
            next.ServeHTTP(w, r)
        }
    })
}


var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
  ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
    return mySigningKey, nil
  },
  SigningMethod: jwt.SigningMethodHS256,
})