package lib

type TransitionInput struct {
	state string
	input string
}

type AFD struct {
	initialState     string
	transitions      map[TransitionInput]string
	acceptanceStates []string
}
