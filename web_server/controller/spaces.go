package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/PranoSA/samba_share_backend/web_server/models"
	"github.com/julienschmidt/httprouter"
)

func (ar AppRouter) CreateSpace(w http.ResponseWriter, r *http.Request, pa httprouter.Params) {

	w.Header().Add("Content-Type", "application/json")

	user, ok := r.Context().Value("Authorization").(string)

	if !ok {
		log.Fatal("Invalid Type Casting at Spaces.go /Create Space")
	}

	megabytes, e := strconv.Atoi(r.URL.Query().Get("megabytes"))

	if e != nil {
		json.NewEncoder(w)

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if megabytes < 10 || megabytes > 100_000 {

		var ErrorResponse BadRequestResponse = BadRequestResponse{
			{
				ParameterName: "megabytes",
				Param_Type:    "Query",
				Value_Type:    "int64",
				Message:       "Values Between 10 and 100,000 needed ",
			},
		}

		json.NewEncoder(w).Encode(&ErrorResponse)
		return
	}

	res, err := ar.Models.Spaces.CreateSpace(models.SpaceRequest{
		Owner:     user,
		Megabytes: int64(megabytes),
	})

	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
	}

	json.NewEncoder(w).Encode(&res)

}

func (ar AppRouter) DeleteSpace(w http.ResponseWriter, r *http.Request, ap httprouter.Params) {
	spaceid := ap.ByName("spaceid")
	email := r.Context().Value("Authorization")
	if email == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	res, err := ar.Models.Spaces.GetSpaceById(models.DeleteSpaceRequest{
		Owner:    email.(string),
		Space_id: spaceid,
	})
	if err == models.ErrorEntryDoesNotExist {
		w.WriteHeader(http.StatusForbidden)
	}

	json.NewEncoder(w).Encode(res)
}
