import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { Search, Plus, Loader2, MoreVertical, Edit, Trash2 } from "lucide-react";
import { toast } from "sonner";
import { format } from "date-fns";

import { getCustomers, createCustomer, updateCustomer, deleteCustomer } from "@/lib/customers.api";
import type { Customer } from "@/types";

import { Button } from "@/components/ui/Button";
import { Input } from "@/components/ui/Input";
import { Card, CardContent } from "@/components/ui/Card";
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from "@/components/ui/DropdownMenu";

export default function CustomersPage() {
  const queryClient = useQueryClient();
  const [search, setSearch] = useState("");
  const [page, setPage] = useState(1);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingCustomer, setEditingCustomer] = useState<Customer | null>(null);

  const [formData, setFormData] = useState({
    name: "",
    phone: "",
    email: "",
    address: "",
  });

  const { data, isLoading } = useQuery({
    queryKey: ["customers", search, page],
    queryFn: () => getCustomers({ search, limit: 20, offset: (page - 1) * 20 }),
  });

  const createMutation = useMutation({
    mutationFn: createCustomer,
    onSuccess: () => {
      toast.success("Customer added successfully");
      queryClient.invalidateQueries({ queryKey: ["customers"] });
      closeModal();
    },
    onError: (err: any) => toast.error(err.message || "Failed to add customer"),
  });

  const updateMutation = useMutation({
    mutationFn: updateCustomer,
    onSuccess: () => {
      toast.success("Customer updated successfully");
      queryClient.invalidateQueries({ queryKey: ["customers"] });
      closeModal();
    },
    onError: (err: any) => toast.error(err.message || "Failed to update customer"),
  });

  const deleteMutation = useMutation({
    mutationFn: deleteCustomer,
    onSuccess: () => {
      toast.success("Customer deleted");
      queryClient.invalidateQueries({ queryKey: ["customers"] });
    },
    onError: (err: any) => toast.error(err.message || "Failed to delete customer"),
  });

  const openModalForNew = () => {
    setEditingCustomer(null);
    setFormData({ name: "", phone: "", email: "", address: "" });
    setIsModalOpen(true);
  };

  const openModalForEdit = (customer: Customer) => {
    setEditingCustomer(customer);
    setFormData({
      name: customer.name,
      phone: customer.phone,
      email: customer.email,
      address: customer.address,
    });
    setIsModalOpen(true);
  };

  const closeModal = () => {
    setIsModalOpen(false);
    setEditingCustomer(null);
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (editingCustomer) {
      updateMutation.mutate({ id: editingCustomer.id, payload: formData });
    } else {
      createMutation.mutate(formData);
    }
  };

  const handleDelete = (id: string) => {
    if (window.confirm("Are you sure you want to delete this customer?")) {
      deleteMutation.mutate(id);
    }
  };

  return (
    <div className="p-6 max-w-7xl mx-auto space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-3xl font-bold tracking-tight text-[#c6a962]">
          Customers
        </h1>
        <Button onClick={openModalForNew} className="gap-2">
          <Plus className="h-4 w-4" /> Add Customer
        </Button>
      </div>

      <div className="flex items-center gap-4 bg-card p-4 rounded-lg border border-[#c6a962]/20 backdrop-blur-md">
        <div className="relative flex-1 max-w-md">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
          <Input
            placeholder="Search by name or phone..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="pl-9 bg-muted border-border text-foreground"
          />
        </div>
      </div>

      <Card className="bg-card border-[#c6a962]/20 backdrop-blur-md">
        <CardContent className="p-0">
          <div className="overflow-x-auto">
            <table className="w-full text-sm text-left text-foreground/90">
              <thead className="text-xs uppercase bg-[#c6a962]/10 text-[#c6a962] border-b border-[#c6a962]/20">
                <tr>
                  <th className="px-6 py-4 font-medium">Name</th>
                  <th className="px-6 py-4 font-medium">Contact</th>
                  <th className="px-6 py-4 font-medium">Total Spent</th>
                  <th className="px-6 py-4 font-medium">Joined</th>
                  <th className="px-6 py-4 text-right font-medium">Actions</th>
                </tr>
              </thead>
              <tbody>
                {isLoading ? (
                  <tr>
                    <td colSpan={5} className="px-6 py-12 text-center">
                      <Loader2 className="h-6 w-6 animate-spin mx-auto text-[#c6a962]" />
                    </td>
                  </tr>
                ) : data?.customers.length === 0 ? (
                  <tr>
                    <td colSpan={5} className="px-6 py-12 text-center text-muted-foreground">
                      No customers found. Add your first customer!
                    </td>
                  </tr>
                ) : (
                  data?.customers.map((customer) => (
                    <tr key={customer.id} className="border-b border-border/50 hover:bg-white/5 transition-colors">
                      <td className="px-6 py-4 font-medium text-foreground">{customer.name}</td>
                      <td className="px-6 py-4">
                        <div className="text-foreground/90">{customer.phone}</div>
                        {customer.email && <div className="text-xs text-muted-foreground">{customer.email}</div>}
                      </td>
                      <td className="px-6 py-4 font-medium text-[#c6a962]">
                        ₹{customer.total_purchases.toLocaleString("en-IN", { minimumFractionDigits: 2 })}
                      </td>
                      <td className="px-6 py-4 text-muted-foreground">
                        {format(new Date(customer.created_at), "MMM d, yyyy")}
                      </td>
                      <td className="px-6 py-4 text-right">
                        <DropdownMenu>
                          <DropdownMenuTrigger asChild>
                            <Button variant="ghost" className="h-8 w-8 p-0">
                              <MoreVertical className="h-4 w-4" />
                            </Button>
                          </DropdownMenuTrigger>
                          <DropdownMenuContent align="end" className="bg-card border-border">
                            <DropdownMenuItem onClick={() => openModalForEdit(customer)} className="text-foreground/90 hover:text-white cursor-pointer">
                              <Edit className="h-4 w-4 mr-2" /> Edit
                            </DropdownMenuItem>
                            <DropdownMenuItem onClick={() => handleDelete(customer.id)} className="text-red-400 hover:text-red-300 cursor-pointer">
                              <Trash2 className="h-4 w-4 mr-2" /> Delete
                            </DropdownMenuItem>
                          </DropdownMenuContent>
                        </DropdownMenu>
                      </td>
                    </tr>
                  ))
                )}
              </tbody>
            </table>
          </div>
          {data && data.total > 20 && (
            <div className="flex items-center justify-between border-t border-border/50 pt-4 mt-4">
              <div className="text-sm text-muted-foreground">
                Showing {((page - 1) * 20) + 1} to {Math.min(page * 20, data.total)} of {data.total} customers
              </div>
              <div className="flex items-center space-x-2">
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => setPage(p => Math.max(1, p - 1))}
                  disabled={page === 1}
                  className="h-8 border-border hover:bg-[#c6a962]/10"
                >
                  Previous
                </Button>
                <div className="text-sm font-medium text-foreground px-2">
                  Page {page} of {Math.ceil(data.total / 20)}
                </div>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => setPage(p => Math.min(Math.ceil(data.total / 20), p + 1))}
                  disabled={page >= Math.ceil(data.total / 20)}
                  className="h-8 border-border hover:bg-[#c6a962]/10"
                >
                  Next
                </Button>
              </div>
            </div>
          )}
        </CardContent>
      </Card>

      {/* Basic Modal */}
      {isModalOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm p-4">
          <div className="bg-card border border-border rounded-xl shadow-2xl w-full max-w-md overflow-hidden">
            <div className="p-6 border-b border-border">
              <h2 className="text-xl font-semibold text-foreground">
                {editingCustomer ? "Edit Customer" : "Add New Customer"}
              </h2>
            </div>
            <form onSubmit={handleSubmit} className="p-6 space-y-4">
              <div className="space-y-2">
                <label className="text-sm text-muted-foreground">Customer Name *</label>
                <Input
                  required
                  value={formData.name}
                  onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                  className="bg-muted border-border text-white"
                />
              </div>
              <div className="space-y-2">
                <label className="text-sm text-muted-foreground">Phone Number *</label>
                <Input
                  required
                  value={formData.phone}
                  onChange={(e) => setFormData({ ...formData, phone: e.target.value })}
                  className="bg-muted border-border text-white"
                />
              </div>
              <div className="space-y-2">
                <label className="text-sm text-muted-foreground">Email Address (Optional)</label>
                <Input
                  type="email"
                  value={formData.email}
                  onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                  className="bg-muted border-border text-white"
                />
              </div>
              <div className="space-y-2">
                <label className="text-sm text-muted-foreground">Address (Optional)</label>
                <Input
                  value={formData.address}
                  onChange={(e) => setFormData({ ...formData, address: e.target.value })}
                  className="bg-muted border-border text-white"
                />
              </div>
              <div className="pt-4 flex justify-end gap-3">
                <Button type="button" variant="outline" onClick={closeModal} className="border-border text-foreground/90">
                  Cancel
                </Button>
                <Button type="submit" disabled={createMutation.isPending || updateMutation.isPending}>
                  {createMutation.isPending || updateMutation.isPending ? "Saving..." : "Save Customer"}
                </Button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
