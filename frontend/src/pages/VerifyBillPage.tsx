import { useEffect, useState } from "react";
import { useParams, Link } from "react-router-dom";
import { CheckCircle, XCircle, AlertTriangle, ArrowLeft } from "lucide-react";
import { verifyInvoice } from "@/lib/bills.api";
import type { PublicVerificationResponse } from "@/lib/bills.api";
import { formatCurrency } from "@/lib/utils";

export default function VerifyBillPage() {
  const { id } = useParams<{ id: string }>(); // the token
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [data, setData] = useState<PublicVerificationResponse | null>(null);

  useEffect(() => {
    if (!id) return;
    verifyInvoice(id)
      .then((res) => {
        setData(res);
        setLoading(false);
      })
      .catch((err) => {
        setError(err.message || "Failed to verify invoice");
        setLoading(false);
      });
  }, [id]);

  if (loading) {
    return (
      <div className="flex h-screen items-center justify-center bg-gray-50 dark:bg-zinc-900">
        <div className="text-center">
          <div className="mb-4 h-12 w-12 animate-spin rounded-full border-4 border-primary border-t-transparent mx-auto"></div>
          <p className="text-lg font-medium text-muted-foreground">Verifying secure invoice...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-zinc-950 py-12 px-4 sm:px-6 lg:px-8 flex items-center justify-center">
      <div className="max-w-md w-full space-y-8 bg-white dark:bg-zinc-900 p-8 rounded-xl shadow-2xl border border-gray-100 dark:border-zinc-800 relative overflow-hidden">
        
        {/* Background gradient blur */}
        <div className="absolute -top-24 -left-24 w-48 h-48 bg-primary/20 rounded-full blur-3xl pointer-events-none"></div>
        <div className="absolute -bottom-24 -right-24 w-48 h-48 bg-primary/10 rounded-full blur-3xl pointer-events-none"></div>

        <div className="relative z-10">
          {error ? (
            <div className="text-center space-y-4">
              <XCircle className="mx-auto h-16 w-16 text-red-500" />
              <h2 className="text-2xl font-bold text-gray-900 dark:text-gray-100">Verification Failed</h2>
              <p className="text-red-500 bg-red-50 dark:bg-red-500/10 p-3 rounded-lg border border-red-100 dark:border-red-500/20">{error}</p>
            </div>
          ) : data ? (
            <div className="text-center space-y-6">
              
              {data.verification_status === "VERIFIED" ? (
                <div className="space-y-2">
                  <CheckCircle className="mx-auto h-20 w-20 text-green-500 animate-in zoom-in duration-500" />
                  <h2 className="text-2xl font-bold text-green-600 dark:text-green-400">Authentic Invoice</h2>
                  <p className="text-sm text-gray-500 dark:text-gray-400">This invoice has been securely verified against cryptographic records.</p>
                </div>
              ) : data.verification_status === "TAMPERED" ? (
                <div className="space-y-2">
                  <XCircle className="mx-auto h-20 w-20 text-red-500 animate-in zoom-in duration-500" />
                  <h2 className="text-2xl font-bold text-red-600 dark:text-red-400">Tampered Invoice</h2>
                  <p className="text-sm text-red-500">Warning: The details on this invoice do not match our secure cryptographic records. This document may have been forged.</p>
                </div>
              ) : (
                <div className="space-y-2">
                  <AlertTriangle className="mx-auto h-20 w-20 text-amber-500 animate-in zoom-in duration-500" />
                  <h2 className="text-2xl font-bold text-amber-600 dark:text-amber-400">Invoice {data.verification_status}</h2>
                  <p className="text-sm text-amber-500">This invoice exists but its status is {data.verification_status}.</p>
                </div>
              )}

              <div className="bg-gray-50 dark:bg-zinc-800/50 rounded-xl p-6 text-left border border-gray-100 dark:border-zinc-800 shadow-inner">
                <div className="space-y-4">
                  <div>
                    <p className="text-xs font-semibold text-gray-400 uppercase tracking-wider mb-1">Issuer / Shop</p>
                    <p className="text-lg font-medium text-gray-900 dark:text-gray-100">{data.shop_name}</p>
                  </div>
                  <div className="grid grid-cols-2 gap-4">
                    <div>
                      <p className="text-xs font-semibold text-gray-400 uppercase tracking-wider mb-1">Invoice No</p>
                      <p className="font-mono text-gray-900 dark:text-gray-100">{data.invoice_number}</p>
                    </div>
                    <div>
                      <p className="text-xs font-semibold text-gray-400 uppercase tracking-wider mb-1">Date</p>
                      <p className="text-gray-900 dark:text-gray-100">{data.invoice_date}</p>
                    </div>
                  </div>
                  <div className="pt-4 border-t border-gray-200 dark:border-zinc-700">
                    <p className="text-xs font-semibold text-gray-400 uppercase tracking-wider mb-1">Grand Total</p>
                    <p className="text-2xl font-bold text-gray-900 dark:text-gray-100">{formatCurrency(data.grand_total)}</p>
                  </div>
                  <div className="pt-4 border-t border-gray-200 dark:border-zinc-700 grid grid-cols-2 gap-4">
                    <div>
                      <p className="text-xs font-semibold text-gray-400 uppercase tracking-wider mb-1">Payment Status</p>
                      <p className={`font-semibold ${data.balance_due <= 0 ? "text-green-600 dark:text-green-400" : "text-amber-600 dark:text-amber-400"}`}>
                        {data.balance_due <= 0 ? "Fully Paid" : "Pending Balance"}
                      </p>
                    </div>
                    {data.balance_due > 0 && (
                      <div>
                        <p className="text-xs font-semibold text-gray-400 uppercase tracking-wider mb-1">Remaining</p>
                        <p className="text-lg font-bold text-red-600 dark:text-red-400">{formatCurrency(data.balance_due)}</p>
                      </div>
                    )}
                  </div>
                </div>
              </div>

            </div>
          ) : null}

          <div className="mt-8 text-center">
             <Link to="/" className="inline-flex items-center text-sm font-medium text-primary hover:text-primary/80 transition-colors">
               <ArrowLeft className="mr-2 h-4 w-4" /> Return to Home
             </Link>
          </div>
        </div>
      </div>
    </div>
  );
}
