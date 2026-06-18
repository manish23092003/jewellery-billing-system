// ═══════════════════════════════════════════════════════════════════════
// Jewellery Billing — TypeScript Type Definitions
// ═══════════════════════════════════════════════════════════════════════

// ── User ─────────────────────────────────────────────────────────────

export type UserRole = "admin" | "staff";

export interface User {
  id: string;
  name: string;
  email: string;
  role: UserRole;
  created_at: string;
  updated_at: string;
}

export interface CreateUserRequest {
  name: string;
  email: string;
  password: string;
  role: UserRole;
}

export interface UpdateUserRequest {
  name?: string;
  email?: string;
  password?: string;
  role?: UserRole;
}

// ── Auth ─────────────────────────────────────────────────────────────

export interface LoginRequest {
  email: string;
  password: string;
}

export interface AuthResponse {
  access_token: string;
  refresh_token: string;
  user: User;
}

// ── API Envelope ─────────────────────────────────────────────────────

export interface APIResponse<T = unknown> {
  success: boolean;
  data?: T;
  error?: string;
  meta?: APIMeta;
}

export interface APIMeta {
  page: number;
  per_page: number;
  total: number;
  total_pages: number;
}
