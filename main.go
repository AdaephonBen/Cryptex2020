
package main    
import (
    "strconv"
    "encoding/json"
    "time"
    "io"
    "compress/gzip"
    "github.com/graphql-go/graphql"
    // "github.com/graphql-go/handler"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/bson"
    "context"
    "fmt"
    "net/http"
    "strings"
    "errors"
    "log"
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
    Username string  `json:"username"`
    Level int `json:"level"`
}

type DatabaseUserObject struct {
    Secret string `json:"secret"`
    ClientID string `json:"clientID"`
    Username string `json:"username"`
    Level int `json:"level"`
}

type LevelResponse struct {
    Level int 
    URL string
}

var answers map[string]string


// var context.TODO(), _ = context.WithTimeout(context.Background(), 10*time.Second)
// var client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
// var collection = client.Database("cryptex").Collection("users")

var collection *mongo.Collection 

func main() {
    fmt.Println("Server started... ")
    fmt.Println("To do : Protect all endpoints with JWT Auth")
    fmt.Println("Change level type to int. It's string rn. ")
    answers = make(map[string]string)
    answers["0"] = "triskaidekaphobia"
    answers["1"] = "nerdfameagain"
    answers["2"] = "ireland"
    answers["3"] = "beatles"
    answers["4"] = "magic"
    answers["5"] = "pabloescobar"
    // Set client options
    clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
    // Connect to MongoDB
    client, err := mongo.Connect(context.TODO(), clientOptions)
    if err != nil {
        log.Fatal(err)
    }
    // Check the connection
    err = client.Ping(context.TODO(), nil)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Connected to MongoDB!")
    collection = client.Database("Cryptex").Collection("users")
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
    router.HandleFunc("/acceptedrules/{secret}", AcceptedRules)
    router.HandleFunc("/answer/{secret}/{level}/{answer}", AnswerQuestion)
    router.HandleFunc("/level/{secret}", LevelHandler)
    router.HandleFunc("/leaderboard", LeaderboardHandler)
    router.HandleFunc("/leaderboardtable", LeaderboardTableHandler)
    router.HandleFunc("/css", CSSHandler)
    router.HandleFunc("/midi.mid", MIDIHandler)
    router.HandleFunc("/rules", RulesHandler)
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
                    // Querying for the right user
                    filter := bson.M{"clientID" : p.Args["clientID"].(string)}
                    var result map[string]interface{}
                    _ = collection.FindOne(context.TODO(), filter).Decode(&result)
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
                    // Querying for the right user
                    filter := bson.M{"username" : p.Args["username"].(string)}
                    var result map[string]interface{}
                    _ = collection.FindOne(context.TODO(), filter).Decode(&result)
                    // Returning the level of the queried user
                    if result["level"] == nil {
                        return false, nil
                    }
                    return true, nil
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

    router.PathPrefix("/cryptex/").Handler(http.StripPrefix("/cryptex/",http.FileServer(http.Dir("./dist/"))))
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
//     context.TODO(), _ = context.WithTimeout(context.Background(), 5*time.Second)
//     res, err := collection.InsertOne(context.TODO(), bson.M{"secret": secret, "username": username, "level": -1})
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
    find, _ := collection.Find(context.TODO(), bson.M{"clientID": vars["ID"]})
    JSOND, _ := json.Marshal(find.Next(context.TODO()))
    UserStatus := string(JSOND)
    if strings.Compare(UserStatus, "false") == 0 {
        res, _ := collection.InsertOne(context.TODO(), bson.M{"clientID":vars["ID"], "username":vars["username"], "level": -1, "secret": vars["secret"][0:378], "lastModified": time.Now().UTC()})
        fmt.Println("Added a new user to MongoDB")
        fmt.Println("MongoDB ID ")
        fmt.Println(res.InsertedID)
    }
}
func AcceptedRules(w http.ResponseWriter, request *http.Request) {
    vars := mux.Vars(request)
    fmt.Println(vars["secret"][0:378]);
    filter := bson.D{{"secret", vars["secret"][0:378]}}
    update := bson.D{
        {"$set", bson.D{
            {"level", 0},
        }},
    }
    _, err := collection.UpdateOne(context.TODO(), filter, update)
    if err != nil {
        log.Fatal(err)
    }
}

func AnswerQuestion(w http.ResponseWriter, request *http.Request) {
    vars := mux.Vars(request)
    find, _ := collection.Find(context.TODO(), bson.M{"secret": vars["secret"][0:378]})
    JSOND, _ := json.Marshal(find.Next(context.TODO()))
    if strings.Compare(string(JSOND), "true") == 0 {
        if val, ok := answers[vars["level"]]; ok {
            var current DatabaseUserObject
            err := find.Decode(&current)
            if (err != nil) {
                fmt.Println("Error decoding database object ", err)
            }
            if strings.Compare(strconv.Itoa(current.Level), vars["level"]) == 0 {
                if strings.Compare(val, vars["answer"]) == 0 {
                    filter := bson.D{{"secret", vars["secret"][0:378]}}
                    update := bson.D{
                        {"$inc", bson.D {
                            {"level", 1},
                        }},
                    }
                    _, err := collection.UpdateOne(context.TODO(), filter, update)
                    update = bson.D{
                        {"$set", bson.D {
                            {"lastModified", time.Now().UTC()},
                        }},
                    }
                    _, err = collection.UpdateOne(context.TODO(), filter, update)
                    if err != nil {
                        fmt.Println("Error updating ", err)
                        responseJSON("DatabaseError", w, http.StatusInternalServerError)
                    } else {
                        responseJSON("Correct", w, http.StatusOK)
                    }
                } else {
                    responseJSON("Wrong", w, http.StatusOK)
                }
            } else {
                responseJSON("LevelNoMatch", w, http.StatusOK)
            }
        } else {
            responseJSON("InvalidLevel", w, http.StatusOK)
        }
    } else {
        responseJSON("InvalidToken", w, http.StatusOK)
    }
}

