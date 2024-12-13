import { Component } from "solid-js"

const LoadingSpinner: Component = () => {
  return (
    <div class={"h-screen flex justify-center items-center"}>
      <div
        class={"w-16 h-16 border-4 border-indigo-600 border-solid rounded-full animate-spin border-t-transparent"}></div>
    </div>
  )
}

export default LoadingSpinner