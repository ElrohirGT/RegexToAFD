package main

import (
	"fmt"
	"log"
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
	stateEvaluationRecord := lib.Set[lib.AFDState]{}
	pairStatesEvaluationTable := lib.AFDStateTable[bool]{}

	// Now fuse all non distinguishable pairs...
	for _, aState := range states {
		nameBuilder := strings.Builder{}
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

			if bState != aState {
				stateEvaluationRecord.Add(bState)
				nameBuilder.WriteString(bState)
			}
		}

		if stateEvaluationRecord.Add(aState) {
			stateName := nameBuilder.String()
			log.Default().Printf("Creating out state: %s", stateName)
			outStates = append(outStates, stateName)
			outAFD.Transitions[stateName] = map[lib.AlphabetInput]lib.AFDState{}
		}
	}

	return outAFD
}
