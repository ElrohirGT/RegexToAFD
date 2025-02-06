package lib

import (
    "testing"
)

// General Test for BST
func TestBST(t *testing.T) {
    // Node Creation
    nodes := []*BSTNode{
        {Key: 1, Val: CreateOperatorToken(AND)},
        {Key: 2, Val: CreateValueToken('a')},
        {Key: 3, Val: CreateValueToken('b')},
    }

    // Creates tree
    tree := new(BST)

    // Insertar nodos
    for _, node := range nodes {
        tree.Insert(node)
    }

    // in-order transverse
    got := tree.List()

    // Expected nodes 
    expectedKeys := []int{3,1,2}
    expectedVals := []RX_Token{CreateValueToken('b'),CreateOperatorToken(AND),CreateValueToken('a')}

    // Verifies total nodes
    if len(got) != len(expectedKeys) {
        t.Fatalf("Número incorrecto de nodos. Esperado %d, pero obtuvo %d", len(expectedKeys), len(got))
    }

    // Verifies each node 
    for i, node := range got {
       if node.Key != expectedKeys[i] {
        t.Errorf("Nodo incorrecto en posición %d: esperado (%d) pero obtuvo (%d)", 
            i, expectedKeys[i], node.Key)
        }

        if node.Val.value != nil && expectedVals[i].value != nil {
            if *node.Val.value != *expectedVals[i].value {
                t.Errorf("Nodo incorrecto en posición %d: esperado (%c) pero obtuvo (%c)", 
                    i, *expectedVals[i].value, *node.Val.value)
            }
        } else if node.Val.operator != nil && expectedVals[i].operator != nil {
            if *node.Val.operator != *expectedVals[i].operator {
                t.Errorf("Nodo incorrecto en posición %d: esperado (%d) pero obtuvo (%d)", 
                    i, *expectedVals[i].operator, *node.Val.operator)
            }
        } else {
            t.Errorf("Nodo incorrecto en posición %d: los tipos de valor no coinciden", i)
        }
    }
}

// Class example
func TestExampleBST(t *testing.T) {
    // Node Creation
    nodes := []*BSTNode{
        {Key: 0, Val: CreateOperatorToken(AND)},
        {Key: 1, Val: CreateValueToken('#')},
        {Key: 2, Val: CreateOperatorToken(AND)},
        {Key: 3, Val: CreateValueToken('b')},
        {Key: 4, Val: CreateOperatorToken(AND)},
        {Key: 5, Val: CreateValueToken('b')},
        {Key: 6, Val: CreateOperatorToken(AND)},
        {Key: 7, Val: CreateValueToken('a')},
        {Key: 8, Val: CreateOperatorToken(ZERO_OR_MANY)},
        {Key: 9, Val: CreateOperatorToken(OR)},
        {Key: 10, Val: CreateValueToken('a')},
        {Key: 11, Val: CreateValueToken('b')},
    }

    // Creates tree
    tree := new(BST)

    // Insertar nodos
    for _, node := range nodes {
        tree.Insert(node)
    }

    // in-order transverse
    got := tree.List()

    // Expected nodes 
    expectedKeys := []int{11,9,10,8,6,7,4,5,2,3,0,1}
    expectedVals := []RX_Token{
        CreateValueToken('b'),
        CreateOperatorToken(OR),
        CreateValueToken('a'),
        CreateOperatorToken(ZERO_OR_MANY),
        CreateOperatorToken(AND),
        CreateValueToken('a'),
        CreateOperatorToken(AND),
        CreateValueToken('b'),
        CreateOperatorToken(AND),
        CreateValueToken('b'),
        CreateOperatorToken(AND),
        CreateValueToken('#'),
    }

    // Verifies total nodes
    if len(got) != len(expectedKeys) {
        t.Fatalf("Número incorrecto de nodos. Esperado %d, pero obtuvo %d", len(expectedKeys), len(got))
    }

    // Verifies each node 
    for i, node := range got {
       if node.Key != expectedKeys[i] {
        t.Errorf("Nodo incorrecto en posición %d: esperado (%d) pero obtuvo (%d)", 
            i, expectedKeys[i], node.Key)
        }

        if node.Val.value != nil && expectedVals[i].value != nil {
            if *node.Val.value != *expectedVals[i].value {
                t.Errorf("Nodo incorrecto en posición %d: esperado (%c) pero obtuvo (%c)", 
                    i, *expectedVals[i].value, *node.Val.value)
            }
        } else if node.Val.operator != nil && expectedVals[i].operator != nil {
            if *node.Val.operator != *expectedVals[i].operator {
                t.Errorf("Nodo incorrecto en posición %d: esperado (%d) pero obtuvo (%d)", 
                    i, *expectedVals[i].operator, *node.Val.operator)
            }
        } else {
            t.Errorf("Nodo incorrecto en posición %d: los tipos de valor no coinciden", i)
        }
    }
}

