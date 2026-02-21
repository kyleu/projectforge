package mvc

import "{{{ .Package }}}/app/util"

type TransitionType string

const (
	TransitionStay    TransitionType = "stay"
	TransitionPush    TransitionType = "push"
	TransitionReplace TransitionType = "replace"
	TransitionPop     TransitionType = "pop"
	TransitionQuit    TransitionType = "quit"
	TransitionRoute   TransitionType = "route"
)

type Transition struct {
	Type   TransitionType
	Screen string
	Data   util.ValueMap
}

func Stay() Transition {
	return Transition{Type: TransitionStay}
}

func Push(screen string, data util.ValueMap) Transition {
	return Transition{Type: TransitionPush, Screen: screen, Data: data}
}

func Replace(screen string, data util.ValueMap) Transition {
	return Transition{Type: TransitionReplace, Screen: screen, Data: data}
}

func Pop() Transition {
	return Transition{Type: TransitionPop}
}

func Quit() Transition {
	return Transition{Type: TransitionQuit}
}

func Route(screen string, data util.ValueMap) Transition {
	return Transition{Type: TransitionRoute, Screen: screen, Data: data}
}
