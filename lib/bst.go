package lib

type BSTNode struct {
	Key int
	Val RX_Token

	father *BSTNode

	left  *BSTNode
	right *BSTNode

	extraProperties TableRow
}

type BST struct {
	root *BSTNode
}

type TableRow struct {
	nullable  bool
	firstpos  []int
	lastpos   []int
	followpos []int
	simbol    string
}

func (b *BSTNode) IsNullable() bool {
	return b.extraProperties.nullable
}

func (b *BSTNode) IsLeaf() bool {
	return b.left == nil && b.right == nil
}

func (b *BST) Insert(n *BSTNode) {
	if b.root == nil {
		b.root = n
		return
	}
	b.root.insert(n)
}

func (b *BSTNode) insert(n *BSTNode) *BSTNode {
	if b == nil {
		return n
	}

	if n.Val.operator != nil && *n.Val.operator == OR {
		b.right = b.right.insert(n)
		return b
	}

	if b.Val.operator != nil && *b.Val.operator == OR && n.Val.value != nil {
		if b.left == nil {
			b.left = n
			return b
		}
		if b.right == nil {
			b.right = n
			return b
		}
	}

	if n.Val.operator != nil || b.Val.operator == nil {
		b.left = b.left.insert(n)
		return b
	}

	if b.right == nil && *b.Val.operator != ZERO_OR_MANY {
		b.right = b.right.insert(n)
	} else {
		b.left = b.left.insert(n)
	}

	return b
}

func (b *BST) List() []*BSTNode {
	if b.root == nil {
		return nil
	}

	// FIXME: This operation changes the tree everytime because the values are references!
	son := b.root
	for son.left != nil {
		son = son.left
	}

	result := []*BSTNode{}

	for son != nil {
		result = append(result, son)

		if son.father != nil {
			brother := son.father.right
			if brother != nil {
				result = append(result, son.father.right)
			}
		}

		son = son.father

	}

	return result
}

func (b *BST) Insertion(postfix []RX_Token) {
	var stack Stack[BSTNode]
	for i, v := range postfix {
		node := BSTNode{Key: i, Val: v}

		if v.IsOperator() {
			op := *v.GetOperator()
			switch op {
			case AND, OR:
				right := stack.Pop().GetValue()
				left := stack.Pop().GetValue()

				node.left = &left
				node.right = &right

				left.father = &node
				right.father = &node

			case ZERO_OR_MANY:
				left := stack.Pop().GetValue()

				node.left = &left
				left.father = &node
			}
		}

		stack.Push(node)
	}

	root := stack.Pop().GetValue()
	b.root = &root
}

