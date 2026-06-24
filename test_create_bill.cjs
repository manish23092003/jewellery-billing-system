const axios = require('axios');

async function test() {
  try {
    // 1. Login
    const loginRes = await axios.post('http://localhost:8080/api/auth/login', {
      email: 'admin@example.com',
      password: 'password123'
    });
    
    const token = loginRes.data.data.access_token;
    console.log("Logged in:", token);

    // 2. Create Bill
    const payload = {
      type: "invoice",
      status: "completed",
      advance_amount: 0,
      invoice_date: new Date().toISOString(),
      customer_name: "Test Customer",
      customer_phone: "1234567890",
      payment_method: "cash",
      notes: "",
      items: [
        {
          item_name: "Gold Ring",
          metal_type: "gold",
          purity: "22K",
          weight: 10,
          rate_per_gram: 5000,
          making_charge: 500,
          gst_percentage: 3,
          quantity: 1,
          charges: []
        }
      ],
      old_gold_items: []
    };

    const res = await axios.post('http://localhost:8080/api/bills', payload, {
      headers: {
        Authorization: `Bearer ${token}`
      }
    });

    console.log("Success:", res.data);
  } catch (err) {
    if (err.response) {
      console.error("Error Response:", err.response.data);
    } else {
      console.error("Error:", err.message);
    }
  }
}

test();
