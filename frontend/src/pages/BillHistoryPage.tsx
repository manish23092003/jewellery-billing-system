import { useState } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/Card";
import { Button } from "@/components/ui/Button";
import { Input } from "@/components/ui/Input";
import { useQuery } from "@tanstack/react-query";
import { getBills } from "@/lib/bills.api";
import api from "@/lib/api";
import { Loader2 } from "lucide-react";
import { toast } from "sonner";

export default function BillHistoryPage() {
  const [searchTerm, setSearchTerm] = useState("");

  const { data: bills, isLoading } = useQuery({
    queryKey: ["bills"],
    queryFn: getBills,
  });

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
            className="w-64 bg-black/50 border-gray-700 text-gray-200" 
          />
        </div>
      </div>

      <Card className="bg-black/40 backdrop-blur-md border-[#c6a962]/20">
        <CardHeader>
          <CardTitle className="text-gray-200">Recent Transactions</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="overflow-x-auto">
            <table className="w-full text-sm text-left">
              <thead className="text-xs text-gray-400 uppercase bg-black/50 border-b border-[#c6a962]/20">
                <tr>
                  <th className="px-6 py-4 font-medium">Invoice No.</th>
                  <th className="px-6 py-4 font-medium">Date</th>
                  <th className="px-6 py-4 font-medium">Customer</th>
                  <th className="px-6 py-4 font-medium">Payment</th>
                  <th className="px-6 py-4 font-medium text-right">Amount</th>
                  <th className="px-6 py-4 font-medium text-center">Action</th>
                </tr>
              </thead>
              <tbody>
                {filteredBills.length === 0 && (
                  <tr>
                    <td colSpan={6} className="text-center py-8 text-gray-500">No invoices found</td>
                  </tr>
                )}
                {filteredBills.map((bill: any) => (
                  <tr key={bill.id} className="border-b border-gray-800 hover:bg-[#c6a962]/5 transition-colors">
                    <td className="px-6 py-4 font-medium text-[#c6a962]">{bill.invoice_number}</td>
                    <td className="px-6 py-4 text-gray-300">{bill.invoice_date.split('T')[0]}</td>
                    <td className="px-6 py-4 text-gray-300">{bill.customer_name}</td>
                    <td className="px-6 py-4 text-gray-300 capitalize">{bill.payment_method}</td>
                    <td className="px-6 py-4 text-right font-bold text-gray-200">₹{bill.grand_total.toLocaleString()}</td>
                    <td className="px-6 py-4 text-center">
                      <Button 
                        variant="outline" 
                        size="sm" 
                        onClick={() => downloadPDF(bill.id, bill.invoice_number)}
                        className="h-8 border-gray-600 text-gray-300 hover:text-white"
                      >
                        View PDF
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
  );
}
