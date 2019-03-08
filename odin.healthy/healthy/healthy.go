package healthy

import (
	"io/ioutil"
	"sync"

	"fmt"
	"os"
	"sync/atomic"

	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

var access_cnt int64

type odinHealthyHandler struct {
	lock     sync.RWMutex
	liveness map[string]HealthyCheck
	readness map[string]HealthyCheck
}

func NewHealthy() *odinHealthyHandler {
	inst := &odinHealthyHandler{
		lock:     sync.RWMutex{},
		liveness: make(map[string]HealthyCheck),
		readness: make(map[string]HealthyCheck),
	}

	return inst
}

func (h *odinHealthyHandler) AddLiveness(name string, check HealthyCheck) error {
	h.lock.Lock()
	defer h.lock.Unlock()

	if _, ok := h.liveness[name]; ok {
		return fmt.Errorf("%s has exist in liveness", name)
	}

	h.liveness[name] = check

	return nil
}

func (h *odinHealthyHandler) AddReadness(name string, check HealthyCheck) error {
	h.lock.Lock()
	defer h.lock.Unlock()

	if _, ok := h.readness[name]; ok {
		return fmt.Errorf("%s has exist in readness", name)
	}

	h.readness[name] = check

	return nil
}

func (h *odinHealthyHandler) aggregation() *HealthyResult {
	result := NewHealthyResult()
	h.lock.RLock()
	defer h.lock.RUnlock()

	for name, check := range h.liveness {
		if err := check(); err != nil {
			result.PushError(name, err)
			return result
		}

		result.PushOK(name)
	}

	for name, check := range h.readness {
		if err := check(); err != nil {
			result.PushError(name, err)
			return result
		}

		result.PushOK(name)
	}

	return result
}

func (h *odinHealthyHandler) HealthyEndpoint(ctx *routing.Context) error {
	ctx.SetContentType("application/json")

	atomic.AddInt64(&access_cnt, 1)

	result := h.aggregation()
	if result.Code != 0 {
		ctx.Response.SetStatusCode(fasthttp.StatusServiceUnavailable)
	}

	ctx.Write(result.Json())

	return nil
}

func AccessHandler(ctx *routing.Context) error {

	ctx.WriteString(fmt.Sprintf("access count: %d\n", atomic.LoadInt64(&access_cnt)))

	return nil
}

func DebugOKHandler(ctx *routing.Context) error {
	f, err := os.OpenFile("/tmp/abc.txt", os.O_RDWR|os.O_CREATE, 0664)
	if err != nil {
		ctx.WriteString("openfile error")
		return nil
	}
	f.WriteString("123456789")
	f.Close()

	ctx.WriteString("ok")

	return nil
}

func DebugErrHandler(ctx *routing.Context) error {
	f, err := os.OpenFile("/tmp/abc.txt", os.O_RDWR|os.O_CREATE, 0664)
	if err != nil {
		ctx.WriteString("openfile error")
		return nil
	}

	f.WriteString("000000000")
	f.Close()

	ctx.WriteString("ok")

	return nil
}

func CatFileHandler(ctx *routing.Context) error {
	f, err := os.OpenFile("/tmp/hello.txt", os.O_RDWR|os.O_CREATE, 0664)
	if err != nil {
		ctx.WriteString("openfile error")
		return nil
	}

	defer f.Close()

	b, _ := ioutil.ReadAll(f)

	ctx.Write(b)

	return nil
}
