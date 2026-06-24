import { useState } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/Card";
import { Button } from "@/components/ui/Button";
import { Input } from "@/components/ui/Input";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { getBills, addPayment } from "@/lib/bills.api";
import api from "@/lib/api";
import { Loader2, Trash2 } from "lucide-react";
import { toast } from "sonner";
import { useNavigate } from "react-router-dom";

export default function BillHistoryPage() {
  const navigate = useNavigate();
  const [searchTerm, setSearchTerm] = useState("");
  const [isPaymentModalOpen, setIsPaymentModalOpen] = useState(false);
  const [selectedBillForPayment, setSelectedBillForPayment] = useState<any>(null);
  const [paymentAmount, setPaymentAmount] = useState<string>("");

  const queryClient = useQueryClient();

  const { data: bills, isLoading } = useQuery({
    queryKey: ["bills"],
    queryFn: getBills,
  });

  const deleteMutation = useMutation({
    mutationFn: async (id: string) => {
      await api.delete(`/bills/${id}`);
    },
    onSuccess: () => {
      toast.success("Invoice deleted successfully");
      queryClient.invalidateQueries({ queryKey: ["bills"] });
    },
    onError: () => {
      toast.error("Failed to delete invoice");
    }
  });

  const paymentMutation = useMutation({
    mutationFn: addPayment,
    onSuccess: () => {
      toast.success("Payment added successfully");
      queryClient.invalidateQueries({ queryKey: ["bills"] });
      closePaymentModal();
    },
    onError: (err: any) => {
      toast.error(err.message || "Failed to add payment");
    }
  });

  const openPaymentModal = (bill: any) => {
    setSelectedBillForPayment(bill);
    setPaymentAmount(bill.balance_due.toString());
    setIsPaymentModalOpen(true);
  };

  const closePaymentModal = () => {
    setIsPaymentModalOpen(false);
    setSelectedBillForPayment(null);
    setPaymentAmount("");
  };

  const handlePaymentSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedBillForPayment) return;
    const amount = parseFloat(paymentAmount);
    if (isNaN(amount) || amount <= 0) {
      toast.error("Please enter a valid amount");
      return;
    }
    paymentMutation.mutate({ id: selectedBillForPayment.id, amount });
  };

  const handleDelete = (id: string) => {
    if (window.confirm("Are you sure you want to delete this invoice? This action cannot be undone.")) {
      deleteMutation.mutate(id);
    }
  };

  const downloadPDF = async (id: string, invoiceNum: string) => {
    try {
      toast.loading("Generating PDF...", { id: "pdf" });
      const response = await api.get(`/bills/${id}/pdf`, { responseType: 'blob' });
      const url = window.URL.createObjectURL(new Blob([response.data], { type: 'application/pdf' }));
      const link = document.createElement('a');
      link.href = url;
      link.setAttribute('download', `${invoiceNum}.pdf`);
      document.body.appendChild(link);
      link.click();
      link.parentNode?.removeChild(link);
      toast.success("PDF downloaded successfully", { id: "pdf" });
    } catch (error) {
      toast.error("Failed to download PDF", { id: "pdf" });
    }
  };

  const shareWhatsApp = async (bill: any) => {
    try {
      toast.loading("Preparing PDF for sharing...", { id: "share" });
      const response = await api.get(`/bills/${bill.id}/pdf`, { responseType: 'blob' });
      const file = new File([response.data], `Invoice_${bill.invoice_number}.pdf`, { type: 'application/pdf' });
      
      const text = `Hello ${bill.customer_name},\n\nThank you for your purchase! Your Invoice No is *${bill.invoice_number}* for an amount of *₹${bill.grand_total.toLocaleString()}*.\n\nHave a great day!`;

      // Attempt native file sharing (Mobile / Supported Desktops)
      if (navigator.canShare && navigator.canShare({ files: [file] })) {
        toast.dismiss("share");
        await navigator.share({
          title: `Invoice ${bill.invoice_number}`,
          text: text,
          files: [file]
        });
        return;
      }

      // Fallback for browsers that don't support native file sharing
      toast.dismiss("share");
      toast.info("Downloading PDF. Please attach it in WhatsApp.", { duration: 4000 });
      
      // Trigger download
      const urlBlob = window.URL.createObjectURL(new Blob([response.data], { type: 'application/pdf' }));
      const link = document.createElement('a');
      link.href = urlBlob;
      link.setAttribute('download', `Invoice_${bill.invoice_number}.pdf`);
      document.body.appendChild(link);
      link.click();
      link.parentNode?.removeChild(link);

      // Open WhatsApp web with text
      const phone = bill.customer_phone ? bill.customer_phone.replace(/\D/g, '') : '';
      const fallbackText = text + "\n\n(Please see the attached PDF)";
      const urlWa = `https://wa.me/${phone}?text=${encodeURIComponent(fallbackText)}`;
      setTimeout(() => window.open(urlWa, '_blank'), 1000);

    } catch (error) {
      toast.error("Failed to share", { id: "share" });
    }
  };

  const filteredBills = bills?.filter((b: any) => 
    b.invoice_number.toLowerCase().includes(searchTerm.toLowerCase()) ||
    b.customer_name.toLowerCase().includes(searchTerm.toLowerCase())
  ) || [];

  if (isLoading) {
    return <div className="flex h-[50vh] items-center justify-center"><Loader2 className="h-8 w-8 animate-spin text-[#c6a962]" /></div>;
  }

  return (
    <div className="p-6 max-w-7xl mx-auto space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-3xl font-bold tracking-tight text-[#c6a962]">
          Invoice History
        </h1>
        <div className="flex space-x-4">
          <Input 
            placeholder="Search invoice or customer..." 
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className="w-64 bg-muted border-border text-foreground" 
          />
        </div>
      </div>

      <Card className="bg-card backdrop-blur-md border-[#c6a962]/20">
        <CardHeader>
          <CardTitle className="text-foreground">Recent Transactions</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="overflow-x-auto">
            <table className="w-full min-w-[1000px] text-sm text-left">
              <thead className="text-xs text-muted-foreground uppercase bg-muted border-b border-[#c6a962]/20">
                <tr>
                  <th className="px-6 py-4 font-medium">Document No.</th>
                  <th className="px-6 py-4 font-medium">Date</th>
                  <th className="px-6 py-4 font-medium">Customer</th>
                  <th className="px-6 py-4 font-medium">Payment</th>
                  <th className="px-6 py-4 font-medium text-right">Amount</th>
                  <th className="px-6 py-4 font-medium text-right">Balance Due</th>
                  <th className="px-6 py-4 font-medium text-center">Status</th>
                  <th className="px-6 py-4 pr-10 font-medium text-right">Action</th>
                </tr>
              </thead>
              <tbody>
                {filteredBills.length === 0 && (
                  <tr>
                    <td colSpan={8} className="text-center py-8 text-muted-foreground">No invoices found</td>
                  </tr>
                )}
                {filteredBills.map((bill: any) => (
                  <tr key={bill.id} className="border-b border-border hover:bg-[#c6a962]/5 transition-colors">
                    <td className="px-6 py-4 font-medium text-[#c6a962]">
                      {bill.invoice_number}
                      {bill.type === "estimate" && (
                        <span className="ml-2 text-[10px] uppercase bg-blue-500/20 text-blue-300 px-2 py-0.5 rounded border border-blue-500/30">
                          Estimate
                        </span>
                      )}
                    </td>
                    <td className="px-6 py-4 text-foreground/90">{bill.invoice_date.split('T')[0]}</td>
                    <td className="px-6 py-4 text-foreground/90">{bill.customer_name}</td>
                    <td className="px-6 py-4 text-foreground/90 capitalize">{bill.payment_method}</td>
                    <td className="px-6 py-4 text-right font-bold text-foreground">₹{bill.grand_total.toLocaleString()}</td>
                    <td className="px-6 py-4 text-right font-bold text-red-400">
                      ₹{bill.balance_due.toLocaleString(undefined, {minimumFractionDigits: 2, maximumFractionDigits: 2})}
                    </td>
                    <td className="px-6 py-4 text-center">
                      {bill.balance_due <= 0 ? (
                        <span className="text-xs uppercase bg-green-500/20 text-green-400 px-2 py-1 rounded-full border border-green-500/30 font-semibold tracking-wide">
                          Paid
                        </span>
                      ) : (
                        <span className="text-xs uppercase bg-yellow-500/20 text-yellow-400 px-2 py-1 rounded-full border border-yellow-500/30 font-semibold tracking-wide">
                          Pending
                        </span>
                      )}
                    </td>
                    <td className="px-6 py-4">
                      <div className="flex items-center justify-end gap-2 pr-4">
                      {bill.balance_due > 0 && bill.type !== "estimate" && (
                        <Button 
                          variant="outline" 
                          size="sm" 
                          onClick={() => openPaymentModal(bill)}
                          className="h-8 border-[#c6a962]/50 text-[#c6a962] hover:bg-[#c6a962]/10"
                        >
                          Pay
                        </Button>
                      )}
                      {bill.type === "estimate" && (
                        <Button 
                          variant="outline" 
                          size="sm" 
                          onClick={() => navigate(`/bills/new?convertFrom=${bill.id}`)}
                          className="h-8 border-blue-500/50 text-blue-400 hover:bg-blue-500/10"
                        >
                          Convert
                        </Button>
                      )}
                      <Button 
                        variant="outline" 
                        size="sm" 
                        onClick={() => downloadPDF(bill.id, bill.invoice_number)}
                        className="h-8 border-gray-600 text-foreground/90 hover:text-white"
                      >
                        PDF
                      </Button>
                      <Button 
                        variant="outline" 
                        size="sm" 
                        onClick={() => shareWhatsApp(bill)}
                        className="h-8 border-green-600/50 text-green-400 hover:bg-green-500/10 hover:text-green-300"
                        title="Share on WhatsApp"
                      >
                        WhatsApp
                      </Button>
                      <Button 
                        variant="ghost" 
                        size="icon" 
                        onClick={() => handleDelete(bill.id)}
                        disabled={deleteMutation.isPending}
                        className="h-8 w-8 text-red-400 hover:text-red-500 hover:bg-red-500/10"
                        title="Delete Invoice"
                      >
                        <Trash2 className="h-4 w-4" />
                      </Button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </CardContent>
      </Card>

      {isPaymentModalOpen && selectedBillForPayment && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm p-4">
          <div className="bg-card border border-border rounded-xl shadow-2xl w-full max-w-sm overflow-hidden">
            <div className="p-6 border-b border-border">
              <h2 className="text-xl font-semibold text-foreground">
                Add Payment
              </h2>
              <p className="text-sm text-muted-foreground mt-1">
                {selectedBillForPayment.invoice_number} - {selectedBillForPayment.customer_name}
              </p>
            </div>
            <form onSubmit={handlePaymentSubmit} className="p-6 space-y-4">
              <div className="space-y-2">
                <label className="text-sm text-muted-foreground">Amount (Max: ₹{selectedBillForPayment.balance_due})</label>
                <Input
                  required
                  type="number"
                  step="0.01"
                  min="0.01"
                  max={selectedBillForPayment.balance_due}
                  value={paymentAmount}
                  onChange={(e) => setPaymentAmount(e.target.value)}
                  className="bg-muted border-border text-white text-lg font-bold"
                />
              </div>
              <div className="pt-4 flex justify-end gap-3">
                <Button type="button" variant="outline" onClick={closePaymentModal} className="border-border text-foreground/90">
                  Cancel
                </Button>
                <Button type="submit" disabled={paymentMutation.isPending} className="bg-[#c6a962] text-black hover:bg-[#b0965a]">
                  {paymentMutation.isPending ? "Saving..." : "Add Payment"}
                </Button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
