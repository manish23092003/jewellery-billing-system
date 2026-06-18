import api from "./api";
import type { APIResponse } from "@/types";

export interface BillItem {
  id?: string;
  item_name: string;
  metal_type: string;
  purity: string;
  weight: number;
  rate_per_gram: number;
  making_charge: number;
  gst_percentage: number;
  quantity: number;
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
  invoice_date: string;
  customer_name: string;
  customer_phone: string;
  payment_method: string;
  notes: string;
  items: BillItem[];
}

export const getBills = async (): Promise<Bill[]> => {
  const { data } = await api.get<APIResponse<Bill[]>>("/bills");
  if (!data.success) throw new Error(data.error);
  return data.data || [];
};

export const createBill = async (payload: CreateBillRequest): Promise<Bill> => {
  const { data } = await api.post<APIResponse<Bill>>("/bills", payload);
  if (!data.success) throw new Error(data.error);
  return data.data!;
};
