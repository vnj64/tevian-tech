package face_cloud

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sync"
	"tevian/domain/models"
)

const (
	minWorkers = 4
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

func (s service) processImage(token string, imageAddress string) (map[string]interface{}, error) {
	file, err := os.Open(imageAddress)
	if err != nil {
		return nil, fmt.Errorf("can't open file %s: %v", imageAddress, err)
	}
	defer file.Close()

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/detect?%s", s.cfg.BaseFaceCloudUrl(), queryParameters()), file)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Add("Content-Type", "image/jpeg")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var resultObj map[string]interface{}
	err = json.Unmarshal(result, &resultObj)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	return resultObj, nil
}

func (s service) Detect(token string, images []models.Image) (models.DetectResult, error) {
	resultsChan := make(chan map[string]interface{})
	errChan := make(chan error)
	taskQueue := make(chan models.Image, len(images))

	workerCount := minWorkers
	if len(images) < minWorkers {
		workerCount = len(images)
	}
	var wg sync.WaitGroup

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for image := range taskQueue {
				result, err := s.processImage(token, image.ImageAddress)
				if err != nil {
					errChan <- err
					continue
				}
				result["imageName"] = image.ImageName
				resultsChan <- result
			}
		}()
	}

	for _, image := range images {
		taskQueue <- image
	}
	close(taskQueue)

	go func() {
		wg.Wait()
		close(resultsChan)
		close(errChan)
	}()

	var finalResult models.DetectResult
	var finalError error

	for result := range resultsChan {
		imageName := result["imageName"].(string)
		finalResult.ImageData = append(finalResult.ImageData, parseImageResult(result, imageName))
	}

	select {
	case err := <-errChan:
		finalError = err
	default:
		finalResult.Status = models.StatusCompleted
	}

	return finalResult, finalError
}

func parseImageResult(info map[string]interface{}, imageName string) models.ImageData {
	data := info["data"].([]interface{})

	var faces []models.Faces

	for _, entry := range data {
		person := entry.(map[string]interface{})
		bbox := person["bbox"].(map[string]interface{})
		demographics := person["demographics"].(map[string]interface{})
		ageMap := demographics["age"].(map[string]interface{})
		age := ageMap["mean"].(float64)
		gender := demographics["gender"].(string)

		face := models.Faces{
			BoundingBox: models.BoundingBox{
				X:      int(bbox["x"].(float64)),
				Y:      int(bbox["y"].(float64)),
				Width:  int(bbox["width"].(float64)),
				Height: int(bbox["height"].(float64)),
			},
			Gender: gender,
			Age:    age,
		}

		faces = append(faces, face)
	}

	return models.ImageData{
		Name:  imageName,
		Faces: faces,
	}
}
