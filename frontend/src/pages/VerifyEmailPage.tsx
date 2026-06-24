import { useState, useEffect } from "react";
import { Link, useSearchParams } from "react-router-dom";
import { Gem, Loader2, CheckCircle, XCircle } from "lucide-react";
import { verifyEmail } from "@/lib/auth.api";

export default function VerifyEmailPage() {
  const [status, setStatus] = useState<"loading" | "success" | "error">("loading");
  const [errorMsg, setErrorMsg] = useState("");
  const [searchParams] = useSearchParams();
  const token = searchParams.get("token");

  useEffect(() => {
    if (!token) {
      setStatus("error");
      setErrorMsg("No verification token found.");
      return;
    }

    const verify = async () => {
      try {
        await verifyEmail({ token });
        setStatus("success");
      } catch (error: any) {
        setStatus("error");
        setErrorMsg(error.message || "Failed to verify email");
      }
    };

    verify();
  }, [token]);

  return (
    <div className="min-h-screen bg-black flex flex-col justify-center py-12 sm:px-6 lg:px-8 relative overflow-hidden">
      <div className="absolute inset-0 bg-[url('https://images.unsplash.com/photo-1599643478524-fb66f70d00f8?q=80&w=2938&auto=format&fit=crop')] bg-cover bg-center opacity-10"></div>
      <div className="absolute inset-0 bg-gradient-to-t from-black via-black/80 to-transparent"></div>

      <div className="sm:mx-auto sm:w-full sm:max-w-md relative z-10">
        <div className="flex justify-center">
          <div className="w-16 h-16 rounded-2xl bg-gradient-to-tr from-amber-500 to-yellow-300 p-0.5 shadow-lg shadow-amber-500/20">
            <div className="w-full h-full bg-black rounded-2xl flex items-center justify-center">
              <Gem className="w-8 h-8 text-amber-500" />
            </div>
          </div>
        </div>
        
        <div className="mt-8 bg-zinc-900/50 backdrop-blur-xl py-8 px-4 shadow-2xl sm:rounded-3xl sm:px-10 border border-white/10 text-center">
          {status === "loading" && (
            <div className="py-8">
              <Loader2 className="w-12 h-12 text-amber-500 animate-spin mx-auto mb-4" />
              <h2 className="text-xl font-bold text-white">Verifying your email...</h2>
              <p className="text-gray-400 mt-2">Please wait a moment.</p>
            </div>
          )}

          {status === "success" && (
            <div className="py-8">
              <CheckCircle className="w-16 h-16 text-green-500 mx-auto mb-4" />
              <h2 className="text-2xl font-bold text-white mb-2">Email Verified!</h2>
              <p className="text-gray-400 mb-6">Your email has been successfully verified. You can now access all features.</p>
              <Link
                to="/"
                className="w-full flex justify-center py-3 px-4 border border-transparent rounded-xl shadow-sm text-sm font-medium text-black bg-gradient-to-r from-amber-500 to-yellow-400 hover:from-amber-400 hover:to-yellow-300 transition-all"
              >
                Continue to Dashboard
              </Link>
            </div>
          )}

          {status === "error" && (
            <div className="py-8">
              <XCircle className="w-16 h-16 text-red-500 mx-auto mb-4" />
              <h2 className="text-2xl font-bold text-white mb-2">Verification Failed</h2>
              <p className="text-gray-400 mb-6">{errorMsg}</p>
              <Link
                to="/login"
                className="w-full flex justify-center py-3 px-4 border border-transparent rounded-xl shadow-sm text-sm font-medium text-black bg-gradient-to-r from-amber-500 to-yellow-400 hover:from-amber-400 hover:to-yellow-300 transition-all"
              >
                Back to Login
              </Link>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
