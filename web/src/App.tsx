import type { Component } from "solid-js"
import DashboardPage from "./pages/dashboard"
import { Toaster } from "solid-toast"

const DemoApp: Component = () => {
  return (
    <div class={"bg-gray-50 h-screen w-full"}>
      <div class={"text-lg text-red-700 w-full text-center"}>App is running in demo mode. No authentication and no VPN</div>
      <DashboardPage/>
      <Toaster />
    </div>
  )
}

const App: Component = () => {
  if (import.meta.env.MODE == "demo") {
    return <DemoApp/>
  }

  return (
   <>
     <DashboardPage/>
     <Toaster/>
   </>
  )
}

export default App




