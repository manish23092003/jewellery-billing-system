import api from "./api";
import type { APIResponse } from "@/types";

export interface Expense {
  id: string;
  category: string;
  amount: number;
  description: string;
  expense_date: string;
}

export const getExpenses = async (): Promise<Expense[]> => {
  const { data } = await api.get<APIResponse<Expense[]>>("/expenses");
  if (!data.success) throw new Error(data.error);
  return data.data || [];
};

export const createExpense = async (payload: {
  category: string;
  amount: number;
  description: string;
  expense_date: string;
}): Promise<Expense> => {
  const { data } = await api.post<APIResponse<Expense>>("/expenses", payload);
  if (!data.success) throw new Error(data.error);
  return data.data!;
};

export const deleteExpense = async (id: string): Promise<void> => {
  const { data } = await api.delete<APIResponse<void>>(`/expenses/${id}`);
  if (!data.success) throw new Error(data.error);
};
