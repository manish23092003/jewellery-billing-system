import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { AuthProvider } from "@/context/AuthContext";
import LoginPage from "@/pages/LoginPage";
import RegisterPage from "@/pages/RegisterPage";
import ForgotPasswordPage from "@/pages/ForgotPasswordPage";
import ResetPasswordPage from "@/pages/ResetPasswordPage";
import VerifyEmailPage from "@/pages/VerifyEmailPage";
import DashboardPage from "@/pages/DashboardPage";
import VerifyBillPage from "@/pages/VerifyBillPage";
import MetalRatesPage from "@/pages/MetalRatesPage";
import CreateBillPage from "@/pages/CreateBillPage";
import BillHistoryPage from "@/pages/BillHistoryPage";
import ExpensesPage from "@/pages/ExpensesPage";
import CustomersPage from "@/pages/CustomersPage";
import SettingsPage from "@/pages/SettingsPage";
import Layout from "@/components/layout/Layout";
import { ThemeProvider } from "@/components/ThemeProvider";
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
    <ThemeProvider defaultTheme="system" storageKey="jewellery-theme">
      <QueryClientProvider client={queryClient}>
        <Toaster position="top-right" richColors theme="system" />
      <AuthProvider>
        <BrowserRouter>
          <Routes>
            {/* Public Auth Routes */}
            <Route path="/login" element={<LoginPage />} />
            <Route path="/register" element={<RegisterPage />} />
            <Route path="/forgot-password" element={<ForgotPasswordPage />} />
            <Route path="/reset-password" element={<ResetPasswordPage />} />
            <Route path="/verify-email" element={<VerifyEmailPage />} />
            <Route path="/verify/:id" element={<VerifyBillPage />} />

            {/* Protected routes wrapped in Layout */}
            <Route element={<Layout />}>
              <Route path="/" element={<DashboardPage />} />
              <Route path="/metal-rates" element={<MetalRatesPage />} />
              <Route path="/bills/new" element={<CreateBillPage />} />
              <Route path="/bills/history" element={<BillHistoryPage />} />
              <Route path="/customers" element={<CustomersPage />} />
              <Route path="/expenses" element={<ExpensesPage />} />
              <Route path="/settings" element={<SettingsPage />} />
            </Route>

            <Route path="*" element={<Navigate to="/" replace />} />
          </Routes>
        </BrowserRouter>
      </AuthProvider>
    </QueryClientProvider>
    </ThemeProvider>
  );
}
