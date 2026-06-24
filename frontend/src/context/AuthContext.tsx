import {
  createContext,
  useContext,
  useState,
  useEffect,
  useCallback,
  type ReactNode,
} from "react";
import api from "@/lib/api";
import type { User, Organization, AuthResponse, APIResponse, RegisterRequest } from "@/types";
// ── Context shape ──────────────────────────────────────────────────

interface AuthContextType {
  user: User | null;
  organization: Organization | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (email: string, password: string) => Promise<void>;
  register: (data: RegisterRequest) => Promise<void>;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

// ── Provider ───────────────────────────────────────────────────────

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [organization, setOrganization] = useState<Organization | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  // On mount, restore user and org from localStorage if tokens exist.
  useEffect(() => {
    const storedUser = localStorage.getItem("user");
    const storedOrg = localStorage.getItem("organization");
    const token = localStorage.getItem("access_token");

    if (storedUser && storedOrg && token) {
      try {
        setUser(JSON.parse(storedUser));
        setOrganization(JSON.parse(storedOrg));
      } catch {
        // Corrupted data — clear everything.
        localStorage.removeItem("user");
        localStorage.removeItem("organization");
        localStorage.removeItem("access_token");
        localStorage.removeItem("refresh_token");
      }
    }
    setIsLoading(false);
  }, []);

  const login = useCallback(async (email: string, password: string) => {
    const { data } = await api.post<APIResponse<AuthResponse>>("/auth/login", {
      email,
      password,
    });

    if (!data.success || !data.data) {
      throw new Error(data.error || "Login failed");
    }

    const { access_token, refresh_token, user: userData, organization: orgData } = data.data;

    localStorage.setItem("access_token", access_token);
    localStorage.setItem("refresh_token", refresh_token);
    localStorage.setItem("user", JSON.stringify(userData));
    localStorage.setItem("organization", JSON.stringify(orgData));

    setUser(userData);
    setOrganization(orgData);
  }, []);

  const register = useCallback(async (reqData: RegisterRequest) => {
    const { data } = await api.post<APIResponse<AuthResponse>>("/auth/register", reqData);

    if (!data.success || !data.data) {
      throw new Error(data.error || "Registration failed");
    }

    const { access_token, refresh_token, user: userData, organization: orgData } = data.data;

    localStorage.setItem("access_token", access_token);
    localStorage.setItem("refresh_token", refresh_token);
    localStorage.setItem("user", JSON.stringify(userData));
    localStorage.setItem("organization", JSON.stringify(orgData));

    setUser(userData);
    setOrganization(orgData);
  }, []);

  const logout = useCallback(() => {
    // Fire-and-forget server-side logout.
    api.post("/auth/logout").catch(() => {});

    localStorage.removeItem("access_token");
    localStorage.removeItem("refresh_token");
    localStorage.removeItem("user");
    localStorage.removeItem("organization");
    setUser(null);
    setOrganization(null);
  }, []);

  return (
    <AuthContext.Provider
      value={{
        user,
        organization,
        isAuthenticated: !!user,
        isLoading,
        login,
        register,
        logout,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

// ── Hook ───────────────────────────────────────────────────────────

export function useAuth(): AuthContextType {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
}
