package util

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	duo "github.com/duosecurity/duo_api_golang"
	"github.com/duosecurity/duo_api_golang/authapi"
)

type duoResponse struct {
	Stat     string
	Response struct {
		Result    string
		Status    string
		StatusMsg string
	}
}

// DuoPush sends a Duo push to the committer
func DuoPush(committer string, committerEmail string,
	group string, repo string) error {

	ikey, err := GetSecret("duoIntegrationKey")
	if err != nil {
		return err
	}

	skey, err := GetSecret("duoSecretKey")
	if err != nil {
		return err
	}

	host, err := GetSecret("duoHost")
	if err != nil {
		return err
	}

	duoClient := duo.NewDuoApi(
		ikey,
		skey,
		host,
		"",
	)

	d := authapi.NewAuthApi(*duoClient)
	check, err := d.Check()
	if err != nil {
		return err
	}
	if check.StatResult.Stat != "OK" {
		return fmt.Errorf("Could not connect to Duo: %s (%s)", *check.StatResult.Message, *check.StatResult.Message_Detail)
	}
	params, err := url.ParseQuery("")
	if err != nil {
		return err
	}

	username := strings.Split(committerEmail, "@")[0]
	params.Add("username", username)
	params.Add("factor", "push")
	params.Add("device", "auto")
	params.Add("type", "Deploy")
	params.Add("display_username", committer)
	params.Add("pushinfo", fmt.Sprintf("Deploy=%s/%s", group, repo))

	_, body, err := d.SignedCall("POST", "/auth/v2/auth", params)
	if err != nil {
		return err
	}

	var duoResp duoResponse
	err = json.Unmarshal(body, &duoResp)
	if err != nil {
		return err
	}

	if duoResp.Response.Result == "allow" {
		return nil
	}

	return fmt.Errorf("Duo 2FA responed with '%s'", string(body))
}
