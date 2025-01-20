package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strconv"

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

// validateCustomer checks if a customer exists in the database
func validateCustomer(customerID int) (bool, error) {
    var exists bool
    query := "SELECT EXISTS(SELECT 1 FROM customer WHERE id=?)"
    err := database.QueryRow(query, customerID).Scan(&exists)
    if err != nil {
        return false, err
    }
    return exists, nil
}

// APIcustget handles the GET request for fetching a customer
func APIcustget(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    cID, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid customer ID", http.StatusBadRequest)
        return
    }

    exists, err := validateCustomer(cID)
    if err != nil {
        http.Error(w, "Error validating customer", http.StatusInternalServerError)
        log.Println("Error validating customer:", err)
        return
    }

    if !exists {
        http.Error(w, "Customer not found", http.StatusNotFound)
        return
    }

    cus := Customer{}
    err = database.QueryRow("SELECT id, name, email FROM customer WHERE id=?", cID).
        Scan(&cus.ID, &cus.Name, &cus.Email)
    if err != nil {
        http.Error(w, "Customer not found", http.StatusNotFound)
        log.Println("Customer not found:", cID)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(cus)
}

func main() {
    initDB()
    defer database.Close()

    r := mux.NewRouter()

    r.HandleFunc("/api/cust/{id:[0-9]+}", APIcustget).Methods("GET").Schemes("http")

    log.Printf("Server running on http://localhost%s", PORT)
    log.Fatal(http.ListenAndServe(PORT, r))
}