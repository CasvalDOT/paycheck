package svcloud

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

var authCookieName string = ".AspNet.ApplicationCookie"

type svcloud struct {
	baseURL    string
	mainURL    string
	authCookie *http.Cookie
}

// Document structure  ...
type Document struct {
	Name string
	Path string
}

// ISvcloud main interface
type ISvcloud interface {
	Login(string, string) (*http.Cookie, error)
	DownloadDocument(Document, string) error
	ListDocuments() ([]Document, error)
}

func (s *svcloud) fetchDocumentPage() ([]byte, error) {
	var bodyAsBytes []byte
	client := http.Client{}

	request, err := http.NewRequest("GET", s.mainURL, nil)
	if err != nil {
		return bodyAsBytes, err
	}

	request.AddCookie(s.authCookie)

	response, err := client.Do(request)
	if err != nil {
		return bodyAsBytes, err
	}

	defer response.Body.Close()

	if response.StatusCode >= 400 {
		return bodyAsBytes, errors.New(response.Status)
	}

	bodyAsBytes, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return bodyAsBytes, err
	}

	return bodyAsBytes, nil
}

func (s *svcloud) ListDocuments() ([]Document, error) {
	var documents []Document

	documentPageAsByte, err := s.fetchDocumentPage()
	if err != nil {
		return documents, err
	}

	reader := bytes.NewReader(documentPageAsByte)

	document, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return documents, err
	}

	documentsInDOM := document.Find("[data-name='Name'] a")

	var wg sync.WaitGroup
	channel := make(chan (Document), documentsInDOM.Length())

	wg.Add(documentsInDOM.Length())

	go func() {
		for item := range channel {
			documents = append(documents, item)
			wg.Done()
		}
	}()

	documentsInDOM.Each(func(index int, selection *goquery.Selection) {
		var fileName string
		var href string

		defer func() {
			channel <- Document{
				Name: fileName,
				Path: href,
			}
		}()

		href, _ = selection.Attr("href")
		fileName = selection.Text()
	})

	wg.Wait()

	return documents, nil
}

func (s *svcloud) DownloadDocument(target Document, destination string) error {
	client := http.Client{}

	endpoint := s.baseURL + target.Path + "&download=True"

	request, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return err
	}

	request.AddCookie(s.authCookie)

	response, err := client.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("NOT_OK")
	}

	file, err := os.Create(destination)
	if err != nil {
		return err
	}

	file.Chmod(0o600)

	io.Copy(file, response.Body)

	return nil
}

func (s *svcloud) Login(username string, password string) (*http.Cookie, error) {
	loginPayload, err := json.Marshal(map[string]string{
		"Username": username,
		"Password": password,
	})
	if err != nil {
		return nil, err
	}

	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	request, err := http.NewRequest("POST", s.mainURL, bytes.NewBuffer(loginPayload))
	if err != nil {
		return nil, err
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; rv:91.0) Gecko/20100101 Firefox/91.0")
	request.Header.Add("Origin", s.baseURL)
	request.Header.Add("Referer", s.mainURL)
	request.Header.Add("Connection", "keep-alive")

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode >= 400 {
		return nil, errors.New(response.Status)
	}

	if len(response.Header["Set-Cookie"]) == 0 {
		return nil, errors.New("NO_SET_COOKIE")
	}

	cookie := response.Header["Set-Cookie"][0]
	cookieValueRegex := regexp.MustCompile("^" + authCookieName + "=(.*?);")

	matchs := cookieValueRegex.FindAllStringSubmatch(cookie, -1)

	if len(matchs) == 0 {
		return nil, errors.New("NO_MATCH_COOKIE_PATTERN")
	}

	if len(matchs[0]) == 0 {
		return nil, errors.New("NO_MATCH_COOKIE_PATTERN")
	}

	cookieToSet := &http.Cookie{
		Name:  authCookieName,
		Value: matchs[0][1],
	}

	s.authCookie = cookieToSet

	return cookieToSet, nil
}

// New instance
func New(baseURL string, endpoint string) ISvcloud {
	return &svcloud{
		baseURL: baseURL,
		mainURL: baseURL + endpoint,
	}
}
