package stitch

type Node struct {
	Name    string
	Inputs  []string
	Outputs []string
}

var ImportedNodes map[string][]string = map[string][]string{
	"snmp": []string{"get", "walk"},
}

var StandardNodes map[string]map[string]Node = map[string]map[string]Node{
	"snmp": map[string]Node{
		"get": Node{Name: "get", Inputs: []string{"Input"}, Outputs: []string{"Output", "Error"}},
	},
}
