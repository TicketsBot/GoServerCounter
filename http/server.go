package http

import (
	"encoding/json"
	"fmt"
	"github.com/TicketsBot/GoServerCounter/config"
	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
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

	UpdateBody struct {
		Key string `json:"key"`
		Shard int `json:"shard"`
		ServerCount int `json:"serverCount"`
	}
)

var(
	counts = make(map[int]int)
	lock sync.Mutex
)

func StartServer() {
	router := routing.New()

	// /total
	totalHandler := fasthttp.CompressHandler(TotalHandler)
	router.Get("/total", func(ctx *routing.Context) error {
		totalHandler(ctx.RequestCtx)
		return nil
	})

	// /update
	updateHandler := fasthttp.CompressHandler(UpdateHandler)
	router.Post("/update", func(ctx *routing.Context) error {
		updateHandler(ctx.RequestCtx)
		return nil
	})

	err := fasthttp.ListenAndServe(config.Conf.Host, router.HandleRequest); if err != nil {
		panic(err)
	}
}

func TotalHandler(ctx *fasthttp.RequestCtx) {
	total := 0

	lock.Lock()
	defer lock.Unlock()
	for _, count := range counts {
		total += count
	}

	Respond(ctx, 200, Total{
		Success: true,
		Count: total,
	})
}

func UpdateHandler(ctx *fasthttp.RequestCtx) {
	var body UpdateBody
	err := json.Unmarshal(ctx.PostBody(), body); if err != nil {
		Respond(ctx, 400, GenericResponse{Success:false})
		return
	}

	lock.Lock()
	defer lock.Unlock()
	counts[body.Shard] = body.ServerCount

	fmt.Println(fmt.Sprintf("%d -> %d", body.Shard, body.ServerCount))

	Respond(ctx, 200, GenericResponse{Success:true})
}

func Respond(ctx *fasthttp.RequestCtx, responseCode int, response interface{}) {
	ctx.SetContentType("application/json; charset=utf8")
	ctx.Response.SetStatusCode(responseCode)

	marshalled, err := json.Marshal(response)

	if err != nil {
		fmt.Println(err.Error())
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
