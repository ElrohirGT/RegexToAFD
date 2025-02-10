package main

import (
	"fmt"
	"slices"
	"strings"

	"github.com/ElrohirGT/RegexToAFD/lib"
)

func getNonAcceptanceStates(afd lib.AFD) []string {
	nonAcceptedStates := []string{}
	for state := range afd.Transitions {
		_, found := afd.AcceptanceStates[state]

		if !found {
			nonAcceptedStates = append(nonAcceptedStates, state)
		}
	}

	return nonAcceptedStates
}

func MinimizeAFD(afd lib.AFD) lib.AFD {

	// acceptedStates := afd.AcceptanceStates
	// nonAcceptedStates := getNonAcceptanceStates(afd)
	stateEquivalenceTable := lib.AFDStateTable[lib.AFDPairType]{}

	states := []lib.AFDState{}
	for state := range afd.Transitions {
		states = append(states, state)
	}
	slices.Sort(states)

	// First obtain all pairs that are distinguishable between each other...
	for _, aState := range states {
		for _, bState := range states {
			if stateEquivalenceTable.PairAlreadyExists(&aState, &bState) {
				continue
			}

			afd.MarkIfDistinguishable(&aState, &bState, &stateEquivalenceTable)
		}
	}

	outAFD := lib.AFD{
		Transitions: make(map[lib.AFDState]map[lib.AlphabetInput]lib.AFDState),
	}
	outStates := []string{}
	statesAlreadyFused := lib.Set[lib.AFDState]{}
	pairStatesEvaluationTable := lib.AFDStateTable[bool]{}

	// Now fuse all non distinguishable pairs...
	for _, aState := range states {
		if statesAlreadyFused.Contains(aState) {
			continue
		}

		nameBuilder := strings.Builder{}
		nameBuilder.WriteString("|")
		nameBuilder.WriteString(aState)

		equivalentStates := []lib.AFDState{}
		for _, bState := range states {
			pairType, found := stateEquivalenceTable.Get(&aState, &bState)
			if !found {
				panic(fmt.Sprintf("Pair was not found in table! (`%s`,`%s`)", aState, bState))
			}

			if pairType == lib.DISTINCT {
				continue
			}

			equivalentStates = append(equivalentStates, bState)
		}

		eqStatesCombinations := lib.Combinations(equivalentStates, 2)

		for _, states := range eqStatesCombinations {
			if len(states) != 2 {
				break
			}

			aState := states[0]
			bState := states[1]

			if alreadyEvaluated, found := pairStatesEvaluationTable.Get(&aState, &bState); found && alreadyEvaluated {
				continue
			} else {
				pairStatesEvaluationTable.AddOrUpdate(aState, bState, true)
			}

			if bState != aState && !statesAlreadyFused.Contains(bState) {
				statesAlreadyFused.Add(bState)
				nameBuilder.WriteString("|")
				nameBuilder.WriteString(bState)
			}
		}

		if statesAlreadyFused.Add(aState) {
			nameBuilder.WriteString("|")
			stateName := nameBuilder.String()
			outStates = append(outStates, stateName)
			outAFD.Transitions[stateName] = make(map[lib.AlphabetInput]lib.AFDState)
		}
	}

	// Now get all transitions...
	newAFDStates := outAFD.GetAllStates()
	for _, combinedStates := range newAFDStates {
		parts := strings.Split(combinedStates, "|")

		for idx := 1; idx < len(parts); idx += 2 {
			state := parts[idx]
			originalStateTransitions := afd.Transitions[state]

			for input, outState := range originalStateTransitions {
				// Find the state from the OutAFD that contains the outState
				for _, newState := range newAFDStates {
					if strings.Contains(newState, outState) {
						outAFD.Transitions[combinedStates][input] = newState
						break
					}
				}
			}
		}
	}

	// Now we find the initial state...
	for _, combinedStates := range newAFDStates {
		initialStateFound := false
		parts := strings.Split(combinedStates, "|")

		for idx := 1; idx < len(parts); idx += 2 {
			state := parts[idx]

			if state == afd.InitialState {
				outAFD.InitialState = combinedStates
				initialStateFound = true
				break
			}
		}

		// Initial state found...
		if initialStateFound {
			break
		}
	}

	// Finding accepted states...
	outAFD.AcceptanceStates = lib.Set[lib.AFDState]{}
	for _, combinedStates := range newAFDStates {
		parts := strings.Split(combinedStates, "|")

		for idx := 1; idx < len(parts); idx += 2 {
			state := parts[idx]

			if afd.AcceptanceStates.Contains(state) {
				outAFD.AcceptanceStates.Add(combinedStates)
				break
			}
		}

	}

	return outAFD
}
