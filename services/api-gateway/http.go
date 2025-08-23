package main

import (
	"encoding/json"
	"log"
	"net/http"
	"ride-sharing/services/api-gateway/grpc_clients"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/tracing"
)

var tracer = tracing.GetTracer("api-gateway")

func handleTripStart(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "handleTripStart")
	defer span.End()

	var reqBody startTripRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "failed to parse JSON data", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	tripService, err := grpc_clients.NewTripServiceClient()
	if err != nil {
		log.Fatal(err)
	}
	defer tripService.Close()

	trip, err := tripService.Client.CreateTrip(ctx, reqBody.toProto())

	if err != nil {
		log.Printf("Failed to start a trip: %v", err)
		http.Error(w, "failed to start trip", http.StatusInternalServerError)
		return
	}

	response := contracts.APIResponse{
		Data: trip,
	}

	writeJSON(w, http.StatusCreated, response)
}

func handleTripPreview(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "handleTripPreview")
	defer span.End()

	var reqBody previewTripRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "failed to parse JSON data", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if reqBody.UserID == "" {
		http.Error(w, "user ID is required", http.StatusBadRequest)
		return
	}

	tripService, err := grpc_clients.NewTripServiceClient()
	if err != nil {
		log.Fatal(err)
	}
	defer tripService.Close()

	tripPreview, err := tripService.Client.PreviewTrip(ctx, reqBody.toProto())
	if err != nil {
		log.Printf("Failed to preview a trip: %v", err)
		http.Error(w, "failed to preview trip", http.StatusInternalServerError)
		return
	}

	response := contracts.APIResponse{
		Data: tripPreview,
	}

	writeJSON(w, http.StatusCreated, response)
}
