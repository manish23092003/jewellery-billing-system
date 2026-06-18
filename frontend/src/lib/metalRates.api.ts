import api from "./api";
import type { APIResponse } from "@/types";

export interface MetalRate {
  id: string;
  metal_type: string;
  purity: string;
  rate_per_gram: number;
  effective_date: string;
}

export const getLatestRates = async (): Promise<MetalRate[]> => {
  const { data } = await api.get<APIResponse<MetalRate[]>>("/metal-rates/latest");
  if (!data.success) throw new Error(data.error);
  return data.data || [];
};

export const createMetalRate = async (payload: {
  metal_type: string;
  purity: string;
  rate_per_gram: number;
}): Promise<MetalRate> => {
  const { data } = await api.post<APIResponse<MetalRate>>("/metal-rates", payload);
  if (!data.success) throw new Error(data.error);
  return data.data!;
};
