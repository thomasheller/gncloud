package gncloud

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/thomasheller/ghttp"
)

type Nextcloud struct {
	BaseURL string
	User    string
	Token   string
}

func (n Nextcloud) Delete(path string, v interface{}) error {
	return n.getOrDelete(http.MethodDelete, path, v)
}

func (n Nextcloud) Get(path string, v interface{}) error {
	return n.getOrDelete(http.MethodGet, path, v)
}

func (n Nextcloud) Post(path string, data url.Values, v interface{}) error {
	return n.postOrPut(http.MethodPost, path, data, v)
}

func (n Nextcloud) Put(path string, data url.Values, v interface{}) error {
	return n.postOrPut(http.MethodPut, path, data, v)
}

func (n Nextcloud) getOrDelete(method string, path string, v interface{}) error {
	req, err := n.buildRequest(method, path)
	if err != nil {
		return err
	}

	var r Response
	err = ghttp.JSON(req, &r)
	if err != nil {
		return err
	}

	return n.decodeResponse(r, v)
}

func (n Nextcloud) postOrPut(method string, path string, data url.Values, v interface{}) error {
	req, err := n.buildRequest(method, path)
	if err != nil {
		return err
	}

	var r Response
	err = ghttp.FormJSON(req, data, &r)
	if err != nil {
		return err
	}

	return n.decodeResponse(r, v)
}

func (n Nextcloud) buildRequest(method string, path string) (*http.Request, error) {
	req, err := http.NewRequest(method, n.url(path), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("OCS-APIRequest", "true")
	req.SetBasicAuth(n.User, n.Token)

	return req, nil
}

func (n Nextcloud) decodeResponse(r Response, v interface{}) error {
	if !ghttp.Success(r.Ocs.Meta.Statuscode) {
		return fmt.Errorf("Error from Nextcloud API (%d %s): %s", r.Ocs.Meta.Statuscode, r.Ocs.Meta.Status, r.Ocs.Meta.Message)
	}

	if v == nil {
		return nil
	}

	return json.Unmarshal(r.Ocs.Data, v)
}

func (n Nextcloud) url(path string) string {
	return fmt.Sprintf("%s/ocs/v2.php/%s?format=json", n.BaseURL, path)
}
