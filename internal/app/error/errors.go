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
	ErrNoSuchStation      = Error{"NO_SUCH_STATION", 404, "입력한 역명으로 조회 가능한 역을 찾을 수 없습니다."}
	ErrNoData             = Error{"NO_DATA", 404, "혼잡도 정보가 없습니다."}
	ErrInternal           = Error{"COMMON_INTERNAL", 500, "관리자에게 문의하세요."}
)
