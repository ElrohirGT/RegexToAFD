package main

import (
	"fmt"
	"testing"

	"github.com/ElrohirGT/RegexToAFD/lib"
)

func compareAFDs(t *testing.T, expected *lib.AFD, result *lib.AFD) {
	errorMsgFormat := "%s\nExpected: %#v\nResult: %#v"
	if expected.InitialState != result.InitialState {
		t.Fatal(fmt.Sprintf(errorMsgFormat, "The initial state is not the same!", *expected, *result))
	}

	if len(expected.AcceptanceStates) != len(result.AcceptanceStates) {
		t.Fatal(fmt.Sprintf(errorMsgFormat, "The number of acceptance states is not the same!", *expected, *result))
	}

	for state := range result.AcceptanceStates {
		if _, found := expected.AcceptanceStates[state]; !found {
			t.Fatal(fmt.Sprintf(errorMsgFormat,
				fmt.Sprintf("The state `%s` was not found inside the expected acceptance states set!\nExpected: %#v", state, expected.AcceptanceStates),
				*expected,
				*result))
		}
	}

	for state, transitions := range expected.Transitions {
		resultTransitions, found := result.Transitions[state]
		if !found {
			t.Fatal(fmt.Sprintf(errorMsgFormat,
				fmt.Sprintf("Transitions for state `%s` not found in result AFD!", state),
				*expected,
				*result))
		}

		for input, outState := range transitions {
			resultOutState, found := resultTransitions[input]

			if !found {
				t.Fatal(fmt.Sprintf(errorMsgFormat,
					fmt.Sprintf("The input `%s` for state: `%s` was not found on result AFD!", input, state),
					*expected, *result))
			}

			if outState != resultOutState {
				t.Fatal(fmt.Sprintf(errorMsgFormat,
					fmt.Sprintf("The transition from `%s` state with input `%s` doesn't match! (`%s` != `%s`)", state, input, outState, resultOutState),
					*expected, *result))
			}
		}
	}
}

func TestVideoAFD(t *testing.T) {
	originalAFD := lib.AFD{
		InitialState: "1",
		Transitions: map[lib.AFDState]map[lib.AlphabetInput]lib.AFDState{
			"1": {
				"a": "2",
				"b": "4",
			},
			"2": {
				"a": "4",
				"b": "3",
			},
			"3": {
				"a": "3",
				"b": "3",
			},
			"4": {
				"a": "4",
				"b": "5",
			},
			"5": {
				"a": "5",
				"b": "5",
			},
			"6": {
				"a": "6",
				"b": "5",
			},
		},
		AcceptanceStates: lib.Set[lib.AFDState]{"3": struct{}{}, "5": struct{}{}},
	}
	expectedAFD := lib.AFD{
		InitialState: "|1|",
		Transitions: map[lib.AFDState]map[lib.AlphabetInput]lib.AFDState{
			"|1|": {
				"a": "|2|4|6|",
				"b": "|2|4|6|",
			},

			"|2|4|6|": {
				"a": "|2|4|6|",
				"b": "|3|5|",
			},

			"|3|5|": {
				"a": "|3|5|",
				"b": "|3|5|",
			},
		},
		AcceptanceStates: lib.Set[lib.AFDState]{"|3|5|": struct{}{}},
	}

	result := MinimizeAFD(originalAFD)

	compareAFDs(t, &expectedAFD, &result)
}
