// MessageType expression operator type
type MessageType uint8

// operator constant
const (
	_ MessageType = iota
    {{range $operator := . -}}
    MessageType{{$operator}}
    {{end}}
)

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
