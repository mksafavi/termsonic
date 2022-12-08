package music

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"strings"

	"github.com/delucks/go-subsonic"
	"github.com/jfbus/httprs"
)

// Stream2 patches subsonic.Client.Stream to return a ReadCloser for use with "beep"
func Stream2(s *subsonic.Client, id string, parameters map[string]string) (io.ReadCloser, error) {
	params := url.Values{}
	params.Add("id", id)
	for k, v := range parameters {
		params.Add(k, v)
	}
	response, err := s.Request("GET", "stream", params)
	if err != nil {
		return nil, err
	}
	contentType := response.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "text/xml") || strings.HasPrefix(contentType, "application/xml") {
		// An error was returned
		responseBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
		resp := subsonic.Response{}
		err = xml.Unmarshal(responseBody, &resp)
		if err != nil {
			return nil, err
		}
		if resp.Error != nil {
			err = fmt.Errorf("Error #%d: %s\n", resp.Error.Code, resp.Error.Message)
		} else {
			err = fmt.Errorf("An error occurred: %#v\n", resp)
		}
		return nil, err
	}

	return httprs.NewHttpReadSeeker(response), nil
}

func Download2(s *subsonic.Client, id string) (io.ReadCloser, error) {
	params := url.Values{}
	params.Add("id", id)
	response, err := s.Request("GET", "download", params)
	if err != nil {
		return nil, err
	}
	contentType := response.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "text/xml") || strings.HasPrefix(contentType, "application/xml") {
		// An error was returned
		responseBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
		resp := subsonic.Response{}
		err = xml.Unmarshal(responseBody, &resp)
		if err != nil {
			return nil, err
		}
		if resp.Error != nil {
			err = fmt.Errorf("Error #%d: %s\n", resp.Error.Code, resp.Error.Message)
		} else {
			err = fmt.Errorf("An error occurred: %#v\n", resp)
		}
		return nil, err
	}
	return httprs.NewHttpReadSeeker(response), nil
}
