package main

import (
	"fmt"
	"log"
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

	// First obtain all pairs that are distinguishable between each other...
	for aState := range afd.Transitions {
		for bState := range afd.Transitions {
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
	stateEvaluationTable := lib.AFDStateTable[bool]{}

	// Now fuse all non distinguishable pairs...
	for aState := range afd.Transitions {
		nameBuilder := strings.Builder{}
		nameBuilder.WriteString(aState)

		for bState := range afd.Transitions {
			if alreadyEvaluated, found := stateEvaluationTable.Get(&aState, &bState); found && alreadyEvaluated {
				continue
			} else {
				stateEvaluationTable.AddOrUpdate(aState, bState, true)
			}

			pairType, found := stateEquivalenceTable.Get(&aState, &bState)
			if !found {
				panic(fmt.Sprintf("Pair was not found in table! (`%s`,`%s`)", aState, bState))
			}

			if pairType == lib.DISTINCT {
				continue
			}

			nameBuilder.WriteString(bState)
		}

		stateName := nameBuilder.String()
		log.Default().Printf("Creating out state: %s", stateName)
		outStates = append(outStates, stateName)
		outAFD.Transitions[stateName] = map[lib.AlphabetInput]lib.AFDState{}
	}

	return outAFD
}
