package controller

import (
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
	a.Models.Samba_Shares.DeleteShare(models.SambaShareResponse{
		Email:   r.Context().Value("Authorization").(string),
		Shareid: pa.ByName("shareid"),
	})

}

func (a AppRouter) CreateShare(w http.ResponseWriter, r *http.Request, pa httprouter.Params) {
	//It will want Space ID Instead
	a.Models.Samba_Shares.AddShare(models.SambaShareResponse{})
}

func (a AppRouter) InviteUsers(w http.ResponseWriter, r *http.Request, pa httprouter.Params) {

}

func (a AppRouter) AcceptInvite(w http.ResponseWriter, r *http.Request, pa httprouter.Params) {

}
