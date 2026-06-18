import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/Card";
import { AreaChart, Area, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';
import { useQuery } from "@tanstack/react-query";
import { getDashboard } from "@/lib/dashboard.api";
import { Loader2 } from "lucide-react";

export default function DashboardPage() {
  const { data, isLoading, isError } = useQuery({
    queryKey: ["dashboard"],
    queryFn: getDashboard,
  });

  if (isLoading) {
    return (
      <div className="flex h-[50vh] items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-[#c6a962]" />
      </div>
    );
  }

  if (isError || !data) {
    return <div className="text-red-500 text-center p-6">Failed to load dashboard data</div>;
  }

  const { metrics, trends } = data;

  return (
    <div className="p-6 max-w-7xl mx-auto space-y-8">
      <h1 className="text-3xl font-bold tracking-tight text-[#c6a962]">
        Dashboard Overview
      </h1>

      {/* Metrics Row 1 - Today */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <Card className="bg-black/40 backdrop-blur-md border-[#c6a962]/20 shadow-[0_0_15px_rgba(198,169,98,0.05)]">
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-gray-400">Today's Revenue</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-3xl font-bold text-gray-200">₹{metrics.today_sales.toLocaleString()}</div>
          </CardContent>
        </Card>
        
        <Card className="bg-black/40 backdrop-blur-md border-[#c6a962]/20 shadow-[0_0_15px_rgba(198,169,98,0.05)]">
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-gray-400">Today's Expenses</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-3xl font-bold text-red-400">₹{metrics.today_expenses.toLocaleString()}</div>
          </CardContent>
        </Card>

        <Card className="bg-gradient-to-br from-[#c6a962]/10 to-transparent backdrop-blur-md border-[#c6a962]/40 shadow-[0_0_15px_rgba(198,169,98,0.1)]">
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-[#c6a962]">Today's Profit</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-3xl font-bold text-[#c6a962]">₹{metrics.today_profit.toLocaleString()}</div>
          </CardContent>
        </Card>
      </div>

      {/* Metrics Row 2 - Monthly */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <Card className="bg-black/40 backdrop-blur-md border-gray-800">
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-gray-400">Monthly Revenue</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-gray-300">₹{metrics.monthly_sales.toLocaleString()}</div>
          </CardContent>
        </Card>
        
        <Card className="bg-black/40 backdrop-blur-md border-gray-800">
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-gray-400">Monthly Expenses</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-red-400/80">₹{metrics.monthly_expenses.toLocaleString()}</div>
          </CardContent>
        </Card>

        <Card className="bg-black/40 backdrop-blur-md border-gray-800">
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-gray-400">Monthly Profit</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-green-500/80">₹{metrics.monthly_profit.toLocaleString()}</div>
          </CardContent>
        </Card>
      </div>

      {/* Charts */}
      <Card className="bg-black/40 backdrop-blur-md border-[#c6a962]/20 p-2">
        <CardHeader>
          <CardTitle className="text-gray-200">Revenue vs Expense Trend</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="h-[400px] w-full mt-4">
            <ResponsiveContainer width="100%" height="100%">
              <AreaChart data={trends} margin={{ top: 10, right: 30, left: 0, bottom: 0 }}>
                <defs>
                  <linearGradient id="colorSales" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="5%" stopColor="#c6a962" stopOpacity={0.3}/>
                    <stop offset="95%" stopColor="#c6a962" stopOpacity={0}/>
                  </linearGradient>
                  <linearGradient id="colorExpense" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="5%" stopColor="#f87171" stopOpacity={0.3}/>
                    <stop offset="95%" stopColor="#f87171" stopOpacity={0}/>
                  </linearGradient>
                </defs>
                <CartesianGrid strokeDasharray="3 3" stroke="#333" vertical={false} />
                <XAxis dataKey="date" stroke="#666" tick={{fill: '#888', fontSize: 12}} />
                <YAxis stroke="#666" tick={{fill: '#888', fontSize: 12}} tickFormatter={(value) => `₹${value/1000}k`} />
                <Tooltip 
                  contentStyle={{ backgroundColor: '#111', borderColor: '#333', color: '#fff', borderRadius: '8px' }}
                  itemStyle={{ color: '#ccc' }}
                />
                <Area type="monotone" dataKey="sales" stroke="#c6a962" strokeWidth={2} fillOpacity={1} fill="url(#colorSales)" name="Revenue" />
                <Area type="monotone" dataKey="expenses" stroke="#f87171" strokeWidth={2} fillOpacity={1} fill="url(#colorExpense)" name="Expenses" />
              </AreaChart>
            </ResponsiveContainer>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
