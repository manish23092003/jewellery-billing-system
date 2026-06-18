import api from "./api";
import type { APIResponse } from "@/types";

export interface ShopSettings {
  id?: number;
  shop_name: string;
  gstin: string;
  phone: string;
  address: string;
  invoice_prefix: string;
  logo_path?: string;
}

export const getSettings = async (): Promise<ShopSettings> => {
  const { data } = await api.get<APIResponse<ShopSettings>>("/settings");
  if (!data.success) throw new Error(data.error);
  return data.data!;
};

export const updateSettings = async (payload: ShopSettings): Promise<ShopSettings> => {
  const { data } = await api.put<APIResponse<ShopSettings>>("/settings", payload);
  if (!data.success) throw new Error(data.error);
  return data.data!;
};

export const uploadLogo = async (file: File): Promise<ShopSettings> => {
  const formData = new FormData();
  formData.append("logo", file);
  
  const { data } = await api.post<APIResponse<ShopSettings>>("/settings/logo", formData, {
    headers: {
      "Content-Type": "multipart/form-data",
    },
  });
  if (!data.success) throw new Error(data.error);
  return data.data!;
};
