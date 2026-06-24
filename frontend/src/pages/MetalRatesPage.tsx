import { useState } from "react";
import { Button } from "@/components/ui/Button";
import { Input } from "@/components/ui/Input";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/Card";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { getLatestRates, createMetalRate } from "@/lib/metalRates.api";
import { Loader2 } from "lucide-react";
import { toast } from "sonner";

export default function MetalRatesPage() {
  const queryClient = useQueryClient();
  const [formData, setFormData] = useState({ metal_type: "", purity: "", rate_per_gram: "" });

  const { data: rates, isLoading } = useQuery({
    queryKey: ["metalRates"],
    queryFn: getLatestRates,
  });

  const mutation = useMutation({
    mutationFn: createMetalRate,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["metalRates"] });
      toast.success("Metal rate updated successfully");
      setFormData({ metal_type: "", purity: "", rate_per_gram: "" });
    },
    onError: (err: any) => {
      toast.error(err.message || "Failed to update metal rate");
    }
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!formData.metal_type || !formData.purity || !formData.rate_per_gram) return;
    
    const today = new Date().toISOString().split("T")[0];
    
    mutation.mutate({
      metal_type: formData.metal_type,
      purity: formData.purity,
      rate_per_gram: parseFloat(formData.rate_per_gram),
      effective_date: today,
    });
  };

  if (isLoading) {
    return <div className="flex h-[50vh] items-center justify-center"><Loader2 className="h-8 w-8 animate-spin text-[#c6a962]" /></div>;
  }

  return (
    <div className="p-6 max-w-6xl mx-auto space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-3xl font-bold tracking-tight text-[#c6a962]">
          Daily Metal Rates
        </h1>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {rates?.map((rate: any, i: number) => (
          <Card key={i} className="hover:shadow-lg transition-all border-[#c6a962]/20 bg-card backdrop-blur-md">
            <CardHeader className="pb-2">
              <CardTitle className="capitalize text-lg text-foreground">
                {rate.metal_type} - {rate.purity}
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="text-4xl font-bold text-[#c6a962]">
                ₹{rate.rate_per_gram}
                <span className="text-sm font-normal text-muted-foreground ml-1">/ g</span>
              </div>
              <p className="text-sm text-muted-foreground mt-2">
                Effective: {rate.effective_date}
              </p>
            </CardContent>
          </Card>
        ))}
      </div>

      <Card className="mt-8 bg-card backdrop-blur-md border-[#c6a962]/20">
        <CardHeader>
          <CardTitle className="text-foreground">Add / Update Rate</CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="grid grid-cols-1 md:grid-cols-4 gap-4 items-end">
            <div className="space-y-2">
              <label className="text-sm font-medium text-foreground/90">Metal Type</label>
              <Input 
                required
                placeholder="e.g. gold" 
                value={formData.metal_type} 
                onChange={(e) => setFormData({...formData, metal_type: e.target.value})} 
                className="text-foreground border-border bg-muted" 
              />
            </div>
            <div className="space-y-2">
              <label className="text-sm font-medium text-foreground/90">Purity</label>
              <Input 
                required
                placeholder="e.g. 22K" 
                value={formData.purity} 
                onChange={(e) => setFormData({...formData, purity: e.target.value})} 
                className="text-foreground border-border bg-muted" 
              />
            </div>
            <div className="space-y-2">
              <label className="text-sm font-medium text-foreground/90">Rate / Gram (₹)</label>
              <Input 
                required
                type="number" 
                placeholder="0.00" 
                value={formData.rate_per_gram} 
                onChange={(e) => setFormData({...formData, rate_per_gram: e.target.value})} 
                className="text-foreground border-border bg-muted" 
              />
            </div>
            <Button type="submit" disabled={mutation.isPending} className="w-full h-9">
              {mutation.isPending ? "Saving..." : "Save Rate"}
            </Button>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
