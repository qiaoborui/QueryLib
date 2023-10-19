package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"sync"
)

type studentIDs struct {
	IDs []string `json:"ids"`
}

type singleStudentID struct {
	ID string `json:"id"`
}

type result struct {
	ID         string `json:"id"`
	SpaceName  string `json:"spaceName"`
	NameMerge  string `json:"nameMerge"`
	StatusName string `json:"statusName"`
	Date       string `json:"date"`
	Error      string `json:"error,omitempty"`
}

type JSONData struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
	Data   []Data `json:"data"`
}

type BeginTime struct {
	Date         string `json:"date"`
	Timezone     string `json:"timezone"`
	TimezoneType int    `json:"timezone_type"`
}

type Data struct {
	BeginTime  BeginTime `json:"beginTime"`
	SpaceName  string    `json:"spaceName"`
	NameMerge  string    `json:"nameMerge"`
	StatusName string    `json:"statusname"`
}

func QueryUsers(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error reading request body",
		})
		return
	}
	var s studentIDs
	err = json.Unmarshal(body, &s)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error unmarshalling request body",
		})
		return
	}
	if s.IDs == nil {
		var ss singleStudentID
		err = json.Unmarshal(body, &ss)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error unmarshalling request body",
			})
			return
		}
		s.IDs = append(s.IDs, ss.ID)
	}
	c.JSON(http.StatusOK, getStudentData(s.IDs))

}

func QueryUser(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No id provided",
		})
		return
	}
	c.JSON(http.StatusOK, getStudentData([]string{id}))
}

func getStudentData(ids []string) []result {
	concurrencyLimit := 20
	results := make([]result, 0)
	sem := make(chan bool, concurrencyLimit)
	var wg sync.WaitGroup
	var appendwg sync.WaitGroup
	resCh := make(chan result)
	appendwg.Add(1)
	// This goroutine is responsible for appending to the results slice
	go func() {
		for r := range resCh {
			results = append(results, r)
		}
		appendwg.Done()
	}()
	for _, id := range ids {
		wg.Add(1)
		sem <- true
		go func(id string) {
			defer wg.Done()
			defer func() { <-sem }()
			resp, err := http.Get(fmt.Sprintf("http://rg.lib.xauat.edu.cn/api.php/currentuse?user=%s", id))
			if err != nil {
				resCh <- result{ID: id, Error: "Error: " + err.Error()}
			} else {
				defer resp.Body.Close()
				var data JSONData
				err = json.NewDecoder(resp.Body).Decode(&data)
				fmt.Printf("id: %s %+v\n", id, data)
				if err != nil {
					resCh <- result{ID: id, Error: "Error: " + err.Error()}
				} else {
					if data.Status == 1 {
						if len(data.Data) == 0 {
							resCh <- result{ID: id, Error: "Error: No Data", StatusName: "未使用"}
						}
						for _, d := range data.Data {
							resCh <- result{
								ID:         id,
								SpaceName:  d.SpaceName,
								NameMerge:  d.NameMerge,
								StatusName: d.StatusName,
								Date:       d.BeginTime.Date,
							}
						}
					} else {
						resCh <- result{ID: id, Error: "Error: " + data.Msg}
					}
				}
			}
		}(id)
	}
	wg.Wait()
	close(resCh) // Close the channel after all goroutines have finished
	appendwg.Wait()
	return results
}
