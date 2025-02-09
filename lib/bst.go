package lib

type BSTNode struct {
	Key   int
	Val   RX_Token
	left  *BSTNode
	right *BSTNode
}

type BST struct {
	root *BSTNode
}

type TableRow struct {
    nullable bool
    firtspos []int
    lastpos []int
    followpos []int
    simbol *rune
}

func (b *BST) Insert(n *BSTNode) {
	b.root = b.root.insert(n)
}

func (b *BSTNode) insert(n *BSTNode) *BSTNode {
	if b == nil {
		return n
	}
	if n.Val.value == nil {
		b.left = b.left.insert(n)
	} else {
		if b.right == nil && *b.Val.operator != ZERO_OR_MANY {
		    b.right = b.right.insert(n)
		} else {
		    b.left = b.left.insert(n)
		}
            
	}
	return b
}

func (b *BST) List() []*BSTNode {
    if b.root == nil {
        return nil
    }

    result := []*BSTNode{}
    stack := []*BSTNode{b.root}

    for len(stack) > 0 {
        current := stack[len(stack)-1]
        stack = stack[:len(stack)-1]

        result = append(result, current)

        if current.left != nil {
            stack = append(stack, current.left)
        }
        if current.right != nil {
            stack = append(stack, current.right)
        }
    }

    for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
        result[i], result[j] = result[j], result[i]
    }

    return result
}

func (b *BST) insertion(postfix []RX_Token) {
	for i, v := range postfix {
		node := &BSTNode{Key: i, Val: v}
		b.Insert(node)
	}
}

func convertTreeToTable(nodes []*BSTNode) []*TableRow {
    table := []*TableRow{}

    // sets Leaf i first
    for i, v := range nodes {
        newRow := new(TableRow)

        if v.Val.value != nil && v.Val.value.HasValue() {
            // nullable
            newRow.nullable = false

            // firstpos
            newRow.firtspos = append(newRow.firtspos, i)

            // lastpos
            newRow.lastpos = append(newRow.lastpos, i)

        } else if v.Val.value != nil && !v.Val.value.HasValue() {
            newRow.nullable = true
        } else if *v.Val.operator == AND {
            //nullable
            newRow.nullable = table[i-2].nullable == true && newRow.nullable == true

            // firstpos
            if table[i-2].nullable == true {
                union_slice := append(table[i-2].firtspos, table[i-1].firtspos ...)
                newRow.firtspos = append(newRow.firtspos, union_slice ...)
            } else {
                newRow.firtspos = append(newRow.firtspos, table[i-2].firtspos ...)
            }

            // lastpos
            if table[i-1].nullable == true {
                union_slice := append(table[i-2].lastpos, table[i-1].lastpos ...)
                newRow.lastpos = append(newRow.lastpos, union_slice ...)
            } else {
                newRow.lastpos = append(newRow.lastpos, table[i-1].lastpos ...)
            }

            // followpos
            for _, pos := range table[i-2].lastpos {
                table[pos].followpos = append(table[pos].followpos, table[i-1].firtspos ...)
            }
 
        } else if *v.Val.operator == OR {
            // nullable
            newRow.nullable = table[i-2].nullable == true || newRow.nullable == true 
            
            // firtspos
            union_slice := append(table[i-2].firtspos, table[i-1].firtspos ...)
            newRow.firtspos = append(newRow.firtspos, union_slice ...)

            // lastpos
            union_slice = append(table[i-2].lastpos, table[i-1].lastpos ...)
            newRow.lastpos = append(newRow.lastpos, union_slice ...)

        } else {
            // nullable
            newRow.nullable = true

            // firstpos
            newRow.firtspos = append(newRow.firtspos, table[i-1].firtspos ...)

            // lastpos
            newRow.lastpos = append(newRow.lastpos, table[i-1].lastpos ...)


            // followpos
            for _, pos := range newRow.lastpos {
                table[pos].followpos = append(table[pos].followpos, newRow.firtspos ...)
            }
        }

        table = append(table, newRow)
    }

    return table
}


