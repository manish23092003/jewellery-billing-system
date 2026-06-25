import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Link, useNavigate } from "react-router-dom";
import { Gem, Loader2, Store, User as UserIcon, Mail, Phone, Lock, ArrowRight, ArrowLeft, CheckCircle2 } from "lucide-react";
import { registerSchema, type RegisterFormValues } from "@/schemas/register.schema";
import { useAuth } from "@/context/AuthContext";
import { toast } from "sonner";

export default function RegisterPage() {
  const [step, setStep] = useState<1 | 2>(1);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const { register: registerBusiness } = useAuth();
  const navigate = useNavigate();

  const {
    register,
    handleSubmit,
    trigger,
    formState: { errors },
  } = useForm<RegisterFormValues>({
    resolver: zodResolver(registerSchema),
    mode: "onChange",
  });

  const onNextStep = async () => {
    // Validate Step 1 fields before proceeding
    const isStep1Valid = await trigger(["business_name", "phone"]);
    if (isStep1Valid) {
      setStep(2);
    }
  };

  const onSubmit = async (data: RegisterFormValues) => {
    setIsSubmitting(true);
    try {
      await registerBusiness(data);
      toast.success("Account created successfully! Welcome aboard.");
      navigate("/");
    } catch (error: any) {
      // Extract the real server error message from the axios response
      const message = error?.response?.data?.error || error?.message || "Registration failed";
      toast.error(message);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="flex min-h-screen bg-background">
      {/* ── Left Panel: Branding & Features ──────────────────────────────────── */}
      <div className="hidden lg:flex lg:w-[45%] relative overflow-hidden bg-gradient-to-br from-background via-muted/50 to-muted">
        <div className="absolute inset-0 opacity-10">
          <div
            className="absolute inset-0"
            style={{
              backgroundImage: `radial-gradient(circle at 25% 25%, #c6a962 1px, transparent 1px),
                               radial-gradient(circle at 75% 75%, #c6a962 1px, transparent 1px)`,
              backgroundSize: "60px 60px",
            }}
          />
        </div>
        <div className="absolute top-1/4 left-1/2 -translate-x-1/2 w-96 h-96 rounded-full bg-gradient-to-br from-[#c6a962]/20 to-[#e5c77d]/10 blur-3xl" />

        <div className="relative z-10 flex flex-col justify-between w-full p-12 lg:px-16 animate-fade-in h-full">
          <div>
            <div className="flex items-center gap-3">
              <div className="p-3 rounded-xl bg-gradient-to-br from-[#c6a962]/20 to-[#c6a962]/5 border border-[#c6a962]/20">
                <Gem className="h-8 w-8 text-[#e5c77d]" strokeWidth={1.5} />
              </div>
              <span className="text-xl font-bold tracking-widest text-[#e5c77d] uppercase" style={{ fontFamily: "var(--font-heading)" }}>
                Jewellery Billing
              </span>
            </div>
            
            <h1 className="text-4xl xl:text-5xl font-bold mt-16 leading-tight tracking-tight text-foreground" style={{ fontFamily: "var(--font-heading)" }}>
              Elevate your <br/>
              <span className="text-transparent bg-clip-text bg-gradient-to-r from-[#c6a962] to-[#f3e5c0]">
                Jewellery Business
              </span>
            </h1>
            <p className="text-muted-foreground mt-6 text-lg max-w-md leading-relaxed">
              Join thousands of leading jewellers who trust our platform to manage their billing, inventory, and analytics with unparalleled elegance.
            </p>
          </div>

          <div className="space-y-6 pb-8">
            {[
              "Instant GST Invoicing & Billing",
              "Real-time Multi-tenant Data Isolation",
              "Advanced Analytics & Daily Metal Rates",
            ].map((feature, idx) => (
              <div key={idx} className="flex items-center gap-4">
                <div className="flex-shrink-0 w-6 h-6 rounded-full bg-[#c6a962]/20 flex items-center justify-center border border-[#c6a962]/30">
                  <CheckCircle2 className="w-4 h-4 text-[#e5c77d]" />
                </div>
                <span className="text-muted-foreground font-medium">{feature}</span>
              </div>
            ))}
          </div>
        </div>
      </div>

      {/* ── Right Panel: Registration Form ───────────────────────────────── */}
      <div className="flex-1 flex flex-col justify-center px-6 py-12 lg:px-20 xl:px-32 bg-background relative">
        <div className="absolute top-0 right-0 w-[500px] h-[500px] bg-gradient-to-bl from-[#c6a962]/5 to-transparent rounded-bl-full pointer-events-none blur-3xl" />
        
        <div className="w-full max-w-md mx-auto relative z-10 animate-fade-in">
          {/* Mobile header */}
          <div className="lg:hidden flex items-center gap-3 mb-10 justify-center">
            <div className="p-2 rounded-xl bg-[#c6a962]/10 border border-[#c6a962]/20">
              <Gem className="h-6 w-6 text-[#e5c77d]" />
            </div>
            <span className="text-xl font-bold text-[#e5c77d] uppercase" style={{ fontFamily: "var(--font-heading)" }}>
              Jewellery Billing
            </span>
          </div>

          <div className="mb-10">
            <h2 className="text-3xl font-bold text-foreground mb-2" style={{ fontFamily: "var(--font-heading)" }}>
              {step === 1 ? "Create your workspace" : "Admin credentials"}
            </h2>
            <p className="text-muted-foreground">
              {step === 1 ? "Let's start with your business details." : "Secure your account with an admin login."}
            </p>
          </div>

          {/* Stepper Indicator */}
          <div className="flex items-center gap-2 mb-10">
            <div className={`h-1.5 flex-1 rounded-full transition-all duration-500 ${step >= 1 ? "bg-gradient-to-r from-[#c6a962] to-[#e5c77d]" : "bg-muted"}`} />
            <div className={`h-1.5 flex-1 rounded-full transition-all duration-500 ${step >= 2 ? "bg-gradient-to-r from-[#e5c77d] to-[#c6a962]" : "bg-muted"}`} />
          </div>

          <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
            <div className={`space-y-6 transition-all duration-500 ${step === 1 ? "block opacity-100" : "hidden opacity-0"}`}>
              {/* Business Name */}
              <div className="space-y-2">
                <label className="text-sm font-medium text-muted-foreground ml-1">Business Name <span className="text-red-500">*</span></label>
                <div className="relative group">
                  <div className="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
                    <Store className="h-5 w-5 text-muted-foreground group-focus-within:text-[#c6a962] transition-colors" />
                  </div>
                  <input
                    {...register("business_name")}
                    type="text"
                    className={`w-full h-12 pl-12 rounded-xl bg-white/5 border ${errors.business_name ? 'border-red-500/50' : 'border-white/10'} text-foreground placeholder-muted-foreground focus:outline-none focus:border-[#c6a962] focus:bg-white/10 transition-all`}
                    placeholder="E.g. Krishna Jewellers Pvt Ltd"
                  />
                </div>
                {errors.business_name && <p className="text-xs text-red-500 ml-1 mt-1">{errors.business_name.message}</p>}
              </div>

              {/* Phone */}
              <div className="space-y-2">
                <label className="text-sm font-medium text-muted-foreground ml-1">Business Phone <span className="text-red-500">*</span></label>
                <div className="relative group">
                  <div className="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
                    <Phone className="h-5 w-5 text-muted-foreground group-focus-within:text-[#c6a962] transition-colors" />
                  </div>
                  <input
                    {...register("phone")}
                    type="tel"
                    className={`w-full h-12 pl-12 rounded-xl bg-white/5 border ${errors.phone ? 'border-red-500/50' : 'border-white/10'} text-foreground placeholder-muted-foreground focus:outline-none focus:border-[#c6a962] focus:bg-white/10 transition-all`}
                    placeholder="+91 98765 43210"
                  />
                </div>
                {errors.phone && <p className="text-xs text-red-500 ml-1 mt-1">{errors.phone.message}</p>}
              </div>

              <button
                type="button"
                onClick={onNextStep}
                className="w-full h-12 mt-4 flex items-center justify-center gap-2 rounded-xl font-semibold bg-white text-black hover:bg-gray-200 transition-colors"
              >
                Continue <ArrowRight className="w-4 h-4" />
              </button>
            </div>

            <div className={`space-y-6 transition-all duration-500 ${step === 2 ? "block opacity-100 animate-fade-in" : "hidden opacity-0"}`}>
              {/* Owner Name */}
              <div className="space-y-2">
                <label className="text-sm font-medium text-muted-foreground ml-1">Your Full Name <span className="text-red-500">*</span></label>
                <div className="relative group">
                  <div className="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
                    <UserIcon className="h-5 w-5 text-muted-foreground group-focus-within:text-[#c6a962] transition-colors" />
                  </div>
                  <input
                    {...register("owner_name")}
                    type="text"
                    className={`w-full h-12 pl-12 rounded-xl bg-white/5 border ${errors.owner_name ? 'border-red-500/50' : 'border-white/10'} text-foreground placeholder-muted-foreground focus:outline-none focus:border-[#c6a962] focus:bg-white/10 transition-all`}
                    placeholder="John Doe"
                  />
                </div>
                {errors.owner_name && <p className="text-xs text-red-500 ml-1 mt-1">{errors.owner_name.message}</p>}
              </div>

              {/* Email */}
              <div className="space-y-2">
                <label className="text-sm font-medium text-muted-foreground ml-1">Email Address <span className="text-red-500">*</span></label>
                <div className="relative group">
                  <div className="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
                    <Mail className="h-5 w-5 text-muted-foreground group-focus-within:text-[#c6a962] transition-colors" />
                  </div>
                  <input
                    {...register("email")}
                    type="email"
                    className={`w-full h-12 pl-12 rounded-xl bg-white/5 border ${errors.email ? 'border-red-500/50' : 'border-white/10'} text-foreground placeholder-muted-foreground focus:outline-none focus:border-[#c6a962] focus:bg-white/10 transition-all`}
                    placeholder="admin@jewellery.com"
                  />
                </div>
                {errors.email && <p className="text-xs text-red-500 ml-1 mt-1">{errors.email.message}</p>}
              </div>

              {/* Password */}
              <div className="space-y-2">
                <label className="text-sm font-medium text-muted-foreground ml-1">Password <span className="text-red-500">*</span></label>
                <div className="relative group">
                  <div className="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
                    <Lock className="h-5 w-5 text-muted-foreground group-focus-within:text-[#c6a962] transition-colors" />
                  </div>
                  <input
                    {...register("password")}
                    type="password"
                    className={`w-full h-12 pl-12 rounded-xl bg-white/5 border ${errors.password ? 'border-red-500/50' : 'border-white/10'} text-foreground placeholder-muted-foreground focus:outline-none focus:border-[#c6a962] focus:bg-white/10 transition-all`}
                    placeholder="••••••••"
                  />
                </div>
                {errors.password && <p className="text-xs text-red-500 ml-1 mt-1">{errors.password.message}</p>}
              </div>

              {/* Confirm Password */}
              <div className="space-y-2">
                <label className="text-sm font-medium text-muted-foreground ml-1">Confirm Password <span className="text-red-500">*</span></label>
                <div className="relative group">
                  <div className="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
                    <Lock className="h-5 w-5 text-muted-foreground group-focus-within:text-[#c6a962] transition-colors" />
                  </div>
                  <input
                    {...register("confirm_password")}
                    type="password"
                    className={`w-full h-12 pl-12 rounded-xl bg-white/5 border ${errors.confirm_password ? 'border-red-500/50' : 'border-white/10'} text-foreground placeholder-muted-foreground focus:outline-none focus:border-[#c6a962] focus:bg-white/10 transition-all`}
                    placeholder="••••••••"
                  />
                </div>
                {errors.confirm_password && <p className="text-xs text-red-500 ml-1 mt-1">{errors.confirm_password.message}</p>}
              </div>

              <div className="flex gap-4 mt-8">
                <button
                  type="button"
                  onClick={() => setStep(1)}
                  className="w-12 h-12 flex items-center justify-center rounded-xl bg-muted border-border border text-foreground transition-all"
                >
                  <ArrowLeft className="w-5 h-5" />
                </button>
                <button
                  type="submit"
                  disabled={isSubmitting}
                  className="flex-1 h-12 flex items-center justify-center rounded-xl font-bold bg-gradient-to-r from-[#c6a962] to-[#b8882a] text-foreground hover:from-[#b8882a] hover:to-[#996824] shadow-lg shadow-[#c6a962]/20 transition-all disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  {isSubmitting ? (
                    <span className="flex items-center gap-2"><Loader2 className="w-5 h-5 animate-spin" /> Creating...</span>
                  ) : (
                    "Launch Workspace"
                  )}
                </button>
              </div>
            </div>
          </form>

          <div className="mt-8 text-center">
            <p className="text-sm text-muted-foreground">
              Already have a workspace?{" "}
              <Link to="/login" className="font-semibold text-[#c6a962] hover:text-[#e5c77d] transition-colors">
                Sign in here
              </Link>
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}
