package testcc

type Information struct {
	TraceNo        string        `json:"traceNo"`
	Data	Date	`json:"data"`
	IssuerId	string `json:"issuer_id"`
}

type Date struct {
	TraceNo	string	`json:"trace_no"`
	TargetProperty	string	`json:"target_property"`
	TargetSize	string	`json:"target_size"`
	PictureHash	 string	`json:"picture_hash"`
	HashAlgorithm	string	`json:"hash_algorithm"`
}