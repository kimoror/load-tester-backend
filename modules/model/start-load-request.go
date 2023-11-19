package model

type StartLoadRequest struct {
	StepConfig []stepConfig `json:"stepConfig"`
	Url        string       `json:"url"`
	TestBody   TestBody     `json:"testBody"`
}

type stepConfig struct {
	Rate            int `json:"rate"`
	DurationSeconds int `json:"durationSeconds"`
}

type TestBody struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
}
