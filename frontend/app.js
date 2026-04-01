const API_URL = "http://localhost:8080";

function showResponse(data) {
    document.getElementById("output").textContent =
        JSON.stringify(data, null, 2);
}

async function createOrder() {
    const body = {
        customer_id: document.getElementById("customerId").value,
        item_name: document.getElementById("itemName").value,
        amount: Number(document.getElementById("amount").value)
    };

    const response = await fetch(`${API_URL}/orders`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(body)
    });

    const data = await response.json();
    showResponse(data);

    if (data.id) {
        document.getElementById("orderId").value = data.id;
    }
}

async function getOrder() {
    const id = document.getElementById("orderId").value;

    const response = await fetch(`${API_URL}/orders/${id}`);
    const data = await response.json();

    showResponse(data);
}

async function cancelOrder() {
    const id = document.getElementById("orderId").value;

    const response = await fetch(`${API_URL}/orders/${id}/cancel`, {
        method: "PATCH"
    });

    const data = await response.json();
    showResponse(data);
}