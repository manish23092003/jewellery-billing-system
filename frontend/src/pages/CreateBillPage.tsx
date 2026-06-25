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

type Charge = { id: number; charge_name: string; amount: string };
type Item = {
  id: number;
  item_name: string;
  hsn_code: string;
  metal_type: string;
  purity: string;
  weight: string;
  rate_per_gram: string;
  making_charge: string;
  gst_percentage: string;
  quantity: string;
  charges: Charge[];
};

const GOLD_PURITIES = ["24K", "22K", "18K", "14K", "10K"];
const SILVER_PURITIES = ["999", "925", "800", "Sterling"];

function newGoldItem(): Item {
  return { id: Date.now(), item_name: "", hsn_code: "71131910", metal_type: "gold", purity: "22K", weight: "", rate_per_gram: "", making_charge: "", gst_percentage: "3", quantity: "1", charges: [] };
}
function newSilverItem(): Item {
  return { id: Date.now(), item_name: "", hsn_code: "71131120", metal_type: "silver", purity: "925", weight: "", rate_per_gram: "", making_charge: "", gst_percentage: "3", quantity: "1", charges: [] };
}

// ── Reusable Item Row Component ────────────────────────────────────────────────
function ItemRow({ item, purities, onUpdate, onRemove, onAddCharge, onRemoveCharge, onUpdateCharge, index }: {
  item: Item;
  purities: string[];
  onUpdate: (id: number, field: string, value: string) => void;
  onRemove: (id: number) => void;
  onAddCharge: (id: number) => void;
  onRemoveCharge: (itemId: number, chargeId: number) => void;
  onUpdateCharge: (itemId: number, chargeId: number, field: string, value: string) => void;
  index: number;
}) {
  const subtotal = ((parseFloat(item.weight) || 0) * (parseFloat(item.rate_per_gram) || 0) + (parseFloat(item.making_charge) || 0) + item.charges.reduce((a, c) => a + (parseFloat(c.amount) || 0), 0)) * (parseInt(item.quantity) || 1);

  return (
    <div className="border border-border rounded-lg p-4 space-y-3 relative bg-muted/30">
      <div className="absolute -top-2.5 left-3 text-xs text-muted-foreground bg-card px-2">Item {index + 1}</div>
      {/* Row 1 */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
        <div className="space-y-1">
          <label className="text-xs text-muted-foreground">Item Name</label>
          <Input value={item.item_name} onChange={e => onUpdate(item.id, "item_name", e.target.value)} placeholder="e.g. Gold Ring" className="h-8 text-sm bg-muted border-border text-foreground" />
        </div>
        <div className="space-y-1">
          <label className="text-xs text-muted-foreground">HSN Code</label>
          <Input value={item.hsn_code} onChange={e => onUpdate(item.id, "hsn_code", e.target.value)} placeholder="7113..." className="h-8 text-sm bg-muted border-border text-foreground" />
        </div>
        <div className="space-y-1">
          <label className="text-xs text-muted-foreground">Purity</label>
          <select
            value={item.purity}
            onChange={e => onUpdate(item.id, "purity", e.target.value)}
            className="flex h-8 w-full rounded-md border border-border bg-muted px-3 py-1 text-sm text-foreground"
          >
            {purities.map(p => <option key={p} value={p}>{p}</option>)}
          </select>
        </div>
        <div className="space-y-1">
          <label className="text-xs text-muted-foreground">Quantity</label>
          <Input type="number" min="1" value={item.quantity} onChange={e => onUpdate(item.id, "quantity", e.target.value)} placeholder="1" className="h-8 text-sm bg-muted border-border text-foreground" />
        </div>
      </div>
      {/* Row 2 */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
        <div className="space-y-1">
          <label className="text-xs text-muted-foreground">Weight (g)</label>
          <Input type="number" min="0" step="0.001" value={item.weight} onChange={e => onUpdate(item.id, "weight", e.target.value)} placeholder="0.000" className="h-8 text-sm bg-muted border-border text-foreground" />
        </div>
        <div className="space-y-1">
          <label className="text-xs text-muted-foreground">Rate / gram (₹)</label>
          <Input type="number" min="0" step="0.01" value={item.rate_per_gram} onChange={e => onUpdate(item.id, "rate_per_gram", e.target.value)} placeholder="0.00" className="h-8 text-sm bg-muted border-border text-foreground" />
        </div>
        <div className="space-y-1">
          <label className="text-xs text-muted-foreground">Making Charge (₹)</label>
          <Input type="number" min="0" step="0.01" value={item.making_charge} onChange={e => onUpdate(item.id, "making_charge", e.target.value)} placeholder="0.00" className="h-8 text-sm bg-muted border-border text-foreground" />
        </div>
        <div className="space-y-1">
          <label className="text-xs text-muted-foreground">GST %</label>
          <Input type="number" min="0" step="0.5" value={item.gst_percentage} onChange={e => onUpdate(item.id, "gst_percentage", e.target.value)} placeholder="3" className="h-8 text-sm bg-muted border-border text-foreground" />
        </div>
      </div>

      {/* Additional Charges */}
      {item.charges.map(charge => (
        <div key={charge.id} className="flex gap-2 items-center ml-2">
          <Input value={charge.charge_name} onChange={e => onUpdateCharge(item.id, charge.id, "charge_name", e.target.value)} placeholder="Charge name (e.g. Stone)" className="h-7 text-xs bg-muted border-border text-foreground max-w-[180px]" />
          <Input type="number" value={charge.amount} onChange={e => onUpdateCharge(item.id, charge.id, "amount", e.target.value)} placeholder="₹ Amount" className="h-7 text-xs bg-muted border-border text-foreground max-w-[120px]" />
          <Button variant="ghost" size="sm" onClick={() => onRemoveCharge(item.id, charge.id)} className="h-7 text-xs text-red-400 hover:text-red-300">✕</Button>
        </div>
      ))}

      {/* Footer */}
      <div className="flex justify-between items-center pt-1 border-t border-border/50">
        <Button variant="ghost" size="sm" onClick={() => onAddCharge(item.id)} className="text-xs text-[#c6a962] h-7">+ Add Charge</Button>
        <div className="flex items-center gap-4">
          <span className="text-xs text-muted-foreground">Subtotal: ₹{subtotal.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}</span>
          <Button variant="ghost" size="sm" onClick={() => onRemove(item.id)} className="h-7 text-xs text-red-400 hover:text-red-300 hover:bg-red-500/10">Remove</Button>
        </div>
      </div>
    </div>
  );
}

// ── Main Page ──────────────────────────────────────────────────────────────────
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

  // Separate Gold and Silver item lists
  const [goldItems, setGoldItems] = useState<Item[]>([newGoldItem()]);
  const [silverItems, setSilverItems] = useState<Item[]>([]);

  // ── Helpers ──
  const updateItems = (setter: React.Dispatch<React.SetStateAction<Item[]>>) =>
    (id: number, field: string, value: string) =>
      setter(prev => prev.map(item => item.id === id ? { ...item, [field]: value } : item));

  const removeItem = (setter: React.Dispatch<React.SetStateAction<Item[]>>) =>
    (id: number) => setter(prev => prev.filter(i => i.id !== id));

  const addCharge = (setter: React.Dispatch<React.SetStateAction<Item[]>>) =>
    (itemId: number) => setter(prev => prev.map(item => item.id === itemId ? { ...item, charges: [...item.charges, { id: Date.now(), charge_name: "", amount: "" }] } : item));

  const removeCharge = (setter: React.Dispatch<React.SetStateAction<Item[]>>) =>
    (itemId: number, chargeId: number) => setter(prev => prev.map(item => item.id === itemId ? { ...item, charges: item.charges.filter(c => c.id !== chargeId) } : item));

  const updateCharge = (setter: React.Dispatch<React.SetStateAction<Item[]>>) =>
    (itemId: number, chargeId: number, field: string, value: string) =>
      setter(prev => prev.map(item => item.id === itemId ? { ...item, charges: item.charges.map(c => c.id === chargeId ? { ...c, [field]: value } : c) } : item));

  // ── Load from estimate ──
  useEffect(() => {
    if (convertFromId) {
      getBillById(convertFromId).then((billData) => {
        setCustomer({ id: "", name: billData.customer_name || "", phone: billData.customer_phone || "", payment_method: billData.payment_method || "cash" });
        setBillType("invoice");
        setAdvanceAmount(billData.advance_amount ? billData.advance_amount.toString() : "");
        if (billData.items && billData.items.length > 0) {
          const mapped: Item[] = billData.items.map((i: any, idx: number) => ({
            id: Date.now() + idx, item_name: i.item_name, hsn_code: i.hsn_code || "",
            metal_type: i.metal_type, purity: i.purity, weight: i.weight.toString(),
            rate_per_gram: i.rate_per_gram.toString(), making_charge: i.making_charge.toString(),
            gst_percentage: i.gst_percentage.toString(), quantity: i.quantity.toString(),
            charges: (i.charges || []).map((c: any, cIdx: number) => ({ id: Date.now() + idx + cIdx, charge_name: c.charge_name, amount: c.amount.toString() }))
          }));
          setGoldItems(mapped.filter(i => i.metal_type === "gold"));
          setSilverItems(mapped.filter(i => i.metal_type === "silver"));
        }
        if (billData.old_gold_items && billData.old_gold_items.length > 0) {
          setOldGoldItems(billData.old_gold_items.map((og: any, idx: number) => ({
            id: Date.now() + idx, name: og.name, weight: og.weight.toString(),
            purity: og.purity, melting_loss_percentage: og.melting_loss_percentage.toString(), rate_per_gram: og.rate_per_gram.toString()
          })));
        }
        toast.success("Estimate details loaded.");
      }).catch(() => toast.error("Failed to load estimate for conversion."));
    }
  }, [convertFromId]);

  // ── Calculations ──
  const calcItemsTotal = (items: Item[]) => {
    let sub = 0, gst = 0;
    items.forEach(item => {
      const w = parseFloat(item.weight) || 0;
      const r = parseFloat(item.rate_per_gram) || 0;
      const m = parseFloat(item.making_charge) || 0;
      const cSum = item.charges.reduce((a, c) => a + (parseFloat(c.amount) || 0), 0);
      const q = parseInt(item.quantity) || 1;
      const gstPct = parseFloat(item.gst_percentage) || 3;
      const lineSub = (w * r + m + cSum) * q;
      sub += lineSub;
      gst += lineSub * (gstPct / 100);
    });
    return { sub, gst };
  };

  const gold = calcItemsTotal(goldItems);
  const silver = calcItemsTotal(silverItems);
  const subtotal = gold.sub + silver.sub;
  const totalGst = gold.gst + silver.gst;
  const grandTotal = subtotal + totalGst;

  let totalOldGold = 0;
  oldGoldItems.forEach(og => {
    const w = parseFloat(og.weight) || 0;
    const r = parseFloat(og.rate_per_gram) || 0;
    const loss = parseFloat(og.melting_loss_percentage) || 0;
    totalOldGold += w * r * (1 - loss / 100);
  });

  const adv = parseFloat(advanceAmount) || 0;
  const balanceDue = Math.max(0, grandTotal - totalOldGold - adv);

  // ── Submit ──
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
      const msg = err.response?.data?.error || err.message || "Failed to create bill";
      toast.error(msg);
    }
  });

  const handleSubmit = () => {
    if (!customer.name) return toast.error("Customer name is required");
    const allItems = [...goldItems, ...silverItems];
    if (allItems.length === 0) return toast.error("Add at least one Gold or Silver item");
    for (const item of allItems) {
      if (!parseFloat(item.weight) || parseFloat(item.weight) <= 0) return toast.error("Weight must be > 0 for all items");
      if (!parseFloat(item.rate_per_gram) || parseFloat(item.rate_per_gram) <= 0) return toast.error("Rate/gram must be > 0 for all items");
    }
    mutation.mutate({
      type: billType,
      status: billType === "estimate" ? "pending" : "completed",
      advance_amount: parseFloat(advanceAmount) || 0,
      invoice_date: new Date().toISOString(),
      customer_name: customer.name,
      customer_phone: customer.phone,
      payment_method: customer.payment_method,
      notes: "",
      convert_from_id: convertFromId || undefined,
      items: allItems.map(item => ({
        item_name: item.item_name || "Item",
        hsn_code: item.hsn_code,
        metal_type: item.metal_type,
        purity: item.purity,
        weight: parseFloat(item.weight) || 0,
        rate_per_gram: parseFloat(item.rate_per_gram) || 0,
        making_charge: parseFloat(item.making_charge) || 0,
        gst_percentage: parseFloat(item.gst_percentage) || 3,
        quantity: parseInt(item.quantity) || 1,
        charges: item.charges.map(c => ({ charge_name: c.charge_name || "Charge", amount: parseFloat(c.amount) || 0 }))
      })),
      old_gold_items: oldGoldItems.map(og => ({
        name: og.name || "Old Gold",
        weight: parseFloat(og.weight) || 0,
        purity: og.purity,
        melting_loss_percentage: parseFloat(og.melting_loss_percentage) || 0,
        rate_per_gram: parseFloat(og.rate_per_gram) || 0
      }))
    });
  };

  return (
    <div className="p-6 max-w-7xl mx-auto space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <h1 className="text-3xl font-bold tracking-tight text-[#c6a962]">Create New Document</h1>
        <div className="flex space-x-2 bg-card p-1 rounded-lg border border-border">
          <Button variant="ghost" size="sm" onClick={() => setBillType("invoice")} className={billType === "invoice" ? "bg-[#c6a962] text-black hover:bg-[#b0965a] hover:text-black" : "text-muted-foreground hover:text-white"}>Tax Invoice</Button>
          <Button variant="ghost" size="sm" onClick={() => setBillType("estimate")} className={billType === "estimate" ? "bg-[#c6a962] text-black hover:bg-[#b0965a] hover:text-black" : "text-muted-foreground hover:text-white"}>Estimate / Order</Button>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-2 space-y-6">

          {/* Customer Details */}
          <Card className="bg-card backdrop-blur-md border-[#c6a962]/20">
            <CardHeader className="flex flex-row items-center justify-between">
              <CardTitle className="text-foreground">Customer Details</CardTitle>
              {customersData?.customers && customersData.customers.length > 0 && (
                <select
                  className="bg-muted border border-border text-gray-300 text-sm rounded-md px-2 py-1 max-w-[200px]"
                  onChange={(e) => {
                    const c = customersData.customers.find((x: any) => x.id === e.target.value);
                    if (c) setCustomer(prev => ({ ...prev, id: c.id, name: c.name, phone: c.phone }));
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
                <Input value={customer.name} onChange={e => setCustomer({ ...customer, name: e.target.value, id: "" })} placeholder="Enter name" className="text-foreground bg-muted border-border" />
              </div>
              <div className="space-y-2">
                <label className="text-sm text-muted-foreground">Phone Number</label>
                <Input value={customer.phone} onChange={e => setCustomer({ ...customer, phone: e.target.value, id: "" })} placeholder="Enter phone" className="text-foreground bg-muted border-border" />
              </div>
              <div className="space-y-2">
                <label className="text-sm text-muted-foreground">Payment Method</label>
                <select value={customer.payment_method} onChange={e => setCustomer({ ...customer, payment_method: e.target.value })} className="flex h-9 w-full rounded-md border border-border bg-muted px-3 py-1 text-sm text-foreground">
                  <option value="cash">Cash</option>
                  <option value="card">Card</option>
                  <option value="upi">UPI</option>
                  <option value="bank_transfer">Bank Transfer</option>
                </select>
              </div>
            </CardContent>
          </Card>

          {/* ── GOLD ITEMS CARD ── */}
          <Card className="bg-card backdrop-blur-md border-yellow-500/30">
            <CardHeader className="flex flex-row justify-between items-center border-b border-yellow-500/20 pb-4">
              <div className="flex items-center gap-3">
                <div className="w-3 h-3 rounded-full bg-yellow-400 shadow-[0_0_8px_2px_rgba(250,204,21,0.4)]" />
                <CardTitle className="text-yellow-400 text-lg">Gold Items</CardTitle>
                <span className="text-xs text-muted-foreground bg-yellow-400/10 border border-yellow-400/20 px-2 py-0.5 rounded-full">HSN: 71131910</span>
              </div>
              <Button variant="outline" size="sm" onClick={() => setGoldItems(prev => [...prev, newGoldItem()])} className="border-yellow-500/40 text-yellow-400 hover:bg-yellow-400/10 hover:text-yellow-300">
                + Add Gold Item
              </Button>
            </CardHeader>
            <CardContent className="space-y-4 pt-4">
              {goldItems.length === 0 && (
                <p className="text-center text-muted-foreground text-sm py-4">No gold items. Click "Add Gold Item" to begin.</p>
              )}
              {goldItems.map((item, index) => (
                <ItemRow
                  key={item.id}
                  item={item}
                  index={index}
                  purities={GOLD_PURITIES}
                  onUpdate={updateItems(setGoldItems)}
                  onRemove={removeItem(setGoldItems)}
                  onAddCharge={addCharge(setGoldItems)}
                  onRemoveCharge={removeCharge(setGoldItems)}
                  onUpdateCharge={updateCharge(setGoldItems)}
                />
              ))}
              {goldItems.length > 0 && (
                <div className="flex justify-end pt-2 text-sm">
                  <span className="text-yellow-400/70">Gold Subtotal: ₹{(gold.sub + gold.gst).toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}</span>
                </div>
              )}
            </CardContent>
          </Card>

          {/* ── SILVER ITEMS CARD ── */}
          <Card className="bg-card backdrop-blur-md border-slate-400/30">
            <CardHeader className="flex flex-row justify-between items-center border-b border-slate-400/20 pb-4">
              <div className="flex items-center gap-3">
                <div className="w-3 h-3 rounded-full bg-slate-300 shadow-[0_0_8px_2px_rgba(203,213,225,0.4)]" />
                <CardTitle className="text-slate-300 text-lg">Silver Items</CardTitle>
                <span className="text-xs text-muted-foreground bg-slate-400/10 border border-slate-400/20 px-2 py-0.5 rounded-full">HSN: 71131120</span>
              </div>
              <Button variant="outline" size="sm" onClick={() => setSilverItems(prev => [...prev, newSilverItem()])} className="border-slate-400/40 text-slate-300 hover:bg-slate-400/10 hover:text-slate-200">
                + Add Silver Item
              </Button>
            </CardHeader>
            <CardContent className="space-y-4 pt-4">
              {silverItems.length === 0 && (
                <p className="text-center text-muted-foreground text-sm py-4">No silver items. Click "Add Silver Item" to begin.</p>
              )}
              {silverItems.map((item, index) => (
                <ItemRow
                  key={item.id}
                  item={item}
                  index={index}
                  purities={SILVER_PURITIES}
                  onUpdate={updateItems(setSilverItems)}
                  onRemove={removeItem(setSilverItems)}
                  onAddCharge={addCharge(setSilverItems)}
                  onRemoveCharge={removeCharge(setSilverItems)}
                  onUpdateCharge={updateCharge(setSilverItems)}
                />
              ))}
              {silverItems.length > 0 && (
                <div className="flex justify-end pt-2 text-sm">
                  <span className="text-slate-300/70">Silver Subtotal: ₹{(silver.sub + silver.gst).toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}</span>
                </div>
              )}
            </CardContent>
          </Card>

          {/* ── OLD GOLD EXCHANGE ── */}
          <Card className="bg-card backdrop-blur-md border-[#c6a962]/20">
            <CardHeader className="flex flex-row justify-between items-center">
              <CardTitle className="text-foreground">Old Gold / Exchange</CardTitle>
              <Button variant="outline" size="sm" onClick={() => setOldGoldItems(prev => [...prev, { id: Date.now(), name: "", weight: "", purity: "22K", melting_loss_percentage: "0", rate_per_gram: "" }])} className="border-[#c6a962]/50 text-[#c6a962]">
                + Add Old Gold
              </Button>
            </CardHeader>
            <CardContent className="space-y-4">
              {oldGoldItems.length === 0 && <p className="text-center text-muted-foreground text-sm py-2">No old gold added.</p>}
              {oldGoldItems.map((og, index) => (
                <div key={og.id} className="border border-border rounded-lg p-4 relative bg-muted/30">
                  <div className="absolute -top-2.5 left-3 text-xs text-muted-foreground bg-card px-2">Old Gold {index + 1}</div>
                  <div className="grid grid-cols-2 md:grid-cols-5 gap-3 items-end">
                    <div className="space-y-1">
                      <label className="text-xs text-muted-foreground">Item Name</label>
                      <Input value={og.name} onChange={e => setOldGoldItems(prev => prev.map(i => i.id === og.id ? { ...i, name: e.target.value } : i))} placeholder="e.g. Broken Chain" className="h-8 text-sm bg-muted border-border text-foreground" />
                    </div>
                    <div className="space-y-1">
                      <label className="text-xs text-muted-foreground">Purity</label>
                      <Input value={og.purity} onChange={e => setOldGoldItems(prev => prev.map(i => i.id === og.id ? { ...i, purity: e.target.value } : i))} placeholder="22K" className="h-8 text-sm bg-muted border-border text-foreground" />
                    </div>
                    <div className="space-y-1">
                      <label className="text-xs text-muted-foreground">Weight (g)</label>
                      <Input type="number" value={og.weight} onChange={e => setOldGoldItems(prev => prev.map(i => i.id === og.id ? { ...i, weight: e.target.value } : i))} placeholder="0.000" className="h-8 text-sm bg-muted border-border text-foreground" />
                    </div>
                    <div className="space-y-1">
                      <label className="text-xs text-muted-foreground">Melting Loss %</label>
                      <Input type="number" value={og.melting_loss_percentage} onChange={e => setOldGoldItems(prev => prev.map(i => i.id === og.id ? { ...i, melting_loss_percentage: e.target.value } : i))} placeholder="0" className="h-8 text-sm bg-muted border-border text-foreground" />
                    </div>
                    <div className="space-y-1">
                      <label className="text-xs text-muted-foreground">Rate/g (₹)</label>
                      <Input type="number" value={og.rate_per_gram} onChange={e => setOldGoldItems(prev => prev.map(i => i.id === og.id ? { ...i, rate_per_gram: e.target.value } : i))} placeholder="0.00" className="h-8 text-sm bg-muted border-border text-foreground" />
                    </div>
                  </div>
                  <div className="flex justify-between items-center mt-3 pt-2 border-t border-border/50">
                    <span className="text-xs text-muted-foreground">Exchange Value: ₹{((parseFloat(og.weight) || 0) * (1 - (parseFloat(og.melting_loss_percentage) || 0) / 100) * (parseFloat(og.rate_per_gram) || 0)).toFixed(2)}</span>
                    <Button variant="ghost" size="sm" onClick={() => setOldGoldItems(prev => prev.filter(i => i.id !== og.id))} className="h-7 text-xs text-red-400 hover:text-red-300">Remove</Button>
                  </div>
                </div>
              ))}
            </CardContent>
          </Card>
        </div>

        {/* ── INVOICE SUMMARY ── */}
        <div>
          <Card className="bg-card border-[#c6a962] shadow-sm sticky top-6">
            <CardHeader>
              <CardTitle className="text-[#c6a962]">Invoice Summary</CardTitle>
            </CardHeader>
            <CardContent className="space-y-3">
              {/* Gold breakdown */}
              {goldItems.length > 0 && (
                <div className="flex justify-between text-yellow-400/80 text-sm">
                  <span>Gold ({goldItems.length} item{goldItems.length > 1 ? "s" : ""})</span>
                  <span>₹{gold.sub.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}</span>
                </div>
              )}
              {/* Silver breakdown */}
              {silverItems.length > 0 && (
                <div className="flex justify-between text-slate-300/80 text-sm">
                  <span>Silver ({silverItems.length} item{silverItems.length > 1 ? "s" : ""})</span>
                  <span>₹{silver.sub.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}</span>
                </div>
              )}
              <div className="border-t border-border/50 pt-2 space-y-2">
                <div className="flex justify-between text-muted-foreground text-sm">
                  <span>Gross Subtotal</span>
                  <span>₹{subtotal.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}</span>
                </div>
                <div className="flex justify-between text-muted-foreground text-sm">
                  <span>CGST</span>
                  <span>₹{(totalGst / 2).toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}</span>
                </div>
                <div className="flex justify-between text-muted-foreground text-sm">
                  <span>SGST</span>
                  <span>₹{(totalGst / 2).toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}</span>
                </div>
                <div className="flex justify-between text-foreground font-medium pt-1">
                  <span>Total Value</span>
                  <span>₹{grandTotal.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}</span>
                </div>
              </div>

              {oldGoldItems.length > 0 && (
                <div className="flex justify-between text-red-400 text-sm">
                  <span>– Old Gold Exchange</span>
                  <span>–₹{totalOldGold.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}</span>
                </div>
              )}

              <div className="flex justify-between items-center text-green-400 text-sm">
                <span>{billType === "estimate" ? "Advance Payment" : "Amount Paid"} (₹)</span>
                <Input type="number" value={advanceAmount} onChange={(e) => setAdvanceAmount(e.target.value)} placeholder="0.00" className="w-24 h-7 bg-muted border-border text-right text-green-400 text-sm" />
              </div>

              <div className="pt-3 border-t border-border flex justify-between items-center">
                <span className="text-lg font-bold text-foreground">Balance Due</span>
                <span className="text-2xl font-bold text-[#c6a962]">₹{balanceDue.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}</span>
              </div>

              <Button onClick={handleSubmit} disabled={mutation.isPending || (goldItems.length === 0 && silverItems.length === 0)} className="w-full mt-4 h-12 text-base font-semibold">
                {mutation.isPending ? <><Loader2 className="h-5 w-5 animate-spin mr-2" />Generating...</> : "Generate Bill"}
              </Button>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}
