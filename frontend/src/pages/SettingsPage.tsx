import { useState, useEffect, useRef } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/Card";
import { Button } from "@/components/ui/Button";
import { Input } from "@/components/ui/Input";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { getSettings, updateSettings, uploadLogo } from "@/lib/settings.api";
import { Loader2, UploadCloud, Image as ImageIcon } from "lucide-react";
import { toast } from "sonner";

export default function SettingsPage() {
  const queryClient = useQueryClient();
  const fileInputRef = useRef<HTMLInputElement>(null);
  
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
        Shop Settings
      </h1>

      <Card className="bg-black/40 backdrop-blur-md border-[#c6a962]/20">
        <CardHeader>
          <CardTitle className="text-gray-200">Shop Logo</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex items-center space-x-6">
            <div className="h-24 w-24 rounded-lg bg-black/50 border border-gray-700 flex items-center justify-center overflow-hidden">
              {settings?.logo_path ? (
                <img 
                  src={`http://localhost:8080${settings.logo_path}`} 
                  alt="Shop Logo" 
                  className="h-full w-full object-contain"
                />
              ) : (
                <ImageIcon className="h-8 w-8 text-gray-500" />
              )}
            </div>
            <div className="space-y-2">
              <p className="text-sm text-gray-400">
                Upload a logo to display on the dashboard and your PDF invoices. <br />
                Recommended format: PNG, JPG, or JPEG.
              </p>
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
            </div>
          </div>
        </CardContent>
      </Card>

      <Card className="bg-black/40 backdrop-blur-md border-[#c6a962]/20">
        <CardHeader>
          <CardTitle className="text-gray-200">General Details</CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSave} className="space-y-6">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div className="space-y-2">
                <label className="text-sm text-gray-400">Shop Name</label>
                <Input 
                  required
                  value={formData.shop_name} 
                  onChange={(e) => setFormData({...formData, shop_name: e.target.value})}
                  className="bg-black/50 border-gray-700 text-gray-200" 
                />
              </div>
              <div className="space-y-2">
                <label className="text-sm text-gray-400">GSTIN</label>
                <Input 
                  value={formData.gstin} 
                  onChange={(e) => setFormData({...formData, gstin: e.target.value})}
                  className="bg-black/50 border-gray-700 text-gray-200" 
                />
              </div>
              <div className="space-y-2">
                <label className="text-sm text-gray-400">Phone Number</label>
                <Input 
                  required
                  value={formData.phone} 
                  onChange={(e) => setFormData({...formData, phone: e.target.value})}
                  className="bg-black/50 border-gray-700 text-gray-200" 
                />
              </div>
              <div className="space-y-2">
                <label className="text-sm text-gray-400">Invoice Prefix</label>
                <Input 
                  required
                  value={formData.invoice_prefix} 
                  onChange={(e) => setFormData({...formData, invoice_prefix: e.target.value})}
                  className="bg-black/50 border-gray-700 text-gray-200" 
                />
              </div>
              <div className="md:col-span-2 space-y-2">
                <label className="text-sm text-gray-400">Shop Address</label>
                <Input 
                  required
                  value={formData.address} 
                  onChange={(e) => setFormData({...formData, address: e.target.value})}
                  className="bg-black/50 border-gray-700 text-gray-200" 
                />
              </div>
            </div>
            
            <div className="flex justify-end">
              <Button type="submit" disabled={mutation.isPending}>
                {mutation.isPending ? "Saving..." : "Save Changes"}
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
