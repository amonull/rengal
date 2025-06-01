package anilist

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/samber/lo"

	"github.com/amonull/rengal/log"
	"github.com/amonull/rengal/network"
)

// login to Anilist
func (a *Anilist) login() error {
	log.Info("Logging in to Anilist")

	if a.id() == "" {
		e := errors.New("no ID set")
		log.Error(e)
		return e
	}
	if a.secret() == "" {
		e := errors.New("no secret set")
		log.Error(e)
		return e
	}
	if a.code() == "" {
		e := errors.New("no code set")
		log.Error(e)
		return e
	}

	// anilist body for login
	body := map[string]any{
		"client_id":     a.id(),
		"client_secret": a.secret(),
		"grant_type":    "authorization_code",
		"redirect_uri":  "https://anilist.co/api/v2/oauth/pin",
		"code":          a.code(),
	}

	// encode body
	jsonBody := lo.Must(json.Marshal(&body))

	// create request
	log.Info("Sending login request to Anilist")
	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		"https://anilist.co/api/v2/oauth/token",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		log.Error(err)
		return err
	}

	// set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// send request
	resp, err := network.Client.Do(req)
	if err != nil {
		log.Error(err)
		return err
	}

	defer func() {
		// naked return used to return error from defer
		err = resp.Body.Close()
	}()

	// check response code
	if resp.StatusCode != http.StatusOK {
		log.Info("Request failed with status code: " + strconv.Itoa(resp.StatusCode))
		return fmt.Errorf("invalid response code %d", resp.StatusCode)
	}

	// decode response
	log.Info("Decoding response from Anilist")
	var response struct {
		AccessToken string `json:"access_token"`
	}

	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Error(err)
		return err
	}

	// set token
	log.Info("Logged in Anilist")
	a.token = response.AccessToken

	return nil
}
