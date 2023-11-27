package kvstorage

import "fmt"

type btreeLayer interface {
	GetRootPageNumber() pageNumber
	Get(key []byte) ([]byte, error)
	Set(key []byte, value []byte) error
}

type btreeLayerImpl struct {
	dataAccessLayer dataAccessLayer
	rootPageNumber  pageNumber
	minFillPercent  float64
	maxFillPercent  float64
}

func newBtreeLayer(dataAccessLayer dataAccessLayer, rootPageNumber pageNumber, minFillPercent float64, maxFillPercent float64) (btreeLayer, error) {
	btl := &btreeLayerImpl{
		dataAccessLayer: dataAccessLayer,
		rootPageNumber:  rootPageNumber,
		minFillPercent:  minFillPercent,
		maxFillPercent:  maxFillPercent,
	}

	if _, err := btl.getRootNode(); err != nil {
		return nil, fmt.Errorf("error while getting root node: %s", err.Error())
	}

	return btl, nil
}

func (b *btreeLayerImpl) GetRootPageNumber() pageNumber {
	return b.rootPageNumber
}

func (b *btreeLayerImpl) Get(key []byte) ([]byte, error) {
	node, index, _, err := b.findKey(key, true)
	if err != nil {
		return nil, fmt.Errorf("error while finding key: %s", err.Error())
	}

	item := node.GetItem(index)
	if item == nil {
		return nil, fmt.Errorf("error while getting item")
	}

	return item.value, nil
}

func (b *btreeLayerImpl) Set(key []byte, value []byte) error {
	nodeToInsertIn, insertionIndex, ancestors, err := b.findKey(key, false)
	if err != nil {
		return fmt.Errorf("error while finding key: %s", err.Error())
	}

	if index, ok := nodeToInsertIn.FindIndex(key); ok {
		nodeToInsertIn.GetItem(index).value = value
	} else {
		nodeToInsertIn.AddItemAt(insertionIndex, &item{
			key:   key,
			value: value,
		})
	}

	err = b.writeNode(nodeToInsertIn)
	if err != nil {
		return fmt.Errorf("error while writing node: %s", err.Error())
	}

	err = b.insertionFixup(ancestors)
	if err != nil {
		return fmt.Errorf("error while fixing up insertion: %s", err.Error())
	}

	return nil
}

func (b *btreeLayerImpl) insertionFixup(ancestors []int) error {
	nodes, err := b.getNodesFromAncestors(ancestors)
	if err != nil {
		return fmt.Errorf("error while getting nodes from ancestors: %s", err.Error())
	}

	for i := len(nodes) - 2; i >= 0; i-- {
		parentNode := nodes[i]
		node := nodes[i+1]
		nodeIndex := ancestors[i]

		if b.isOverpopulated(node) {
			err := b.split(parentNode, node, nodeIndex)
			if err != nil {
				return fmt.Errorf("error while splitting node: %s", err.Error())
			}
		}
	}

	rootNode := nodes[0]
	if b.isOverpopulated(rootNode) {
		newRoot := b.allocateNewNode(nil, []pageNumber{rootNode.GetPageNumber()})
		err := b.split(newRoot, rootNode, 0)
		if err != nil {
			return fmt.Errorf("error while splitting root node: %s", err.Error())
		}

		err = b.writeNode(newRoot)
		if err != nil {
			return fmt.Errorf("error while writing new root node: %s", err.Error())
		}
		b.rootPageNumber = newRoot.GetPageNumber()
	}

	return nil
}

func (b *btreeLayerImpl) maxThreshold() float64 {
	return float64(b.dataAccessLayer.GetPageSize()) * b.maxFillPercent
}

func (b *btreeLayerImpl) isOverpopulated(node node) bool {
	return float64(node.GetSize()) > b.maxThreshold()
}

func (b *btreeLayerImpl) minThreshold() float64 {
	return float64(b.dataAccessLayer.GetPageSize()) * b.minFillPercent
}

func (b *btreeLayerImpl) isUnderpopulated(node node) bool {
	return float64(node.GetSize()) < b.minThreshold()
}

func (b *btreeLayerImpl) split(node node, nodeToSplit node, index int) error {
	splitIndex := nodeToSplit.GetSplitIndex(b.minThreshold())
	splitItem := nodeToSplit.GetItem(splitIndex)

	newNode, err := b.splitNodeIntoNewNode(nodeToSplit, splitIndex)
	if err != nil {
		return fmt.Errorf("error while splitting node: %s", err.Error())
	}

	node.AddItemAt(index, splitItem)
	node.SetChildPageNumberAt(index+1, newNode.GetPageNumber())

	err = b.writeNodes(node, nodeToSplit)
	if err != nil {
		return fmt.Errorf("error while writing nodes: %s", err.Error())
	}

	return nil
}

