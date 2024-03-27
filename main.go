package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type HealthCheckResponse struct {
  Status string `json:"status"`
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
  response := HealthCheckResponse{Status: "ok"}

  // This will pipe the JSON byte directly to the ResponseWriter, significantly optimized for performance
  json.NewEncoder(w).Encode(&response)

  
  // This method will create a JSON and marshal it to a byte array before writing it to the ResponseWriter
  // Not a good practice for performance in comparison to the above method
  // data, err := json.Marshal(response)
  // if err != nil {
  //   http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
  //   return
  // }
  // w.Header().Set("Content-Type", "application/json")
  // w.WriteHeader(http.StatusOK)
  // w.Write([]byte(data))
}

func main() {
  httpMux := http.NewServeMux()

  httpMux.HandleFunc("/health", HealthCheckHandler)

  server := &http.Server{
    Addr: ":8080",
    Handler: httpMux,
  }


  fmt.Printf("Server is running on port %s\n", server.Addr)
  log.Fatal(server.ListenAndServe())
}
