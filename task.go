package udnssdk

import (
	"fmt"
	"log"
	"time"
)

// TasksService provides access to the tasks resources
type TasksService struct {
	client *Client
}

// Task wraps a task response
type Task struct {
	TaskID         string `json:"taskId"`
	TaskStatusCode string `json:"taskStatusCode"`
	Message        string `json:"message"`
	ResultURI      string `json:"resultUri"`
}

// TaskListDTO wraps a list of Task resources, from an HTTP response
type TaskListDTO struct {
	Tasks      []Task     `json:"tasks"`
	Queryinfo  QueryInfo  `json:"queryInfo"`
	Resultinfo ResultInfo `json:"resultInfo"`
}

type taskWrapper struct {
	Task Task `json:"task"`
}

// taskResultPath links to the task result url.
func taskResultPath(tid string) string {
	return fmt.Sprintf("tasks/%s/result", tid)
}

// taskPath links to the task url.
func taskPath(tid string) string {
	return fmt.Sprintf("tasks/%s", tid)
}

func taskQueryPath(query string, offset int) string {
	if query != "" {
		return fmt.Sprintf("tasks?sort=NAME&query=%s&offset=%d", query, offset)
	}
	return fmt.Sprintf("tasks?offset=%d", offset)
}

// GetTaskStatus Get the status of a task.
func (s *TasksService) GetTaskStatus(tid string) (Task, *Response, error) {
	uri := taskPath(tid)
	var t Task
	res, err := s.client.get(uri, &t)
	return t, res, err
}

// GetTaskResultByURI requests a task by its URI
func (s *TasksService) GetTaskResultByURI(uri string) (*Response, error) {
	req, err := s.client.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	res, err := s.client.HTTPClient.Do(req)

	if err != nil {
		return &Response{Response: res}, err
	}
	return &Response{Response: res}, err
}

// GetTaskResultByID  requests a task by its task id
func (s *TasksService) GetTaskResultByID(tid string) (*Response, error) {
	uri := taskResultPath(tid)
	return s.client.GetResultByURI(uri)
}

// GetTaskResultByTask  requests a task by the provided task's result uri
func (s *TasksService) GetTaskResultByTask(t Task) (*Response, error) {
	return s.client.GetResultByURI(t.ResultURI)
}

// ListAllTasks requests all tasks, list
func (s *TasksService) ListAllTasks(query string) ([]Task, error) {
	// TODO: Sane Configuration for timeouts / retries
	maxerrs := 5
	waittime := 5 * time.Second

	// init accumulators
	dtos := []Task{}
	offset := 0
	errcnt := 0

	for {
		reqDtos, ri, res, err := s.ListTasks(query, offset)
		if err != nil {
			if res.StatusCode >= 500 {
				errcnt = errcnt + 1
				if errcnt < maxerrs {
					time.Sleep(waittime)
					continue
				}
			}
			return dtos, err
		}

		log.Printf("[DEBUG] ResultInfo: %+v\n", ri)
		for _, d := range reqDtos {
			dtos = append(dtos, d)
		}
		if ri.ReturnedCount+ri.Offset >= ri.TotalCount {
			return dtos, nil
		}
		offset = ri.ReturnedCount + ri.Offset
		continue
	}
}

// ListTasks request tasks by query & offset, list them also returning list metadata, the actual response, or an error
func (s *TasksService) ListTasks(query string, offset int) ([]Task, ResultInfo, *Response, error) {
	var tld TaskListDTO

	uri := taskQueryPath(query, offset)
	res, err := s.client.get(uri, &tld)

	ts := []Task{}
	for _, t := range tld.Tasks {
		ts = append(ts, t)
	}
	return ts, tld.Resultinfo, res, err
}

// DeleteTask deletes a task.
func (s *TasksService) DeleteTask(tid string) (*Response, error) {
	path := taskPath(tid)
	return s.client.delete(path, nil)
}
