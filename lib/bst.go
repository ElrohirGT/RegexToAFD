package lib

type BSTNode struct {
    Key int 
    Val interface{}
    left *BSTNode
    right *BSTNode
}

type BST struct {
    root *BSTNode
}

func (b *BST) Insert(n *BSTNode) {
    b.root = b.root.insert(n)
}

func (b *BSTNode) insert(n *BSTNode) {
    if b == nil {
        reutrn n
    }
    if n.Key < b.Key {
        b.left = b.left.insert(n)
    } else if n.Key > b.key {
        b.right = b.right.insert(n)
    }
    return b
}

func (b *BST) List() []*Node {
    ret := []*BSTNode{}
    stack := []*BSTNode{}
    current := b.root
    for len(stack) != 0 || current != nil {
        if current != nil {
            stack = append(stack, current)
            current = current.left
        } else {
            current = stack[len(stack)-1]
            stack = stack[:len(stack)-1]
            ret = append(ret, current)
            current = current.right
        }
    }
    return ret
}


