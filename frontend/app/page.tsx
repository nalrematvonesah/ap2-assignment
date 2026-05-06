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
      const res = await fetch("http://localhost:8080/payments", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          order_id: Number(orderId),
          amount: Number(amount),
          email,
        }),
      })

      const data = await res.json()

      setStatus(`Payment ${data.status}`)
    } catch (err) {
      setStatus("Failed to process payment")
    }

    setLoading(false)
  }

  return (
    <main className="min-h-screen relative overflow-hidden bg-black flex items-center justify-center px-4">

      {/* Background Glow */}
      <div className="absolute w-[500px] h-[500px] bg-blue-500/20 blur-3xl rounded-full top-[-100px] left-[-100px]" />
      <div className="absolute w-[500px] h-[500px] bg-purple-500/20 blur-3xl rounded-full bottom-[-100px] right-[-100px]" />

      {/* Card */}
      <div className="relative z-10 w-full max-w-md backdrop-blur-xl bg-white/10 border border-white/20 rounded-3xl shadow-2xl p-10">

        {/* Header */}
        <div className="mb-8 text-center">
          <h1 className="text-4xl font-bold text-white mb-2">
            Payment System
          </h1>

          <p className="text-zinc-300">
            Event-Driven Microservices Platform
          </p>
        </div>

        {/* Form */}
        <div className="flex flex-col gap-4">

          <input
            type="number"
            placeholder="Order ID"
            value={orderId}
            onChange={(e) => setOrderId(e.target.value)}
            className="bg-white/10 border border-white/20 rounded-xl p-4 text-white placeholder:text-zinc-400 outline-none focus:border-blue-400 transition"
          />

          <input
            type="number"
            placeholder="Amount"
            value={amount}
            onChange={(e) => setAmount(e.target.value)}
            className="bg-white/10 border border-white/20 rounded-xl p-4 text-white placeholder:text-zinc-400 outline-none focus:border-blue-400 transition"
          />

          <input
            type="email"
            placeholder="Email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            className="bg-white/10 border border-white/20 rounded-xl p-4 text-white placeholder:text-zinc-400 outline-none focus:border-blue-400 transition"
          />

          {/* Button */}
          <button
            onClick={handlePayment}
            disabled={loading}
            className="mt-2 bg-gradient-to-r from-blue-500 to-purple-600 hover:scale-[1.02] active:scale-[0.98] transition-all duration-200 rounded-xl p-4 font-bold text-white shadow-lg disabled:opacity-50"
          >
            {loading ? "Processing..." : "Process Payment"}
          </button>

          {/* Status */}
          {status && (
            <div className={`mt-4 rounded-xl p-4 text-center font-semibold border ${
              status.includes("Failed")
                ? "bg-red-500/20 border-red-500/30 text-red-300"
                : "bg-green-500/20 border-green-500/30 text-green-300"
            }`}>
              {status.includes("Failed") ? "❌ " : "✅ "}
              {status}
            </div>
          )}
        </div>

        {/* Footer */}
        <div className="mt-8 text-center text-sm text-zinc-400">
          Go • gRPC • RabbitMQ • PostgreSQL • Docker
        </div>
      </div>
    </main>
  )
}