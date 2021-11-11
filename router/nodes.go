package router

import (
	"errors"
	"fmt"
	"strings"
)

const (
	//
	charDelimiter string = `/`
	charPlus      string = `+`
	charTail      string = `#`
	// недопустимые символя для обычного текстового узла
	invalidChars string = "#"
	//

	// типы узлов
	nodeTypeString NodeType = iota
	nodeTypePrefix          // /str+/
	nodeTypeSuffix          // /+str/
	nodeTypePlus            // /+/
	nodeTypeTail            // /#

)

type (
	NodeType int

	INode interface {
		Type() NodeType
		Validate(p string) bool
		String() string
	}
	//
	nodeString struct {
		v string
	}

	nodePrefix struct {
		v string
	}

	nodeSuffix struct {
		v string
	}

	nodePlus struct{}

	nodeTail struct{}
)

func newNodeError(detail string) error {
	return errors.New(fmt.Sprintf("Node error : %s. ", detail))
}

func (n *nodeString) Validate(s string) bool {
	return n.v == s
}

func (n *nodeString) Type() NodeType {
	return nodeTypeString
}

func (n *nodeString) String() string {
	return n.v
}

func (n *nodePrefix) Validate(s string) bool {
	return strings.HasPrefix(s, n.v)
}

func (n *nodePrefix) Type() NodeType {
	return nodeTypePrefix
}

func (n *nodePrefix) String() string {
	return n.v
}

func (n *nodeSuffix) Validate(s string) bool {
	return strings.HasSuffix(s, n.v)
}

func (n *nodeSuffix) Type() NodeType {
	return nodeTypeSuffix
}

func (n *nodeSuffix) String() string {
	return n.v
}

func (n *nodePlus) Validate(s string) bool {
	return len(s) > 0
}

func (n *nodePlus) Type() NodeType {
	return nodeTypePlus
}

func (n *nodePlus) String() string {
	return "+"
}

func (n *nodeTail) Validate(s string) bool {
	return len(s) > 0
}

func (n *nodeTail) Type() NodeType {
	return nodeTypeTail
}

func (n *nodeTail) String() string {
	return "#"
}

//
// createNodeValue - создает узел нужного типа
//
func createNodeValue(s string) (INode, error) {

	countPlus := strings.Count(s, charPlus)

	if len(s) > 0 && countPlus < 2 {

		if s == charTail {
			return &nodeTail{}, nil
		}

		if strings.Contains(s, charTail) {
			return nil, newNodeError(s)
		}

		if s == charPlus {
			return &nodePlus{}, nil
		}

		if countPlus == 0 {
			if len(s) < 1 || strings.ContainsAny(s, invalidChars) {
				return nil, newNodeError(s)
			}
			return &nodeString{v: s}, nil
		}

		if strings.HasPrefix(s, charPlus) {
			temps := strings.TrimPrefix(s, charPlus)
			if strings.ContainsAny(temps, charPlus+charTail) {
				return nil, newNodeError(s)
			}
			return &nodeSuffix{v: temps}, nil
		}

		if strings.HasSuffix(s, charPlus) {
			temps := strings.TrimSuffix(s, charPlus)
			if strings.ContainsAny(temps, charPlus+charTail) {
				return nil, newNodeError(s)
			}
			return &nodePrefix{v: temps}, nil
		}
	}
	return nil, newNodeError(s)
}
