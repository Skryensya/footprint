package completions

import "github.com/footprint-tools/footprint-cli/internal/dispatchers"

var commandTree *dispatchers.DispatchNode

// RegisterCommandTree stores the command tree for later use by completion generators
// This should be called from main.go after building the tree
func RegisterCommandTree(root *dispatchers.DispatchNode) {
	commandTree = root
}

// GetCommandTree returns the registered command tree
func GetCommandTree() *dispatchers.DispatchNode {
	return commandTree
}
