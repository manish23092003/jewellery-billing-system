import { useState } from "react";
import { useNavigate, useLocation } from "react-router-dom";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Eye, EyeOff, Gem, Loader2 } from "lucide-react";
import { useAuth } from "@/context/AuthContext";
import { loginSchema, type LoginFormData } from "@/schemas/auth.schema";

/**
 * LoginPage — premium jewellery-themed login with split layout.
 *
 * Left:  Branding panel with animated gold ornament pattern
 * Right: Clean login form with validation
 */
export default function LoginPage() {
  const navigate = useNavigate();
  const location = useLocation();
  const { login } = useAuth();

  const [showPassword, setShowPassword] = useState(false);
  const [serverError, setServerError] = useState("");

  const from = (location.state as { from?: { pathname: string } })?.from?.pathname || "/";

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<LoginFormData>({
    resolver: zodResolver(loginSchema),
    defaultValues: { email: "", password: "" },
  });

  const onSubmit = async (data: LoginFormData) => {
    setServerError("");
    try {
      await login(data.email, data.password);
      navigate(from, { replace: true });
    } catch (err: unknown) {
      if (err && typeof err === "object" && "response" in err) {
        const axiosErr = err as { response?: { data?: { error?: string } } };
        setServerError(axiosErr.response?.data?.error || "Login failed");
      } else {
        setServerError("Unable to connect to server");
      }
    }
  };

  return (
    <div className="flex min-h-screen">
      {/* ── Left Panel: Branding ──────────────────────────────────── */}
      <div className="hidden lg:flex lg:w-1/2 relative overflow-hidden bg-gradient-to-br from-background via-muted/50 to-muted">
        {/* Animated gold pattern overlay */}
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

        {/* Glowing orb */}
        <div className="absolute top-1/4 left-1/2 -translate-x-1/2 w-72 h-72 rounded-full bg-gradient-to-br from-[#c6a962]/20 to-[#e5c77d]/10 blur-3xl" />

        {/* Content */}
        <div className="relative z-10 flex flex-col items-center justify-center w-full px-12 animate-fade-in">
          {/* Diamond icon */}
          <div className="mb-8 p-5 rounded-2xl bg-gradient-to-br from-[#c6a962]/20 to-[#c6a962]/5 border border-[#c6a962]/20">
            <Gem className="h-14 w-14 text-[#e5c77d]" strokeWidth={1.5} />
          </div>

          <h1
            className="text-5xl font-bold text-center mb-4 tracking-tight"
            style={{ fontFamily: "var(--font-heading)" }}
          >
            <span className="text-gold-gradient">Jewellery</span>
            <br />
            <span className="text-white/90">Billing System</span>
          </h1>

          <p className="text-white/50 text-center max-w-sm text-lg leading-relaxed mt-2">
            Manage your billing, expenses, and accounts with precision and
            elegance.
          </p>

          {/* Feature badges */}
          <div className="mt-10 flex flex-wrap justify-center gap-3">
            {["Billing", "Expenses", "Reports", "Metal Rates"].map((feature) => (
              <span
                key={feature}
                className="px-4 py-1.5 rounded-full text-xs font-medium tracking-wide
                           bg-[#c6a962]/10 text-[#e5c77d]/80 border border-[#c6a962]/20"
              >
                {feature}
              </span>
            ))}
          </div>
        </div>
      </div>

      {/* ── Right Panel: Login Form ───────────────────────────────── */}
      <div className="flex w-full lg:w-1/2 items-center justify-center px-6 py-12 bg-background">
        <div className="w-full max-w-md animate-fade-in">
          {/* Mobile header (hidden on desktop where left panel shows) */}
          <div className="lg:hidden flex flex-col items-center mb-10">
            <div className="p-3 rounded-xl bg-primary/10 mb-4">
              <Gem className="h-8 w-8 text-primary" strokeWidth={1.5} />
            </div>
            <h1
              className="text-2xl font-bold text-gold-gradient"
              style={{ fontFamily: "var(--font-heading)" }}
            >
              Jewellery Billing
            </h1>
          </div>

          {/* Form header */}
          <div className="mb-8">
            <h2
              className="text-3xl font-bold text-foreground"
              style={{ fontFamily: "var(--font-heading)" }}
            >
              Welcome back
            </h2>
            <p className="text-muted-foreground mt-2">
              Sign in to your account to continue
            </p>
          </div>

          {/* Server error banner */}
          {serverError && (
            <div className="mb-6 p-4 rounded-lg bg-destructive/10 border border-destructive/20 text-destructive text-sm animate-fade-in">
              {serverError}
            </div>
          )}

          {/* Login form */}
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-5">
            {/* Email */}
            <div className="space-y-2">
              <label
                htmlFor="email"
                className="block text-sm font-medium text-foreground"
              >
                Email Address
              </label>
              <input
                id="email"
                type="email"
                autoComplete="email"
                placeholder="admin@jewellery.com"
                className={`w-full h-11 px-4 rounded-lg border bg-background text-foreground
                  placeholder:text-muted-foreground/50 transition-all duration-200
                  focus:outline-none focus:ring-2 focus:ring-primary/40 focus:border-primary
                  ${errors.email ? "border-destructive ring-1 ring-destructive/30" : "border-input"}`}
                {...register("email")}
              />
              {errors.email && (
                <p className="text-xs text-destructive mt-1">
                  {errors.email.message}
                </p>
              )}
            </div>

            {/* Password */}
            <div className="space-y-2">
              <div className="flex items-center justify-between">
                <label
                  htmlFor="password"
                  className="block text-sm font-medium text-foreground"
                >
                  Password
                </label>
                <button
                  type="button"
                  onClick={() => navigate("/forgot-password")}
                  className="text-xs font-medium text-[#c6a962] hover:text-[#e5c77d] transition-colors"
                >
                  Forgot password?
                </button>
              </div>
              <div className="relative">
                <input
                  id="password"
                  type={showPassword ? "text" : "password"}
                  autoComplete="current-password"
                  placeholder="••••••••"
                  className={`w-full h-11 px-4 pr-11 rounded-lg border bg-background text-foreground
                    placeholder:text-muted-foreground/50 transition-all duration-200
                    focus:outline-none focus:ring-2 focus:ring-primary/40 focus:border-primary
                    ${errors.password ? "border-destructive ring-1 ring-destructive/30" : "border-input"}`}
                  {...register("password")}
                />
                <button
                  type="button"
                  tabIndex={-1}
                  onClick={() => setShowPassword((v) => !v)}
                  className="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground transition-colors"
                >
                  {showPassword ? (
                    <EyeOff className="h-4 w-4" />
                  ) : (
                    <Eye className="h-4 w-4" />
                  )}
                </button>
              </div>
              {errors.password && (
                <p className="text-xs text-destructive mt-1">
                  {errors.password.message}
                </p>
              )}
            </div>

            {/* Submit button */}
            <button
              type="submit"
              disabled={isSubmitting}
              className="w-full h-11 rounded-lg font-semibold text-sm tracking-wide
                bg-gradient-to-r from-[#c6a962] to-[#b8882a] text-white
                shadow-lg shadow-[#c6a962]/20
                hover:from-[#b8882a] hover:to-[#996824]
                disabled:opacity-60 disabled:cursor-not-allowed
                transition-all duration-300 transform hover:scale-[1.01] active:scale-[0.99]"
            >
              {isSubmitting ? (
                <span className="flex items-center justify-center gap-2">
                  <Loader2 className="h-4 w-4 animate-spin" />
                  Signing in…
                </span>
              ) : (
                "Sign In"
              )}
            </button>
          </form>

          {/* Footer links */}
          <div className="mt-8 text-center space-y-4">
            <p className="text-sm text-muted-foreground">
              Don't have an account?{" "}
              <button
                onClick={() => navigate("/register")}
                className="font-medium text-[#c6a962] hover:text-[#e5c77d] transition-colors"
              >
                Register your business
              </button>
            </p>
            <p className="text-xs text-muted-foreground/60">
              © {new Date().getFullYear()} Jewellery Billing System
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}
