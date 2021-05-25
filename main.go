package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	cmdClear = flag.NewFlagSet("convert", flag.ExitOnError)
	cmdList  = flag.NewFlagSet("colors", flag.ExitOnError)
)

var usage = `Usage: td [command] [options...] string
	Command: pop
	Description: Finish the current todo
	Command: clear
	Description: Clears the todolist
	Command: ls
	Description: Shows your todolist
	Command: rm
	Description: Removes a single todo from the list
`

func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, usage)
	}

	if len(os.Args) <= 1 {
		listFirstCmd()
		return
	}
	if len(os.Args) < 2 {
		usageAndExit("no command selected")
	}
	switch os.Args[1] {
	case "clear":
		clearCmd()
	case "ls":
		listCmd()
	case "pop":
		popCmd()
	default:
		addCmd(os.Args[1], 0)
	}
}

func listCmd() {
	list, err := Open()
	if err != nil {
		usageAndExit(err.Error())
	}

	if len(list.Todos) == 0 {
		usageAndExit("List is empty!")
	}

	for _, item := range list.Todos {
		fmt.Println(item.Title)
	}
}

func addCmd(todo string, priority int) {
	list, err := Open()
	if err != nil {
		usageAndExit(err.Error())
	}

	list.Push(newTodo(todo, priority))

	if err = ToFile(list); err != nil {
		usageAndExit(err.Error())
	}
}

func listFirstCmd() {
	list, err := Open()
	if err != nil {
		usageAndExit(err.Error())
	}

	first := list.First()
	fmt.Println(first.Title)
}

func popCmd() {
	list, err := Open()
	if err != nil {
		usageAndExit(err.Error())
	}

	list.Pop()

	if err = ToFile(list); err != nil {
		usageAndExit(err.Error())
	}
}

func clearCmd() {
	list, err := Open()
	if err != nil {
		usageAndExit(err.Error())
	}

	list.Clear()

	if err = ToFile(list); err != nil {
		usageAndExit(err.Error())
	}
}

func Open() (*Todolist, error) {
	var list *Todolist
	var err error
	if IsConfigExists() {
		list, err = FromFile()
		if err != nil {
			return list, err
		}
	} else {
		list = &Todolist{
			Todos: []*Todo{},
		}
		if err = ToFile(list); err != nil {
			return list, err
		}
	}
	return list, err
}

func usageAndExit(msg string) {
	if msg != "" {
		fmt.Fprintf(os.Stderr, msg)
		fmt.Fprintf(os.Stderr, "\n\n")
	}
	flag.Usage()
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(1)
}
