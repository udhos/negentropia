/*
I want to make a tree from a table. the table is as following:

OrgID   OrgName        parentID
A001    Dept           0
A002    subDept1        A001
A003    sub_subDept    A002
A006    gran_subDept   A003
A004    subDept2        A001

and i want the result is as following,how to do it using go:

Dept

--subDept1

----sub_subDept

------gran_subDept

--subDept2

http://stackoverflow.com/questions/22957638/make-a-tree-from-a-table-using-golang
*/

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type Node struct {
	name     string
	children []*Node
}

var (
	nodeTable = map[string]*Node{}
	root      *Node
)

func add(id, name, parentId string) {
	fmt.Printf("add: id=%v name=%v parentId=%v\n", id, name, parentId)

	node := &Node{name: name, children: []*Node{}}

	if parentId == "0" {
		root = node
	} else {

		parent, ok := nodeTable[parentId]
		if !ok {
			fmt.Printf("add: parentId=%v: not found\n", parentId)
			return
		}

		parent.children = append(parent.children, node)
	}

	nodeTable[id] = node
}

func scan() {
	input := os.Stdin
	reader := bufio.NewReader(input)
	lineCount := 0
	for {
		lineCount++
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("error reading lines: %v\n", err)
			return
		}
		tokens := strings.Fields(line)
		if t := len(tokens); t != 3 {
			fmt.Printf("bad input line %v: tokens=%d [%v]\n", lineCount, t, line)
			continue
		}
		add(tokens[0], tokens[1], tokens[2])
	}
}

func showNode(node *Node, prefix string) {
	if prefix == "" {
		fmt.Printf("%v\n\n", node.name)
	} else {
		fmt.Printf("%v %v\n\n", prefix, node.name)
	}
	for _, n := range node.children {
		showNode(n, prefix+"--")
	}
}

func show() {
	if root == nil {
		fmt.Printf("show: root node not found\n")
		return
	}
	fmt.Printf("RESULT:\n")
	showNode(root, "")
}

func main() {
	fmt.Printf("main: reading input from stdin\n")
	scan()
	fmt.Printf("main: reading input from stdin -- done\n")
	show()
	fmt.Printf("main: end\n")
}
