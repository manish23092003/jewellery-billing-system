fetch('https://jewellery-billing-system.onrender.com/api/auth/register', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    business_name: "Test Shop 999",
    owner_name: "Test Owner",
    email: "newtest999xyz@gmail.com",
    phone: "9876543210",
    password: "password123",
    confirm_password: "password123"
  })
}).then(r => r.json()).then(d => console.log(JSON.stringify(d, null, 2))).catch(e => console.error(e));
