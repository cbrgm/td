package main

import (
	"sort"
)

type Todolist struct {
	Todos []*Todo `json:"todos"`
}

func (l *Todolist) Pop() {
	if len(l.Todos) > 0 {
		l.Todos = removeItem(l.Todos, 0)
	}
}

func removeItem(s []*Todo, index int) []*Todo {
	return append(s[:index], s[index+1:]...)
}

func (l *Todolist) First() (todo *Todo) {
	if len(l.Todos) > 0 {
		return l.Todos[0]
	}
	return newTodo("Empty!", 0)
}

func (l *Todolist) Push(todo *Todo) {
	if todo.Priority == 0 {
		todo.Priority = l.highestPriority()
	}
	l.Todos = append(l.Todos, todo)

	By(func(p1, p2 *Todo) bool {
		return p1.Priority > p2.Priority
	}).Sort(l.Todos)
}

func (l *Todolist) Clear() {
	l.Todos = []*Todo{}
}

func (l *Todolist) highestPriority() int {
	highest := 0
	for _, todo := range l.Todos {
		if todo.Priority >= highest {
			highest = todo.Priority + 10
		}
	}
	return highest
}

type Todo struct {
	Title    string `json:"title"`
	Priority int    `json:"priority"`
}

func newTodo(title string, priority int) *Todo {
	return &Todo{
		Title:    title,
		Priority: priority,
	}
}

type By func(p1, p2 *Todo) bool

func (by By) Sort(vars []*Todo) {
	ps := &resultSorter{
		vars: vars,
		by:   by,
	}
	sort.Sort(ps)
}

type resultSorter struct {
	vars []*Todo
	by   func(p1, p2 *Todo) bool
}

func (s *resultSorter) Len() int {
	return len(s.vars)
}

func (s *resultSorter) Swap(i, j int) {
	s.vars[i], s.vars[j] = s.vars[j], s.vars[i]
}

func (s *resultSorter) Less(i, j int) bool {
	return s.by(s.vars[i], s.vars[j])
}
