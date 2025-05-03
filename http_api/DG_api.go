package http_api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// 定义 API 基础 URL 常量
const BaseURL = "http://127.0.0.1:8920" // 这里替换为实际的 API URL

// 游戏信息响应结构体
type GameInfoResponse struct {
	Status         int    `json:"status"`
	Code           string `json:"code"`
	StrengthConfig struct {
		Strength       int `json:"strength"`
		RandomStrength int `json:"randomStrength"`
	} `json:"strengthConfig"`
	GameConfig struct {
		StrengthChangeInterval     []int   `json:"strengthChangeInterval"`
		EnableBChannel             bool    `json:"enableBChannel"`
		BChannelStrengthMultiplier float64 `json:"bChannelStrengthMultiplier"`
		PulseId                    string  `json:"pulseId"`
		PulseMode                  string  `json:"pulseMode"`
		PulseChangeInterval        int     `json:"pulseChangeInterval"`
	} `json:"gameConfig"`
	ClientStrength struct {
		Strength int `json:"strength"`
		Limit    int `json:"limit"`
	} `json:"clientStrength"`
	CurrentPulseId string `json:"currentPulseId"`
}

// 获取游戏信息函数
func GetGameInfo(clientId string) (*GameInfoResponse, error) {
	url := fmt.Sprintf("%s/api/v2/game/%s", BaseURL, clientId) // 使用常量 BaseURL
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer resp.Body.Close()

	var gameInfo GameInfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&gameInfo); err != nil {
		return nil, err
	}

	if gameInfo.Status != 1 {
		return nil, fmt.Errorf("Failed to get game info: %s", gameInfo.Code)
	}

	return &gameInfo, nil
}

// 波形列表响应结构体
type PulseListResponse struct {
	Status    int    `json:"status"`
	Code      string `json:"code"`
	PulseList []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"pulseList"`
}

// 获取波形列表函数
func GetPulseList() (*PulseListResponse, error) {
	url := fmt.Sprintf("%s/api/v2/pulse_list", BaseURL) // 使用常量 BaseURL
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer resp.Body.Close()

	var pulseList PulseListResponse
	if err := json.NewDecoder(resp.Body).Decode(&pulseList); err != nil {
		return nil, err
	}

	if pulseList.Status != 1 {
		return nil, fmt.Errorf("Failed to get pulse list: %s", pulseList.Code)
	}

	return &pulseList, nil
}

// 设置强度请求的结构体
type SetStrengthConfigRequest struct {
	Strength struct {
		Add *int `json:"add,omitempty"`
		Sub *int `json:"sub,omitempty"`
		Set *int `json:"set,omitempty"`
	} `json:"strength,omitempty"`
	RandomStrength struct {
		Add *int `json:"add,omitempty"`
		Sub *int `json:"sub,omitempty"`
		Set *int `json:"set,omitempty"`
	} `json:"randomStrength,omitempty"`
}

// 设置强度函数
func SetStrength(clientId string, config SetStrengthConfigRequest) error {
	url := fmt.Sprintf("%s/api/v2/game/%s/strength", BaseURL, clientId) // 使用常量 BaseURL

	body, err := json.Marshal(config)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer resp.Body.Close()

	var response struct {
		Status           int      `json:"status"`
		Code             string   `json:"code"`
		Message          string   `json:"message"`
		SuccessClientIds []string `json:"successClientIds"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return err
	}

	if response.Status != 1 {
		return fmt.Errorf("Failed to set strength: %s", response.Message)
	}

	return nil
}

// 一键开火请求结构体
type FireRequest struct {
	Strength int    `json:"strength"`
	Time     int    `json:"time"`
	Override bool   `json:"override"`
	PulseId  string `json:"pulseId"`
}

// 一键开火函数
func Fire(clientId string, request FireRequest) error {
	url := fmt.Sprintf("%s/api/v2/game/%s/action/fire", BaseURL, clientId) // 使用常量 BaseURL

	body, err := json.Marshal(request)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer resp.Body.Close()

	var response struct {
		Status           int      `json:"status"`
		Code             string   `json:"code"`
		Message          string   `json:"message"`
		SuccessClientIds []string `json:"successClientIds"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return err
	}

	if response.Status != 1 {
		return fmt.Errorf("Failed to trigger fire action: %s", response.Message)
	}

	return nil
}
