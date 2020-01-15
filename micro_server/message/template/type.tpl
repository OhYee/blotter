// MessageType expression operator type
type MessageType uint8

// operator constant
const (
	_ MessageType = iota
    {{range $operator := . -}}
    MessageType{{$operator}}
    {{end}}
)

{{range $operator := . -}}
// MessageType{{$operator}}Handle function to solve {{$operator}} type message
type MessageType{{$operator}}Handle func (data {{$operator}}) (err error)
{{end}}

var typeName = [...]string{
    "Unknown", 
    {{range $operator := . -}}
    "{{$operator}}", 
    {{end}}
}

func (t MessageType) String() string {
	if int(t) > len(typeName) {
		return "Unknown"
	}
	return typeName[t]
}
