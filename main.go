package main

import (
	"fmt"
	"time"
)

type FlowComponent interface {
	downTo(lower any) FlowComponent
}

type OnePath struct {
	Id        string
	Path      chan any
	Converter func(any) any
}

func (component *OnePath) downTo(lower any) FlowComponent {
	onePathLower, ok := lower.(*OnePath)
	if ok {

		go func() {
			for {
				//fmt.Println(component.Id, ":: getting in")
				data := <-component.Path

				//fmt.Println(component.Id, ":: converting")
				convertedData := component.Converter(data)

				//fmt.Println(component.Id, ":: getting out")
				onePathLower.Path <- convertedData
			}
		}()

		return onePathLower
	}

	onePathBranchLower, ok := lower.(*OnePathBranch)
	if ok {

		go func() {
			for {
				//fmt.Println(component.Id, ":: getting in")
				data := <-component.Path

				//fmt.Println(component.Id, ":: converting")
				convertedData := component.Converter(data)

				//fmt.Println(component.Id, ":: getting out")
				onePathBranchLower.Path <- convertedData
			}
		}()

		return onePathBranchLower
	}

	return nil
}

type OnePathBranch struct {
	Id        string
	Path      chan any
	Predicate func(any) (any, bool)
	positive  *OnePath
	negative  *OnePath
}

func (component *OnePathBranch) Positive(lower *OnePath) *OnePath {
	go func() {
		for {
			//fmt.Println(component.Id, ":: getting in")
			data := <-component.Path

			//fmt.Println(component.Id, ":: predicating")
			convertedData, neetToTransmit := component.Predicate(data)

			if neetToTransmit {
				if component.positive != nil {
					//fmt.Println(component.Id, ":: getting out to positive")
					component.positive.Path <- convertedData
				}
			} else {
				if component.negative != nil {
					//fmt.Println(component.Id, ":: getting out to negative")
					component.negative.Path <- convertedData
				}
			}
		}

	}()

	component.positive = lower
	return lower
}

func (component *OnePathBranch) Negative(lower *OnePath) *OnePath {
	component.negative = lower
	return lower
}

func (component *OnePathBranch) downTo(lower any) FlowComponent {
	return nil
}

func main() {
	app := &OnePath{
		Id:   "app",
		Path: make(chan any, 1),
		Converter: func(data any) any {
			return data
		},
	}

	pasrer1 := &OnePath{
		Id:   "parser1",
		Path: make(chan any, 1),
		Converter: func(data any) any {
			return data
		},
	}

	pasrer2 := &OnePath{
		Id:   "parser2",
		Path: make(chan any, 1),
		Converter: func(data any) any {
			return data
		},
	}

	pasrer3 := &OnePath{
		Id:   "parser3",
		Path: make(chan any, 1),
		Converter: func(data any) any {
			return data
		},
	}

	solution := &OnePathBranch{
		Id:   "solution",
		Path: make(chan any, 1),
		Predicate: func(data any) (any, bool) {
			return data, true
		},
	}

	result1 := &OnePath{
		Id:   "result1",
		Path: make(chan any, 1),
		Converter: func(data any) any {
			return data
		},
	}

	result2 := &OnePath{
		Id:   "result2",
		Path: make(chan any, 1),
		Converter: func(data any) any {
			return data
		},
	}

	// app.downTo(pasrer1).downTo(pasrer2).downTo(pasrer3).downTo(result)

	app.
		downTo(pasrer1).
		downTo(pasrer2).
		downTo(pasrer3).
		downTo(solution)

	solution.
		Positive(result1)

	solution.
		Negative(result2)

	for {
		app.Path <- "hello"

		go func() {
			for {
				data := <-result1.Path
				fmt.Println("resul1 ::", data)
			}
		}()

		go func() {
			for {
				data := <-result2.Path
				fmt.Println("result2 ::", data)
			}
		}()

		time.Sleep(1 * time.Second)
	}
}
