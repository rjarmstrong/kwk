package runtime

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	//"github.com/mitchellh/go-ps"
	"github.com/rjarmstrong/kwk-types"
	"github.com/rjarmstrong/kwk/src/out"
)

const processNodeEnvName = "PROCESS_NODE"

// ProcessNode stores all the metadata about the currently running snippet and its process.
type ProcessNode struct {
	URI     string       `msg:"ali" json:"uri"`
	Level   int64        `msg:"l" json:"lvl"`
	Args    []string     `msg:"args" json:"args"`
	Caller  *ProcessNode `msg:"c" json:"c"`
	Runner  string       `msg:"rnr" json:"rnr"`
	PPid    int          `msg:"ppid" json:"ppid"` //Parent pid
	PRunner string       `msg:"prnr" json:"prnr"` //Parent exe
	Pid     int          `msg:"pid" json:"pid"`   //Can only be retrospectively set.
}

var nodes = []*ProcessNode{}

var emptyCaller = &ProcessNode{URI: "-"}

// NewProcessNode creates a new instance of a ProcessNode. Used immediately before execution.
func NewProcessNode(a types.Alias, runner string, args []string, caller *ProcessNode) *ProcessNode {
	//printTree(os.Getpid(), args)
	exe, _ := os.Executable()
	n := &ProcessNode{URI: a.VersionURI(), Runner: runner, Args: args, PPid: os.Getpid(), PRunner: exe}
	if caller != nil {
		n.Caller = caller
		n.Level = caller.Level + 1
	} else {
		n.Caller = emptyCaller
		n.Level = 1
	}
	nodes = append(nodes, n)
	return n
}

// Complete captures the pid to store it in the ProcessNode.
func (node *ProcessNode) Complete(pid int) {
	node.Pid = pid
	out.Debug("NODE: %s", node.URI)
}

//func printTree(pid int, args []string) {
//	p, err := ps.FindProcess(pid)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	if p == nil || p.Pid() == 0 {
//		fmt.Println("DONE")
//		return
//	}
//	fmt.Println(p.Executable(), p.Pid(), args)
//	printTree(p.PPid(), nil)
//}

func getCurrentNode(a types.Alias, runner string, args []string, c *exec.Cmd) (*ProcessNode, error) {
	caller, err := GetCallerNode()
	if err != nil {
		return nil, err
	}
	node := NewProcessNode(a, runner, args, caller)
	b, _ := json.Marshal(node)
	nodeString := fmt.Sprintf("%s=%s", processNodeEnvName, b)
	c.Env = append(os.Environ(), nodeString)
	return node, nil
}

// GetCallerNode get the parent process metadata for the currently executing Snippet.
func GetCallerNode() (*ProcessNode, error) {
	callerString, ok := os.LookupEnv(processNodeEnvName)
	var caller *ProcessNode
	if ok && callerString != "" {
		caller = &ProcessNode{}
		err := json.Unmarshal([]byte(callerString), caller)
		if err != nil {
			return nil, err
		}
		if caller.Level == 0 {
			caller = nil
		}
	}
	return caller, nil
}
