package domain

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"wireguard-ui/internal"
)

const PeerTemplate = `[Interface]
Address = %s
PrivateKey = %s
ListenPort = 51820
DNS = %s

[Peer]
PublicKey = %s
PresharedKey = %s
Endpoint = %s
AllowedIPs = %s`

type Peer struct {
	peerName         string
	peerIP           string
	peerPrivateKey   string
	peerPublicKey    string
	peerPresharedKey string
}

func CreatePeer(name, peerIP, serverPublicKey string) (*Peer, error) {
	err := createPeerDirectory(name)

	if err != nil {
		return nil, err
	}
	private, public, shared, err := createPeerKeys(name)

	if err != nil {
		_ = DeletePeer(name)
		return nil, err
	}

	endpoint, dns, allowedIPs, err := getPeerConfig()

	if err != nil {
		_ = DeletePeer(name)
		return nil, err
	}

	configPath := "/config/" + name + "/" + name + ".conf"

	if internal.DevMode {
		configPath = "./wgconf/" + name + "/" + name + ".conf"
	}

	config := fmt.Sprintf(PeerTemplate, peerIP, private, dns, serverPublicKey, shared, endpoint, allowedIPs)

	if err := os.WriteFile(configPath, []byte(config), 0640); err != nil {
		_ = DeletePeer(name)
		return nil, err
	}

	if err = os.Chown(configPath, 1000, 1000); err != nil {
		_ = DeletePeer(name)
		return nil, err
	}

	err = createPeerQR(name)

	if err != nil {
		_ = DeletePeer(name)
		return nil, err
	}

	return &Peer{
		peerName:         name,
		peerIP:           peerIP,
		peerPrivateKey:   private,
		peerPublicKey:    public,
		peerPresharedKey: shared,
	}, nil
}

func DeletePeer(name string) error {
	if internal.DevMode {
		return os.RemoveAll("./wgconf/" + name)
	}

	return os.RemoveAll("/config/" + name)
}

func getPeerConfig() (string, string, string, error) {
	serverUrl := os.Getenv("SERVERURL")

	if serverUrl == "" {
		return "", "", "", fmt.Errorf("SERVERURL env variable not set")
	}

	serverPort := os.Getenv("SERVERPORT")

	if serverPort == "" {
		return "", "", "", fmt.Errorf("SERVERPORT env variable not set")
	}

	dns := os.Getenv("PEERDNS")

	if dns == "" {
		dns = "1.1.1.1"
	}

	allowedIPs := os.Getenv("ALLOWEDIPS")

	if allowedIPs == "" {
		allowedIPs = "0.0.0.0/0, ::/0"
	}

	return serverUrl + ":" + serverPort, dns, allowedIPs, nil
}

func createPeerQR(peer string) error {
	qrPath := "/config/" + peer + "/" + peer + ".png"
	configPath := "/config/" + peer + "/" + peer + ".conf"
	qrencode := exec.Command("qrencode", "-o", qrPath, "-r", configPath)

	if internal.DevMode {
		qrencode = exec.Command("docker", "exec", "wireguard", "qrencode", "-o", qrPath, "-r", configPath)
	}

	return qrencode.Run()
}

func createPeerKeys(peer string) (string, string, string, error) {
	peerPath := "/config/" + peer + "/"

	if internal.DevMode {
		peerPath = "./wgconf/" + peer + "/"
	}

	private, err := createPeerKey("genkey")

	if err != nil {
		return "", "", "", err
	}
	public, err := createPeerPublicKey(private)

	if err != nil {
		return "", "", "", err
	}
	shared, err := createPeerKey("genpsk")

	if err != nil {
		return "", "", "", err
	}

	if err := os.WriteFile(peerPath+"privatekey-"+peer, []byte(private), 0600); err != nil {
		return "", "", "", err
	}

	if err := os.WriteFile(peerPath+"publickey-"+peer, []byte(public), 0600); err != nil {
		return "", "", "", err
	}

	if err := os.WriteFile(peerPath+"presharedkey-"+peer, []byte(shared), 0600); err != nil {
		return "", "", "", err
	}

	return private, public, shared, nil
}

func createPeerKey(key string) (string, error) {
	buffer := bytes.Buffer{}

	genkey := exec.Command("wg", key)

	if internal.DevMode {
		genkey = exec.Command("docker", "exec", "wireguard", "wg", key)
	}
	genkey.Stdout = &buffer

	if err := genkey.Run(); err != nil {
		return "", err
	}

	return strings.ReplaceAll(buffer.String(), "\n", ""), nil
}

func createPeerPublicKey(privateKey string) (string, error) {
	buffer := bytes.Buffer{}

	if internal.DevMode {
		genkey := exec.Command("docker", "exec", "-i", "wireguard", "wg", "pubkey")
		genkey.Stdout = &buffer
		genkey.Stdin = bytes.NewBufferString(privateKey)

		if err := genkey.Run(); err != nil {
			return "", err
		}

		return strings.ReplaceAll(buffer.String(), "\n", ""), nil
	}
	genkey := exec.Command("wg", "pubkey")
	genkey.Stdout = &buffer
	genkey.Stdin = bytes.NewBufferString(privateKey)

	if err := genkey.Run(); err != nil {
		return "", err
	}

	return strings.ReplaceAll(buffer.String(), "\n", ""), nil
}

func createPeerDirectory(peer string) error {
	peerPath := "/config/" + peer

	if internal.DevMode {
		peerPath = "./wgconf/" + peer
	}

	_, err := os.Stat(peerPath)

	if os.IsExist(err) {
		return fmt.Errorf("peer %s already exists", peerPath)
	}

	if err = os.Mkdir(peerPath, 0750); err != nil {
		return err
	}

	return os.Chown(peerPath, 1000, 1000)
}
