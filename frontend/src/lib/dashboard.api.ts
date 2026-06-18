import api from "./api";
import type { APIResponse } from "@/types";

export interface DashboardMetrics {
  today_sales: number;
  today_expenses: number;
  today_profit: number;
  monthly_sales: number;
  monthly_expenses: number;
  monthly_profit: number;
}

export interface DailyTrend {
  date: string;
  sales: number;
  expenses: number;
  profit: number;
}

export interface DashboardData {
  metrics: DashboardMetrics;
  trends: DailyTrend[];
}

export const getDashboard = async (): Promise<DashboardData> => {
  const { data } = await api.get<APIResponse<DashboardData>>(`/analytics/dashboard?t=${new Date().getTime()}`);
  if (!data.success) throw new Error(data.error);
  return data.data!;
};
