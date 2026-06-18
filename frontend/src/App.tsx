import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { AuthProvider } from "@/context/AuthContext";
import LoginPage from "@/pages/LoginPage";
import DashboardPage from "@/pages/DashboardPage";
import MetalRatesPage from "@/pages/MetalRatesPage";
import CreateBillPage from "@/pages/CreateBillPage";
import BillHistoryPage from "@/pages/BillHistoryPage";
import ExpensesPage from "@/pages/ExpensesPage";
import SettingsPage from "@/pages/SettingsPage";
import Layout from "@/components/layout/Layout";
import { Toaster } from "sonner";

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: 1,
      refetchOnWindowFocus: false,
      staleTime: 5 * 60 * 1000,
    },
  },
});

export default function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <Toaster position="top-right" richColors theme="dark" />
      <AuthProvider>
        <BrowserRouter>
          <Routes>
            <Route path="/login" element={<LoginPage />} />

            {/* Protected routes wrapped in Layout */}
            <Route element={<Layout />}>
              <Route path="/" element={<DashboardPage />} />
              <Route path="/metal-rates" element={<MetalRatesPage />} />
              <Route path="/bills/new" element={<CreateBillPage />} />
              <Route path="/bills/history" element={<BillHistoryPage />} />
              <Route path="/expenses" element={<ExpensesPage />} />
              <Route path="/settings" element={<SettingsPage />} />
            </Route>

            <Route path="*" element={<Navigate to="/" replace />} />
          </Routes>
        </BrowserRouter>
      </AuthProvider>
    </QueryClientProvider>
  );
}
