
package main    
import (
    "encoding/json"
    "io"
    "compress/gzip"
    "time"
    "github.com/graphql-go/graphql"
    // "github.com/graphql-go/handler"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/bson"
    "context"
    "fmt"
    "net/http"
    "strings"
    "errors"
    // "log"
    "os"
    "github.com/gorilla/handlers"
    // "github.com/codegangsta/negroni"
    // "github.com/auth0/go-jwt-middleware"
    "github.com/dgrijalva/jwt-go"
    "github.com/gorilla/mux"
    "go.mongodb.org/mongo-driver/mongo"
    // "github.com/joho/godotenv"
)

type Response struct {
    Message string `json:"message"`
}

type Jwks struct {
    Keys []JSONWebKeys `json:"keys"`
}

type JSONWebKeys struct {
    Kty string `json:"kty"`
    Kid string `json:"kid"`
    Use string `json:"use"`
    N string `json:"n"`
    E string `json:"e"`
    X5c []string `json:"x5c"`
}

type user struct {
    clientID string `json:"clientID"`
    username string  `json:"username"`
    level int `json:"level"`
}

// var ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
// var client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
// var collection = client.Database("cryptex").Collection("users")

func main() {
    fmt.Println("Server started... ")
    fmt.Println("To do : Protect all endpoints with JWT Auth")
    fmt.Println("Change level type to int. It's string rn. ")
    // jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options {
    //     ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
    //         // Verify 'aud' claim
    //         aud := "https://cryptex2020.auth0.com/api/v2/"
    //         checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
    //         if !checkAud {
    //             return token, errors.New("Invalid audience.")
    //         }
    //         // Verify 'iss' claim
    //         iss := "https://cryptex2020.auth0.com/"
    //         checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
    //         if !checkIss {
    //             return token, errors.New("Invalid issuer.")
    //         }

    //         cert, err := getPemCert(token)
    //         if err != nil {
    //             panic(err.Error())
    //         }

    //         result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
    //         return result, nil
    //     },
    //     SigningMethod: jwt.SigningMethodRS256,
    // })
    router := mux.NewRouter()

    router.Handle("/", http.FileServer(http.Dir("./dist/")))
    // router.Handle("/callback", http.ServeFile())
    // Without JWT middleware check
    // router.Handle("/things", ThingsHandler).Methods("GET")

    // ALL API CALLS (GraphQL) are defined here

    router.HandleFunc("/adduser/{ID}/{username}/{secret}", AddUser)
    // Define GraphQL User Type :
    // userType := graphql.NewObject(graphql.ObjectConfig{
    //     Name: "User", 
    //     Fields : graphql.Fields{
    //         "clientID":&graphql.Field{
    //             Type: graphql.String,
    //         },
    //         "username":&graphql.Field{
    //             Type: graphql.String,
    //         },
    //         "level":&graphql.Field{
    //             Type: graphql.Int,
    //         },
    //     },
    // })
    // Define GraphQL Root Query : (Every field in this RootQuery represents a possible query)
    rootQuery := graphql.NewObject(graphql.ObjectConfig{
        Name: "Query",
        Fields: graphql.Fields{
            "level":&graphql.Field{
                Type: graphql.String,
                Args: graphql.FieldConfigArgument {
                    "clientID": &graphql.ArgumentConfig {
                        Type: graphql.String,
                    },
                },
                Resolve: func(p graphql.ResolveParams) (interface {}, error) {
                    // MongoDB Initial Code
                    client, _ := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
                    ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
                    _ = client.Connect(ctx)
                    collection := client.Database("Cryptex").Collection("users")
                    ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
                    // Querying for the right user
                    filter := bson.M{"clientID" : p.Args["clientID"].(string)}
                    var result map[string]interface{}
                    _ = collection.FindOne(ctx, filter).Decode(&result)
                    // Returning the level of the queried user
                    if result["level"] == nil {
                        return "-2", nil
                    }
                    return result["level"], nil
                },
            },
            "doesUsernameExist":&graphql.Field{
                Type: graphql.Boolean,
                Args: graphql.FieldConfigArgument {
                    "username":&graphql.ArgumentConfig {
                        Type: graphql.String,
                    },
                },
                Resolve: func(p graphql.ResolveParams) (interface {}, error) {
                    // MongoDB Initial Code
                    client, _ := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
                    ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
                    _ = client.Connect(ctx)
                    collection := client.Database("Cryptex").Collection("users")
                    ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
                    // Querying for the right user
                    filter := bson.M{"username" : p.Args["username"].(string)}
                    var result map[string]interface{}
                    _ = collection.FindOne(ctx, filter).Decode(&result)
                    // Returning the level of the queried user
                    if result["level"] == nil {
                        return false, nil
                    }
                    return true, nil
                },
            },
            "level1":&graphql.Field{
                Type: graphql.String,
                Args: graphql.FieldConfigArgument{
                    "ID":&graphql.ArgumentConfig{
                        Type: graphql.String,                  
                    },
                },
            },
        },
    })
    // Create schema with Root Query and Mutator
    var schema, _ = graphql.NewSchema(
        graphql.SchemaConfig{
            Query: rootQuery,
        },
    )

    router.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
        result := executeQuery(r.URL.Query().Get("query"), schema)
        json.NewEncoder(w).Encode(result)
    })
    // END OF BLOCK DEFINING API CALLS

    // router.Handle("/api/private", negroni.New(
    //     negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
    //     negroni.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    //         message := "Hello from a private endpoint! You need to be authenticated to see this."
    //         responseJSON(message, w, http.StatusOK)
    // }))))

    // Not necessary when wired up to Auth0 to get tokens
    // router.Handle("/get-token", GetTokenHandler).Methods("GET")

    router.PathPrefix("/").Handler(http.FileServer(http.Dir("./dist/")))
    http.ListenAndServe(":8080", gzipHandler(handlers.LoggingHandler(os.Stdout, router)))
}

