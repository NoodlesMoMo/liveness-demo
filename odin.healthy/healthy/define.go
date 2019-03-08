package healthy

import (
	"encoding/json"

	"github.com/qiangxue/fasthttp-routing"
)

type ResultItem struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type HealthyResult struct {
	Code int          `json:"code"`
	Msg  string       `json:"msg"`
	Data []ResultItem `json:"data"`
}

func NewHealthyResult() *HealthyResult {
	return &HealthyResult{
		Data: make([]ResultItem, 0),
	}
}

func (r *HealthyResult) PushOK(name string) {
	r.Data = append(r.Data, ResultItem{name, "ok"})
}

func (r *HealthyResult) PushError(name string, e error) {
	r.Code += 1
	r.Msg += e.Error() + "\n"
}

func (r *HealthyResult) Json() []byte {
	b, _ := json.Marshal(r)

	return b
}

type HealthyCheck = func() error

type HealthyHandler interface {
	AddLiveness(name string, check HealthyCheck) error
	AddReadness(name string, check HealthyCheck) error

	HealthyEndpoint(ctx *routing.Context) error
}
