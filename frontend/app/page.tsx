"use client"

import { useState } from "react"

export default function Home() {
  const [orderId, setOrderId] = useState("")
  const [amount, setAmount] = useState("")
  const [email, setEmail] = useState("")
  const [status, setStatus] = useState("")
  const [loading, setLoading] = useState(false)

  const handlePayment = async () => {
    setLoading(true)
    setStatus("")

    try {
      const res = await fetch(
        "http://localhost:8080/payments",
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            order_id: Number(orderId),
            amount: Number(amount),
            email,
          }),
        }
      )

      const data = await res.json()

      setStatus(`Payment ${data.status}`)
    } catch {
      setStatus("Failed to process payment")
    }

    setLoading(false)
  }

  return (
    <main className="relative min-h-screen overflow-hidden bg-black text-white">

      {/* GRID */}
      <div className="absolute inset-0 bg-[linear-gradient(rgba(255,255,255,0.03)_1px,transparent_1px),linear-gradient(90deg,rgba(255,255,255,0.03)_1px,transparent_1px)] bg-[size:70px_70px]" />

      {/* GLOW */}
      <div className="absolute top-[-200px] left-[-200px] h-[500px] w-[500px] rounded-full bg-blue-500/20 blur-3xl animate-pulse" />

      <div className="absolute bottom-[-200px] right-[-200px] h-[500px] w-[500px] rounded-full bg-purple-500/20 blur-3xl animate-pulse" />

      {/* HERO */}
      <section className="relative z-10 mx-auto flex min-h-screen max-w-7xl flex-col items-center justify-center px-6 py-20">

        {/* TOP */}
        <div className="mb-12 text-center">

          <div className="mb-4 inline-flex items-center rounded-full border border-white/10 bg-white/5 px-4 py-2 text-sm text-zinc-300 backdrop-blur-xl">
            Production-Style Event-Driven Architecture
          </div>

          <h1 className="max-w-4xl bg-gradient-to-r from-white via-zinc-200 to-zinc-500 bg-clip-text text-6xl font-black leading-tight text-transparent md:text-7xl">
            Distributed Payment Platform
          </h1>

          <p className="mx-auto mt-6 max-w-2xl text-lg text-zinc-400">
            Built with Go, gRPC, RabbitMQ, Redis,
            PostgreSQL, Docker and Next.js
          </p>
        </div>

        {/* STATS */}
        <div className="mb-10 grid w-full max-w-4xl grid-cols-2 gap-4 md:grid-cols-4">

          {[
            "gRPC",
            "RabbitMQ",
            "Redis",
            "Docker",
          ].map((item) => (
            <div
              key={item}
              className="rounded-2xl border border-white/10 bg-white/5 p-6 text-center backdrop-blur-xl"
            >
              <div className="text-2xl font-bold">
                ⚡
              </div>

              <div className="mt-2 text-sm text-zinc-300">
                {item}
              </div>
            </div>
          ))}
        </div>

        {/* PAYMENT CARD */}
        <div className="w-full max-w-xl rounded-3xl border border-white/10 bg-white/5 p-8 shadow-2xl backdrop-blur-2xl">

          <div className="mb-8">
            <h2 className="text-3xl font-bold">
              Process Payment
            </h2>

            <p className="mt-2 text-zinc-400">
              Send a payment request through the distributed event pipeline.
            </p>
          </div>

          <div className="space-y-4">

            <input
              type="number"
              placeholder="Order ID"
              value={orderId}
              onChange={(e) =>
                setOrderId(e.target.value)
              }
              className="w-full rounded-2xl border border-white/10 bg-black/30 p-4 text-white outline-none transition focus:border-blue-500"
            />

            <input
              type="number"
              placeholder="Amount"
              value={amount}
              onChange={(e) =>
                setAmount(e.target.value)
              }
              className="w-full rounded-2xl border border-white/10 bg-black/30 p-4 text-white outline-none transition focus:border-blue-500"
            />

            <input
              type="email"
              placeholder="Customer Email"
              value={email}
              onChange={(e) =>
                setEmail(e.target.value)
              }
              className="w-full rounded-2xl border border-white/10 bg-black/30 p-4 text-white outline-none transition focus:border-blue-500"
            />

            <button
              onClick={handlePayment}
              disabled={loading}
              className="w-full rounded-2xl bg-gradient-to-r from-blue-500 to-purple-600 p-4 text-lg font-bold transition-all duration-300 hover:scale-[1.02] active:scale-[0.98] disabled:opacity-50"
            >
              {loading
                ? "Processing..."
                : "Process Payment"}
            </button>

            {status && (
              <div
                className={`rounded-2xl border p-4 text-center font-semibold ${
                  status.includes("Failed")
                    ? "border-red-500/30 bg-red-500/10 text-red-300"
                    : "border-green-500/30 bg-green-500/10 text-green-300"
                }`}
              >
                {status.includes("Failed")
                  ? "❌ "
                  : "✅ "}
                {status}
              </div>
            )}
          </div>
        </div>

        {/* FOOTER */}
        <div className="mt-10 text-center text-sm text-zinc-500">
          Event-Driven Microservices • CQRS Style Messaging • Reliable Workers
        </div>
      </section>
    </main>
  )
}