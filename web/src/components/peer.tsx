import { Component, createSignal } from "solid-js"
import { DeleteIcon, DownloadIcon, QrIcon } from "./icons"
import Modal from "@lutaok/solid-modal"
import toast from "solid-toast"

type PeerCardProps = {
  peer: Peer,
  update: () => void
}

const PeerCard: Component<PeerCardProps> = ({ peer, update }) => {
  const [qrModalOpened, setQrModalOpened] = createSignal<boolean>(false)

  const deletePeer = async () => {
    const response = await fetch(`/api/v1/peer/${peer.name}`, {
      method: "DELETE"
    })

    if (!response.ok) {
      const error = await response.text()
      toast.error(`Failed to delete peer: ${error}`)
    }
    update()
  }

  return (
    <div class={"w-full rounded-md border-2 border-indigo-600 border-solid flex flex-row justify-between p-2"}>
      <div>
        <p>Name: {peer.name}</p>
        <p>Remote IP: {peer.remote_ip}</p>
        <p>Virtual IP: {peer.local_ip}</p>
      </div>
      <div>
        <div onClick={() => setQrModalOpened(true)} class={"hover:cursor-pointer"}>
          <QrIcon />
        </div>
        <a href={`/api/v1/peer/${peer.name}/config`} target={"_blank"} rel={"noreferrer"}
           class={"hover:cursor-pointer my-4"}>
          <DownloadIcon />
        </a>
        <div onClick={deletePeer} class={"hover:cursor-pointer"}>
          <DeleteIcon />
        </div>
        <Modal isOpen={qrModalOpened()} onCloseRequest={() => setQrModalOpened(false)} closeOnOutsideClick>
          <img src={`/api/v1/peer/${peer.name}/qr`} alt={"QR"} width={320} height={320} />
          <button onClick={() => setQrModalOpened(false)}
                  class={"relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"}>
            Close
          </button>
        </Modal>
      </div>
    </div>
  )
}

export default PeerCard