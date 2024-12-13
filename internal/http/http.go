package http

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"wireguard-ui/internal"
	"wireguard-ui/internal/domain"
	"wireguard-ui/internal/service"
)

func Configure(app *fiber.App) {
	apiGroup := app.Group("/api/v1")

	var config *domain.WireguardConfig

	if internal.DevMode {
		config, _ = domain.LoadWireguardConfig("./wgconf/wg_confs/wg0.conf", "./wgconf/server/publickey-server")
	} else {
		cfg, err := domain.LoadWireguardConfig("/config/wg_confs/wg0.conf", "/config/server/publickey-server")

		if err != nil {
			log.Fatal(err)
		}
		config = cfg
	}

	peerService := service.NewPeerService(config)

	peerController := NewPeerController(peerService)
	peerController.Configure(apiGroup)

	if internal.DevMode {
		proxyURL, _ := url.Parse("http://0.0.0.0:3001")

		proxy := httputil.ReverseProxy{Director: func(req *http.Request) {
			req.URL.Scheme = proxyURL.Scheme
			req.URL.Host = proxyURL.Host
			req.URL.Path = proxyURL.Path + req.URL.Path
			req.Host = proxyURL.Host
		}}

		app.Get("/*", func(ctx fiber.Ctx) error {
			fasthttpadaptor.NewFastHTTPHandler(&proxy)(ctx.Context())
			return nil
		})
	} else {
		app.Get("/*", static.New("./dist"))
	}
}
