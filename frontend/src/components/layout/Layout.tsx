import { Navigate, Outlet, Link, useLocation } from "react-router-dom";
import { useAuth } from "@/context/AuthContext";
import { Button } from "@/components/ui/Button";
import { useQuery } from "@tanstack/react-query";
import { getSettings } from "@/lib/settings.api";

export default function Layout() {
  const { user, logout } = useAuth();
  const location = useLocation();

  const { data: settings } = useQuery({
    queryKey: ["settings"],
    queryFn: getSettings,
  });

  const shopName = settings?.shop_name || "Aura Jewels";
  const shopShortName = shopName.split(" ")[0] || "Aura";

  if (!user) {
    return <Navigate to="/login" replace />;
  }

  const navItems = [
    { name: "Dashboard", path: "/" },
    { name: "Create Invoice", path: "/bills/new" },
    { name: "Invoice History", path: "/bills/history" },
    { name: "Metal Rates", path: "/metal-rates" },
    { name: "Expenses", path: "/expenses" },
    { name: "Settings", path: "/settings" },
  ];

  return (
    <div className="min-h-screen bg-[#0a0a0a] text-gray-100 flex">
      {/* Sidebar */}
      <aside className="w-64 border-r border-[#c6a962]/20 bg-black/40 backdrop-blur-xl flex flex-col hidden md:flex">
        <div className="h-16 flex items-center px-6 border-b border-[#c6a962]/20 space-x-3">
          {settings?.logo_path && (
            <img 
              src={`http://localhost:8080${settings.logo_path}`} 
              alt="Logo" 
              className="h-8 w-8 object-contain"
            />
          )}
          <span className="text-xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-[#c6a962] to-[#f3e5c0] truncate">
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
                    : "text-gray-400 hover:text-gray-200 hover:bg-white/5"
                }`}
              >
                {item.name}
              </Link>
            );
          })}
        </div>
        <div className="p-4 border-t border-[#c6a962]/20">
          <div className="flex items-center space-x-3 mb-4 px-2">
            <div className="w-8 h-8 rounded-full bg-[#c6a962] text-black flex items-center justify-center font-bold">
              {user.name.charAt(0).toUpperCase()}
            </div>
            <div>
              <p className="text-sm font-medium text-gray-200">{user.name}</p>
              <p className="text-xs text-gray-500 capitalize">{user.role}</p>
            </div>
          </div>
          <Button
            variant="outline"
            className="w-full border-gray-700 text-gray-300 bg-transparent hover:bg-red-500/10 hover:text-red-400 hover:border-red-500/50"
            onClick={logout}
          >
            Logout
          </Button>
        </div>
      </aside>

      {/* Main Content */}
      <main className="flex-1 flex flex-col h-screen overflow-hidden">
        {/* Mobile Header */}
        <header className="h-16 border-b border-[#c6a962]/20 bg-black/40 backdrop-blur-md flex items-center justify-between px-6 md:hidden">
          <div className="flex items-center space-x-2">
            {settings?.logo_path && (
              <img 
                src={`http://localhost:8080${settings.logo_path}`} 
                alt="Logo" 
                className="h-8 w-8 object-contain"
              />
            )}
            <span className="text-xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-[#c6a962] to-[#f3e5c0] truncate">
              {shopShortName}
            </span>
          </div>
          <Button variant="ghost" onClick={logout} className="text-gray-400">
            Logout
          </Button>
        </header>

        <div className="flex-1 overflow-y-auto relative">
          {/* Subtle background glow */}
          <div className="absolute top-0 left-1/4 w-96 h-96 bg-[#c6a962]/5 rounded-full blur-[120px] pointer-events-none"></div>
          
          <div className="relative z-10">
            <Outlet />
          </div>
        </div>
      </main>
    </div>
  );
}
