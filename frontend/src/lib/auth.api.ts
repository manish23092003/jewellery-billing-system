import api from "./api";
import type {
  LoginRequest,
  AuthResponse,
  APIResponse,
  RegisterRequest,
  ForgotPasswordRequest,
  ResetPasswordRequest,
  VerifyEmailRequest,
} from "@/types";

export const login = async (data: LoginRequest): Promise<AuthResponse> => {
  const response = await api.post<APIResponse<AuthResponse>>("/auth/login", data);
  if (!response.data.success || !response.data.data) {
    throw new Error(response.data.error || "Login failed");
  }
  return response.data.data;
};

export const registerBusiness = async (data: RegisterRequest): Promise<AuthResponse> => {
  const response = await api.post<APIResponse<AuthResponse>>("/auth/register", data);
  if (!response.data.success || !response.data.data) {
    throw new Error(response.data.error || "Registration failed");
  }
  return response.data.data;
};

export const forgotPassword = async (data: ForgotPasswordRequest): Promise<void> => {
  const response = await api.post<APIResponse>("/auth/forgot-password", data);
  if (!response.data.success) {
    throw new Error(response.data.error || "Failed to send reset link");
  }
};

export const resetPassword = async (data: ResetPasswordRequest): Promise<void> => {
  const response = await api.post<APIResponse>("/auth/reset-password", data);
  if (!response.data.success) {
    throw new Error(response.data.error || "Failed to reset password");
  }
};

export const verifyEmail = async (data: VerifyEmailRequest): Promise<void> => {
  const response = await api.post<APIResponse>("/auth/verify-email", data);
  if (!response.data.success) {
    throw new Error(response.data.error || "Email verification failed");
  }
};

export const resendVerification = async (): Promise<void> => {
  const response = await api.post<APIResponse>("/auth/resend-verification");
  if (!response.data.success) {
    throw new Error(response.data.error || "Failed to resend verification email");
  }
};

export const logoutUser = async (): Promise<void> => {
  try {
    await api.post("/auth/logout");
  } catch (error) {
    // Ignore logout errors
  }
};
