// ORIGINAL: java/DomWalker.java

package domutil

import (
	"golang.org/x/net/html"
)

// WalkNodes used to walk the subtree of the DOM rooted at a particular root. It has two
// function parameters, i.e. fnVisit and fnExit :
// - fnVisit is called when we reach a node during the walk. If it returns false, children
//   of the node will be skipped and fnExit won't be called for this node.
// - fnExit is called when exiting a node, after visiting all of its children.
func WalkNodes(root *html.Node, fnVisit func(*html.Node) bool, fnExit func(*html.Node)) {
	if root == nil {
		return
	}

	visitChildren := false
	if fnVisit != nil {
		visitChildren = fnVisit(root)
	}

	if !visitChildren {
		return
	}

	for child := root.FirstChild; child != nil; child = child.NextSibling {
		WalkNodes(child, fnVisit, fnExit)
	}

	if fnExit != nil {
		fnExit(root)
	}
}
