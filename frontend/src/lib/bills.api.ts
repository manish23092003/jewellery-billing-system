import api from "./api";
import type { APIResponse } from "@/types";

export interface BillItem {
  id?: string;
  item_name: string;
  hsn_code?: string;
  metal_type: string;
  purity: string;
  weight: number;
  rate_per_gram: number;
  making_charge: number;
  gst_percentage: number;
  quantity: number;
  charges?: { charge_name: string; amount: number }[];
  line_total?: number;
}

export interface Bill {
  id: string;
  invoice_number: string;
  invoice_date: string;
  customer_name: string;
  customer_phone: string;
  subtotal: number;
  gst_amount: number;
  grand_total: number;
  payment_method: string;
  items: BillItem[];
}

export interface CreateBillRequest {
  type?: string;
  status?: string;
  advance_amount?: number;
  invoice_date: string;
  customer_name: string;
  customer_phone: string;
  payment_method: string;
  notes: string;
  items: BillItem[];
  convert_from_id?: string;
  old_gold_items?: any[];
}

export const getBills = async (): Promise<Bill[]> => {
  const { data } = await api.get<APIResponse<Bill[]>>("/bills");
  if (!data.success) throw new Error(data.error);
  return data.data || [];
};

export const getBillById = async (id: string): Promise<any> => {
  const { data } = await api.get<APIResponse<any>>(`/bills/${id}`);
  if (!data.success) throw new Error(data.error);
  return data.data;
};

export const createBill = async (payload: CreateBillRequest): Promise<Bill> => {
  const { data } = await api.post<APIResponse<Bill>>("/bills", payload);
  if (!data.success) throw new Error(data.error);
  return data.data!;
};

export const addPayment = async ({ id, amount }: { id: string; amount: number }): Promise<void> => {
  const { data } = await api.post<APIResponse<void>>(`/bills/${id}/payment`, { amount });
  if (!data.success) throw new Error(data.error);
};

export interface PublicVerificationResponse {
  shop_name: string;
  invoice_number: string;
  invoice_date: string;
  grand_total: number;
  balance_due: number;
  verification_status: "VERIFIED" | "TAMPERED" | "VOID" | "CANCELLED" | "ACTIVE";
}

export const verifyInvoice = async (token: string): Promise<PublicVerificationResponse> => {
  const { data } = await api.get<APIResponse<PublicVerificationResponse>>(`/public/verify/${token}`);
  if (!data.success) throw new Error(data.error);
  return data.data!;
};
