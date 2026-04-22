const API = "http://localhost:8080";

let interval;

const statusEl = document.getElementById("status");

function setStatus(status) {
  statusEl.textContent = status;
  statusEl.className = "badge";

  if (status === "Pending") statusEl.classList.add("pending");
  if (status === "Paid") statusEl.classList.add("paid");
  if (status === "Failed") statusEl.classList.add("failed");
}

// CREATE ORDER
document.getElementById("orderForm").addEventListener("submit", async (e) => {
  e.preventDefault();

  const customer_id = document.getElementById("customerId").value;
  const item_name = document.getElementById("itemName").value;
  const amount = parseInt(document.getElementById("amount").value);

  const res = await fetch(`${API}/orders`, {
    method: "POST",
    headers: {"Content-Type": "application/json"},
    body: JSON.stringify({ customer_id, item_name, amount })
  });

  if (!res.ok) {
    alert("Ошибка создания заказа");
    return;
  }

  const data = await res.json();

  document.getElementById("orderId").textContent = data.id;
  document.getElementById("orderCard").classList.remove("hidden");

  setStatus(data.status);

  startPolling(data.id);
});

// LIVE STATUS
function startPolling(id) {
  clearInterval(interval);

  interval = setInterval(async () => {
    const res = await fetch(`${API}/orders/${id}`);
    const data = await res.json();

    setStatus(data.status);

    if (data.status === "Paid" || data.status === "Failed") {
      clearInterval(interval);
    }
  }, 1500);
}