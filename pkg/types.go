package pkg

type CtxType int

const (
	_proxyType CtxType = iota
	_request
	_user
	_uuid
	_report
	_response
)
