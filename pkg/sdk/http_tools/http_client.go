package http_tools

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/tencentad/marketing-api-go-sdk/pkg/config"
	"golang.org/x/oauth2"
)

// HttpClient
type HttpClient struct {
	client *http.Client
	Config *HttpConfig

	//
	Ctx context.Context
	Method string
	Path   string
	HeaderParams map[string]string
	QueryParams map[string]string
	PostBody interface{}
	formParams map[string]string
	FileName string
	FileBytes []byte
	FileKey string
}

// DoProcess 发送请求
func (t *HttpClient) DoProcess() (string, error) {
	// init http config
	t.setHttpConfig()
	// make request
	request, err := t.prepareRequest()
	//r, err := t.client.prepareRequest(ctx, localVarPath, localVarHttpMethod, localVarPostBody, localVarHeaderParams,
	//	localVarQueryParams, localVarFormParams, localVarFileName, localVarFileBytes, localVarFileKey)
	response, err := t.client.Do(request)
	if err != nil || response == nil {
		return "", err
	}

	localVarBody, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		return "", err
	}

	if response.StatusCode < 300 {
		// If we succeed, return the data, otherwise pass on to decode error.
		return string(localVarBody), nil
	}

	if response.StatusCode >= 300 {
		return "", fmt.Errorf("http response: code = %d, body = %s", response.StatusCode, string(localVarBody))
	}
}

func (t *HttpClient) setHttpConfig() {
	t.client = &http.Client{
		Timeout: time.Duration(t.Config.Timeout),
	}
	// set default Content-Type
	if len(t.HeaderParams["Content-Type"]) == 0 {
		t.HeaderParams["Content-Type"] = "text/plain"
	}

	// set default Accept
	if len(t.HeaderParams["Accept"]) == 0 {
		t.HeaderParams["Accept"] = "application/json"
	}
}

func (t *HttpClient) SetMethod(method string) {
	t.Method = strings.ToUpper("method")
}

func (t *HttpClient) getHttpPath() string {
	path := t.Config.BasePath
	if len(t.Config.apiVersion) > 0 {
		path = path + "/" + t.Config.apiVersion
	}
	return path + "/" + t.Method
}

func (t *HttpClient) AddHeaderParam (key string, value string) {
	t.HeaderParams[key] = value
}

func (t *HttpClient) AddQueryParam (key string, value string) {
	t.QueryParams[key] = value
}

