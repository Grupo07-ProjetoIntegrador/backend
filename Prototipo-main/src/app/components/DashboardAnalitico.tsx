import { useState } from "react";
import { 
  Card, 
  CardContent, 
  CardHeader, 
  CardTitle 
} from "./ui/card";
import { 
  Select, 
  SelectContent, 
  SelectItem, 
  SelectTrigger, 
  SelectValue 
} from "./ui/select";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "./ui/popover";
import { Button } from "./ui/button";
import { Command, CommandEmpty, CommandGroup, CommandInput, CommandItem, CommandList } from "./ui/command";
import { Check, ChevronsUpDown, Calendar as CalendarIcon, TrendingUp } from "lucide-react";
import { cn } from "./ui/utils";

const stores = [
  { value: "all", label: "Todas as Lojas" },
  { value: "zara", label: "Zara" },
  { value: "outback", label: "Outback" },
  { value: "renner", label: "Renner" },
  { value: "cacaushow", label: "Cacau Show" },
];

const dateRanges = [
  { value: "this_month", label: "Este Mês" },
  { value: "last_month", label: "Mês Anterior" },
  { value: "this_year", label: "Este Ano" },
];

export function DashboardAnalitico() {
  const [storeOpen, setStoreOpen] = useState(false);
  const [selectedStore, setSelectedStore] = useState("all");
  const [selectedDate, setSelectedDate] = useState("this_month");

  return (
    <div className="flex-1 p-6 md:p-8 flex flex-col gap-6">
      {/* Page Header & Filters */}
      <div className="flex flex-col md:flex-row justify-between items-start md:items-center gap-4">
        <div>
          <h1 className="text-2xl font-semibold text-gray-900">Dashboard de Indicadores</h1>
          <p className="text-sm text-gray-500 mt-1">Visão geral do desempenho de treinamentos.</p>
        </div>

        <div className="flex flex-col sm:flex-row items-center gap-3 w-full md:w-auto">
          {/* Loja Select (Combobox) */}
          <Popover open={storeOpen} onOpenChange={setStoreOpen}>
            <PopoverTrigger asChild>
              <Button
                variant="outline"
                role="combobox"
                aria-expanded={storeOpen}
                className="w-full sm:w-[220px] justify-between bg-white"
              >
                {selectedStore
                  ? stores.find((store) => store.value === selectedStore)?.label
                  : "Selecione a loja..."}
                <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
              </Button>
            </PopoverTrigger>
            <PopoverContent className="w-[220px] p-0">
              <Command>
                <CommandInput placeholder="Buscar loja..." />
                <CommandList>
                  <CommandEmpty>Nenhuma loja encontrada.</CommandEmpty>
                  <CommandGroup>
                    {stores.map((store) => (
                      <CommandItem
                        key={store.value}
                        value={store.value}
                        onSelect={(currentValue) => {
                          setSelectedStore(currentValue === selectedStore ? "all" : currentValue);
                          setStoreOpen(false);
                        }}
                      >
                        <Check
                          className={cn(
                            "mr-2 h-4 w-4",
                            selectedStore === store.value ? "opacity-100" : "opacity-0"
                          )}
                        />
                        {store.label}
                      </CommandItem>
                    ))}
                  </CommandGroup>
                </CommandList>
              </Command>
            </PopoverContent>
          </Popover>

          {/* Date Picker (Select) */}
          <Select value={selectedDate} onValueChange={setSelectedDate}>
            <SelectTrigger className="w-full sm:w-[180px] bg-white">
              <div className="flex items-center gap-2">
                <CalendarIcon className="w-4 h-4 text-gray-500" />
                <SelectValue placeholder="Período" />
              </div>
            </SelectTrigger>
            <SelectContent>
              {dateRanges.map((range) => (
                <SelectItem key={range.value} value={range.value}>
                  {range.label}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
      </div>

      {/* KPI Grid */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        {/* Card 1: Taxa de Presença */}
        <Card className="border-gray-200 shadow-sm">
          <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
            <CardTitle className="text-sm font-medium text-gray-600">Taxa de Presença</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-gray-900">85%</div>
            <p className="text-xs text-gray-500 mt-1">
              vs. 15% de absenteísmo
            </p>
          </CardContent>
        </Card>

        {/* Card 2: Total de Capacitados */}
        <Card className="border-gray-200 shadow-sm">
          <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
            <CardTitle className="text-sm font-medium text-gray-600">Total de Capacitados</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-gray-900">1.240</div>
            <p className="text-xs text-gray-500 mt-1">
              Funcionários únicos
            </p>
          </CardContent>
        </Card>

        {/* Card 3: Horas de Treinamento */}
        <Card className="border-gray-200 shadow-sm">
          <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
            <CardTitle className="text-sm font-medium text-gray-600">Horas de Treinamento</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-gray-900">450h</div>
            <p className="text-xs text-gray-500 mt-1">
              Investidas no período
            </p>
          </CardContent>
        </Card>

        {/* Card 4: NPS Médio */}
        <Card className="border-gray-200 shadow-sm">
          <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
            <CardTitle className="text-sm font-medium text-gray-600">NPS Médio</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-gray-900 flex items-center gap-2">
              9.2<span className="text-sm font-normal text-gray-500">/10</span>
            </div>
            <div className="flex items-center gap-1 mt-1 text-xs text-green-600 font-medium">
              <TrendingUp className="w-3 h-3" />
              <span>+0.4 em relação ao período anterior</span>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}