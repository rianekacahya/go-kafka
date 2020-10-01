package crashy

import cm "github.com/orcaman/concurrent-map"

const (
	//ErrCodeUnexpected code for generic error for unrecognized cause
	ErrCodeUnexpected = "1000"
	//ErrCodeNetBuild code for resource connection build issue
	ErrCodeNetBuild = "1001"
	//ErrCodeNetConnect code for resource connection establish issue
	ErrCodeNetConnect = "1002"
	//ErrCodeValidation code for any data validation issues
	ErrCodeValidation = "1003"
	//ErrCodeFormatting code for any formatting issue(s) includes marshalling and unmarshalling
	ErrCodeFormatting = "1004"
	//ErrCodeDataRead code for any storage read issue(s)
	ErrCodeDataRead = "1005"
	//ErrCodeDataWrite code for any storage write issue(s)
	ErrCodeDataWrite = "1006"
	//ErrCodeNoResult code when data provider given no result for any query
	ErrCodeNoResult = "1007"
	//ErrCodeUnauthorized code when any access doesnt contains enough authorization
	ErrCodeUnauthorized = "1008"
	//ErrCodeExpired code when authorization given is already expired
	ErrCodeExpired = "1009"
	//ErrCodeForbidden code when access to an endpoint without proper permission
	ErrCodeForbidden = "1010"
	//ErrCodeTooManyRequest code when a user token used more than given rate
	ErrCodeTooManyRequest = "1011"
	//ErrCodeDataIncomplete code when data integrity expected but does not satisfied
	ErrCodeDataIncomplete = "1012"
	//ErrCodeSend code for any error due to failed send message(s)
	ErrCodeSend = "1013"
)

var (
	mapper = cm.New()
)

func init() {
	for k, v := range map[string]string{
		ErrCodeUnexpected:     "unexpected error occurred while processing request",
		ErrCodeNetBuild:       "failed to build connection to targeted source",
		ErrCodeNetConnect:     "failed to establish connection to targeted source",
		ErrCodeValidation:     "request contains invalid data",
		ErrCodeFormatting:     "an error occurred while formatting data",
		ErrCodeDataRead:       "failed to read data from data provider",
		ErrCodeDataWrite:      "failed to write data into provider",
		ErrCodeNoResult:       "no result found match criteria",
		ErrCodeUnauthorized:   "unauthorized access",
		ErrCodeForbidden:      "forbidden access",
		ErrCodeExpired:        "code expired",
		ErrCodeTooManyRequest: "request limit exceeded",
		ErrCodeDataIncomplete: "stored data incomplete",
		ErrCodeSend:           "unable to send code due to upstream failure",
	} {
		mapper.Set(k, v)
	}
}

func Message(code string) string {
	if s, ok := mapper.Get(code); ok {
		return s.(string)
	}
	if s, ok := mapper.Get(ErrCodeUnexpected); ok {
		return s.(string)
	}
	return "unexpected error occurred"
}

func Register(c string, s string) {
	mapper.Set(c, s)
}

func RegisterMap(m map[string]string) {
	for k, v := range m {
		Register(k, v)
	}
}

func Messages() map[string]string {
	var r = make(map[string]string)
	mapper.IterCb(func(key string, v interface{}) {
		r[key] = v.(string)
	})
	return r
}