func LeaderboardHandler (w http.ResponseWriter, request *http.Request) {
    options := options.Find()
    options.SetSort(bson.D{{"level", -1}, {"lastModified", 1}})    
    find, _ := collection.Find(context.TODO(), bson.M{}, options)
    var results []user
    for find.Next(context.TODO()) {
        // create a value into which the single document can be decoded
        var elem user
        err := find.Decode(&elem)
        fmt.Println(elem)
        if err != nil {
            fmt.Println("Error decoding leaderboard item")
        }
        results = append(results, elem)
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    jData, _ := json.Marshal(results)
    w.Write(jData)
}

func LevelHandler (w http.ResponseWriter, request *http.Request) {
    fmt.Println("Here")
    vars := mux.Vars(request)
    find, _ := collection.Find(context.TODO(), bson.M{"secret": vars["secret"][0:378]})
    JSOND, _ := json.Marshal(find.Next(context.TODO()))
    if strings.Compare(string(JSOND), "true") == 0 {
        var current DatabaseUserObject
        err := find.Decode(&current)
        if (err != nil) {
            fmt.Println("Not able to read database object")
            responseJSON("DatabaseError", w, http.StatusInternalServerError)
        } else {
            var resp LevelResponse
            if (current.Level == 0) {
                resp = LevelResponse{0, "https://res.cloudinary.com/drgddftct/image/upload/v1547292346/QPADBgJd8EkeBut6.png"}
            } else if (current.Level == 1) {
                resp = LevelResponse{1, "https://res.cloudinary.com/dmridruee/image/upload/v1547295044/qsQK5bRhRvgXjh378d5J/7yXw9wkWaTMXafsC7USs.png"}
            } else if (current.Level == 2) {
                resp = LevelResponse{2, "169B62169B62169B62FFFFFFFFFFFFFFFFFFFF883EFF883EFF883E169B62169B62169B62FFFFFFFFFFFFFFFFFFFF883EFF883EFF883E169B62169B62169B62FFFFFFFFFFFFFFFFFFFF883EFF883EFF883E169B62169B62169B62FFFFFFFFFFFFFFFFFFFF883EFF883EFF883E169B62169B62169B62FFFFFFFFFFFFFFFFFFFF883EFF883EFF883E169B62169B62169B62FFFFFFFFFFFFFFFFFFFF883EFF883EFF883E169B62169B62169B62FFFFFFFFFFFFFFFFFFFF883EFF883EFF883E169B62169B62169B62FFFFFFFFFFFFFFFFFFFF883EFF883EFF883E169B62169B62169B62FFFFFFFFFFFFFFFFFFFF883EFF883EFF883E169B62169B62169B62FFFFFFFFFFFFFFFFFFFF883EFF883EFF883E169B62169B62169B62FFFFFFFFFFFFFFFFFFFF883EFF883EFF883E169B62169B62169B62FFFFFFFFFFFFFFFFFFFF883EFF883EFF883E169B62169B62169B62FFFFFFFFFFFFFFFFFFFF883EFF883EFF883E"}
            } else if (current.Level == 3) {
                resp = LevelResponse{3, "/midi.mid"}
            } else if (current.Level == 4) {
                resp = LevelResponse{4, "https://res.cloudinary.com/do3uy82tk/image/upload/v1564096693/asdfasdf.jpg"}
            } else if (current.Level == 5) {
                resp = LevelResponse{5, "https://res.cloudinary.com/dmridruee/image/upload/v1547211291/0PNQNGAOck2NQwyb6hQV.png"}
            } else {
                resp = LevelResponse{6, "Won"}
            }
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusOK)
            jData, _ := json.Marshal(resp)
            w.Write(jData)
        }
    }
}

func LeaderboardTableHandler(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "leaderboard.html")
}

func CSSHandler(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "./prerenderedviews/css/index.css")
}

func RulesHandler(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "rules.html")
}

func MIDIHandler(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "cryptex.mid")
}

// func submitAnswer(w http.ResponseWriter, request *http.Request) {
//     vars := mux.Vars(request)
//     client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
//     context.TODO(), _ := context.WithTimeout(context.Background(), 10*time.Second)
//     _ = client.Connect(context.TODO())
//     collection := client.Database("Cryptex").Collection("users")
//     context.TODO(), _ = context.WithTimeout(context.Background(), 5*time.Second)
    
// }
// Function is obsoelete, implemented using GraphQL in main()
// func RetrieveLevel(w http.ResponseWriter, request *http.Request) {
//     vars := mux.Vars(request)
//     client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
//     context.TODO(), _ := context.WithTimeout(context.Background(), 10*time.Second)
//     err = client.Connect(context.TODO())
//     fmt.Println(err)
//     collection := client.Database("Cryptex").Collection("users")
//     context.TODO(), _ = context.WithTimeout(context.Background(), 5*time.Second)
//     filter := bson.M{"clientID" : vars["ID"]}
//     var result map[string]interface{}
//     err = collection.FindOne(context.TODO(), filter).Decode(&result)

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