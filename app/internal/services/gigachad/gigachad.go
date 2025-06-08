package gigachad

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

type GigaChatApi struct {
	AccessToken string
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

func NewApi() (*GigaChatApi, error) {
	urls := "https://ngw.devices.sberbank.ru:9443/api/v2/oauth"

	clientID := os.Getenv("GIGACHAT_CLIENT_ID")
	clientSecret := os.Getenv("GIGACHAT_CLIENT_SECRET")

	log.Println(clientID + "     " + clientSecret)
	log.Println("TEST")

	authStr := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))

	form := url.Values{}
	form.Add("scope", "GIGACHAT_API_PERS") // или GIGACHAT_API_B2B / GIGACHAT_API_CORP

	req, err := http.NewRequest("POST", urls, strings.NewReader(form.Encode()))

	if err != nil {
		log.Println("failed to create request:", err)
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")
	req.Header.Set("Authorization", "Basic "+authStr)
	req.Header.Set("RqUID", uuid.New().String()) // UUID v4 для жур

	client := http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // отключает проверку сертификата (только для теста)
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("request failed:", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("failed to get token: %s", resp.Status)
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var tokenResp TokenResponse
	err = json.NewDecoder(resp.Body).Decode(&tokenResp)
	if err != nil {
		log.Println("failed to decode response:", err)
		return nil, err
	}

	log.Println("success token:", tokenResp.AccessToken)

	// Создаём и возвращаем экземпляр API с полученным токеном (пример)
	api := &GigaChatApi{
		AccessToken: tokenResp.AccessToken,
		// другие поля инициализируй по необходимости
	}

	return api, nil
}
