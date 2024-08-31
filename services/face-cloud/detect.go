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
	maxWorkers = 16
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

func (s service) Detect(token string, task *models.Task) (models.DetectResult, error) {
	if task.ImageAddress == nil {
		return models.DetectResult{}, fmt.Errorf("task has no image address")
	}

	resultsChan := make(chan map[string]interface{})
	errChan := make(chan error)
	taskQueue := make(chan string, 1)

	workerCount := minWorkers
	var wg sync.WaitGroup
	// тут подумать по реализации gracefully shutdown
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for imageAddress := range taskQueue {
				result, err := s.processImage(token, imageAddress)
				if err != nil {
					errChan <- err
					continue
				}
				resultsChan <- result
			}
		}()
	}

	taskQueue <- *task.ImageAddress

	go func() {
		wg.Wait()
		close(resultsChan)
		close(errChan)
		close(taskQueue)
	}()

	var finalResult models.DetectResult
	var finalError error

	select {
	case result := <-resultsChan:
		fmt.Println(result)
		finalResult = parseResponse(result, task)
	case err := <-errChan:
		finalError = err
	}

	fmt.Println(finalResult)
	return finalResult, finalError
}

func parseResponse(info map[string]interface{}, task *models.Task) models.DetectResult {
	data := info["data"].([]interface{})
	var images []models.ImageData
	var stats models.Statistics
	var totalMaleAge, totalFemaleAge float64
	var totalMales, totalFemales int

	for _, entry := range data {
		person := entry.(map[string]interface{})
		imageName := task.ImageName
		bbox := person["bbox"].(map[string]interface{})

		demographics := person["demographics"].(map[string]interface{})
		age := demographics["age"].(map[string]interface{})["mean"].(float64)
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

		stats.TotalFaces++
		if gender == "male" {
			totalMales++
			totalMaleAge += age
		} else if gender == "female" {
			totalFemales++
			totalFemaleAge += age
		}

		imageFound := false
		for i := range images {
			if images[i].Name == *imageName {
				images[i].Faces = append(images[i].Faces, face)
				imageFound = true
				break
			}
		}

		if !imageFound {
			images = append(images, models.ImageData{
				Name:  *imageName,
				Faces: []models.Faces{face},
			})
		}
	}

	if totalMales > 0 {
		stats.AverageMaleAge = totalMaleAge / float64(totalMales)
	}
	if totalFemales > 0 {
		stats.AverageFemaleAge = totalFemaleAge / float64(totalFemales)
	}
	stats.TotalMales = totalMales
	stats.TotalFemales = totalFemales

	return models.DetectResult{
		TaskId:     task.Id,
		Status:     models.StatusCompleted,
		ImageData:  images,
		Statistics: stats,
	}
}
