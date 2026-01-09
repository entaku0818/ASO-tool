package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	sqladmin "google.golang.org/api/sqladmin/v1"
)

var (
	projectID    = os.Getenv("GCP_PROJECT_ID")
	instanceName = "aso-db"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func main() {
	if projectID == "" {
		projectID = "voicevox-custom"
	}

	http.HandleFunc("/", handleRequest)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("DB Manager starting on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	// CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	action := r.URL.Query().Get("action")
	if action == "" {
		action = "status"
	}

	ctx := context.Background()
	service, err := sqladmin.NewService(ctx)
	if err != nil {
		respondError(w, "Failed to create SQL Admin client: "+err.Error())
		return
	}

	switch action {
	case "start":
		startInstance(w, service)
	case "stop":
		stopInstance(w, service)
	case "status":
		getStatus(w, service)
	default:
		respondError(w, "Invalid action. Use: start, stop, or status")
	}
}

func startInstance(w http.ResponseWriter, service *sqladmin.Service) {
	instance, err := service.Instances.Get(projectID, instanceName).Do()
	if err != nil {
		respondError(w, "Failed to get instance: "+err.Error())
		return
	}

	if instance.Settings.ActivationPolicy == "ALWAYS" {
		respond(w, Response{Status: "already_running", Message: "Database is already running"})
		return
	}

	instance.Settings.ActivationPolicy = "ALWAYS"
	_, err = service.Instances.Update(projectID, instanceName, instance).Do()
	if err != nil {
		respondError(w, "Failed to start instance: "+err.Error())
		return
	}

	respond(w, Response{Status: "starting", Message: "Database start initiated"})
}

func stopInstance(w http.ResponseWriter, service *sqladmin.Service) {
	instance, err := service.Instances.Get(projectID, instanceName).Do()
	if err != nil {
		respondError(w, "Failed to get instance: "+err.Error())
		return
	}

	if instance.Settings.ActivationPolicy == "NEVER" {
		respond(w, Response{Status: "already_stopped", Message: "Database is already stopped"})
		return
	}

	instance.Settings.ActivationPolicy = "NEVER"
	_, err = service.Instances.Update(projectID, instanceName, instance).Do()
	if err != nil {
		respondError(w, "Failed to stop instance: "+err.Error())
		return
	}

	respond(w, Response{Status: "stopping", Message: "Database stop initiated"})
}

func getStatus(w http.ResponseWriter, service *sqladmin.Service) {
	instance, err := service.Instances.Get(projectID, instanceName).Do()
	if err != nil {
		respondError(w, "Failed to get instance: "+err.Error())
		return
	}

	status := "stopped"
	if instance.Settings.ActivationPolicy == "ALWAYS" {
		status = "running"
	}

	respond(w, Response{Status: status, Message: "Instance state: " + instance.State})
}

func respond(w http.ResponseWriter, resp Response) {
	json.NewEncoder(w).Encode(resp)
}

func respondError(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(Response{Status: "error", Message: message})
}
