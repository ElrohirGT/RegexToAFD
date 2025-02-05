package lib

import "testing"

func TestBST(t *testing.T) {
    r1, r2, r3 := '1', 'a', 'c'
    nodes := []*BSTNode{
        {Key: 1, Val: &r1},
        {Key: 2, Val: &r2},
        {Key: 3, Val: &r3},
    }

    // Crear el árbol
    tree := new(BST)

    // Insertar nodos
    for _, node := range nodes {
        tree.Insert(node)
    }

    // Obtener los nodos ordenados en in-order traversal
    got := tree.List()

    // Nodos esperados después del in-order traversal
    expectedKeys := []int{1, 2, 3}
    expectedVals := []rune{'1', 'a', 'c'}

    // Verificar que la cantidad de nodos es correcta
    if len(got) != len(expectedKeys) {
        t.Fatalf("Número incorrecto de nodos. Esperado %d, pero obtuvo %d", len(expectedKeys), len(got))
    }

    // Verificar cada nodo
    for i, node := range got {
        if node.Key != expectedKeys[i] || *node.Val != expectedVals[i] {
            t.Errorf("Nodo incorrecto en posición %d: esperado (%d, %c) pero obtuvo (%d, %c)", 
                i, expectedKeys[i], expectedVals[i], node.Key, *node.Val)
        }
    }
}



