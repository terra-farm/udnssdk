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
	t := TaskID(tid)
	return t.ResultURI()
}

// taskPath links to the task url.
func taskPath(tid string) string {
	t := TaskID(tid)
	return t.URI()
}

func taskQueryPath(query string, offset int) string {
	return TasksQueryURI(query, offset)
}

// GetTaskStatus Get the status of a task.
func (s *TasksService) GetTaskStatus(tid string) (Task, *Response, error) {
	t := TaskID(tid)
	return s.Find(t)
}

// GetTaskResultByURI requests a task by its URI
func (s *TasksService) GetTaskResultByURI(uri string) (*Response, error) {
	return s.client.GetResultByURI(uri)
}

// GetTaskResultByID  requests a task by its task id
func (s *TasksService) GetTaskResultByID(tid string) (*Response, error) {
	t := TaskID(tid)
	return s.FindResult(t)
}

// GetTaskResultByTask  requests a task by the provided task's result uri
func (s *TasksService) GetTaskResultByTask(t Task) (*Response, error) {
	return s.FindResultByTask(t)
}

// ListAllTasks requests all tasks, list
func (s *TasksService) ListAllTasks(query string) ([]Task, error) {
	return s.Select(query)
}

// ListTasks request tasks by query & offset, list them also returning list metadata, the actual response, or an error
func (s *TasksService) ListTasks(query string, offset int) ([]Task, ResultInfo, *Response, error) {
	return s.SelectWithOffset(query, offset)
}

// DeleteTask deletes a task.
func (s *TasksService) DeleteTask(tid string) (*Response, error) {
	t := TaskID(tid)
	return s.Delete(t)
}

// ===== //

type TaskID string

// URI generates URI for the task result
func (t TaskID) ResultURI() string {
	return fmt.Sprintf("%s/result", t.URI())
}
// taskPath links to the task url.
func (t TaskID) URI() string {
	return fmt.Sprintf("tasks/%s", t)
}

func TasksQueryURI(query string, offset int) string {
	if query != "" {
		return fmt.Sprintf("tasks?sort=NAME&query=%s&offset=%d", query, offset)
	}
	return fmt.Sprintf("tasks?offset=%d", offset)
}

// Select requests all tasks, with pagination
func (s *TasksService) Select(query string) ([]Task, error) {
	// TODO: Sane Configuration for timeouts / retries
	maxerrs := 5
	waittime := 5 * time.Second

	// init accumulators
	dtos := []Task{}
	offset := 0
	errcnt := 0

	for {
		reqDtos, ri, res, err := s.SelectWithOffset(query, offset)
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

// SelectWithOffset request tasks by query & offset, list them also returning list metadata, the actual response, or an error
func (s *TasksService) SelectWithOffset(query string, offset int) ([]Task, ResultInfo, *Response, error) {
	var tld TaskListDTO

	uri := TasksQueryURI(query, offset)
	res, err := s.client.get(uri, &tld)

	ts := []Task{}
	for _, t := range tld.Tasks {
		ts = append(ts, t)
	}
	return ts, tld.Resultinfo, res, err
}

// Find Get the status of a task.
func (s *TasksService) Find(t TaskID) (Task, *Response, error) {
	var tv Task
	res, err := s.client.get(t.URI(), &tv)
	return tv, res, err
}

// FindResult requests
func (s *TasksService) FindResult(t TaskID) (*Response, error) {
	return s.client.GetResultByURI(t.ResultURI())
}

// FindResultByTask  requests a task by the provided task's result uri
func (s *TasksService) FindResultByTask(t Task) (*Response, error) {
	return s.client.GetResultByURI(t.ResultURI)
}

// Delete requests deletions
func (s *TasksService) Delete(t TaskID) (*Response, error) {
	return s.client.delete(t.URI(), nil)
}