// prepareRequest build the request
func (t *HttpClient) prepareRequest() (request *http.Request, err error) {

	var body *bytes.Buffer

	// Detect postBody type and post.
	if t.PostBody != nil {
		contentType := t.HeaderParams["Content-Type"]
		if contentType == "" {
			contentType = detectContentType(t.PostBody)
			t.HeaderParams["Content-Type"] = contentType
		}

		body, err = setBody(t.PostBody, contentType)
		if err != nil {
			return nil, err
		}
	}

	// add form parameters and file if available.
	if strings.HasPrefix(t.HeaderParams["Content-Type"], "multipart/form-data") && len(formParams) > 0 || (len(
		t.FileBytes) > 0 && t.FileName != "") {
		if body != nil {
			return nil, errors.New("Cannot specify postBody and multipart form at the same time.")
		}
		body = &bytes.Buffer{}
		w := multipart.NewWriter(body)

		for k, v := range formParams {
			for _, iv := range v {
				if strings.HasPrefix(k, "@") { // file
					err = addFile(w, k[1:], iv)
					if err != nil {
						return nil, err
					}
				} else { // form value
					w.WriteField(k, iv)
				}
			}
		}
		if len(fileBytes) > 0 && fileName != "" {
			w.Boundary()
			//_, fileNm := filepath.Split(fileName)
			part, err := w.CreateFormFile(fileKey, filepath.Base(fileName))
			if err != nil {
				return nil, err
			}
			_, err = part.Write(fileBytes)
			if err != nil {
				return nil, err
			}
		}
		// Set the Boundary in the Content-Type
		headerParams["Content-Type"] = w.FormDataContentType()

		// Set Content-Length
		headerParams["Content-Length"] = fmt.Sprintf("%d", body.Len())
		w.Close()
	}

	if strings.HasPrefix(headerParams["Content-Type"], "application/x-www-form-urlencoded") && len(formParams) > 0 {
		if body != nil {
			return nil, errors.New("Cannot specify postBody and x-www-form-urlencoded form at the same time.")
		}
		body = &bytes.Buffer{}
		body.WriteString(formParams.Encode())
		// Set Content-Length
		headerParams["Content-Length"] = fmt.Sprintf("%d", body.Len())
	}

	// Setup path and query parameters
	url, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	// Adding Query Param
	query := url.Query()
	for k, v := range queryParams {
		for _, iv := range v {
			query.Add(k, iv)
		}
	}

	// Encode the parameters.
	url.RawQuery = query.Encode()

	// Generate a new request
	if body != nil {
		localVarRequest, err = http.NewRequest(method, url.String(), body)
	} else {
		localVarRequest, err = http.NewRequest(method, url.String(), nil)
	}
	if err != nil {
		return nil, err
	}

	// add header parameters, if any
	if len(headerParams) > 0 {
		headers := http.Header{}
		for h, v := range headerParams {
			headers.Set(h, v)
		}
		localVarRequest.Header = headers
	}

	// Override request host, if applicable
	if c.Cfg.Host != "" {
		localVarRequest.Host = c.Cfg.Host
	}

	// Add the user agent to the request.
	localVarRequest.Header.Add("User-Agent", c.Cfg.UserAgent)

	if ctx != nil {
		// add context to the request
		localVarRequest = localVarRequest.WithContext(ctx)

		// Walk through any authentication.

		// OAuth2 authentication
		if tok, ok := ctx.Value(config.ContextOAuth2).(oauth2.TokenSource); ok {
			// We were able to grab an oauth2 token from the context
			var latestToken *oauth2.Token
			if latestToken, err = tok.Token(); err != nil {
				return nil, err
			}

			latestToken.SetAuthHeader(localVarRequest)
		}

		// Basic HTTP Authentication
		if auth, ok := ctx.Value(config.ContextBasicAuth).(config.BasicAuth); ok {
			localVarRequest.SetBasicAuth(auth.UserName, auth.Password)
		}

		// AccessToken Authentication
		if auth, ok := ctx.Value(config.ContextAccessToken).(string); ok {
			localVarRequest.Header.Add("Authorization", "Bearer "+auth)
		}
	}

	for header, value := range c.Cfg.DefaultHeader {
		localVarRequest.Header.Add(header, value)
	}

	return localVarRequest, nil
}

func Decode(v interface{}, b []byte, contentType string) (err error) {
	if strings.Contains(contentType, "application/xml") {
		if err = xml.Unmarshal(b, v); err != nil {
			return err
		}
		return nil
	} else if strings.Contains(contentType, "application/json") {
		dec := json.NewDecoder(bytes.NewReader(b))
		if err := dec.Decode(v); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("undefined response type")
}

// detectContentType method is used to figure out `Request.Body` content type for request header
func detectContentType(body interface{}) string {
	contentType := "text/plain; charset=utf-8"
	kind := reflect.TypeOf(body).Kind()

	switch kind {
	case reflect.Struct, reflect.Map, reflect.Ptr:
		contentType = "application/json; charset=utf-8"
	case reflect.String:
		contentType = "text/plain; charset=utf-8"
	default:
		if b, ok := body.([]byte); ok {
			contentType = http.DetectContentType(b)
		} else if kind == reflect.Slice {
			contentType = "application/json; charset=utf-8"
		}
	}

	return contentType
}

// Set request body from an interface{}
func setBody(body interface{}, contentType string) (bodyBuf *bytes.Buffer, err error) {
	if bodyBuf == nil {
		bodyBuf = &bytes.Buffer{}
	}

	if reader, ok := body.(io.Reader); ok {
		_, err = bodyBuf.ReadFrom(reader)
	} else if b, ok := body.([]byte); ok {
		_, err = bodyBuf.Write(b)
	} else if s, ok := body.(string); ok {
		_, err = bodyBuf.WriteString(s)
	} else if s, ok := body.(*string); ok {
		_, err = bodyBuf.WriteString(*s)
	} else if jsonCheck.MatchString(contentType) {
		err = json.NewEncoder(bodyBuf).Encode(body)
	} else if xmlCheck.MatchString(contentType) {
		xml.NewEncoder(bodyBuf).Encode(body)
	}

	if err != nil {
		return nil, err
	}

	if bodyBuf.Len() == 0 {
		err = fmt.Errorf("Invalid body type %s\n", contentType)
		return nil, err
	}
	return bodyBuf, nil
}
