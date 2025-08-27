package gigachad

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
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

type GigaChatRequest struct {
	Model    string            `json:"model"`
	Messages []GigaChatMessage `json:"messages"`
}

type GigaChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GigaChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
			Role    string `json:"role"`
		} `json:"message"`
		Index        int    `json:"index"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Object  string `json:"object"`
	Usage   struct {
		PromptTokens          int `json:"prompt_tokens"`
		CompletionTokens      int `json:"completion_tokens"`
		TotalTokens           int `json:"total_tokens"`
		PrecachedPromptTokens int `json:"precached_prompt_tokens"`
	} `json:"usage"`
}

type GigaChatImageRequest struct {
	Model        string            `json:"model"`
	Messages     []GigaChatMessage `json:"messages"`
	FunctionCall string            `json:"function_call"`
}

type GigaChatImageResponse struct {
	Created int64 `json:"created"`
	Data    []struct {
		URL     string `json:"url,omitempty"`
		B64Json string `json:"b64_json,omitempty"`
	} `json:"data"`
}

func NewApi() (*GigaChatApi, error) {
	clientID := os.Getenv("GIGACHAT_CLIENT_ID")
	clientSecret := os.Getenv("GIGACHAT_CLIENT_SECRET")

	// Формируем Basic auth строку
	auth := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))
	log.Println(auth)
	// Данные формы
	form := url.Values{}
	form.Set("scope", "GIGACHAT_API_PERS")

	// Создаём POST-запрос
	req, err := http.NewRequest("POST", "https://ngw.devices.sberbank.ru:9443/api/v2/oauth",
		strings.NewReader(form.Encode()))
	if err != nil {
		log.Fatalf("Ошибка создания запроса: %v", err)
	}

	// Заголовки
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("RqUID", uuid.New().String())

	// HTTP-клиент с отключенной проверкой TLS (аналог curl -k)
	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // ⚠️ ТОЛЬКО ДЛЯ ТЕСТОВ!
			},
		},
	}

	// Отправляем запрос
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Ошибка отправки запроса: %v", err)
	}
	defer resp.Body.Close()

	// Читаем и выводим тело ответа
	body, _ := io.ReadAll(resp.Body)
	tokenResponse := TokenResponse{}
	json.Unmarshal(body, &tokenResponse)

	return &GigaChatApi{
		tokenResponse.AccessToken,
	}, nil
}

func (gigaChat GigaChatApi) Send(content string) GigaChatResponse {
	url := "https://gigachat.devices.sberbank.ru/api/v1/chat/completions"

	request := GigaChatRequest{
		Model: "GigaChat",
		Messages: []GigaChatMessage{
			{
				Role:    "user",
				Content: content,
			},
		},
	}
	body, err := json.Marshal(request)

	if err != nil {
		log.Println(err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))

	if err != nil {
		log.Println(err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+gigaChat.AccessToken)

	response, err := gigaChat.Request(*req)

	if err != nil {
		log.Println(err)
	}

	data, er := io.ReadAll(response.Body)
	if er != nil {
		log.Println(er)
	}
	
	resp := GigaChatResponse{}
	e := json.Unmarshal(data, &resp)
	if err != nil {
		log.Println(e)
	}

	log.Println(string(data))

	return resp
}

func (GigaChatApi GigaChatApi) Request(request http.Request) (*http.Response, error) {
	client := &http.Client{
		Timeout: 40 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // ⚠️ ТОЛЬКО ДЛЯ ТЕСТОВ!
			},
		},
	}

	resp, err := client.Do(&request)
	if err != nil {

		return nil, err
	}

	return resp, nil
}

func (gigaChat GigaChatApi) GenerateImage(prompt string) ([]byte, error) {
	// 1. Запрос на генерацию изображения
	request := GigaChatImageRequest{
		Model:        "GigaChat",
		Messages:     []GigaChatMessage{{Role: "user", Content: prompt}},
		FunctionCall: "auto",
	}
	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", "https://gigachat.devices.sberbank.ru/api/v1/chat/completions", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+gigaChat.AccessToken)

	resp, err := gigaChat.Request(*req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)

	log.Println(string(respBody), err)
	if err != nil {
		return nil, err
	}
	var chatResp GigaChatResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return nil, err
	}

	// 2. Получаем file_id
	var imageID string
	for _, choice := range chatResp.Choices {
		content := choice.Message.Content
		// Регулярка для поиска uuid4 внутри src=""
		re := regexp.MustCompile(`<img[^>]+src="([a-f0-9\-]{36})"`)
		m := re.FindStringSubmatch(content)
		if len(m) == 2 {
			imageID = m[1]
			break
		}
	}
	if imageID == "" {
		return nil, fmt.Errorf("image uuid not found in response")
	}
	log.Println("Image UUID:", imageID)
	// 3. Скачиваем изображение
	fileUrl := "https://gigachat.devices.sberbank.ru/api/v1/files/" + imageID + "/content"
	fileReq, err := http.NewRequest("GET", fileUrl, nil)
	if err != nil {
		return nil, err
	}
	fileReq.Header.Add("Authorization", "Bearer "+gigaChat.AccessToken)
	fileReq.Header.Add("Accept", "image/jpeg")

	fileResp, err := gigaChat.Request(*fileReq)

	if err != nil {
		return nil, err
	}
	defer fileResp.Body.Close()
	imgData, err := io.ReadAll(fileResp.Body)

	if err != nil {
		return nil, err
	}
	return imgData, nil
}
