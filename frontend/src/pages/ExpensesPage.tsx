import { useState } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/Card";
import { Button } from "@/components/ui/Button";
import { Input } from "@/components/ui/Input";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { getExpenses, createExpense, deleteExpense } from "@/lib/expenses.api";
import { Loader2, Trash2 } from "lucide-react";
import { toast } from "sonner";

export default function ExpensesPage() {
  const queryClient = useQueryClient();
  const [formData, setFormData] = useState({ category: "", amount: "", description: "", expense_date: "" });

  const { data: expenses, isLoading } = useQuery({
    queryKey: ["expenses"],
    queryFn: getExpenses,
  });

  const createMutation = useMutation({
    mutationFn: createExpense,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["expenses"] });
      queryClient.invalidateQueries({ queryKey: ["dashboard"] }); // Refresh dashboard stats
      toast.success("Expense logged successfully!");
      setFormData({ category: "", amount: "", description: "", expense_date: "" });
    },
    onError: (err: any) => {
      toast.error(err.message || "Failed to log expense");
    }
  });

  const deleteMutation = useMutation({
    mutationFn: deleteExpense,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["expenses"] });
      queryClient.invalidateQueries({ queryKey: ["dashboard"] });
      toast.success("Expense deleted");
    }
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!formData.category || !formData.amount || !formData.expense_date) return;
    createMutation.mutate({
      category: formData.category,
      amount: parseFloat(formData.amount),
      description: formData.description,
      expense_date: formData.expense_date,
    });
  };

  if (isLoading) {
    return <div className="flex h-[50vh] items-center justify-center"><Loader2 className="h-8 w-8 animate-spin text-[#c6a962]" /></div>;
  }

  return (
    <div className="p-6 max-w-7xl mx-auto space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-3xl font-bold tracking-tight text-[#c6a962]">
          Expense Management
        </h1>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Add Expense Form */}
        <div className="lg:col-span-1">
          <Card className="bg-black/40 backdrop-blur-md border-[#c6a962]/20 sticky top-6">
            <CardHeader>
              <CardTitle className="text-gray-200">Log New Expense</CardTitle>
            </CardHeader>
            <CardContent>
              <form onSubmit={handleSubmit} className="space-y-4">
                <div className="space-y-2">
                  <label className="text-sm text-gray-400">Category</label>
                  <Input 
                    required
                    placeholder="e.g. Salary, Rent" 
                    value={formData.category}
                    onChange={(e) => setFormData({...formData, category: e.target.value})}
                    className="bg-black/50 border-gray-700 text-gray-200" 
                  />
                </div>
                <div className="space-y-2">
                  <label className="text-sm text-gray-400">Amount (₹)</label>
                  <Input 
                    required
                    type="number" 
                    placeholder="0.00" 
                    value={formData.amount}
                    onChange={(e) => setFormData({...formData, amount: e.target.value})}
                    className="bg-black/50 border-gray-700 text-gray-200" 
                  />
                </div>
                <div className="space-y-2">
                  <label className="text-sm text-gray-400">Description</label>
                  <Input 
                    required
                    placeholder="Brief details" 
                    value={formData.description}
                    onChange={(e) => setFormData({...formData, description: e.target.value})}
                    className="bg-black/50 border-gray-700 text-gray-200" 
                  />
                </div>
                <div className="space-y-2">
                  <label className="text-sm text-gray-400">Date</label>
                  <Input 
                    required
                    type="date" 
                    value={formData.expense_date}
                    onChange={(e) => setFormData({...formData, expense_date: e.target.value})}
                    className="bg-black/50 border-gray-700 text-gray-200" 
                  />
                </div>
                <Button type="submit" disabled={createMutation.isPending} className="w-full mt-4">
                  {createMutation.isPending ? "Saving..." : "Save Expense"}
                </Button>
              </form>
            </CardContent>
          </Card>
        </div>

        {/* Expense History Table */}
        <div className="lg:col-span-2">
          <Card className="bg-black/40 backdrop-blur-md border-[#c6a962]/20">
            <CardHeader className="flex flex-row items-center justify-between">
              <CardTitle className="text-gray-200">Recent Expenses</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="overflow-x-auto">
                <table className="w-full text-sm text-left">
                  <thead className="text-xs text-gray-400 uppercase bg-black/50 border-b border-[#c6a962]/20">
                    <tr>
                      <th className="px-4 py-3 font-medium">Date</th>
                      <th className="px-4 py-3 font-medium">Category</th>
                      <th className="px-4 py-3 font-medium">Description</th>
                      <th className="px-4 py-3 font-medium text-right">Amount</th>
                      <th className="px-4 py-3 font-medium text-center">Action</th>
                    </tr>
                  </thead>
                  <tbody>
                    {expenses?.length === 0 && (
                      <tr>
                        <td colSpan={5} className="text-center py-8 text-gray-500">No expenses found</td>
                      </tr>
                    )}
                    {expenses?.map((exp: any) => (
                      <tr key={exp.id} className="border-b border-gray-800 hover:bg-red-500/5 transition-colors">
                        <td className="px-4 py-3 text-gray-300">{exp.expense_date.split('T')[0]}</td>
                        <td className="px-4 py-3 text-[#c6a962] font-medium">{exp.category}</td>
                        <td className="px-4 py-3 text-gray-300">{exp.description}</td>
                        <td className="px-4 py-3 text-right font-bold text-red-400">-₹{exp.amount.toLocaleString()}</td>
                        <td className="px-4 py-3 text-center">
                          <Button 
                            variant="ghost" 
                            size="icon" 
                            onClick={() => deleteMutation.mutate(exp.id)}
                            className="h-8 w-8 text-gray-500 hover:text-red-500 hover:bg-red-500/10"
                          >
                            <Trash2 className="h-4 w-4" />
                          </Button>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}
