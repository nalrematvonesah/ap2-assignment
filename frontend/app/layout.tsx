import "./globals.css"

export const metadata = {
  title: "Payment Platform",
  description: "Event-Driven Payment Microservices System",
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  )
}