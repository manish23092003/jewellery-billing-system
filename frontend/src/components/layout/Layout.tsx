import { Navigate, Outlet, Link, useLocation } from "react-router-dom";
import { useAuth } from "@/context/AuthContext";
import { Button } from "@/components/ui/Button";
import { useQuery } from "@tanstack/react-query";
import { getSettings } from "@/lib/settings.api";

export default function Layout() {
  const { user, organization, logout, isLoading } = useAuth();
  const location = useLocation();

  const { data: settings } = useQuery({
    queryKey: ["settings"],
    queryFn: getSettings,
  });

  const shopName = settings?.shop_name || organization?.business_name || "Jewellery Billing";
  const shopShortName = shopName.split(" ")[0] || "Jewellery";

  // Wait until auth state is restored from localStorage before deciding to redirect.
  // Without this, the page flashes to login and then back, causing an infinite loop.
  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-background">
        <div className="flex flex-col items-center gap-4">
          <div className="w-10 h-10 rounded-full border-4 border-[#c6a962] border-t-transparent animate-spin" />
          <p className="text-muted-foreground text-sm">Loading…</p>
        </div>
      </div>
    );
  }

  if (!user) {
    return <Navigate to="/login" replace />;
  }

  const navItems = [
    { name: "Dashboard", path: "/" },
    { name: "Create Invoice", path: "/bills/new" },
    { name: "Invoice History", path: "/bills/history" },
    { name: "Customers", path: "/customers" },
    { name: "Metal Rates", path: "/metal-rates" },
    { name: "Expenses", path: "/expenses" },
    { name: "Settings", path: "/settings" },
  ];

  return (
    <div className="min-h-screen bg-background text-foreground flex">
      {/* Sidebar */}
      <aside className="w-64 border-r border-border bg-card/50 backdrop-blur-xl flex flex-col hidden md:flex">
        <div className="h-16 flex items-center px-6 border-b border-border space-x-3">
          {settings?.logo_path && (
            <img 
              src={`http://localhost:8080${settings.logo_path}`} 
              alt="Logo" 
              className="h-8 w-8 object-contain"
            />
          )}
          <span className="text-xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-[#c6a962] to-[#b8882a] dark:to-[#f3e5c0] truncate">
            {shopName}
          </span>
        </div>
        <div className="flex-1 py-6 px-4 space-y-2">
          {navItems.map((item) => {
            const isActive = location.pathname === item.path;
            return (
              <Link
                key={item.path}
                to={item.path}
                className={`block px-4 py-3 rounded-lg transition-all ${
                  isActive
                    ? "bg-gradient-to-r from-[#c6a962]/20 to-transparent text-[#c6a962] font-semibold border-l-2 border-[#c6a962]"
                    : "text-muted-foreground hover:text-foreground hover:bg-black/5 dark:hover:bg-white/5"
                }`}
              >
                {item.name}
              </Link>
            );
          })}
        </div>
        <div className="p-4 border-t border-border">
          <div className="flex items-center space-x-3 mb-4 px-2">
            <div className="w-8 h-8 rounded-full bg-[#c6a962] text-white flex items-center justify-center font-bold">
              {user.name.charAt(0).toUpperCase()}
            </div>
            <div>
              <p className="text-sm font-medium text-foreground">{user.name}</p>
              <p className="text-xs text-muted-foreground capitalize">{user.role}</p>
            </div>
          </div>
          <Button
            variant="outline"
            className="w-full border-border text-foreground bg-transparent hover:bg-destructive/10 hover:text-destructive hover:border-destructive/50"
            onClick={logout}
          >
            Logout
          </Button>
        </div>
      </aside>

      {/* Main Content */}
      <main className="flex-1 flex flex-col h-screen overflow-hidden">
        {/* Mobile Header */}
        <header className="h-16 border-b border-border bg-card/80 backdrop-blur-md flex items-center justify-between px-6 md:hidden">
          <div className="flex items-center space-x-2">
            {settings?.logo_path && (
              <img 
                src={`http://localhost:8080${settings.logo_path}`} 
                alt="Logo" 
                className="h-8 w-8 object-contain"
              />
            )}
            <span className="text-xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-[#c6a962] to-[#b8882a] dark:to-[#f3e5c0] truncate">
              {shopShortName}
            </span>
          </div>
          <Button variant="ghost" onClick={logout} className="text-muted-foreground">
            Logout
          </Button>
        </header>

        <div className="flex-1 overflow-y-auto relative">
          {/* Subtle background glow */}
          <div className="absolute top-0 left-1/4 w-96 h-96 bg-[#c6a962]/5 rounded-full blur-[120px] pointer-events-none"></div>
          
          <div className="relative z-10">
            <div className="p-6">
              <Outlet />
            </div>
          </div>
        </div>
      </main>
    </div>
  );
}
