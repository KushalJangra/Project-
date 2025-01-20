package main

import (
	"fmt"
    "encoding/json"
    "net/http"
    "github.com/gorilla/mux"
    
)



type Order struct {
    ID    uint   `json:"id"`
	Quantity int `json:"quantity"`
    
}
const (
    DBHost  = "127.0.0.1"
    DBUser  = "root"
    DBPass  = "Kush@123456"
    DBDbase = "sdb"
    PORT    = ":3000"
)
func initDB() {
    
	dbConn := fmt.Sprintf("%s:%s@tcp(%s)/%s", DBUser, DBPass, DBHost, DBDbase)
    db, err := sql.Open("mysql", dbConn)
    if err != nil {
        log.Fatalf("Database connection error: %v", err)
    }

    database = db
 
    if err := database.Ping(); err != nil {
        log.Fatalf("Database ping error: %v", err)
    }
    log.Println("Database connected successfully!")
	
}

func createOrder(w http.ResponseWriter, r *http.Request) {
    var order Order
    if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Validate customer
    exists, err := validateCustomer(order.ID)
    if err != nil {
        http.Error(w, "Error validating customer", http.StatusInternalServerError)
        log.Println("Error validating customer:", err)
        return
    }

    if !exists {
        http.Error(w, "Customer not found", http.StatusNotFound)
        return
    }

    // Create order
    if result := db.Create(&order); result.Error != nil {
        http.Error(w, result.Error.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(order)
}

func getOrder(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]
    var order Order
    if result := db.First(&order, id); result.Error != nil {
        http.Error(w, "order not found", http.StatusNotFound)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(order)
}

func main() {
    initDB()

    r := mux.NewRouter()

    r.HandleFunc("/orders", createOrder).Methods("POST")
    r.HandleFunc("/orders/{id}", getOrder).Methods("GET")

    http.ListenAndServe(":8080", r)
}