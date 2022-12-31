package gee

import "fmt"

type node struct {
	prefix      string
	children    []*node
	pattern     string
	middlewares []HandlerFunc
}

type trie struct {
	root *node
}

func newTrie() *trie {
	return &trie{root: &node{}}
}

func (n *node) String() string {
	return fmt.Sprintf("node{pattern=%s, part=%s}", n.pattern, n.prefix)
}

func (n *node) isAny() bool {
	if len(n.prefix) == 0 {
		return false
	}
	return n.prefix[0] == '*' || n.prefix[0] == ':'
}

func (n *node) matchOne(part string) *node {
	for _, node := range n.children {
		if node.prefix == part || node.isAny() {
			return node
		}
	}
	return nil
}

func (n *node) matchAll(part string) (nodes []*node) {
	for _, node := range n.children {
		if node.prefix == part || node.isAny() {
			nodes = append(nodes, node)
		}
	}
	return
}

func (t *trie) insert(pattern string, middlewares []HandlerFunc) {
	parts := parsePattern(pattern)
	cur_node := t.root

	for _, part := range parts {
		new_node := cur_node.matchOne(part)
		if new_node == nil {
			new_node = &node{prefix: part}
			cur_node.children = append(cur_node.children, new_node)
		}
		cur_node = new_node
	}
	cur_node.pattern = pattern
	cur_node.middlewares = middlewares
}

func (n *node) search(parts []string) *node {
	// when the node can match any one node,or the parts has been itered over.
	if n.isAny() || len(parts) == 0 {
		// need to check if current node is an pattern
		if n.pattern == "" {
			return nil
		}
		return n
	}
	part := parts[0]
	nodes := n.matchAll(part)
	for _, node := range nodes {
		res := node.search(parts[1:])
		if res != nil {
			return res
		}
	}
	return nil
}

func (t *trie) search(path string) *node {
	parts := parsePattern(path)
	return t.root.search(parts)
}
