package cmd

import (
	"time"
	"bitbucket.com/sharingmachine/kwkcli/models"
	"os"
	"encoding/json"
	"fmt"
	"os/exec"
	"bitbucket.com/sharingmachine/kwkcli/log"
)

////go:generate

type ProcessNode struct {
	AliasString string `msg:"ali" json:"ali"`
	Level       int `msg:"l" json:"lvl"`
	Args        []string `msg:"args" json:"args"`
	Caller      *ProcessNode `msg:"c" json:"c"`
	NodeStart   int64 `msg:"nst" json:"nst"`
	AppStart    int64 `msg:"ast" json:"ast"`
	AppDuration int64 `msg:"ad" json:"ad"`
	Pid         int `msg:"pid" json:"pid"`
}


func NewProcessNode(a models.Alias, args []string, caller *ProcessNode) *ProcessNode {
	n := &ProcessNode{AliasString : a.String(), Args : args, Pid : os.Getpid() }
	n.NodeStart = int64(time.Now().UnixNano())
	if caller != nil {
		n.Caller = caller
		n.Level = caller.Level + 1
		n.AppStart = caller.AppStart
		n.AppDuration=time.Now().UnixNano()-caller.AppStart
	} else {
		n.AppStart = time.Now().UnixNano()
		n.Caller = nil
		n.Level = 1
	}
	return n
}

func getCurrentNode(a models.Alias, args []string, c *exec.Cmd) (*ProcessNode, error) {
	callerString, ok := os.LookupEnv(PROCESS_NODE)
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
	node := NewProcessNode(a, args, caller)
	b, _ := json.Marshal(node)
	nodeString := fmt.Sprintf("%s=%s", PROCESS_NODE, b)
	c.Env = append(os.Environ(), nodeString)
	log.Debug(nodeString)
	log.Debug("Elapsed: %dms", node.AppDuration/1000000)
	return node, nil
}