package service

import (
	"bytes"
	"os/exec"
	"strings"
	"wireguard-ui/internal"
	"wireguard-ui/internal/domain"
)

type PeerService struct {
	config *domain.WireguardConfig
}

func NewPeerService(config *domain.WireguardConfig) *PeerService {
	return &PeerService{config}
}

type Peer struct {
	Name     string `json:"name"`
	LocalIP  string `json:"local_ip"`
	RemoteIP string `json:"remote_ip"`
}

type AddPeer struct {
	Name    string `json:"name"`
	LocalIP string `json:"local_ip"`
}

func NewPeer(name, localIP, remoteIP string) *Peer {
	if remoteIP == "(none)" {
		remoteIP = "not connected"
	}

	return &Peer{
		Name:     name,
		LocalIP:  localIP,
		RemoteIP: remoteIP,
	}
}

func (s *PeerService) GetPeers() ([]*Peer, error) {
	var out bytes.Buffer

	cmd := exec.Command("wg", "show", "wg0", "dump")

	if internal.DevMode {
		cmd = exec.Command("docker", "exec", "wireguard", "wg", "show", "wg0", "dump")
	}
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	peersRaw := strings.Split(out.String(), "\n")
	peersRaw = peersRaw[1 : len(peersRaw)-1]

	var peers []*Peer

	for _, peer := range peersRaw {
		data := strings.Fields(peer)
		name := s.config.GetPeerName(data[0])

		peers = append(peers, NewPeer(name, data[3], data[2]))
	}

	return peers, nil
}

func (s *PeerService) AddPeer(addPeer *AddPeer) error {
	peer, err := domain.CreatePeer(addPeer.Name, addPeer.LocalIP, s.config.GetPublicKey())

	if err != nil {
		return err
	}

	err = s.config.AddPeer(peer)

	if err != nil {
		return err
	}

	return s.config.Restart()
}

func (s *PeerService) RemovePeer(peerName string) error {
	err := domain.DeletePeer(peerName)

	if err != nil {
		return err
	}

	err = s.config.RemovePeer(peerName)

	if err != nil {
		return err
	}

	return s.config.Restart()
}
