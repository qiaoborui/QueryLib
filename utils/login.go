package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

var jar *cookiejar.Jar

type Session struct {
	client *http.Client
}

func NewSession() *Session {
	jar, _ = cookiejar.New(nil)
	return &Session{
		client: &http.Client{
			Jar: jar,
		},
	}
}

func (s *Session) Login(username, password string) error {
	cap_url := "http://rg.lib.xauat.edu.cn/api.php/login"
	captcha, err := s.getCaptcha()
	if err != nil {
		return fmt.Errorf("failed to get captcha: %v", err)
	}
	data := url.Values{}
	data.Set("username", username)
	data.Set("password", password)
	data.Set("verify", captcha)
	req, err := http.NewRequest("POST", cap_url, strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	res, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send login request: %v", err)
	}
	defer res.Body.Close()

	responseData, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read login response: %v", err)
	}

	var loginResponse struct {
		Status int `json:"status"`
	}
	err = json.Unmarshal(responseData, &loginResponse)
	if err != nil {
		return fmt.Errorf("failed to parse login response: %v", err)
	}

	if loginResponse.Status == 1 {
		fmt.Println("Login successful!")
		return nil
	} else {
		fmt.Println("Login failed!")
		return fmt.Errorf("login request returned non-OK status")
	}
}
func (s *Session) getCaptcha() (string, error) {
	resp, err := s.client.Get("http://rg.lib.xauat.edu.cn/api.php/check")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	// Read image data
	imgData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	// Prepare multipart data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("image", "captcha.jpg")
	if err != nil {
		return "", err
	}
	part.Write(imgData)
	writer.Close()

	// Post data
	req, err := http.NewRequest("POST", "https://captcha-captcha-gzdtwircml.cn-hangzhou.fcapp.run/ocr", body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err = s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read response
	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(res), nil
}
