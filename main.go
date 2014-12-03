package main

import "fmt"

// Tree is a Left-Leaning Red-Black (LLRB) implementation of 2-3 trees
type LLRB struct {
	root  *Node
}
type Node struct {
	Item        []byte
	Left, Right *Node // Pointers to left and right child nodes
	Depth       int
}
type slce []byte

func (x slce) Less(y []byte) bool {
	return len(x) < len(y)
}

//
func less(x, y []byte) bool {
	if &x[0] == &pinf[0] {
		return false
	}
	if &x[0] == &ninf[0] {
		return true
	}
	return slce(x).Less(y)
}

// Inf returns an Item that is "bigger than" any other item, if sign is positive.
// Otherwise it returns an Item that is "smaller than" any other item.
func Inf(sign int) []byte {
	if sign == 0 {
		panic("sign")
	}
	if sign > 0 {
		return pinf
	}
	return ninf
}

var (
	ninf = []byte("nI")
	pinf = []byte("pI")
)

type nInf []byte

func (nInf) Less([]byte) bool {
	return true
}

type pInf []byte

func (pInf) Less([]byte) bool {
	return false
}

// New() allocates a new tree
func New() *LLRB {
	return &LLRB{}
}

// Min returns the minimum element in the tree.
func (t *LLRB) Min() []byte {
	h := t.root
	if h == nil {
		return nil
	}
	for h.Left != nil {
		h = h.Left
	}
	return h.Item
}

// Max returns the maximum element in the tree.
func (t *LLRB) Max() []byte {
	h := t.root
	if h == nil {
		return nil
	}
	for h.Right != nil {
		h = h.Right
	}
	return h.Item
}
func (t *LLRB) InsertNoReplaceBulk(items ...[]byte) {
	for _, i := range items {
		t.InsertNoReplace(0, i)
	}
}

// InsertNoReplace inserts item into the tree. If an existing
// element has the same order, both elements remain in the tree.
func (t *LLRB) InsertNoReplace(off int, item []byte) {
	t.root = t.insertNoReplace(off, t.root, item)
	t.root.Depth |= 1
}
func (t *LLRB) insertNoReplace(off int, h *Node, item []byte) *Node {
	if h == nil {
		return newNode(item)
	}
	h = walkDownRot23(h)
	d := h.Depth >> 8
	if off <= d {
		h.Left = t.insertNoReplace(off, h.Left, item)
//		h.Depth += off >> 8
	} else if off-d < len(h.Item) {
		h.Right = t.insertNoReplace(0, h.Right, h.Item[off-d:])
		h.Right = t.insertNoReplace(0, h.Right, item)
		h.Item = h.Item[:off-d]
//		h.Depth += len(item) >> 8
	} else {
		h.Right = t.insertNoReplace(off-d-len(h.Item), h.Right, item)
	}

	return walkUpRot23(h)
}

// Rotation driver routines for 2-3 algorithm
func walkDownRot23(h *Node) *Node { return h }
func walkUpRot23(h *Node) *Node {
	if isRed(h.Right) && !isRed(h.Left) {
		h = rotateLeft(h)
	}
	if isRed(h.Left) && isRed(h.Left.Left) {
		h = rotateRight(h)
	}
	if isRed(h.Left) && isRed(h.Right) {
		flip(h)
	}
	return h
}

func (t *LLRB) Delete(off, size int) {
	t.delete(off, size, t.root)
}
func (t *LLRB) delete(off, size int, h *Node) {
	d := h.Depth >> 8
	if off <= d {
		if !isRed(h.Left) && !isRed(h.Left.Left) {
			h = moveRedLeft(h)
		}
		t.delete(off,size,h.Left)
	} else if off-d < len(h.Item) {
		print("internal deletion")
	} else {
		if isRed(h.Left) {
			h = rotateRight(h)
		}

		if h.Right != nil && !isRed(h.Right) && !isRed(h.Right.Left) {
			h = moveRedRight(h)
		}
		t.delete(off-d-len(h.Item), size, h.Right)
	}
}

