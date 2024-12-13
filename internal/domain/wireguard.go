package domain

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"wireguard-ui/internal"
)

type wireguardPeer struct {
	name         string
	publicKey    string
	presharedKey string
	ip           string
}

func newWireguardPeer(raw string) *wireguardPeer {
	data := strings.Split(raw, "\n")

	name := strings.Split(data[1], " ")[1]
	publicKey := strings.Split(data[2], " = ")[1]
	presharedKey := strings.Split(data[3], " = ")[1]
	ip := strings.Split(data[4], " = ")[1]

	return &wireguardPeer{
		name:         name,
		publicKey:    publicKey,
		presharedKey: presharedKey,
		ip:           ip,
	}
}

type WireguardConfig struct {
	filePath  string
	publicKey string
	configRaw string
	peers     []*wireguardPeer
}

const WireguardPeerTemplate = `

[Peer]
# %s
PublicKey = %s
PresharedKey = %s
AllowedIPs = %s`

func (c *WireguardConfig) GetPublicKey() string {
	return c.publicKey
}

func (c *WireguardConfig) GetPeerName(publicKey string) string {
	for _, peer := range c.peers {
		if peer.publicKey == publicKey {
			return peer.name
		}
	}

	return ""
}

func (c *WireguardConfig) AddPeer(peer *Peer) error {
	c.peers = append(c.peers, &wireguardPeer{
		name:         peer.peerName,
		publicKey:    peer.peerPublicKey,
		presharedKey: peer.peerPresharedKey,
		ip:           peer.peerIP + "/32",
	})

	return c.save()
}

func (c *WireguardConfig) RemovePeer(peerName string) error {
	for i, peer := range c.peers {
		if peer.name == peerName {
			c.peers = append(c.peers[:i], c.peers[i+1:]...)
			return c.save()
		}
	}
	return fmt.Errorf("peer %s not found", peerName)
}

func (c *WireguardConfig) save() error {
	file, err := os.OpenFile(c.filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)

	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(c.configRaw)

	if err != nil {
		return err
	}

	for _, peer := range c.peers {
		_, err = file.WriteString(fmt.Sprintf(WireguardPeerTemplate, peer.name, peer.publicKey, peer.presharedKey, peer.ip))

		if err != nil {
			return err
		}
	}

	return nil
}

func (c *WireguardConfig) Restart() error {
	down := exec.Command("wg-quick", "down", "wg0")
	up := exec.Command("wg-quick", "up", "wg0")

	if internal.DevMode {
		down = exec.Command("docker", "exec", "wireguard", "wg-quick", "down", "wg0")
		up = exec.Command("docker", "exec", "wireguard", "wg-quick", "up", "wg0")
	}

	if err := down.Run(); err != nil {
		return err
	}

	return up.Run()
}

func LoadWireguardConfig(path string, publicKeyPath string) (*WireguardConfig, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}

	if _, err := os.Stat(publicKeyPath); os.IsNotExist(err) {
		return nil, err
	}

	publicKey, err := os.ReadFile(publicKeyPath)

	if err != nil {
		return nil, err
	}

	configBytes, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	configRaw := string(configBytes)
	peersRaw := strings.Split(configRaw, "\n\n")[1:]
	config := strings.Split(configRaw, "\n\n")[0]
	var peers []*wireguardPeer

	for _, peerRaw := range peersRaw {
		peers = append(peers, newWireguardPeer(peerRaw))
	}

	return &WireguardConfig{
		filePath:  path,
		publicKey: strings.ReplaceAll(string(publicKey), "\n", ""),
		configRaw: config,

		peers: peers,
	}, nil
}
