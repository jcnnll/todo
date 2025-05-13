package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/jcnnll/todo"
)

var todoFileName string

func main() {
	if os.Getenv("TODO_FILENAME") == "" {
		os.Setenv("TODO_FILENAME", ".todo.json")
	}
	todoFileName = os.Getenv("TODO_FILENAME")

	task := flag.String("task", "", "Task to be added to the ToDo list")
	add := flag.Bool("add", false, "Add task to the ToDo list")
	del := flag.Int("del", 0, "Delete a task from the ToDo list")
	list := flag.Bool("list", false, "List all tasks")
	verbose := flag.Bool("verbose", false, "View in verbose mode")
	complete := flag.Int("complete", 0, "Item to be marked complete")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "%s tool. Developed for education purposes\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Copyright 2025\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Usage information\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	l := &todo.List{}

	// load todos from file
	if err := l.Get(todoFileName); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	switch {
	case *list:
		fmt.Print(l)

	case *verbose:
		formatted := ""
		for k, t := range *l {
			prefix := "  "
			if t.Done {
				prefix = "X "
			}

			formatted += fmt.Sprintf("%s%d: %s\t%v\n", prefix, k+1, t.Task, t.CreatedAt.Format("2006-01-02"))
		}
		fmt.Print(formatted)

	case *add:
		t, err := getTask(os.Stdin, flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		l.Add(t)
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

	case *del > 0:
		if err := l.Delete(*del); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// save new list
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

	case *complete > 0:
		if err := l.Complete(*complete); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// save new list
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

	case *task != "":
		l.Add(*task)

		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

	default:
		fmt.Fprintln(os.Stderr, "Invalid option")
		os.Exit(1)
	}
}

func getTask(r io.Reader, args ...string) (string, error) {
	if len(args) > 0 {
		return strings.Join(args, " "), nil
	}

	s := bufio.NewScanner(r)
	s.Scan()
	if err := s.Err(); err != nil {
		return "", err
	}

	if len(s.Text()) == 0 {
		return "", fmt.Errorf("Task cannot be blank")
	}

	return s.Text(), nil
}
