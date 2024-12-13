import { Component, createResource, For } from "solid-js"
import CreatePeerForm from "../components/form"
import Peer from "../components/peer"

const DashboardPage: Component = () => {
  const [peers, { refetch }] = createResource<[Peer]>(async () => {
    const response = await fetch("/api/v1/peer")

    return await response.json()
  })

  return (
    <>
      <CreatePeerForm update={refetch} />
      <div class={"grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3 mt-4 mx-4"}>
        <For each={peers()}>
          {(peer) =>
            <Peer peer={peer} update={refetch} />
          }
        </For>
      </div>
    </>
  )
}

export default DashboardPage