// Internal node manipulation routines
func newNode(item []byte) *Node { return &Node{Item: item} }
func isRed(h *Node) bool {
	if h == nil {
		return false
	}
	return 0 == (h.Depth & 1)
}
func rotateLeft(h *Node) *Node {
	x := h.Right
	if 1 == (x.Depth & 1) {
		panic("rotating a black link")
	}
	h.Right = x.Left
	x.Left = h
	x.Depth = h.Depth
	h.Depth = 0
	return x
}
func rotateRight(h *Node) *Node {
	x := h.Left
	if 1 == (x.Depth & 1) {
		panic("rotating a black link")
	}
	h.Left = x.Right
	x.Right = h
	x.Depth = h.Depth
	h.Depth = 0
	return x
}

func flip(h *Node) {
	h.Depth ^= 1
	h.Left.Depth ^= 1
	h.Right.Depth ^= 1
}

func moveRedLeft(h *Node) *Node {
	flip(h)
	if isRed(h.Right.Left) {
		h.Right = rotateRight(h.Right)
		h = rotateLeft(h)
		flip(h)
	}
	return h
}

func moveRedRight(h *Node) *Node {
	flip(h)
	if isRed(h.Left.Left) {
		h = rotateRight(h)
		flip(h)
	}
	return h
}
func fixUp(h *Node) *Node {
	if isRed(h.Right) {
		h = rotateLeft(h)
	}
	if isRed(h.Left) && isRed(h.Left.Left) {
		h = rotateRight(h)
	}
	if isRed(h.Left) && isRed(h.Right) {
		flip(h)
	}
	return h
}

func getbyte(n *Node, i int) byte {
	return getbyte_(n, i)
}

func getbyte_(n *Node, i int) byte {
	d := n.Depth >> 8
	if i < d && n.Left != nil {
		return getbyte_(n.Left, i)
	}
	i -= d
	if i < len(n.Item) {
		return n.Item[i]
	}
	i -= len(n.Item)
	return getbyte_(n.Right, i)
}
func dmp(n *Node) {
	if n.Left != nil {
		dmp(n.Left)
	}
	fmt.Print(string(n.Item))
	if n.Right != nil {
		dmp(n.Right)
	}
}
func dump(n *Node, depth int) {
	if n.Left != nil {
		dump(n.Left, depth+1)
	}
	fmt.Println(depth, "[", (n.Depth >> 8), "|", (n.Depth & 1), "|", string(n.Item), "]")
	if n.Right != nil {
		dump(n.Right, depth+1)
	}
}

func fixdepth(n *Node) (int, int) {
	new := n.Depth & 1
	rd := 0
	ld := 0
	fl := 0
	fr := 0
	if n.Right != nil {
		rd, fr = fixdepth(n.Right)
	}
	if n.Left != nil {
		ld, fl = fixdepth(n.Left)
	}

	new += ld

	if new != n.Depth {
		n.Depth = new
		return ld + rd + (len(n.Item) << 8), fl + fr + 1
	} else {
		return ld + rd + (len(n.Item) << 8), fl + fr
	}

}

func main() {
	tree := New()
	sli1 := []byte("h||ello")
	sli2 := []byte("to world")
	sli3 := []byte("and you")
	sli4 := []byte("thats cool")
	sli5 := []byte("nice")
	sli6 := []byte("gopher")

	tree.InsertNoReplaceBulk(sli1, sli2, sli3, sli4, sli5, sli6)

	_, xed := fixdepth(tree.root)

	fmt.Println("Fixed:", xed)

	dump(tree.root, 0)

	for i:=0; i < 5; i++ {

	tree.InsertNoReplace(37, []byte("just a test"))

	_, xed := fixdepth(tree.root)

	fmt.Println("\nFixed:", xed)

//	dump(tree.root,0)
	dmp(tree.root)
	}

		for i:=0; i < 5; i++ {

	tree.Delete(0,1)

	_, xed := fixdepth(tree.root)

	fmt.Println("\nFixed:", xed)

//	dump(tree.root,0)
	dmp(tree.root)
	}
}
