// ORIGINAL: java/TreeCloneBuilder.java

package domutil

import (
	"golang.org/x/net/html"
)

// TreeClone takes a list of nodes and returns a clone of the minimum tree in the
// DOM that contains all of them. This is done by going through each node, cloning its
// parent and adding children to that parent until the next node is not contained in
// that parent (originally). The list cannot contain a parent of any of the other nodes.
// Children of the nodes in the provided list are excluded.
//
// This implementation doesn't come from the original dom-distiller code. Instead I
// created it from scratch to make it simpler and more Go idiomatic.
func TreeClone(nodes []*html.Node) *html.Node {
	// Get the nearest ancestor
	allAncestors, nearestAncestor := GetAncestors(nodes...)
	if nearestAncestor == nil {
		return nil
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
			if _, exist := allAncestors[child]; exist {
				clone.AppendChild(fnClone(child))
			}
		}

		return clone
	}

	return fnClone(nearestAncestor)
}
