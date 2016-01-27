package udnssdk

import (
	"fmt"
	"log"
	"time"
)

type TasksService struct {
	client *Client
}

type Task struct {
	TaskId         string `json:"taskId"`
	TaskStatusCode string `json:"taskStatusCode"`
	Message        string `json:"message"`
	ResultUri      string `json:"resultUri"`
}

type TaskListDTO struct {
	Tasks      []Task     `json:"tasks"`
	Queryinfo  QueryInfo  `json:"queryInfo"`
	Resultinfo ResultInfo `json:"resultInfo"`
}
type taskWrapper struct {
	Task Task `json:"task"`
}

// taskPath links to the task url.
func taskResultPath(tid string) string {
	path := fmt.Sprintf("tasks/%s/result", tid)
	/*
		if tasktype != nil {
			path += fmt.Sprintf("/%v", tasktype)
			if task != nil {
				path += fmt.Sprintf("/%v", task)
			}
		}
	*/
	return path
}
func taskPath(tid string) string {
	return fmt.Sprintf("tasks/%s", tid)
}

// Get the status of a task.
func (s *TasksService) GetTaskStatus(tid string) (Task, *Response, error) {
	reqStr := taskPath(tid)
	var t Task
	res, err := s.client.get(reqStr, &t)
	if err != nil {
		return t, res, err
	}
	return t, res, err
}

// HTTP BS to dance around bad program structure
func (s *TasksService) GetTaskResultByURI(uri string) (*Response, error) {
	req, err := s.client.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	res, err := s.client.HttpClient.Do(req)

	if err != nil {
		return &Response{Response: res}, err
	}
	return &Response{Response: res}, err
}

func (s *TasksService) GetTaskResult(tid string) (*Response, error) {
	uri := taskResultPath(tid)

	req, err := s.client.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	res, err := s.client.HttpClient.Do(req)

	if err != nil {
		return &Response{Response: res}, err
	}
	return &Response{Response: res}, err
}

// List tasks
//
func (s *TasksService) ListTasks(query string) ([]Task, *Response, error) {
	reqStr := "tasks"
	var tld TaskListDTO
	offset := 0

	log.Printf("DEBUG - ListTasks: %s\n", reqStr)

	res, err := s.client.get(reqStr, &tld)
	pis := []Task{}
	if query != "" {
		reqStr = fmt.Sprintf("%s?sort=NAME&query=%s&offset=", reqStr, query)
	} else {
		reqStr = fmt.Sprintf("%s?offset=", reqStr)
	}
	// TODO: Sane Configuration for timeouts / retries
	timeout := 5
	waittime := 5 * time.Second
	errcnt := 0
	for true {

		res, err := s.client.get(fmt.Sprintf("%s%d", reqStr, offset), &tld)
		if err != nil {
			if res.StatusCode >= 500 {
				errcnt = errcnt + 1
				if errcnt < timeout {
					time.Sleep(waittime)
					continue
				}
			}
			return pis, res, err

		}
		log.Printf("DEBUG - ResultInfo: %+v\n", tld.Resultinfo)
		for _, pi := range tld.Tasks {
			pis = append(pis, pi)
		}
		if tld.Resultinfo.ReturnedCount+tld.Resultinfo.Offset >= tld.Resultinfo.TotalCount {
			return pis, res, nil
		} else {
			offset = tld.Resultinfo.ReturnedCount + tld.Resultinfo.Offset
			continue
		}
	}
	return pis, res, err
}

// DeleteTask deletes a task.
//
func (s *TasksService) DeleteTask(tid string) (*Response, error) {
	path := taskPath(tid)
	return s.client.delete(path, nil)
}