func ConvertTreeToTable(nodes []*BSTNode) []*TableRow {
	table := []*TableRow{}

	for _, node := range nodes {
		if node.IsLeaf() {
			nullable := !(node.Val.IsValue() && node.Val.GetValue().HasValue())
			firstPos := []int{node.Key}
			lastPos := []int{node.Key}

			var simbol string
			if !nullable {
				simbol = string(node.Val.GetValue().GetValue())
			}

			row := TableRow{
				nullable: nullable, firstpos: firstPos, lastpos: lastPos, simbol: simbol,
			}

			node.extraProperties = row
			table = append(table, &row)
		} else {
			op := *node.Val.GetOperator()
			switch op {
			case AND:
				nullable := node.left.IsNullable() && node.right.IsNullable()
				firstPos := node.left.extraProperties.firstpos
				if node.left.IsNullable() {
					firstPos = append(firstPos, node.right.extraProperties.firstpos...)
				}

				lastPos := node.right.extraProperties.lastpos
				if node.right.IsNullable() {
					lastPos = append(lastPos, node.left.extraProperties.lastpos...)
				}
				table = append(table, &TableRow{nullable: nullable, firstpos: firstPos, lastpos: lastPos})

				for i := range node.left.extraProperties.lastpos {
					node_i := nodes[i]
					if node_i.IsLeaf() {
						node_i.extraProperties.followpos = append(node_i.extraProperties.followpos, node.right.extraProperties.firstpos...)
					}
				}

			case OR:
				nullable := node.left.IsNullable() || node.right.IsNullable()
				firstPos := node.left.extraProperties.firstpos
				firstPos = append(firstPos, node.right.extraProperties.firstpos...)

				lastPos := node.right.extraProperties.lastpos
				lastPos = append(lastPos, node.left.extraProperties.lastpos...)

				table = append(table, &TableRow{nullable: nullable, firstpos: firstPos, lastpos: lastPos})

			case ZERO_OR_MANY:
				nullable := true
				firstPos := node.left.extraProperties.firstpos

				lastPos := node.left.extraProperties.lastpos

				table = append(table, &TableRow{nullable: nullable, firstpos: firstPos, lastpos: lastPos})

				for i := range lastPos {
					node_i := nodes[i]
					if node_i.IsLeaf() {
						node_i.extraProperties.followpos = append(node_i.extraProperties.followpos, firstPos...)
					}
				}
			}
		}
	}

	// sets Leaf i first
	// for i, v := range nodes {
	// 	newRow := new(TableRow)

	// 	if v.Val.value != nil && v.Val.value.HasValue() {
	// 		// nullable
	// 		newRow.nullable = false

	// 		// firstpos
	// 		newRow.firtspos = append(newRow.firtspos, i)

	// 		// lastpos
	// 		newRow.lastpos = append(newRow.lastpos, i)

	// 		// simbol
	// 		newRow.simbol = string(v.Val.value.GetValue())

	// 	} else if v.Val.value != nil && !v.Val.value.HasValue() {
	// 		newRow.nullable = true
	// 	} else if *v.Val.operator == AND {
	// 		var c1 int

	// 		if nodes[i-1].Val.operator != nil {
	// 			if *nodes[i-1].Val.operator == OR {
	// 				c1 = 4
	// 			} else {
	// 				c1 = 2
	// 			}
	// 		} else {
	// 			c1 = 2
	// 		}

	// 		//nullable
	// 		newRow.nullable = table[i-c1].nullable == true && table[i-1].nullable == true

	// 		// firstpos
	// 		if table[i-2].nullable == true {
	// 			union_slice := append(table[i-c1].firtspos, table[i-1].firtspos...)
	// 			newRow.firtspos = append(newRow.firtspos, union_slice...)
	// 		} else {
	// 			newRow.firtspos = append(newRow.firtspos, table[i-c1].firtspos...)
	// 		}

	// 		// lastpos
	// 		if table[i-1].nullable == true {
	// 			union_slice := append(table[i-c1].lastpos, table[i-1].lastpos...)
	// 			newRow.lastpos = append(newRow.lastpos, union_slice...)
	// 		} else {
	// 			newRow.lastpos = append(newRow.lastpos, table[i-1].lastpos...)
	// 		}

	// 		// followpos
	// 		for _, pos := range table[i-c1].lastpos {
	// 			table[pos].followpos = append(table[pos].followpos, table[i-1].firtspos...)
	// 		}

	// 	} else if *v.Val.operator == OR {
	// 		// nullable
	// 		newRow.nullable = table[i-2].nullable == true || table[i-1].nullable == true

	// 		// firtspos
	// 		union_slice := append(table[i-2].firtspos, table[i-1].firtspos...)
	// 		newRow.firtspos = append(newRow.firtspos, union_slice...)

	// 		// lastpos
	// 		union_slice = append(table[i-2].lastpos, table[i-1].lastpos...)
	// 		newRow.lastpos = append(newRow.lastpos, union_slice...)

	// 	} else {
	// 		// nullable
	// 		newRow.nullable = true

	// 		// firstpos
	// 		newRow.firtspos = append(newRow.firtspos, table[i-1].firtspos...)

	// 		// lastpos
	// 		newRow.lastpos = append(newRow.lastpos, table[i-1].lastpos...)

	// 		// followpos
	// 		for _, pos := range newRow.lastpos {
	// 			table[pos].followpos = append(table[pos].followpos, newRow.firtspos...)
	// 		}
	// 	}

	// 	table = append(table, newRow)
	// }

	return table
}
