package hostnames

type sourceType struct {
	url       string
	preClean  cleanLineFunc
	checkLine checkLineFunc
	postClean cleanLineFunc
}
