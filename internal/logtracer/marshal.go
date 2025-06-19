package logtracer

type logentry struct {
	Function     string `json:"Function"`
	Source       string `json:"Source"`
	IncomingText string `json:"incoming_text"`
	Level        string `json:"level"`
	Msg          string `json:"msg"`
	ResultIntent string `json:"result_intent"`
	Time         string `json:"time"`
}
