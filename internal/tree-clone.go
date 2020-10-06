// ORIGINAL: java/TreeCloneBuilder.java

package internal

import (
	"golang.org/x/net/html"
)

// BuildTreeClone takes a list of nodes and returns a clone of the minimum tree in the
// DOM that contains all of them. This is done by going through each node, cloning its
// parent and adding children to that parent until the next node is not contained in
// that parent (originally). The list cannot contain a parent of any of the other nodes.
// Children of the nodes in the provided list are excluded.
//
// This implementation doesn't come from the original dom-distiller code. Instead I
// created it from scratch to make it simpler and more Go idiomatic.
func BuildTreeClone(nodes []*html.Node) *html.Node {
	// Find all ancestors
	ancestors := make(map[*html.Node]int)
	for _, node := range nodes {
		// Include the node itself to list of ancestor
		ancestors[node]++

		// Save parents of node to list ancestor
		parent := node.Parent
		for parent != nil {
			ancestors[parent]++
			parent = parent.Parent
		}
	}

	// Find common ancestor
	nNodes := len(nodes)
	commonAncestors := make(map[*html.Node]struct{})
	for node, count := range ancestors {
		if count == nNodes {
			commonAncestors[node] = struct{}{}
		}
	}

	// If there are no common ancestor found, stop
	if len(commonAncestors) == 0 {
		return nil
	}

	// Find the nearest ancestor
	var nearestAncestor *html.Node
	for node := range commonAncestors {
		childIsAncestor := false
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			if _, exist := commonAncestors[child]; exist {
				childIsAncestor = true
				break
			}
		}

		if !childIsAncestor {
			nearestAncestor = node
		}
	}

	// Clone the ancestor and childrens that required to reach specified nodes
	var fnClone func(src *html.Node) *html.Node
	fnClone = func(src *html.Node) *html.Node {
		clone := &html.Node{
			Type:     src.Type,
			DataAtom: src.DataAtom,
			Data:     src.Data,
			Attr:     append([]html.Attribute{}, src.Attr...),
		}

		for child := src.FirstChild; child != nil; child = child.NextSibling {
			if _, exist := ancestors[child]; exist {
				clone.AppendChild(fnClone(child))
			}
		}

		return clone
	}

	return fnClone(nearestAncestor)
}
