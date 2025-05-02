package rpc

type MessageName string
type ResponseName string
const (
	CLOSE MessageName = "close"
	SHUTDOWN MessageName = "shutdown"
	INDEX MessageName = "index file"
	SEARCH MessageName = "search"

	ERROR ResponseName = "error"
	SUCCESS ResponseName = "success"
)

type BaseMessage struct {
	Method MessageName `json:"method"`
}

type BaseResponse struct {
	Type ResponseName `json:"type"`
}

type ErrorResponse struct {
	BaseResponse
	Error string `json:"error"`
}
func NewErrorResponse(e string) *ErrorResponse {
	r := &ErrorResponse{}
	r.Type = ERROR
	r.Error = e
	return r
}

type SuccessResponse struct {
	BaseResponse
	Value string `json:"value"`
}
func NewSuccessResponse(value string) *SuccessResponse {
	s := &SuccessResponse{}
	s.Type = SUCCESS
	s.Value = value
	return s
}

type IndexMessage struct {
	BaseMessage
	Filename string `json:"filename"`
}
func NewIndexMessage(filename string) *IndexMessage {
	i := &IndexMessage{}
	i.Method = INDEX
	i.Filename = filename
	return i
}

type SearchMessage struct {
	BaseMessage
	Query string `json:"query"`
}
func NewSearchMessage(query string) *SearchMessage {
	s := &SearchMessage{}
	s.Method = SEARCH
	s.Query = query
	return s
}

var CloseMessage BaseMessage = BaseMessage {
	Method: CLOSE,
}
var ShutdownMessage BaseMessage = BaseMessage {
	Method: SHUTDOWN,
}
