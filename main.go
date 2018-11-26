package main

import (
	"go.i3wm.org/i3"
	"log"
	"strconv"
)

var (
	er   *i3.EventReceiver
	tree i3.Tree
	root *i3.Node
	err  error
	m    map[string]string // workspace map
)

func buildmap(n *i3.Node, w *i3.Node) {
	for _, c := range n.Nodes {
		switch c.Type {
		case i3.Con, i3.FloatingCon:
			class := c.WindowProperties.Class
			// w should never be nil
			if w != nil {
				m[w.Name] = m[w.Name] + "_" + string(class)
			}
			buildmap(c, w)
		case i3.WorkspaceNode:
			if c.Name != "__i3_scratch" {
				m[c.Name] = ""
				buildmap(c, c)
			}
		default:
			buildmap(c, w)
		}
	}
}

func renameworkspaces() {
	// getting the workspaces
	ws, err := i3.GetWorkspaces()
	if err != nil {
		log.Fatal(err)
	}

	for _, w := range ws {
		n := strconv.Itoa(int(w.Num))
		i3.RunCommand("rename workspace " + w.Name + " to " + n + ":" + m[w.Name])
	}
}

func main() {
	// initializing map
	m = make(map[string]string)

	// subscribing to window events
	er = i3.Subscribe(i3.WindowEventType)

	// looping events
	for er.Next() {

		ev := er.Event().(*i3.WindowEvent)

		switch ev.Change {
		case "new", "close":
			// getting i3 tree
			if tree, err = i3.GetTree(); err != nil {
				log.Fatal(err)
			}
			root = tree.Root

			buildmap(root, nil)
			renameworkspaces()
		}
	}
	log.Fatal(er.Close())
}
