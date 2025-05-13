package todo_test

import (
	"os"
	"testing"

	"github.com/jcnnll/todo"
)

func TestAdd(t *testing.T) {
	l := todo.List{}

	taskName := "New task"

	l.Add(taskName)

	if l[0].Task != taskName {
		t.Errorf("Expected %q, got %q", taskName, l[0].Task)
	}
}

func TestComplete(t *testing.T) {
	l := todo.List{}

	taskName := "New task"
	l.Add(taskName)

	if l[0].Task != taskName {
		t.Errorf("Expected %q, got %q", taskName, l[0].Task)
	}

	if l[0].Done {
		t.Errorf("New task should not be completed")
	}

	l.Complete(1)

	if !l[0].Done {
		t.Errorf("New task should br completed")
	}
}

func TestDelete(t *testing.T) {
	l := todo.List{}

	tasks := []string{
		"Task 1",
		"Task 2",
		"Task 3",
	}

	for _, v := range tasks {
		l.Add(v)
	}

	for i := range tasks {
		if l[i].Task != tasks[i] {
			t.Errorf("Expected %q, got %q", tasks[i], l[i].Task)
		}
	}

	l.Delete(2)

	if len(l) != 2 {
		t.Errorf("Expected list to be %d, got %d", 2, len(l))
	}

	if l[1].Task != tasks[2] {
		t.Errorf("expected %q, git %q", tasks[2], l[1].Task)
	}
}

func TestSaveGet(t *testing.T) {
	l1 := todo.List{}
	l2 := todo.List{}

	taskName := "New task"

	l1.Add(taskName)

	if l1[0].Task != taskName {
		t.Errorf("Expected %q, got %q", taskName, l1[0].Task)
	}

	tf, err := os.CreateTemp("", "")

	if err != nil {
		t.Fatalf("Error creating temp file: %s", err)
	}

	defer os.Remove(tf.Name())

	if err := l1.Save(tf.Name()); err != nil {
		t.Fatalf("Error saving list to file: %s", err)
	}

	if err := l2.Get(tf.Name()); err != nil {
		t.Fatalf("Error getting list from file: %s", err)
	}

	if l1[0].Task != l2[0].Task {
		t.Errorf("Task %q should match %q", l1[0].Task, l2[0].Task)
	}
}
