package handlers

import (
	"fmt"
	"net/http"
	"os"
	"tectonic-api/database"
	"tectonic-api/utils"
)

func UpdatePoints(w http.ResponseWriter, r *http.Request) {
	status := http.StatusNoContent

	p, err := utils.ParseParametersURL(r, "guild_id", "user_ids")
	if err != nil {
		status = http.StatusBadRequest
		utils.JsonWriter(err).IntoHTTP(status)(w, r)
		return
	}

	// Point parameter key, points is for custom point values
	pkey, err := utils.RequireOne(p, "point_event", "points")
	if err != nil {
		status = http.StatusBadRequest
		utils.JsonWriter(err).IntoHTTP(status)(w, r)
		return
	}

	err = database.UpdatePoints(r.Context(), p, pkey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error updating points: %v\n", err)
		status = http.StatusNotFound
	}

	utils.JsonWriter(http.NoBody).IntoHTTP(status)(w, r)
}
