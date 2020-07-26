package http

import (
	"encoding/json"
	"fmt"
	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"log"
	"os"
	"sync"
)

type(
	Total struct {
		Success bool `json:"success"`
		Count int `json:"count"`
	}

	GenericResponse struct {
		Success bool `json:"success"`
	}
)

var(
	Lock sync.RWMutex
	Count int
)

func StartServer() {
	router := routing.New()

	// /total
	totalHandler := fasthttp.CompressHandler(TotalHandler)
	router.Get("/total", func(ctx *routing.Context) error {
		totalHandler(ctx.RequestCtx)
		return nil
	})

	// /total/prometheus
	prometheusHandler := fasthttp.CompressHandler(PrometheusHandler)
	router.Get("/total/prometheus", func(ctx *routing.Context) error {
		prometheusHandler(ctx.RequestCtx)
		return nil
	})

	err := fasthttp.ListenAndServe(os.Getenv("SERVER_ADDR"), router.HandleRequest); if err != nil {
		panic(err)
	}
}

func TotalHandler(ctx *fasthttp.RequestCtx) {
	Lock.RLock()
	defer Lock.RUnlock()

	Respond(ctx, 200, Total{
		Success: true,
		Count: Count,
	})
}

func PrometheusHandler(ctx *fasthttp.RequestCtx) {
	ctx.Response.SetStatusCode(200)

	Lock.RLock()
	res := fmt.Sprintf("tickets_servercount %d", Count)
	Lock.RUnlock()

	_, err := fmt.Fprintln(ctx, res)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func Respond(ctx *fasthttp.RequestCtx, responseCode int, response interface{}) {
	ctx.SetContentType("application/json; charset=utf8")
	ctx.Response.SetStatusCode(responseCode)

	marshalled, err := json.Marshal(response)

	if err != nil {
		log.Println(err.Error())
		ctx.Response.SetStatusCode(500)
		_, err := fmt.Fprintln(ctx, "An internal server occurred")

		if err != nil {
			fmt.Println(err.Error())
		}

		return
	}

	_, err = fmt.Fprintln(ctx, string(marshalled))
	if err != nil {
		fmt.Println(err.Error())
	}
}
