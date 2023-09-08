package models_test

import (
	"testing"
	"time"

	"github.com/PranoSA/samba_share_backend/web_server/models"
)

func TestInviteFunctionality(t *testing.T) {

	t.Run("Test Interperobability", func(t *testing.T) {

		header, token, expir := models.GenInvite()

		worked, err := models.VerifyInvite(header, token, expir)

		if err != nil {
			t.Error("Error In Verifying Invite")
		}

		if !worked {
			t.Error("Failed verifiying Invite ")
		}
	})

	t.Run("Deliberately Wrong", func(t *testing.T) {
		header, token, expir := models.GenInvite()

		token[25] = 172
		token[26] = 114
		token[1] = 7

		worked, err := models.VerifyInvite(header, token, expir)

		if err != nil {
			t.Error("Generated Unexpected Error")
		}

		if worked {
			t.Error("Unexpectedly Verified Invite")
		}
	})

	t.Run("Deliberately Expired", func(t *testing.T) {
		header, token, expir := models.GenInvite()

		new_expir := expir.Add(-time.Hour * 48)

		worked, err := models.VerifyInvite(header, token, new_expir)

		if err != nil {
			t.Error("Unexpectedly Generated Error")
		}

		if worked {
			t.Error("Expected Invite Expiration But WOrked")
		}

	})

}
