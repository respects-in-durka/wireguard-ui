package http

import (
	"github.com/gofiber/fiber/v3"
	"os"
	"wireguard-ui/internal"
	"wireguard-ui/internal/service"
)

type PeerController struct {
	peerService *service.PeerService
}

func NewPeerController(peerService *service.PeerService) *PeerController {
	return &PeerController{
		peerService: peerService,
	}
}

func (c *PeerController) Configure(router fiber.Router) {
	router.Get("/peer", c.GetPeers)
	router.Post("/peer", c.AddPeer)
	router.Delete("/peer/:name", c.DeletePeer)
	router.Get("/peer/:name/config", c.GetPeerConfig)
	router.Get("/peer/:name/qr", c.GetPeerQR)
}

func (c *PeerController) GetPeers(ctx fiber.Ctx) error {
	peers, err := c.peerService.GetPeers()

	if err != nil {
		return err
	}
	return ctx.JSON(peers)
}

func (c *PeerController) AddPeer(ctx fiber.Ctx) error {
	addPeer := &service.AddPeer{}

	if err := ctx.Bind().Body(addPeer); err != nil {
		return err
	}
	return c.peerService.AddPeer(addPeer)
}

func (c *PeerController) DeletePeer(ctx fiber.Ctx) error {
	peerName := ctx.Params("name")

	if peerName == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Peer name can't be empty")
	}

	return c.peerService.RemovePeer(peerName)
}

func (c *PeerController) GetPeerConfig(ctx fiber.Ctx) error {
	peerName := ctx.Params("name")

	if peerName == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Peer name can't be empty")
	}

	if internal.DevMode {
		path := "./wgconf/" + peerName + "/" + peerName + ".conf"

		if _, err := os.Stat(path); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Peer config file not found")
		}

		return ctx.Download(path)
	}

	path := "/config/" + peerName + "/" + peerName + ".conf"

	if _, err := os.Stat(path); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Peer config file not found")
	}

	return ctx.Download(path)
}

func (c *PeerController) GetPeerQR(ctx fiber.Ctx) error {
	peerName := ctx.Params("name")

	if peerName == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Peer name can't be empty")
	}

	if internal.DevMode {
		path := "./wgconf/" + peerName + "/" + peerName + ".png"

		if _, err := os.Stat(path); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Peer qr file not found")
		}

		return ctx.SendFile(path)
	}

	path := "/config/" + peerName + "/" + peerName + ".png"

	if _, err := os.Stat(path); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Peer qr file not found")
	}

	return ctx.SendFile(path)
}
