package go_yapi

import (
	"encoding/gob"
	"encoding/json"
	"errors"
	"golang.org/x/oauth2"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const VersionAPI = "v6"

// http://www.patorjk.com/software/taag/
//      ___ ___         .__
//     /   |   \   ____ |  | ______   ___________  ______
//    /    ~    \_/ __ \|  | \____ \_/ __ \_  __ \/  ___/
//    \    Y    /\  ___/|  |_|  |_> >  ___/|  | \/\___ \
//     \___|_  /  \___  >____/   __/ \___  >__|  /____  >
//           \/       \/     |__|        \/           \/

// NewOauth2Config ...
func NewOauth2Config(clientID, clientSecret string, scopes []string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://oauth.yandex.ru/authorize",
			TokenURL: "https://oauth.yandex.ru/token",
		},
	}
}

type Parameters map[string][]string

func (p Parameters) String() string {
	if len(p) > 0 {
		pURL := &strings.Builder{}
		first := true
		for k := range p {
			if first {
				pURL.WriteString("?")
				first = false
			} else {
				pURL.WriteString("&")
			}
			d := make([]string, 0, len(p[k]))
			for i := range p[k] {
				d = append(d, url.QueryEscape(p[k][i]))
			}
			pURL.WriteString(k + "=" + strings.Join(d, ","))
		}
		return pURL.String()
	}
	return ""
}

// Get ...
func Get(client *http.Client, url string, params Parameters, header map[string]string, v interface{}) error {
	return Request(client, http.MethodGet, url, params, header, http.StatusOK, nil, v)
}

// Post ...
func Post(client *http.Client, url string, params Parameters, header map[string]string, body io.Reader, v interface{}) error {
	return Request(client, http.MethodPost, url, params, header, http.StatusCreated, body, v)
}

// Patch ...
func Patch(client *http.Client, url string, params Parameters, header map[string]string, body io.Reader, v interface{}) error {
	return Request(client, http.MethodPatch, url, params, header, http.StatusOK, body, v)
}

// Delete ...
func Delete(client *http.Client, url string, params Parameters, header map[string]string) error {
	return Request(client, http.MethodDelete, url, params, header, http.StatusOK, nil, nil)
}

func Request(client *http.Client, method, url string, params Parameters, header map[string]string, expectedStatus int, body io.Reader, v interface{}) error {
	req, err := http.NewRequest(method, url+params.String(), body)
	if err != nil {
		return err
	}
	for k := range header {
		req.Header.Add(k, header[k])
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden {
		return errors.New(resp.Status + " " + resp.Header.Get("WWW-Authenticate"))
	}

	if resp.StatusCode != expectedStatus {
		return errors.New(resp.Status)
	}

	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(v)
		if err != nil {
			return err
		}
	}
	return nil
}

func TokenFromFile(tokenFile string) (*oauth2.Token, error) {
	tokData, err := os.Open(tokenFile)
	if err != nil {
		return &oauth2.Token{}, err
	}
	defer tokData.Close()
	var t oauth2.Token
	if err := gob.NewDecoder(tokData).Decode(&t); err != nil {
		return &oauth2.Token{}, err
	}
	return &t, nil
}

func TokenToFile(tokenFile string, token *oauth2.Token) error {
	f, err := os.Create(tokenFile)
	if err != nil {
		return err
	}
	defer f.Close()
	return gob.NewEncoder(f).Encode(*token)
}

func headerOrgID(id int) map[string]string {
	if id > 0 {
		return map[string]string{"X-Org-ID": strconv.Itoa(id)}
	}
	return map[string]string{}
}

func jsonParam(v interface{}) string {
	s, _ := json.Marshal(v)
	return string(s)
}
