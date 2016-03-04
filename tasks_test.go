package udnssdk

import (
	"testing"
)

func Test_ListTasks(t *testing.T) {
	tasks, err := testClient.Tasks.ListAllTasks("")
	t.Logf("Tasks: %+v \n", tasks)
	if err != nil {
		t.Fatal(err)
	}
}
