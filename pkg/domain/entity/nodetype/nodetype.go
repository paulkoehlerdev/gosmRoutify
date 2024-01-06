package nodetype

type NodeType uint8

const (
	EMPTYNODE        NodeType = iota
	ENDNODE          NodeType = iota
	INTERMEDIATENODE NodeType = iota
	CONNECTIONNODE   NodeType = iota
	JUNCTIONNODE     NodeType = iota
)

func (n NodeType) String() string {
	return [...]string{"EMPTYNODE", "ENDNODE", "INTERMEDIATENODE", "CONNECTIONNODE", "JUNCTIONNODE"}[n]
}

func (n NodeType) IsEmpty() bool {
	return n == EMPTYNODE
}

func (n NodeType) IsTowerNode() bool {
	return n == JUNCTIONNODE || n == CONNECTIONNODE
}

func (n NodeType) IsPillarNode() bool {
	return n == INTERMEDIATENODE || n == ENDNODE
}
