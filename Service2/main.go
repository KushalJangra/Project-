package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "math/rand"
    "net/http"
    "strconv"
    "time"

    _ "github.com/go-sql-driver/mysql"
    "github.com/gorilla/mux"
)

const (
    DBHost  = "127.0.0.1"
    DBUser  = "root"
    DBPass  = "Kush@123456"
    DBDbase = "sdb"
    PORT    = ":3000"
)

var database *sql.DB

type Customer struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

type Product struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}
type JSONResponse struct {
    Fields map[string]string `json:"fields"`
}

func initDB() {
    dbConn := fmt.Sprintf("%s:%s@tcp(%s)/%s", DBUser, DBPass, DBHost, DBDbase)
    db, err := sql.Open("mysql", dbConn)
    if err != nil {
        log.Fatalf("Database connection error: %v", err)
    }

    database = db
    // Test the connection
    if err := database.Ping(); err != nil {
        log.Fatalf("Database ping error: %v", err)
    }
    log.Println("Database connected successfully!")
}

// APICommentPost handles the POST request for adding comments
func APICustPost(w http.ResponseWriter, r *http.Request) {
    var commentAdded bool
    resp := JSONResponse{Fields: make(map[string]string)}

    err := r.ParseForm()
    if err != nil {
        http.Error(w, "Unable to parse form", http.StatusBadRequest)
        return
    }

    name := r.FormValue("name")
    email := r.FormValue("email")
    comments := r.FormValue("comments")
    guid := r.FormValue("guid")

    rand.Seed(time.Now().UnixNano())
    randomInt := rand.Intn(1001)

    res, err := database.Exec("INSERT INTO customer (page_id, comment_name, comment_email, comment_text, comment_guid) VALUES (?, ?, ?, ?, ?)", randomInt, name, email, comments, guid)
    if err != nil {
        http.Error(w, "Failed to add comment", http.StatusInternalServerError)
        log.Println(err)
        return
    }

    id, err := res.LastInsertId()
    if err != nil {
        commentAdded = false
        log.Println(err)
    } else {
        commentAdded = true
    }

    // Prepare JSON response
    resp.Fields["id"] = strconv.FormatInt(id, 10)
    resp.Fields["added"] = strconv.FormatBool(commentAdded)

    w.Header().Set("Content-Type", "application/json")
    jsonResp, _ := json.Marshal(resp)
    w.WriteHeader(http.StatusOK)
    w.Write(jsonResp)
}

// cust get
func APIcustget(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    cID := vars["id"]

    cus := Customer{}
    err := database.QueryRow("SELECT page_id, comment_name, comment_email, comment_text, comment_guid FROM customer WHERE page_id=?", cID).
        Scan(&cus.ID, &cus.Name, &cus.Email)
    if err != nil {
        http.Error(w, " not found", http.StatusNotFound)
        log.Println("not found", cID)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(cus)
}

func APIprod(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    pID := vars["id"]

    prods := Product{}
    err := database.QueryRow("SELECT p_id, o_name FROM prod WHERE p_id=?", pID).
        Scan(&prods.ID, &prods.Name)
    if err != nil {
        http.Error(w, "Order not found", http.StatusNotFound)
        log.Println("Couldn't find order with ID:", pID)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(prods)
}

func APIprodPost(w http.ResponseWriter, r *http.Request) {
    var prodadded bool

    resp := JSONResponse{Fields: make(map[string]string)}

    err := r.ParseForm()
    if err != nil {
        http.Error(w, "Unable to parse form", http.StatusBadRequest)
        return
    }

    name := r.FormValue("name")

    rand.Seed(time.Now().UnixNano())
    randomInt := rand.Intn(1001)

    res, err := database.Exec("INSERT INTO prod (p_id, p_name) VALUES ( ?, ?)", randomInt, name)
    if err != nil {
        http.Error(w, "Failed to add order", http.StatusInternalServerError)
        log.Println(err)
        return
    }

    id, err := res.LastInsertId()
    if err != nil {
        prodadded = false
        log.Println(err)
    } else {
        prodadded = true
    }

    resp.Fields["id"] = strconv.FormatInt(id, 10)
    resp.Fields["added"] = strconv.FormatBool(prodadded)

    w.Header().Set("Content-Type", "application/json")
    jsonResp, _ := json.Marshal(resp)
    w.WriteHeader(http.StatusOK)
    w.Write(jsonResp)
}

func main() {
    
    initDB()
    defer database.Close()

    
    r := mux.NewRouter()

    r.HandleFunc("/api/cust", APICustPost).Methods("POST")
    r.HandleFunc("/api/cust/{id:[0-9]+}", APIcustget).Methods("GET").Schemes("http")

    r.HandleFunc("/api/prod/{id:[0-9]+}", APIprod).Methods("GET").Schemes("http")
    r.HandleFunc("/api/prod", APIprodPost).Methods("POST")

    log.Printf("Server running on http://localhost%s", PORT)
    log.Fatal(http.ListenAndServe(PORT, r))
}