package handlers

import (
	"log"
	"net/http"

	"github.com/kreilt/monster-tracker/internal/model"
	"github.com/kreilt/monster-tracker/internal/repository"
)

func Flavors(flavorsRepo *repository.Flavor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		search := r.URL.Query().Get("search")
		lineup := r.URL.Query().Get("lineup")
		rare := r.URL.Query().Get("rare")
		region := r.URL.Query().Get("region")
		status := r.URL.Query().Get("status")

		if rare != "" && !model.IsValidRarity(rare) {
			writeError(w, "invalid rare value", http.StatusBadRequest)
			return
		}

		if status != "" && !model.IsValidStatuses(status) {
			writeError(w, "invalid status value", http.StatusBadRequest)
			return
		}

		flavorsFilter := repository.FlavorFilter{
			Search: search,
			Lineup: lineup,
			Rare:   rare,
			Region: region,
			Status: status,
		}

		flavors, err := flavorsRepo.List(r.Context(), flavorsFilter)
		if err != nil {
			log.Printf("failed to get flavors, %v", err)
			writeError(w, "internal", http.StatusInternalServerError)
			return
		}
		writeJSON(w, flavors, http.StatusOK)
	}
}
