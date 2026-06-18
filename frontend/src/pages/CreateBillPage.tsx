import { useState } from "react";
import { Button } from "@/components/ui/Button";
import { Input } from "@/components/ui/Input";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/Card";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { createBill } from "@/lib/bills.api";
import { useNavigate } from "react-router-dom";
import { toast } from "sonner";
import { Loader2 } from "lucide-react";

export default function CreateBillPage() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  
  const [customer, setCustomer] = useState({ name: "", phone: "", payment_method: "cash" });
  const [items, setItems] = useState([
    { id: Date.now(), item_name: "", metal_type: "gold", purity: "22K", weight: "", rate_per_gram: "", making_charge: "", gst_percentage: "3", quantity: "1" }
  ]);

  const addItem = () => setItems([...items, { id: Date.now(), item_name: "", metal_type: "gold", purity: "22K", weight: "", rate_per_gram: "", making_charge: "", gst_percentage: "3", quantity: "1" }]);
  const removeItem = (id: number) => setItems(items.filter(i => i.id !== id));

  const updateItem = (id: number, field: string, value: string) => {
    setItems(items.map(item => item.id === id ? { ...item, [field]: value } : item));
  };

  // Preview Calculations
  let subtotal = 0;
  let totalGst = 0;
  items.forEach(item => {
    const w = parseFloat(item.weight) || 0;
    const r = parseFloat(item.rate_per_gram) || 0;
    const m = parseFloat(item.making_charge) || 0;
    const q = parseInt(item.quantity) || 1;
    const gstPct = parseFloat(item.gst_percentage) || 3;
    
    const metalVal = w * r;
    const lineSub = (metalVal + m) * q;
    const lineGst = lineSub * (gstPct / 100);
    
    subtotal += lineSub;
    totalGst += lineGst;
  });
  const grandTotal = subtotal + totalGst;

  const mutation = useMutation({
    mutationFn: createBill,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["bills"] });
      queryClient.invalidateQueries({ queryKey: ["dashboard"] });
      toast.success("Invoice generated successfully");
      navigate("/bills/history");
    },
    onError: (err: any) => {
      toast.error(err.message || "Failed to generate invoice");
    }
  });

  const handleSubmit = () => {
    if (!customer.name) return toast.error("Customer name is required");
    if (items.length === 0) return toast.error("At least one item is required");
    
    const payload = {
      invoice_date: new Date().toISOString(),
      customer_name: customer.name,
      customer_phone: customer.phone,
      payment_method: customer.payment_method,
      notes: "",
      items: items.map(item => ({
        item_name: item.item_name || "Unknown Item",
        metal_type: item.metal_type,
        purity: item.purity,
        weight: parseFloat(item.weight) || 0,
        rate_per_gram: parseFloat(item.rate_per_gram) || 0,
        making_charge: parseFloat(item.making_charge) || 0,
        gst_percentage: parseFloat(item.gst_percentage) || 3,
        quantity: parseInt(item.quantity) || 1
      }))
    };
    
    mutation.mutate(payload);
  };

  return (
    <div className="p-6 max-w-7xl mx-auto space-y-6">
      <h1 className="text-3xl font-bold tracking-tight text-[#c6a962]">
        Create New Invoice
      </h1>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-2 space-y-6">
          <Card className="bg-black/40 backdrop-blur-md border-[#c6a962]/20">
            <CardHeader>
              <CardTitle className="text-gray-200">Customer Details</CardTitle>
            </CardHeader>
            <CardContent className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <label className="text-sm text-gray-400">Customer Name</label>
                <Input value={customer.name} onChange={e => setCustomer({...customer, name: e.target.value})} placeholder="Enter name" className="text-gray-200 bg-black/50 border-gray-700" />
              </div>
              <div className="space-y-2">
                <label className="text-sm text-gray-400">Phone Number</label>
                <Input value={customer.phone} onChange={e => setCustomer({...customer, phone: e.target.value})} placeholder="Enter phone" className="text-gray-200 bg-black/50 border-gray-700" />
              </div>
              <div className="space-y-2">
                <label className="text-sm text-gray-400">Payment Method</label>
                <select 
                  value={customer.payment_method} 
                  onChange={e => setCustomer({...customer, payment_method: e.target.value})}
                  className="flex h-9 w-full rounded-md border border-gray-700 bg-black/50 px-3 py-1 text-sm shadow-sm transition-colors text-gray-200 focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
                >
                  <option value="cash">Cash</option>
                  <option value="card">Card</option>
                  <option value="upi">UPI</option>
                  <option value="bank_transfer">Bank Transfer</option>
                </select>
              </div>
            </CardContent>
          </Card>

          <Card className="bg-black/40 backdrop-blur-md border-[#c6a962]/20">
            <CardHeader className="flex flex-row justify-between items-center">
              <CardTitle className="text-gray-200">Bill Items</CardTitle>
              <Button variant="outline" size="sm" onClick={addItem} className="border-[#c6a962]/50 text-[#c6a962]">
                + Add Item
              </Button>
            </CardHeader>
            <CardContent className="space-y-6">
              {items.map((item, index) => (
                <div key={item.id} className="grid grid-cols-12 gap-3 items-end border-b border-gray-800 pb-6 relative">
                  <div className="col-span-12 absolute -top-4 right-0 text-gray-600 text-xs text-right">Item {index + 1}</div>
                  <div className="col-span-3 space-y-2">
                    <label className="text-xs text-gray-500">Item Name</label>
                    <Input value={item.item_name} onChange={e => updateItem(item.id, 'item_name', e.target.value)} placeholder="e.g. Gold Ring" className="text-sm text-gray-200 bg-black/50 border-gray-700 h-8" />
                  </div>
                  <div className="col-span-2 space-y-2">
                    <label className="text-xs text-gray-500">Metal Type</label>
                    <Input value={item.metal_type} onChange={e => updateItem(item.id, 'metal_type', e.target.value)} placeholder="gold" className="text-sm text-gray-200 bg-black/50 border-gray-700 h-8" />
                  </div>
                  <div className="col-span-1 space-y-2">
                    <label className="text-xs text-gray-500">Purity</label>
                    <Input value={item.purity} onChange={e => updateItem(item.id, 'purity', e.target.value)} placeholder="22K" className="text-sm text-gray-200 bg-black/50 border-gray-700 h-8" />
                  </div>
                  <div className="col-span-2 space-y-2">
                    <label className="text-xs text-gray-500">Weight (g)</label>
                    <Input type="number" value={item.weight} onChange={e => updateItem(item.id, 'weight', e.target.value)} placeholder="0.0" className="text-sm text-gray-200 bg-black/50 border-gray-700 h-8" />
                  </div>
                  <div className="col-span-2 space-y-2">
                    <label className="text-xs text-gray-500">Rate/g (₹)</label>
                    <Input type="number" value={item.rate_per_gram} onChange={e => updateItem(item.id, 'rate_per_gram', e.target.value)} placeholder="0.0" className="text-sm text-gray-200 bg-black/50 border-gray-700 h-8" />
                  </div>
                  <div className="col-span-2 space-y-2">
                    <label className="text-xs text-gray-500">Making (₹)</label>
                    <Input type="number" value={item.making_charge} onChange={e => updateItem(item.id, 'making_charge', e.target.value)} placeholder="0.0" className="text-sm text-gray-200 bg-black/50 border-gray-700 h-8" />
                  </div>
                  <div className="col-span-11 mt-2">
                    <div className="text-right text-xs text-gray-400">
                      Value: ₹{((parseFloat(item.weight)||0) * (parseFloat(item.rate_per_gram)||0) + (parseFloat(item.making_charge)||0)).toFixed(2)}
                    </div>
                  </div>
                  <div className="col-span-1 text-right">
                    <Button variant="ghost" size="icon" onClick={() => removeItem(item.id)} className="h-8 w-8 text-red-400 hover:text-red-500 hover:bg-red-500/10"><span className="text-xs">X</span></Button>
                  </div>
                </div>
              ))}
              {items.length === 0 && <p className="text-center text-gray-500 text-sm">No items added. Add an item to create a bill.</p>}
            </CardContent>
          </Card>
        </div>

        <div className="space-y-6">
          <Card className="bg-[#111] border-[#c6a962] shadow-[0_0_15px_rgba(198,169,98,0.15)] sticky top-6">
            <CardHeader>
              <CardTitle className="text-[#c6a962]">Invoice Summary</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex justify-between text-gray-400">
                <span>Subtotal</span>
                <span>₹{subtotal.toLocaleString(undefined, {minimumFractionDigits: 2, maximumFractionDigits: 2})}</span>
              </div>
              <div className="flex justify-between text-gray-400">
                <span>GST Amount</span>
                <span>₹{totalGst.toLocaleString(undefined, {minimumFractionDigits: 2, maximumFractionDigits: 2})}</span>
              </div>
              <div className="pt-4 border-t border-gray-800 flex justify-between items-center">
                <span className="text-lg font-bold text-gray-200">Grand Total</span>
                <span className="text-2xl font-bold text-[#c6a962]">₹{grandTotal.toLocaleString(undefined, {minimumFractionDigits: 2, maximumFractionDigits: 2})}</span>
              </div>
              <Button onClick={handleSubmit} disabled={mutation.isPending || items.length === 0} className="w-full mt-6 h-12 text-lg">
                {mutation.isPending ? <Loader2 className="h-5 w-5 animate-spin mr-2" /> : null}
                {mutation.isPending ? "Generating..." : "Generate Bill"}
              </Button>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}
