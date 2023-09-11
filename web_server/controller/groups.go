package controller

import (
	"encoding/json"
	"net/http"

	"github.com/PranoSA/samba_share_backend/web_server/models"
	"github.com/julienschmidt/httprouter"
)

/**
 *
 * HTTP ROUTES FOR INVITING, ACCEPTING INVITES, CREATING SAMBA GROUPS
 * AND MANAGING SHARES ASSOCIATED WITH EACH GROUP
 *
 */

func (a AppRouter) DeleteShare(w http.ResponseWriter, r *http.Request, pa httprouter.Params) {

	/**
	 *
	 * Get Share
	 * The Problem With Atomic Operations -> Maybe Not An Issue For Now ...
	 *
	 * Ensure the Correct Owner
	 *
	 * Delete If Correct Owner, Also Delete All Samba Shares
	 *
	 *
	 */
	email := r.Context().Value("Authorization")
	if email == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	shareid := pa.ByName("shareid")

	a.Models.Samba_Shares.DeleteShare(models.SambaShareResponse{
		Email:   email.(string),
		Shareid: shareid,
	})

}

type CreateShareBody struct {
	Name string `json:"Body"` //Ignored For Now ....
}

func (a AppRouter) CreateShare(w http.ResponseWriter, r *http.Request, pa httprouter.Params) {

	email := r.Context().Value("Authorization")
	if email == nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	//It will want Space ID Instead
	res, err := a.Models.Samba_Shares.AddShare(models.SambaShareResponse{
		Email:   email.(string),
		Shareid: "",
	})

	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	json.NewEncoder(w).Encode(res)
}

type InviteResponse struct {
	InviteId   string
	InviteLink string
}

func (a AppRouter) InviteUsers(w http.ResponseWriter, r *http.Request, pa httprouter.Params) {

	//a.CorsMiddleware(&w, r)
	shareid := pa.ByName("shareid")

	email := r.Context().Value("Authorization")
	if email == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	res, err := a.Models.Samba_Shares.CreateInvite(models.ShareInviteRequest{
		Email:   email.(string),
		Shareid: shareid,
	})

	if err == models.ErrorEntryDoesNotExist {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var Body InviteResponse
	Body.InviteId = res.Inviteid
	Body.InviteLink = res.Invite_code

	json.NewEncoder(w).Encode(Body)

}

func (a AppRouter) AcceptInvite(w http.ResponseWriter, r *http.Request, pa httprouter.Params) {

}
