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
	minWorkers = 4  // Минимальное количество горутин
	maxWorkers = 16 // Максимальное количество горутин
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

func (s service) processImage(token string, imageAddress string, wg *sync.WaitGroup, results chan<- *models.ResultData, errChan chan<- error) {
	defer wg.Done()

	file, err := os.Open(imageAddress)
	if err != nil {
		{
			errChan <- fmt.Errorf("cant open file: %v", err)
			return
		}
	}
	defer file.Close()

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/detect?%s", s.cfg.BaseFaceCloudUrl(), queryParameters()), file)
	if err != nil {
		errChan <- fmt.Errorf("error creating request: %v", err)
		return
	}

	req.Header.Add("Content-Type", "image/jpeg")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		errChan <- fmt.Errorf("error making request: %v", err)
		return
	}
	defer resp.Body.Close()

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		errChan <- fmt.Errorf("error reading response body: %v", err)
		return
	}

	var fnResult models.ResultData
	err = json.Unmarshal(result, &fnResult)
	if err != nil {
		errChan <- fmt.Errorf("error unmarshaling: %v", err)
		return
	}

	results <- &fnResult
}

func (s service) Detect(token string, imageAddresses []string) ([]*models.ResultData, error) {
	var wg sync.WaitGroup
	results := make(chan *models.ResultData, len(imageAddresses))
	errChan := make(chan error, len(imageAddresses))

	workerCount := minWorkers
	if len(imageAddresses) < minWorkers {
		workerCount = len(imageAddresses)
	} else if len(imageAddresses) > maxWorkers {
		workerCount = maxWorkers
	}

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, imageAddress := range imageAddresses {
				wg.Add(1)
				go s.processImage(token, imageAddress, &wg, results, errChan)
			}
		}()
	}

	wg.Wait()
	close(results)
	close(errChan)

	var finalResults []*models.ResultData
	for result := range results {
		finalResults = append(finalResults, result)
	}

	var finalError error
	for err := range errChan {
		if finalError == nil {
			finalError = err
		}
	}

	return finalResults, finalError
}
