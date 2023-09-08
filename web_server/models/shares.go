package models

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"time"
)

type SambaShareRequest struct {
	Email    string
	Spaceid  string
	Password string
}

type SambaShareResponse struct {
	Email   string
	Shareid string
}

type ShareInviteRequest struct {
	Email   string
	Shareid string
}

type ShareInviteResponse struct {
	Email       string
	Inviteid    string
	Invite_code string
}

type ShareInviteAccept struct {
	Email       string
	Inviteid    string
	Invite_code string
	Password    string
}

/**
 *
 * Later Get Shares By User
 *
 */

type SambaShareModel interface {
	AddShare(SambaShareResponse) (*SambaShareResponse, error)
	DeleteShare(SambaShareResponse) (*SambaShareResponse, error)
	CreateInvite(ShareInviteRequest) (*ShareInviteResponse, error)
	AcceptInvite(ShareInviteAccept) (*ShareInviteResponse, error)
}

func GenInvite() (string, []byte, time.Time) {

	var TokenBytes []byte = make([]byte, 32)

	var TokenString string = base64.StdEncoding.EncodeToString(TokenBytes)

	hashStore := sha256.New()

	hashStore.Write(TokenBytes)

	hashedToken := hashStore.Sum(nil)

	//var hashedToken []byte = hashStore

	return TokenString, hashedToken, time.Now().Add(time.Hour * 24)
}

func VerifyInvite(string_header string, hash_store []byte, expires time.Time) (bool, error) {

	if time.Now().After(expires) {
		return false, nil
	}

	initialTokenBytes, err := base64.StdEncoding.DecodeString(string_header)
	if err != nil {
		return false, ErrorMalformedInvite
	}

	hashedStore := sha256.New()
	hashedStore.Write(initialTokenBytes)

	hashedChallengeToken := hashedStore.Sum(nil)

	if bytes.Compare(hashedChallengeToken, hash_store) != 0 {
		return false, nil
	}

	return true, nil
}
