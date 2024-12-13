import { Component, createSignal } from "solid-js"
import toast from "solid-toast"

type CreatePeerFormProps = {
  update: () => void
}

const CreatePeerForm: Component<CreatePeerFormProps> = ({ update }) => {
  const [name, setName] = createSignal<string>("")
  const [ip, setIP] = createSignal<string>("")

  const nameRegex = new RegExp("^[a-zA-Z0-9]*$")
  const ipRegex = new RegExp("^((25[0-5]|(2[0-4]|1\\d|[1-9]|)\\d)\\.?\\b){4}$")

  const addPeer = async () => {
    if (name() == "" || ip() == "") {
      return
    }

    const response = await fetch(`/api/v1/peer/`, {
      method: "POST",
      body: JSON.stringify({
        "name": name(),
        "local_ip": ip()
      }),
      headers: {
        "Content-Type": "application/json"
      }
    })

    if (!response.ok) {
      const error = await response.text()
      toast.error(`Failed to add peer: ${error}`)
    }
    setName("")
    setIP("")
    update()
  }

  return (
    <div class="bg-gray-50 w-full">
      <form onSubmit={(e) => e.preventDefault()}>
        <div class="shadow overflow-hidden -space-y-px">
          <div class="flex justify-center items-center px-4 py-3">
            <div class="mx-2">
              <input
                name="name"
                value={name()}
                required
                placeholder="Name"
                class="appearance-none rounded-md relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
                classList={{
                  "text-red-700": !nameRegex.test(name())
                }}
                onInput={(e) => setName(e.target.value)}
              />
            </div>
            <div class="mx-2">
              <input
                name="ip"
                value={ip()}
                required
                placeholder="IP"
                class="appearance-none rounded-md relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
                classList={{
                  "text-red-700": !ipRegex.test(ip())
                }}
                onInput={(e) => setIP(e.target.value)}
              />
            </div>
            <div class="bg-gray-50 text-right mx-2">
              <button
                onClick={addPeer}
                class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">
                Add
              </button>
            </div>
          </div>
        </div>
      </form>
    </div>
  )
}

export default CreatePeerForm


