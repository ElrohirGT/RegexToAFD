package lib

type AFDState = string
type AlphabetInput = string

type TransitionInput struct {
	State AFDState
	Input AlphabetInput
}

type AFD struct {
	InitialState AFDState
	// A dictionary that contains a bunch of states.
	// Each AFD state has another dictionary associated with it.
	// Each key on this second dictionary represents an input from the alphabet,
	// and the value is the new State the automata should transition.
	Transitions      map[AFDState]map[AlphabetInput]AFDState
	AcceptanceStates Set[AFDState]
}
