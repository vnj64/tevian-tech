package face_cloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"tevian/domain/services"
)

type AuthResponse struct {
	Data struct {
		AccessToken string `json:"access_token"`
	} `json:"data"`
	Status int `json:"status"`
}

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func GetAccessToken(cloudLogin, cloudPassword string, cfg services.Config) (string, error) {
	body := User{
		Email:    cloudLogin,
		Password: cloudPassword,
	}
	jsonData, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}

	authUrl := fmt.Sprintf("%s/api/v1/login", cfg.BaseFaceCloudUrl())
	resp, err := http.Post(authUrl, "application/json", bytes.NewReader(jsonData))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var result AuthResponse
	err = json.Unmarshal(bodyBytes, &result)

	return result.Data.AccessToken, nil
}
