package face_cloud

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"tevian/domain/services"
)

func queryParameters() string {
	params := url.Values{}
	params.Add("fd_min_size", "0")
	params.Add("fd_max_size", "0")
	params.Add("fd_threshold", "0.8")
	params.Add("rotate_until_faces_found", "false")
	params.Add("orientation_classifier", "false")
	params.Add("demographics", "true")
	params.Add("attributes", "true")
	params.Add("landmarks", "false")
	params.Add("liveness", "false")
	params.Add("quality", "false")
	params.Add("masks", "false")

	return params.Encode()
}

func Detect(token string, body io.Reader, cfg services.Config) (string, error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/detect?%s", cfg.BaseFaceCloudUrl(), queryParameters()), body)
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}
	defer req.Body.Close()

	req.Header.Add("Content-Type", "image/jpeg")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	return string(result), nil

}
