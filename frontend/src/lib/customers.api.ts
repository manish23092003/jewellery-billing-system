import api from "./api";
import type { Customer, CreateCustomerRequest, UpdateCustomerRequest, APIResponse } from "@/types";

interface GetCustomersParams {
  search?: string;
  limit?: number;
  offset?: number;
}

interface PaginatedCustomers {
  customers: Customer[];
  total: number;
}

export const getCustomers = async (params?: GetCustomersParams): Promise<PaginatedCustomers> => {
  const { data } = await api.get<APIResponse<Customer[]>>("/customers", { params });
  if (!data.success) throw new Error(data.error);
  return {
    customers: data.data || [],
    total: data.meta?.total || 0,
  };
};

export const getCustomer = async (id: string): Promise<Customer> => {
  const { data } = await api.get<APIResponse<Customer>>(`/customers/${id}`);
  if (!data.success) throw new Error(data.error);
  return data.data!;
};

export const createCustomer = async (payload: CreateCustomerRequest): Promise<Customer> => {
  const { data } = await api.post<APIResponse<Customer>>("/customers", payload);
  if (!data.success) throw new Error(data.error);
  return data.data!;
};

export const updateCustomer = async ({ id, payload }: { id: string; payload: UpdateCustomerRequest }): Promise<Customer> => {
  const { data } = await api.put<APIResponse<Customer>>(`/customers/${id}`, payload);
  if (!data.success) throw new Error(data.error);
  return data.data!;
};

export const deleteCustomer = async (id: string): Promise<void> => {
  const { data } = await api.delete<APIResponse>(`/customers/${id}`);
  if (!data.success) throw new Error(data.error);
};
