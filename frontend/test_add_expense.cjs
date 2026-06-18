const axios = require('axios');
async function test() {
  try {
    const api = axios.create({ baseURL: 'http://localhost:8080/api' });
    const loginRes = await api.post('/auth/login', { email: 'admin@jewellery.com', password: 'Admin@123' });
    const token = loginRes.data.data.access_token;
    api.defaults.headers.common['Authorization'] = `Bearer ${token}`;

    console.log("Adding expense...");
    await api.post('/expenses', {
      category: 'test',
      amount: 50,
      description: 'testing dashboard update',
      expense_date: '2026-06-18'
    });

    const res = await api.get('/analytics/dashboard');
    console.log("New today_expenses:", res.data.data.metrics.today_expenses);
  } catch (err) {
    console.error(err.response?.status, err.response?.data);
  }
}
test();