func (b *btreeLayerImpl) splitNodeIntoNewNode(nodeToSplit node, splitIndex int) (node, error) {
	newNodeItems := nodeToSplit.GetItemsAfter(splitIndex)
	newNodeChildren := nodeToSplit.GetChildPageNumbersAfter(splitIndex)
	newNode := b.allocateNewNode(newNodeItems, newNodeChildren)

	err := b.writeNode(newNode)
	if err != nil {
		return nil, fmt.Errorf("error while writing new node: %s", err.Error())
	}

	nodeToSplit.SetItems(nodeToSplit.GetItemsBefore(splitIndex))
	nodeToSplit.SetChildPageNumbers(nodeToSplit.GetChildPageNumbersBefore(splitIndex))

	return newNode, nil
}

func (b *btreeLayerImpl) findKey(key []byte, exact bool) (node, int, []int, error) {
	root, err := b.getRootNode()
	if err != nil {
		return nil, -1, nil, fmt.Errorf("error while getting root node: %s", err.Error())
	}

	ancestors := make([]int, 0)
	node, index, err := b.findKeyHelper(root, key, exact, &ancestors)
	if err != nil {
		return nil, -1, nil, fmt.Errorf("error while finding key: %s", err.Error())
	}

	return node, index, ancestors, err
}

func (b *btreeLayerImpl) findKeyHelper(node node, key []byte, exact bool, ancestors *[]int) (node, int, error) {
	index, found := node.FindIndex(key)
	if found {
		return node, index, nil
	}

	if node.IsLeaf() || node.GetChildPageNumber(index) == -1 {
		if exact {
			return nil, -1, fmt.Errorf("key not found")
		}
		return node, index, nil
	}

	*ancestors = append(*ancestors, index)
	child, err := b.getChild(node, index)
	if err != nil {
		return nil, -1, fmt.Errorf("error while getting child: %s", err.Error())
	}

	return b.findKeyHelper(child, key, exact, ancestors)
}

func (b *btreeLayerImpl) getNodesFromAncestors(ancestors []int) ([]node, error) {
	nodes := make([]node, len(ancestors)+1)

	root, err := b.getRootNode()
	if err != nil {
		return nil, fmt.Errorf("error while getting root node: %s", err.Error())
	}
	nodes[0] = root

	for i := 0; i < len(ancestors); i++ {
		child, err := b.getChild(nodes[i], ancestors[i])
		if err != nil {
			return nil, fmt.Errorf("error while getting child: %s", err.Error())
		}
		nodes[i+1] = child
	}

	return nodes, nil
}

func (b *btreeLayerImpl) getChild(node node, index int) (node, error) {
	childPageNumber := node.GetChildPageNumber(index)
	if childPageNumber == -1 {
		return nil, fmt.Errorf("error while getting child node: child page number is -1")
	}

	child, err := b.getNodeFromPage(childPageNumber)
	if err != nil {
		return nil, fmt.Errorf("error while getting child node: %s", err.Error())
	}

	return child, nil
}

func (b *btreeLayerImpl) getRootNode() (node, error) {
	if b.rootPageNumber == -1 {
		err := b.writeNewRootNode()
		if err != nil {
			return nil, fmt.Errorf("error while writing new root node: %s", err.Error())
		}
	}

	node, err := b.getNodeFromPage(b.rootPageNumber)
	if err != nil {
		return nil, fmt.Errorf("error while getting root node: %s", err.Error())
	}
	return node, nil
}

func (b *btreeLayerImpl) writeNewRootNode() error {
	root := b.allocateNewEmptyNode()

	err := b.writeNode(root)
	if err != nil {
		return fmt.Errorf("error while writing root node: %s", err.Error())
	}

	b.rootPageNumber = root.GetPageNumber()
	return nil
}

func (b *btreeLayerImpl) getNodeFromPage(pageNum pageNumber) (node, error) {
	node := newEmptyNode(pageNum)

	page, err := b.dataAccessLayer.ReadPage(pageNum)
	if err != nil {
		return nil, fmt.Errorf("error while reading page: %s", err.Error())
	}

	err = node.Deserialize(page.buf)
	if err != nil {
		return nil, fmt.Errorf("error while deserializing node: %s", err.Error())
	}

	return node, nil
}

func (b *btreeLayerImpl) writeNodes(nodes ...node) error {
	for index, node := range nodes {
		err := b.writeNode(node)
		if err != nil {
			return fmt.Errorf("error while writing node at %d: %s", index, err.Error())
		}
	}
	return nil
}

func (b *btreeLayerImpl) writeNode(node node) error {
	page := b.dataAccessLayer.AllocateEmptyPage()

	err := node.Serialize(page.buf)
	if err != nil {
		return fmt.Errorf("error while serializing node: %s", err.Error())
	}

	if node.GetPageNumber() == -1 {
		pageNum, err := b.dataAccessLayer.WritePage(page.buf)
		if err != nil {
			return fmt.Errorf("error while writing page: %s", err.Error())
		}

		node.SetPageNumber(pageNum)
		return nil
	}

	err = b.dataAccessLayer.OverwritePage(node.GetPageNumber(), page.buf)
	if err != nil {
		return fmt.Errorf("error while overwriting page: %s", err.Error())
	}
	return nil
}

func (b *btreeLayerImpl) allocateNewEmptyNode() node {
	return newEmptyNode(-1)
}

func (b *btreeLayerImpl) allocateNewNode(items []*item, children []pageNumber) node {
	return newNode(-1, items, children)
}
