const axios = require('axios');
async function test() {
  try {
    const api = axios.create({ baseURL: 'http://localhost:8080/api' });
    const loginRes = await api.post('/auth/login', { email: 'admin@jewellery.com', password: 'Admin@123' });
    const token = loginRes.data.data.access_token;
    api.defaults.headers.common['Authorization'] = `Bearer ${token}`;

    const res = await api.get('/settings');
    console.log(res.data.data);
  } catch (err) {
    console.error(err.response?.status, err.response?.data);
  }
}
test();
