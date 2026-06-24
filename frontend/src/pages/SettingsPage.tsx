import { useState, useEffect, useRef } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/Card";
import { Button } from "@/components/ui/Button";
import { Input } from "@/components/ui/Input";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { getSettings, updateSettings, uploadLogo } from "@/lib/settings.api";
import { Loader2, UploadCloud, Image as ImageIcon, Moon, Sun, Monitor } from "lucide-react";
import { toast } from "sonner";
import { useTheme } from "@/components/ThemeProvider";

export default function SettingsPage() {
  const queryClient = useQueryClient();
  const fileInputRef = useRef<HTMLInputElement>(null);
  const { theme, setTheme } = useTheme();
  
  const [formData, setFormData] = useState({
    shop_name: "",
    gstin: "",
    phone: "",
    address: "",
    invoice_prefix: "INV",
  });

  const { data: settings, isLoading } = useQuery({
    queryKey: ["settings"],
    queryFn: getSettings,
  });

  useEffect(() => {
    if (settings) {
      setFormData({
        shop_name: settings.shop_name,
        gstin: settings.gstin,
        phone: settings.phone,
        address: settings.address,
        invoice_prefix: settings.invoice_prefix,
      });
    }
  }, [settings]);

  const mutation = useMutation({
    mutationFn: updateSettings,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["settings"] });
      toast.success("Settings saved successfully!");
    },
    onError: (err: any) => {
      toast.error(err.message || "Failed to save settings");
    }
  });

  const logoMutation = useMutation({
    mutationFn: uploadLogo,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["settings"] });
      toast.success("Logo uploaded successfully!");
    },
    onError: (err: any) => {
      toast.error(err.message || "Failed to upload logo");
    }
  });

  const handleSave = (e: React.FormEvent) => {
    e.preventDefault();
    mutation.mutate(formData as any);
  };

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files.length > 0) {
      logoMutation.mutate(e.target.files[0]);
    }
  };

  if (isLoading) {
    return <div className="flex h-[50vh] items-center justify-center"><Loader2 className="h-8 w-8 animate-spin text-[#c6a962]" /></div>;
  }

  return (
    <div className="p-6 max-w-4xl mx-auto space-y-6">
      <h1 className="text-3xl font-bold tracking-tight text-[#c6a962]">
        Settings
      </h1>

      <Card className="bg-card backdrop-blur-md border-border shadow-sm">
        <CardHeader>
          <CardTitle className="text-foreground">Appearance</CardTitle>
          <p className="text-sm text-muted-foreground">Customize how the application looks on your device.</p>
        </CardHeader>
        <CardContent>
          <div className="flex flex-col space-y-4">
            <label className="text-sm font-medium text-foreground">Theme Preference</label>
            <div className="flex flex-wrap gap-4">
              <Button 
                variant={theme === "light" ? "default" : "outline"}
                className={theme === "light" ? "bg-[#c6a962] text-white hover:bg-[#b0965a]" : "border-border text-foreground hover:bg-muted"}
                onClick={() => setTheme("light")}
              >
                <Sun className="mr-2 h-4 w-4" />
                Light
              </Button>
              <Button 
                variant={theme === "dark" ? "default" : "outline"}
                className={theme === "dark" ? "bg-[#c6a962] text-white hover:bg-[#b0965a]" : "border-border text-foreground hover:bg-muted"}
                onClick={() => setTheme("dark")}
              >
                <Moon className="mr-2 h-4 w-4" />
                Dark
              </Button>
              <Button 
                variant={theme === "system" ? "default" : "outline"}
                className={theme === "system" ? "bg-[#c6a962] text-white hover:bg-[#b0965a]" : "border-border text-foreground hover:bg-muted"}
                onClick={() => setTheme("system")}
              >
                <Monitor className="mr-2 h-4 w-4" />
                System
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>

      <Card className="bg-card backdrop-blur-md border-[#c6a962]/20 shadow-sm">
        <CardHeader>
          <CardTitle className="text-foreground">Shop Logo</CardTitle>
          <p className="text-sm text-muted-foreground">Upload a logo to display on the dashboard and your PDF invoices.</p>
        </CardHeader>
        <CardContent>
          <div className="flex items-center space-x-6">
            <div className="h-24 w-24 rounded-lg bg-muted border border-border flex items-center justify-center overflow-hidden">
              {settings?.logo_path ? (
                <img 
                  src={`http://localhost:8080${settings.logo_path}`} 
                  alt="Shop Logo" 
                  className="h-full w-full object-contain"
                />
              ) : (
                <ImageIcon className="h-8 w-8 text-muted-foreground" />
              )}
            </div>
            <div className="space-y-2 flex-1">
              <input 
                type="file" 
                accept=".png,.jpg,.jpeg" 
                className="hidden" 
                ref={fileInputRef}
                onChange={handleFileChange}
              />
              <Button 
                variant="outline" 
                onClick={() => fileInputRef.current?.click()}
                disabled={logoMutation.isPending}
                className="border-[#c6a962]/50 text-[#c6a962] hover:bg-[#c6a962]/10"
              >
                {logoMutation.isPending ? (
                  <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                ) : (
                  <UploadCloud className="h-4 w-4 mr-2" />
                )}
                Upload New Logo
              </Button>
              <p className="text-xs text-muted-foreground mt-2">
                Recommended format: PNG, JPG, or JPEG. Max size 2MB.
              </p>
            </div>
          </div>
        </CardContent>
      </Card>

      <Card className="bg-card backdrop-blur-md border-[#c6a962]/20 shadow-sm">
        <CardHeader>
          <CardTitle className="text-foreground">General Details</CardTitle>
          <p className="text-sm text-muted-foreground">Update your shop information and billing details.</p>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSave} className="space-y-6">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div className="space-y-2">
                <label className="text-sm font-medium text-muted-foreground">Shop Name</label>
                <Input 
                  required
                  value={formData.shop_name} 
                  onChange={(e) => setFormData({...formData, shop_name: e.target.value})}
                  className="bg-muted border-border text-foreground" 
                />
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium text-muted-foreground">GSTIN</label>
                <Input 
                  value={formData.gstin} 
                  onChange={(e) => setFormData({...formData, gstin: e.target.value})}
                  className="bg-muted border-border text-foreground" 
                />
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium text-muted-foreground">Phone Number</label>
                <Input 
                  required
                  value={formData.phone} 
                  onChange={(e) => setFormData({...formData, phone: e.target.value})}
                  className="bg-muted border-border text-foreground" 
                />
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium text-muted-foreground">Invoice Prefix</label>
                <Input 
                  required
                  value={formData.invoice_prefix} 
                  onChange={(e) => setFormData({...formData, invoice_prefix: e.target.value})}
                  className="bg-muted border-border text-foreground" 
                />
              </div>
              <div className="md:col-span-2 space-y-2">
                <label className="text-sm font-medium text-muted-foreground">Shop Address</label>
                <Input 
                  required
                  value={formData.address} 
                  onChange={(e) => setFormData({...formData, address: e.target.value})}
                  className="bg-muted border-border text-foreground" 
                />
              </div>
            </div>
            
            <div className="flex justify-end pt-4 border-t border-border mt-6">
              <Button type="submit" disabled={mutation.isPending} className="bg-[#c6a962] hover:bg-[#b0965a] text-black font-semibold">
                {mutation.isPending ? "Saving..." : "Save Changes"}
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
