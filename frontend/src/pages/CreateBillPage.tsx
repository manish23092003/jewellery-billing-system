import { useState } from "react";
import { Button } from "@/components/ui/Button";
import { Input } from "@/components/ui/Input";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/Card";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { createBill, getBillById } from "@/lib/bills.api";
import { getCustomers } from "@/lib/customers.api";
import { useNavigate, useSearchParams } from "react-router-dom";
import { toast } from "sonner";
import { Loader2 } from "lucide-react";
import { useEffect } from "react";

export default function CreateBillPage() {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const convertFromId = searchParams.get("convertFrom");
  const queryClient = useQueryClient();
  const [customer, setCustomer] = useState({ id: "", name: "", phone: "", payment_method: "cash" });
  const [billType, setBillType] = useState("invoice");
  const [advanceAmount, setAdvanceAmount] = useState("");
  const [oldGoldItems, setOldGoldItems] = useState<any[]>([]);
  
  const { data: customersData } = useQuery({
    queryKey: ["customers", "all"],
    queryFn: () => getCustomers({ limit: 100 }),
  });
  const [items, setItems] = useState([
    { id: Date.now(), item_name: "", hsn_code: "71131910", metal_type: "gold", purity: "22K", weight: "", rate_per_gram: "", making_charge: "", gst_percentage: "3", quantity: "1", charges: [] as { id: number, charge_name: string, amount: string }[] }
  ]);

  useEffect(() => {
    if (convertFromId) {
      getBillById(convertFromId).then((billData) => {
        setCustomer({ id: "", name: billData.customer_name || "", phone: billData.customer_phone || "", payment_method: billData.payment_method || "cash" });
        setBillType("invoice"); // Convert to Tax Invoice
        setAdvanceAmount(billData.advance_amount ? billData.advance_amount.toString() : "");
        
        if (billData.items && billData.items.length > 0) {
          setItems(billData.items.map((i: any, index: number) => ({
            id: Date.now() + index,
            item_name: i.item_name,
            hsn_code: i.hsn_code || "",
            metal_type: i.metal_type,
            purity: i.purity,
            weight: i.weight.toString(),
            rate_per_gram: i.rate_per_gram.toString(),
            making_charge: i.making_charge.toString(),
            gst_percentage: i.gst_percentage.toString(),
            quantity: i.quantity.toString(),
            charges: (i.charges || []).map((c: any, cIdx: number) => ({
              id: Date.now() + index + cIdx,
              charge_name: c.charge_name,
              amount: c.amount.toString()
            }))
          })));
        }
        
        if (billData.old_gold_items && billData.old_gold_items.length > 0) {
          setOldGoldItems(billData.old_gold_items.map((og: any, index: number) => ({
            id: Date.now() + index,
            name: og.name,
            weight: og.weight.toString(),
            purity: og.purity,
            melting_loss_percentage: og.melting_loss_percentage.toString(),
            rate_per_gram: og.rate_per_gram.toString(),
          })));
        }
        
        toast.success("Estimate details loaded. You can now adjust weights and prices.");
      }).catch((err) => {
        toast.error("Failed to load estimate for conversion.");
      });
    }
  }, [convertFromId]);

  const getHSN = (metalType: string) => {
    const type = metalType.toLowerCase();
    if (type.includes("silver")) return "71131120";
    if (type.includes("platinum")) return "71131950";
    if (type.includes("diamond")) return "71131930";
    return "71131910"; // Default Gold
  };

  const addItem = () => setItems([...items, { id: Date.now(), item_name: "", hsn_code: "71131910", metal_type: "gold", purity: "22K", weight: "", rate_per_gram: "", making_charge: "", gst_percentage: "3", quantity: "1", charges: [] }]);
  const removeItem = (id: number) => setItems(items.filter(i => i.id !== id));

  const updateItem = (id: number, field: string, value: string) => {
    setItems(items.map(item => {
      if (item.id === id) {
        const updatedItem = { ...item, [field]: value };
        if (field === 'metal_type') {
          updatedItem.hsn_code = getHSN(value);
        }
        return updatedItem;
      }
      return item;
    }));
  };

  const addCharge = (itemId: number) => {
    setItems(items.map(item => item.id === itemId ? { ...item, charges: [...item.charges, { id: Date.now(), charge_name: "", amount: "" }] } : item));
  };
  const removeCharge = (itemId: number, chargeId: number) => {
    setItems(items.map(item => item.id === itemId ? { ...item, charges: item.charges.filter(c => c.id !== chargeId) } : item));
  };
  const updateCharge = (itemId: number, chargeId: number, field: string, value: string) => {
    setItems(items.map(item => item.id === itemId ? {
      ...item, charges: item.charges.map(c => c.id === chargeId ? { ...c, [field]: value } : c)
    } : item));
  };

  const addOldGold = () => setOldGoldItems([...oldGoldItems, { id: Date.now(), name: "", weight: "", purity: "22K", melting_loss_percentage: "0", rate_per_gram: "" }]);
  const removeOldGold = (id: number) => setOldGoldItems(oldGoldItems.filter(i => i.id !== id));
  const updateOldGold = (id: number, field: string, value: string) => {
    setOldGoldItems(oldGoldItems.map(item => item.id === id ? { ...item, [field]: value } : item));
  };

  // Preview Calculations
  let subtotal = 0;
  let totalGst = 0;
  items.forEach(item => {
    const w = parseFloat(item.weight) || 0;
    const r = parseFloat(item.rate_per_gram) || 0;
    const m = parseFloat(item.making_charge) || 0;
    const cSum = item.charges.reduce((acc, c) => acc + (parseFloat(c.amount) || 0), 0);
    const q = parseInt(item.quantity) || 1;
    const gstPct = parseFloat(item.gst_percentage) || 3;
    
    const metalVal = w * r;
    const lineSub = (metalVal + m + cSum) * q;
    const lineGst = lineSub * (gstPct / 100);
    
    subtotal += lineSub;
    totalGst += lineGst;
  });
  const grandTotal = subtotal + totalGst;

  let totalOldGold = 0;
  oldGoldItems.forEach(og => {
    const w = parseFloat(og.weight) || 0;
    const r = parseFloat(og.rate_per_gram) || 0;
    const loss = parseFloat(og.melting_loss_percentage) || 0;
    const grossVal = w * r;
    const netVal = grossVal - (grossVal * (loss / 100));
    totalOldGold += netVal;
  });

  const adv = parseFloat(advanceAmount) || 0;
  let balanceDue = grandTotal - totalOldGold - adv;
  if (balanceDue < 0) balanceDue = 0;

  const mutation = useMutation({
    mutationFn: createBill,
    onSuccess: () => {
      toast.success("Bill generated successfully");
      queryClient.invalidateQueries({ queryKey: ["bills"] });
      queryClient.invalidateQueries({ queryKey: ["dashboard"] });
      queryClient.invalidateQueries({ queryKey: ["customers"] });
      navigate("/bills/history");
    },
    onError: (err: any) => {
      const msg = err.response?.data?.error || err.response?.data?.message || err.message || "Failed to create bill";
      toast.error(msg);
    }
  });

  const handleSubmit = () => {
    if (!customer.name) return toast.error("Customer name is required");
    if (items.length === 0) return toast.error("At least one item is required");
    for (const item of items) {
      if (!parseFloat(item.weight) || parseFloat(item.weight) <= 0) return toast.error("Weight must be greater than 0 for all items");
      if (!parseFloat(item.rate_per_gram) || parseFloat(item.rate_per_gram) <= 0) return toast.error("Rate per gram must be greater than 0 for all items");
    }

    
    const payload = {
      type: billType,
      status: billType === "estimate" ? "pending" : "completed",
      advance_amount: parseFloat(advanceAmount) || 0,
      invoice_date: new Date().toISOString(),
      customer_name: customer.name,
      customer_phone: customer.phone,
      payment_method: customer.payment_method,
      notes: "",
      convert_from_id: convertFromId || undefined,
      items: items.map(item => ({
        item_name: item.item_name || "Unknown Item",
        hsn_code: item.hsn_code || "",
        metal_type: item.metal_type,
        purity: item.purity,
        weight: parseFloat(item.weight) || 0,
        rate_per_gram: parseFloat(item.rate_per_gram) || 0,
        making_charge: parseFloat(item.making_charge) || 0,
        gst_percentage: parseFloat(item.gst_percentage) || 3,
        quantity: parseInt(item.quantity) || 1,
        charges: item.charges.map(c => ({ charge_name: c.charge_name || "Additional Charge", amount: parseFloat(c.amount) || 0 }))
      })),
      old_gold_items: oldGoldItems.map(og => ({
        name: og.name || "Old Gold",
        weight: parseFloat(og.weight) || 0,
        purity: og.purity,
        melting_loss_percentage: parseFloat(og.melting_loss_percentage) || 0,
        rate_per_gram: parseFloat(og.rate_per_gram) || 0
      }))
    };
    
    mutation.mutate(payload);
  };

  return (
    <div className="p-6 max-w-7xl mx-auto space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-3xl font-bold tracking-tight text-[#c6a962]">
          Create New Document
        </h1>
        <div className="flex space-x-2 bg-card p-1 rounded-lg border border-border">
          <Button 
            variant="ghost" 
            size="sm"
            onClick={() => setBillType("invoice")} 
            className={billType === "invoice" ? "bg-[#c6a962] text-black hover:bg-[#b0965a] hover:text-black" : "text-muted-foreground hover:text-white"}
          >
            Tax Invoice
          </Button>
          <Button 
            variant="ghost" 
            size="sm"
            onClick={() => setBillType("estimate")} 
            className={billType === "estimate" ? "bg-[#c6a962] text-black hover:bg-[#b0965a] hover:text-black" : "text-muted-foreground hover:text-white"}
          >
            Estimate / Order
          </Button>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-2 space-y-6">
          <Card className="bg-card backdrop-blur-md border-[#c6a962]/20">
            <CardHeader className="flex flex-row items-center justify-between">
              <CardTitle className="text-foreground">Customer Details</CardTitle>
              {customersData?.customers && customersData.customers.length > 0 && (
                <select 
                  className="bg-muted border border-border text-gray-300 text-sm rounded-md px-2 py-1 max-w-[200px]"
                  onChange={(e) => {
                    const selectedId = e.target.value;
                    if (!selectedId) return;
                    const c = customersData.customers.find((x: any) => x.id === selectedId);
                    if (c) {
                      setCustomer(prev => ({ ...prev, id: c.id, name: c.name, phone: c.phone }));
                    }
                  }}
                  value={customer.id}
                >
                  <option value="">-- Select Existing --</option>
                  {customersData.customers.map((c: any) => (
                    <option key={c.id} value={c.id}>{c.name} ({c.phone})</option>
                  ))}
                </select>
              )}
            </CardHeader>
            <CardContent className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <label className="text-sm text-muted-foreground">Customer Name</label>
                <Input value={customer.name} onChange={e => setCustomer({...customer, name: e.target.value, id: ""})} placeholder="Enter name" className="text-foreground bg-muted border-border" />
              </div>
              <div className="space-y-2">
                <label className="text-sm text-muted-foreground">Phone Number</label>
                <Input value={customer.phone} onChange={e => setCustomer({...customer, phone: e.target.value, id: ""})} placeholder="Enter phone" className="text-foreground bg-muted border-border" />
              </div>
              <div className="space-y-2">
                <label className="text-sm text-muted-foreground">Payment Method</label>
                <select 
                  value={customer.payment_method} 
                  onChange={e => setCustomer({...customer, payment_method: e.target.value})}
                  className="flex h-9 w-full rounded-md border border-border bg-muted px-3 py-1 text-sm shadow-sm transition-colors text-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
                >
                  <option value="cash">Cash</option>
                  <option value="card">Card</option>
                  <option value="upi">UPI</option>
                  <option value="bank_transfer">Bank Transfer</option>
                </select>
              </div>
            </CardContent>
          </Card>

          <Card className="bg-card backdrop-blur-md border-[#c6a962]/20">
            <CardHeader className="flex flex-row justify-between items-center">
              <CardTitle className="text-foreground">Bill Items</CardTitle>
              <Button variant="outline" size="sm" onClick={addItem} className="border-[#c6a962]/50 text-[#c6a962]">
                + Add Item
              </Button>
            </CardHeader>
            <CardContent className="space-y-6 overflow-x-auto">
              <div className="min-w-[800px] lg:min-w-0 space-y-6">
              {items.map((item, index) => (
                <div key={item.id} className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-8 gap-3 items-end border-b border-border pb-6 relative">
                  <div className="col-span-full absolute -top-4 right-0 text-gray-600 text-xs text-right">Item {index + 1}</div>
                  <div className="space-y-2">
                    <label className="text-xs text-muted-foreground">Item Name</label>
                    <Input value={item.item_name} onChange={e => updateItem(item.id, 'item_name', e.target.value)} placeholder="e.g. Gold Ring" className="text-sm text-foreground bg-muted border-border h-8" />
                  </div>
                  <div className="space-y-2">
                    <label className="text-xs text-muted-foreground">HSN</label>
                    <Input value={item.hsn_code} onChange={e => updateItem(item.id, 'hsn_code', e.target.value)} placeholder="7113" className="text-sm text-foreground bg-muted border-border h-8" />
                  </div>
                  <div className="space-y-2">
                    <label className="text-xs text-muted-foreground">Metal Type</label>
                    <select 
                      value={item.metal_type} 
                      onChange={e => updateItem(item.id, 'metal_type', e.target.value)}
                      className="flex h-8 w-full rounded-md border border-border bg-muted px-3 py-1 text-sm shadow-sm transition-colors text-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
                    >
                      <option value="gold">Gold</option>
                      <option value="silver">Silver</option>
                      <option value="platinum">Platinum</option>
                      <option value="diamond">Diamond</option>
                    </select>
                  </div>
                  <div className="space-y-2">
                    <label className="text-xs text-muted-foreground">Purity</label>
                    <Input value={item.purity} onChange={e => updateItem(item.id, 'purity', e.target.value)} placeholder="22K" className="text-sm text-foreground bg-muted border-border h-8" />
                  </div>
                  <div className="space-y-2">
                    <label className="text-xs text-muted-foreground">Wt (g)</label>
                    <Input type="number" min="0" step="0.001" value={item.weight} onChange={e => updateItem(item.id, 'weight', e.target.value)} placeholder="0.0" className="text-sm text-foreground bg-muted border-border h-8" />
                  </div>
                  <div className="space-y-2">
                    <label className="text-xs text-muted-foreground">Qty</label>
                    <Input type="number" min="1" step="1" value={item.quantity} onChange={e => updateItem(item.id, 'quantity', e.target.value)} placeholder="1" className="text-sm text-foreground bg-muted border-border h-8" />
                  </div>
                  <div className="space-y-2">
                    <label className="text-xs text-muted-foreground">Rate/g (₹)</label>
                    <Input type="number" min="0" step="0.01" value={item.rate_per_gram} onChange={e => updateItem(item.id, 'rate_per_gram', e.target.value)} placeholder="0.0" className="text-sm text-foreground bg-muted border-border h-8" />
                  </div>
                  <div className="space-y-2">
                    <label className="text-xs text-muted-foreground">Making (₹)</label>
                    <Input type="number" min="0" step="0.01" value={item.making_charge} onChange={e => updateItem(item.id, 'making_charge', e.target.value)} placeholder="0.0" className="text-sm text-foreground bg-muted border-border h-8" />
                  </div>
                  <div className="col-span-full mt-4 space-y-3">
                    {item.charges.map((charge) => (
                      <div key={charge.id} className="flex gap-2 items-center ml-4 md:ml-12">
                        <Input value={charge.charge_name} onChange={e => updateCharge(item.id, charge.id, 'charge_name', e.target.value)} placeholder="e.g. Stone Charge" className="text-sm text-foreground bg-muted border-border h-8 max-w-[200px]" />
                        <Input type="number" min="0" step="0.01" value={charge.amount} onChange={e => updateCharge(item.id, charge.id, 'amount', e.target.value)} placeholder="Amount (₹)" className="text-sm text-foreground bg-muted border-border h-8 max-w-[120px]" />
                        <Button variant="ghost" size="sm" onClick={() => removeCharge(item.id, charge.id)} className="h-8 text-muted-foreground hover:text-red-400">Remove</Button>
                      </div>
                    ))}
                    <div className="flex justify-between items-center ml-0 md:ml-12 border-t border-border pt-2">
                      <Button variant="ghost" size="sm" onClick={() => addCharge(item.id)} className="text-xs text-[#c6a962]">
                        + Add Charge
                      </Button>
                      <div className="flex items-center gap-4">
                        <div className="text-right text-xs text-muted-foreground font-medium">
                          Item Subtotal: ₹{(((parseFloat(item.weight)||0) * (parseFloat(item.rate_per_gram)||0) + (parseFloat(item.making_charge)||0) + item.charges.reduce((acc, c) => acc + (parseFloat(c.amount) || 0), 0)) * (parseInt(item.quantity)||1)).toFixed(2)}
                        </div>
                        <Button variant="ghost" size="icon" onClick={() => removeItem(item.id)} className="h-8 w-8 text-red-400 hover:text-red-500 hover:bg-red-500/10"><span className="text-xs">X</span></Button>
                      </div>
                    </div>
                  </div>
                </div>
              ))}
              {items.length === 0 && <p className="text-center text-muted-foreground text-sm">No items added. Add an item to create a bill.</p>}
              </div>
            </CardContent>
          </Card>

          <Card className="bg-card backdrop-blur-md border-[#c6a962]/20">
            <CardHeader className="flex flex-row justify-between items-center">
              <CardTitle className="text-foreground">Old Gold / Exchange</CardTitle>
              <Button variant="outline" size="sm" onClick={addOldGold} className="border-[#c6a962]/50 text-[#c6a962]">
                + Add Old Gold
              </Button>
            </CardHeader>
            <CardContent className="space-y-6 overflow-x-auto">
              <div className="min-w-[800px] lg:min-w-0 space-y-6">
              {oldGoldItems.map((og, index) => (
                <div key={og.id} className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-3 items-end border-b border-border pb-6 relative">
                  <div className="col-span-full absolute -top-4 right-0 text-gray-600 text-xs text-right">Item {index + 1}</div>
                  <div className="space-y-2">
                    <label className="text-xs text-muted-foreground">Item Name</label>
                    <Input value={og.name} onChange={e => updateOldGold(og.id, 'name', e.target.value)} placeholder="e.g. Broken Chain" className="text-sm text-foreground bg-muted border-border h-8" />
                  </div>
                  <div className="space-y-2">
                    <label className="text-xs text-muted-foreground">Purity</label>
                    <Input value={og.purity} onChange={e => updateOldGold(og.id, 'purity', e.target.value)} placeholder="22K" className="text-sm text-foreground bg-muted border-border h-8" />
                  </div>
                  <div className="space-y-2">
                    <label className="text-xs text-muted-foreground">Weight (g)</label>
                    <Input type="number" min="0" step="0.001" value={og.weight} onChange={e => updateOldGold(og.id, 'weight', e.target.value)} placeholder="0.0" className="text-sm text-foreground bg-muted border-border h-8" />
                  </div>
                  <div className="space-y-2">
                    <label className="text-xs text-muted-foreground">Melting Loss %</label>
                    <Input type="number" min="0" step="0.01" value={og.melting_loss_percentage} onChange={e => updateOldGold(og.id, 'melting_loss_percentage', e.target.value)} placeholder="0.0" className="text-sm text-foreground bg-muted border-border h-8" />
                  </div>
                  <div className="space-y-2">
                    <label className="text-xs text-muted-foreground">Rate/g (₹)</label>
                    <Input type="number" min="0" step="0.01" value={og.rate_per_gram} onChange={e => updateOldGold(og.id, 'rate_per_gram', e.target.value)} placeholder="0.0" className="text-sm text-foreground bg-muted border-border h-8" />
                  </div>
                  <div className="col-span-full mt-4 flex justify-end items-center gap-4">
                    <div className="text-right text-xs text-muted-foreground font-medium">
                      Exchange Value: ₹{((parseFloat(og.weight)||0) * (1 - (parseFloat(og.melting_loss_percentage)||0)/100) * (parseFloat(og.rate_per_gram)||0)).toFixed(2)}
                    </div>
                    <Button variant="ghost" size="icon" onClick={() => removeOldGold(og.id)} className="h-8 w-8 text-red-400 hover:text-red-500 hover:bg-red-500/10"><span className="text-xs">X</span></Button>
                  </div>
                </div>
              ))}
              {oldGoldItems.length === 0 && <p className="text-center text-muted-foreground text-sm">No old gold added.</p>}
              </div>
            </CardContent>
          </Card>
        </div>

        <div className="space-y-6">
          <Card className="bg-card border-[#c6a962] shadow-sm sticky top-6">
            <CardHeader>
              <CardTitle className="text-[#c6a962]">Invoice Summary</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex justify-between text-muted-foreground">
                <span>Gross Subtotal</span>
                <span>₹{subtotal.toLocaleString(undefined, {minimumFractionDigits: 2, maximumFractionDigits: 2})}</span>
              </div>
              <div className="flex justify-between text-muted-foreground">
                <span>CGST</span>
                <span>₹{(totalGst / 2).toLocaleString(undefined, {minimumFractionDigits: 2, maximumFractionDigits: 2})}</span>
              </div>
              <div className="flex justify-between text-muted-foreground">
                <span>SGST</span>
                <span>₹{(totalGst / 2).toLocaleString(undefined, {minimumFractionDigits: 2, maximumFractionDigits: 2})}</span>
              </div>
              <div className="flex justify-between items-center text-foreground font-medium">
                <span>Total Items Value</span>
                <span>₹{grandTotal.toLocaleString(undefined, {minimumFractionDigits: 2, maximumFractionDigits: 2})}</span>
              </div>
              {oldGoldItems.length > 0 && (
                <div className="flex justify-between text-red-400">
                  <span>- Old Gold Exchange</span>
                  <span>-₹{totalOldGold.toLocaleString(undefined, {minimumFractionDigits: 2, maximumFractionDigits: 2})}</span>
                </div>
              )}
              <div className="flex justify-between items-center text-green-400 pt-2">
                <span>{billType === "estimate" ? "Advance Payment" : "Amount Paid"} (₹)</span>
                <Input 
                  type="number" 
                  value={advanceAmount} 
                  onChange={(e) => setAdvanceAmount(e.target.value)} 
                  placeholder="0.00" 
                  className="w-24 h-8 bg-muted border-border text-right text-green-400"
                />
              </div>
              <div className="pt-4 border-t border-border flex justify-between items-center">
                <span className="text-lg font-bold text-foreground">Balance Due</span>
                <span className="text-2xl font-bold text-[#c6a962]">₹{balanceDue.toLocaleString(undefined, {minimumFractionDigits: 2, maximumFractionDigits: 2})}</span>
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
