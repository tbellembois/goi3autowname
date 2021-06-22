package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
	"go.i3wm.org/i3/v4"
)

var (
	er           *i3.EventReceiver
	tree         i3.Tree
	root         *i3.Node
	err          error
	jsonFile     *os.File          // mapping file
	byteValue    []byte            // mapping file content
	applications Applications      // unmarshalled mapping file
	mm           map[string]string // app class<>name map
	wm           map[string]string // workspace map
)

type Applications struct {
	Applications []Application `json:"applications"`
}

type Application struct {
	Class string `json:"class"`
	Name  string `json:"name"`
}

func buildmap(n *i3.Node, w *i3.Node) {
	for _, c := range n.Nodes {
		switch c.Type {
		case i3.Con, i3.FloatingCon:
			class := string(c.WindowProperties.Class)
			// w should never be nil
			// mapping class<>name?
			newname := class
			if m, ok := mm[class]; ok {
				newname = m
			}
			if w != nil {
				log.WithFields(log.Fields{"n.ID": n.ID, "wm[w.Name]": wm[w.Name]}).Debug("buildmap")
				if wm[w.Name] == "" {
					wm[w.Name] = newname
				} else {
					wm[w.Name] = wm[w.Name] + " " + newname
				}
			}
			buildmap(c, w)
		case i3.WorkspaceNode:
			if c.Name != "__i3_scratch" {
				wm[c.Name] = ""
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
		// skip if using non numerical window names
		if int(w.Num) < 0 {
			log.Info("non integer window ID found: " + n)
			break
		}
		newname := n
		// creating a new workspace m[w.Name] is empty
		if wm[w.Name] != "" {
			newname = n + ":" + wm[w.Name]
		}
		log.WithFields(log.Fields{"n": n, "w.Name": w.Name, "newname": newname}).Debug("renameworkspaces")
		if _, err := i3.RunCommand(fmt.Sprintf(`rename workspace "%s" to "%s"`, w.Name , newname)); err != nil {
			log.WithFields(log.Fields{"n":n, "err":err}).Debug("rename failed")
		}
	}
}

func main() {
	// getting the program parameters
	debug := flag.Bool("debug", false, "debug (verbose log), default is info")
	mapf := flag.String("mapf", "./goi3autowname.json", "json map file full path")
	flag.Parse()

	// initializing maps
	wm = make(map[string]string)
	mm = make(map[string]string)

	// setting the log level
	if *debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	// opening the mapping file
	if jsonFile, err = os.Open(*mapf); err != nil {
		log.Info("no goi3autowname.json mapping file found, running with defaults")
	}
	defer jsonFile.Close()

	// reading the mapping file
	if byteValue, err = ioutil.ReadAll(jsonFile); err != nil {
		log.Info("error reading mapping file, running with defaults - err:" + err.Error())
	}
	if err = json.Unmarshal(byteValue, &applications); err != nil {
		log.Info("error unmarshalling mapping file, running with defaults - err:" + err.Error())
	}

	// loading the mapping into the map
	for i := 0; i < len(applications.Applications); i++ {
		log.Info(applications.Applications[i].Class + "->" + applications.Applications[i].Name)
		mm[applications.Applications[i].Class] = applications.Applications[i].Name
	}

	// subscribing to window events
	er = i3.Subscribe(i3.WindowEventType)

	// looping events
	for er.Next() {

		ev := er.Event().(*i3.WindowEvent)

		switch ev.Change {
		case "new", "close", "title":
			// getting i3 tree
			if tree, err = i3.GetTree(); err != nil {
				log.Fatal(err)
			}
			root = tree.Root

			buildmap(root, nil)
			log.WithFields(log.Fields{"wm": wm}).Debug("buildmap")
			renameworkspaces()
		default:
			log.WithFields(log.Fields{"ev.Change": ev.Change}).Debug("ev.Change")
		}
	}
	log.Fatal(er.Close())
}