type gzipResponseWriter struct {
    io.Writer
    http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {    return w.Writer.Write(b)   }

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
// var newUser = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request)
// {
//     vars := mux.Vars(r)
//     secret := vars["secret"]
//     username := vars["username"]
//     fmt.Println(secret + " : " + username)
//     ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
//     res, err := collection.InsertOne(ctx, bson.M{"secret": secret, "username": username, "level": -1})
//     id := res.InsertedID
//     fmt.Println(id)
//     fmt.Println(err)
// })
func getPemCert(token *jwt.Token) (string, error) {
    cert := ""
    resp, err := http.Get("https://cryptex2020.auth0.com/.well-known/jwks.json")

    if err != nil {
        return cert, err
    }
    defer resp.Body.Close()

    var jwks = Jwks{}
    err = json.NewDecoder(resp.Body).Decode(&jwks)

    if err != nil {
        return cert, err
    }

    for k, _ := range jwks.Keys {
        if token.Header["kid"] == jwks.Keys[k].Kid {
            cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
        }
    }

    if cert == "" {
        err := errors.New("Unable to find appropriate key.")
        return cert, err
    }

    return cert, nil
}
func String(n int32) string {
    buf := [11]byte{}
    pos := len(buf)
    i := int64(n)
    signed := i < 0
    if signed {
        i = -i
    }
    for {
        pos--
        buf[pos], i = '0'+byte(i%10), i/10
        if i == 0 {
            if signed {
                pos--
                buf[pos] = '-'
            }
            return string(buf[pos:])
        }
    }
}
func responseJSON(message string, w http.ResponseWriter, statusCode int) {
    enableCors(&w);
    response := Response{message}

    jsonResponse, err := json.Marshal(response)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    w.Write(jsonResponse)
}
func AddUser(w http.ResponseWriter, request *http.Request) {
    vars := mux.Vars(request)
    // MongoDB
    client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
    ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
    err = client.Connect(ctx)
    fmt.Println(err)
    collection := client.Database("Cryptex").Collection("users")
    ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
    find, err := collection.Find(ctx, bson.M{"clientID": vars["ID"]})
    JSOND, err := json.Marshal(find.Next(ctx))
    UserStatus := string(JSOND)
    if strings.Compare(UserStatus, "false") == 0 {
        res, _ := collection.InsertOne(ctx, bson.M{"clientID":vars["ID"], "username":vars["username"], "level": -1, "secret": vars["secret"][0:378]})
        fmt.Println("Added a new user to MongoDB")
        fmt.Println("MongoDB ID ")
        fmt.Println(res.InsertedID)
    }
}
// func submitAnswer(w http.ResponseWriter, request *http.Request) {
//     vars := mux.Vars(request)
//     client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
//     ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
//     _ = client.Connect(ctx)
//     collection := client.Database("Cryptex").Collection("users")
//     ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
    
// }
// Function is obsoelete, implemented using GraphQL in main()
// func RetrieveLevel(w http.ResponseWriter, request *http.Request) {
//     vars := mux.Vars(request)
//     client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
//     ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
//     err = client.Connect(ctx)
//     fmt.Println(err)
//     collection := client.Database("Cryptex").Collection("users")
//     ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
//     filter := bson.M{"clientID" : vars["ID"]}
//     var result map[string]interface{}
//     err = collection.FindOne(ctx, filter).Decode(&result)

// }
func executeQuery(query string, schema graphql.Schema) *graphql.Result {
    result := graphql.Do(graphql.Params{
        Schema:        schema,
        RequestString: query,
    })
    if len(result.Errors) > 0 {
        fmt.Printf("wrong result, unexpected errors: %v", result.Errors)
    }
    return result
}
// Does not provide any fine tuning. Adjust CORS funciton later. 
func enableCors(w *http.ResponseWriter) {
    (*w).Header().Set("Access-Control-Allow-Origin", "*")
}