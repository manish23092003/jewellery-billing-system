// ═══════════════════════════════════════════════════════════════════════
// Jewellery Billing — TypeScript Type Definitions
// ═══════════════════════════════════════════════════════════════════════

// ── Organization ─────────────────────────────────────────────────────

export interface Organization {
  id: string;
  business_name: string;
  owner_name: string;
  email: string;
  phone: string;
  gstin: string;
  address: string;
  subscription_status: string;
  created_at: string;
  updated_at: string;
}

export interface BillItemCharge {
  id?: string;
  charge_name: string;
  amount: number;
}

export interface BillItem {
  id: string;
  item_name: string;
  metal_type: string;
  purity: string;
  weight: number;
  rate_per_gram: number;
  making_charge: number;
  gst_percentage: number;
  quantity: number;
  line_total: number;
  charges?: BillItemCharge[];
}

export interface CreateBillItemRequest {
  item_name: string;
  metal_type: string;
  purity: string;
  weight: number;
  rate_per_gram: number;
  making_charge: number;
  gst_percentage: number;
  quantity: number;
  charges?: BillItemCharge[];
}

export interface BillOldGold {
  id?: string;
  name: string;
  weight: number;
  purity: string;
  melting_loss_percentage: number;
  rate_per_gram: number;
  total_value: number;
}

export interface Bill {
  id: string;
  invoice_number: string;
  invoice_date: string;
  type: string;
  status: string;
  customer_name: string;
  customer_phone: string;
  subtotal: number;
  gst_amount: number;
  old_gold_amount: number;
  advance_amount: number;
  grand_total: number;
  balance_due: number;
  payment_method: string;
  notes: string;
  items: BillItem[];
  old_gold_items: BillOldGold[];
}

export interface CreateBillRequest {
  type: string;
  status: string;
  advance_amount: number;
  invoice_date: string;
  customer_name: string;
  customer_phone: string;
  payment_method: string;
  notes: string;
  items: CreateBillItemRequest[];
  old_gold_items: Omit<BillOldGold, "id" | "total_value">[];
}

export interface Customer {
  id: string;
  organization_id: string;
  name: string;
  phone: string;
  email: string;
  address: string;
  total_purchases: number;
  created_at: string;
  updated_at: string;
}

export interface CreateCustomerRequest {
  name: string;
  phone: string;
  email: string;
  address: string;
}

export interface UpdateCustomerRequest {
  name: string;
  phone: string;
  email: string;
  address: string;
}

export interface RegisterRequest {
  business_name: string;
  owner_name: string;
  email: string;
  phone: string;
  password?: string;
  confirm_password?: string;
}

// ── User ─────────────────────────────────────────────────────────────

export type UserRole = "admin" | "staff";

export interface User {
  id: string;
  organization_id: string;
  name: string;
  email: string;
  role: UserRole;
  is_active: boolean;
  email_verified: boolean;
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

export interface ForgotPasswordRequest {
  email: string;
}

export interface ResetPasswordRequest {
  token: string;
  password: string;
  confirm_password: string;
}

export interface VerifyEmailRequest {
  token: string;
}

export interface AuthResponse {
  access_token: string;
  refresh_token: string;
  user: User;
  organization: Organization;
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
