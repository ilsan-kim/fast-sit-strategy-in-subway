package serror

type Error struct {
	Code           string
	HTTPStatusCode int `json:"-"`
	Message        string
}

func (e Error) Error() string {
	return e.Message
}

var (
	ErrInvalidRequestTime = Error{"NOT_VALID_REQUEST_TIME", 400, "요청이 불가능한 시간대입니다."}
	ErrExternalService    = Error{"EXTERNAL_SERVICE_NOT_WORKING", 400, "외부 서비스가 작동하지 않습니다."}
)
