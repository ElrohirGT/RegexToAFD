package lib

import "testing"

func TestConvertFromTableToAFD(t *testing.T) {
	// Definir tabla de ejemplo
	table := []TableRow{
		{false, []int{0}, []int{0}, []int{1}, "a"},
		{false, []int{1}, []int{1}, []int{}, "b"},
		{false, []int{0}, []int{1}, []int{}, ""},
	}

	afd := convertFromTableToAFD(table)

	// Verificar estado inicial
	expectedInitial := "0"
	if afd.InitialState != expectedInitial {
		t.Errorf("Expected initial state %s, got %s", expectedInitial, afd.InitialState)
	}

	// Verificar transiciones
	expectedTransitions := map[string]map[string]string{
		"0": {"a": "1"},
	}

	for state, transitions := range expectedTransitions {
		for input, expectedNextState := range transitions {
			if afd.Transitions[state][input] != expectedNextState {
				t.Errorf("Expected transition (%s, %s) -> %s, got %s",
					state, input, expectedNextState, afd.Transitions[state][input])
			}
		}
	}

	// Verificar estados de aceptación
	expectedFinalState := "1"
	if !afd.AcceptanceStates.Contains(expectedFinalState) {
		t.Errorf("Expected final state %s to be in AcceptanceStates", expectedFinalState)
	}
}
