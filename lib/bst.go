package lib

type BSTNode struct {
    Key int 
    Val RX_Token
    left *BSTNode
    right *BSTNode
}

type BST struct {
    root *BSTNode
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
        if b.right == nil {
            b.right = b.right.insert(n)
        } else {
            b.left = b.left.insert(n)
        }
    }
    return b
}

func (b *BST) List() []*BSTNode {
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

func (b *BST) insertion(postfix []RX_Token){
    for i, v := range postfix {
        node := &BSTNode{Key: i, Val: v}
        b.Insert(node)
    }
}
