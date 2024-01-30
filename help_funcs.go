package itswizard_m_berlinPreVersion

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func createUrl(query map[string]string, method string, dev bool, endpoint string) (*url.URL, error) {
	u := new(url.URL)
	var err error

	if dev {
		u, err = url.Parse(endpoint + method)
		if err != nil {
			return nil, err
		}
	} else {
		u, err = url.Parse(endpoint + method)
		if err != nil {
			return nil, err
		}
	}

	u.Scheme = "https"
	q := u.Query()
	for k, v := range query {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	return u, nil
}

/*
Returns all new Persons created in the last 5 minutes
*/
func (p *BlusdConnection) callAPI(t time.Time, sid string, service string, usetime bool) (out []byte, err error) {

	q := make(map[string]string)
	if usetime {
		q["dtLetzteAenderung"] = t.Format(timeLayout)
	}
	if sid != "" {
		q["schuleUID"] = sid
	}

	url, err := createUrl(q, service, p.dev, p.endpoint)
	p.url = url
	if err != nil {
		return
	}
	log.Println(p.url.String())
	req, _ := http.NewRequest("GET", p.url.String(), nil)
	req.Header.Add("Authorization", p.authentificationHeader)
	req.Header.Add("Accept", "application/json")
	resp, err := p.client.Do(req)
	resp_body, _ := ioutil.ReadAll(resp.Body)

	if strings.Contains(string(resp_body), "listError") {
		var errorResp blusdError
		err = json.Unmarshal(resp_body, &errorResp)
		if err != nil {
			return nil, errors.New(err.Error() + " " + string(resp_body))
		}
		var errorString string
		for _, k := range errorResp.ListError {
			errorString = errorString + " --- " + k.ErrorText + " : " + k.ErrorSource + " : " + " : " + k.Reference
		}
		err = resp.Body.Close()
		if err != nil {
			return
		}
		return nil, errors.New(errorString)
	}

	out = resp_body
	err = resp.Body.Close()
	if err != nil {
		return
	}
	return
}
