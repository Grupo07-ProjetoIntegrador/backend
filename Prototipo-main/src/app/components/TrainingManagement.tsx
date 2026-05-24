import { useEffect, useState, useRef } from "react";
import { Plus, ExternalLink, Calendar, BookOpen, Tag, Search, ChevronDown, Filter, ChevronRight, Edit2, Trash2, Users, Clock, Star, Award, AlertTriangle, Store, UserCheck, Activity, Hash, List, ChevronLeft, CheckCircle2, XCircle } from "lucide-react";
import { TrainingForm } from "./TrainingForm";
import { TrainingDetails } from "./TrainingDetails";
import { StoreDetails } from "./StoreDetails";

import { Popover, PopoverContent, PopoverTrigger, PopoverAnchor } from "./ui/popover";
import { Button } from "./ui/button";
import { Input } from "./ui/input";
import { Calendar as CalendarComponent } from "./ui/calendar";
import { Tabs, TabsList, TabsTrigger, TabsContent } from "./ui/tabs";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "./ui/select";
import { Card, CardHeader, CardTitle, CardContent, CardDescription } from "./ui/card";
import { Progress } from "./ui/progress";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription } from "./ui/dialog";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "./ui/table";
import { Badge } from "./ui/badge";
import { format, startOfMonth, endOfMonth, startOfWeek, endOfWeek, eachDayOfInterval, isSameMonth, isSameDay, addMonths, subMonths } from "date-fns";
import { cn } from "./ui/utils";
import { ptBR } from "date-fns/locale";
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, Legend } from "recharts";

const API_BASE_URL = "http://localhost:8080";

type ApiTraining = {
  id: string;
  tema: string;
  segmento?: string;
  data?: string;
  data_hora?: string;
  horario_inicio?: string;
  conteudo?: string;
  status?: string;
};

const formatTrainingDate = (dateTime?: string, fallback = "") => {
  if (!dateTime) return fallback;
  const date = new Date(dateTime);
  if (isNaN(date.getTime())) return fallback || dateTime;
  return format(date, "dd MMM yyyy", { locale: ptBR });
};

const normalizeTrainingFromApi = (training: ApiTraining) => {
  const dataHora = training.data_hora || "";
  const status = (training.status || "").toLowerCase();

  return {
    id: training.id,
    tema: training.tema,
    segmento: training.segmento || "Geral",
    data: formatTrainingDate(dataHora, training.data || ""),
    hora: training.horario_inicio || (dataHora ? format(new Date(dataHora), "HH:mm") : "00:00"),
    dataHora: dataHora || new Date().toISOString(),
    conteudo: training.conteudo || "",
    isCancelado: status === "cancelado",
    attendanceList: [],
  };
};

const chartData = [
  { name: "Jan", inscritos: 120, presentes: 95 },
  { name: "Fev", inscritos: 145, presentes: 110 },
  { name: "Mar", inscritos: 130, presentes: 125 },
  { name: "Abr", inscritos: 160, presentes: 150 },
  { name: "Mai", inscritos: 150, presentes: 140 },
  { name: "Jun", inscritos: 175, presentes: 165 },
];

const mockStoreData = [
  { id: 1, name: "Renner", segment: "Vestuário", manager: "Carlos Silva", luc: "LUC-101" },
  { id: 2, name: "C&A", segment: "Vestuário", manager: "Ana Oliveira", luc: "LUC-102" },
  { id: 3, name: "Centauro", segment: "Artigos Esportivos", manager: "Pedro Santos", luc: "LUC-105" },
  { id: 4, name: "O Boticário", segment: "Cosméticos", manager: "Julia Lima", luc: "LUC-104" },
  { id: 5, name: "Riachuelo", segment: "Vestuário", manager: "Roberto Almeida", luc: "LUC-103" },
  { id: 6, name: "Loja de Informática XYZ", segment: "Tecnologia", manager: "Fernanda Costa", luc: "LUC-106" },
  { id: 7, name: "Moda Íntima ABC", segment: "Vestuário", manager: "Marcos Pereira", luc: "LUC-107" },
  { id: 8, name: "Acessórios e Cia", segment: "Acessórios", manager: "Luciana Gomes", luc: "LUC-108" },
  { id: 9, name: "Ótica Elegance", segment: "Óticas", manager: "Beatriz Nogueira", luc: "LUC-109" },
];

const generateVariedAttendance = (trainingId: number, tema: string) => {
  return mockStoreData.map((store) => {
    let status = "Presente";
    
    if (store.name === "Riachuelo" && tema === "Inovações Tecnológicas em Vendas") {
      status = "Ausente";
    } else if (store.name === "Centauro" && tema === "Inovações Tecnológicas em Vendas") {
      status = "Ausente"; // Just another example to ensure variety
    } else if ((trainingId + store.id) % 4 === 0) {
      status = "Ausente";
    } else if ((trainingId * store.id) % 7 === 0) {
      status = "Ausente";
    }

    return {
      id: store.id,
      luc: store.luc,
      loja: store.name,
      representante: store.manager,
      status: status
    };
  });
};

const initialTrainings = [
  { id: 1, tema: "Técnicas de Vendas e Retenção", segmento: "Geral", data: "10 Jan 2026", hora: "14:00", dataHora: "2026-01-10T14:00:00", conteudo: "Atendimento ao cliente", isCancelado: false },
  { id: 2, tema: "Gestão de Estoque e Visual Merchandising", segmento: "Vestuário", data: "25 Jan 2026", hora: "09:30", dataHora: "2026-01-25T09:30:00", conteudo: "Organização e exposição de produtos", isCancelado: false },
  { id: 3, tema: "Compliance e Boas Práticas Comerciais", segmento: "Geral", data: "10 Fev 2026", hora: "16:00", dataHora: "2026-02-10T16:00:00", conteudo: "Normas internas e regulamentações", isCancelado: false },
  { id: 4, tema: "Liderança e Gestão de Equipes", segmento: "Geral", data: "20 Fev 2026", hora: "10:00", dataHora: "2026-02-20T10:00:00", conteudo: "Desenvolvimento de líderes", isCancelado: false },
  { id: 5, tema: "Experiência do Cliente no Varejo", segmento: "Geral", data: "05 Mar 2026", hora: "15:00", dataHora: "2026-03-05T15:00:00", conteudo: "Jornada de compra e satisfação", isCancelado: false },
  { id: 6, tema: "Marketing Digital para Lojistas", segmento: "Geral", data: "15 Mar 2026", hora: "09:00", dataHora: "2026-03-15T09:00:00", conteudo: "Estratégias para redes sociais", isCancelado: false },
  { id: 7, tema: "Prevenção de Perdas no Shopping", segmento: "Geral", data: "02 Abr 2026", hora: "14:30", dataHora: "2026-04-02T14:30:00", conteudo: "Identificação e mitigação de riscos", isCancelado: false },
  { id: 8, tema: "Inovações Tecnológicas em Vendas", segmento: "Tecnologia e Eletro", data: "18 Abr 2026", hora: "11:00", dataHora: "2026-04-18T11:00:00", conteudo: "Uso de IA e novas ferramentas", isCancelado: false },
  { id: 9, tema: "Gestão Financeira para Pequenos Negócios", segmento: "Geral", data: "25 Abr 2026", hora: "13:00", dataHora: "2026-04-25T13:00:00", conteudo: "Controle de fluxo de caixa", isCancelado: false },
  { id: 10, tema: "Design de Interiores Comerciais", segmento: "Casa e Decoração", data: "02 Mai 2026", hora: "10:30", dataHora: "2026-05-02T10:30:00", conteudo: "Otimização de layout de loja", isCancelado: false },
  { id: 11, tema: "Estratégias de Comunicação Interna", segmento: "Geral", data: "20 Mai 2026", hora: "15:30", dataHora: "2026-05-20T15:30:00", conteudo: "Melhoria no engajamento da equipe", isCancelado: false },
  { id: 12, tema: "Sustentabilidade no Varejo", segmento: "Geral", data: "10 Jun 2026", hora: "09:30", dataHora: "2026-06-10T09:30:00", conteudo: "Práticas ESG em lojas", isCancelado: false },
].map(t => {
  const isAgendado = new Date(t.dataHora) > new Date(2026, 4, 4);
  return {
    ...t,
    attendanceList: isAgendado ? [] : generateVariedAttendance(t.id, t.tema)
  };
});

export function TrainingManagement() {
  const [trainingsList, setTrainingsList] = useState<any[]>([]);
  const [isLoadingTrainings, setIsLoadingTrainings] = useState(true);
  const [trainingsError, setTrainingsError] = useState("");
  const [activeTab, setActiveTab] = useState("list");
  const [viewMode, setViewMode] = useState<"list" | "calendar">("list");
  const [calendarMonth, setCalendarMonth] = useState(new Date(2026, 4, 1)); // May 2026
  const [activeFilter, setActiveFilter] = useState<"todos" | "tema" | "segmento" | "data" | "conteudo" | "status">("todos");
  const [searchQuery, setSearchQuery] = useState("");
  const [dateQuery, setDateQuery] = useState<Date>();
  const [dateInput, setDateInput] = useState("");
  const [isCreating, setIsCreating] = useState(false);
  const [editingTraining, setEditingTraining] = useState<any>(null);
  const [selectedTraining, setSelectedTraining] = useState<any>(null);
  const [trainingToDelete, setTrainingToDelete] = useState<any>(null);
  const [selectedStore, setSelectedStore] = useState<any>(null);
  const [storeSearchQuery, setStoreSearchQuery] = useState("");
  const [lucSearchQuery, setLucSearchQuery] = useState("");
  const [currentPage, setCurrentPage] = useState(1);
  const [dashboardPeriod, setDashboardPeriod] = useState("ano");
  const itemsPerPage = 10;
  const dateInputRef = useRef<HTMLInputElement>(null);

  const carregarTreinamentos = async () => {
    setIsLoadingTrainings(true);
    setTrainingsError("");

    try {
      const response = await fetch(`${API_BASE_URL}/api/treinamentos`);

      if (!response.ok) {
        throw new Error("Erro ao buscar treinamentos");
      }

      const data: ApiTraining[] = await response.json();
      setTrainingsList((data || []).map(normalizeTrainingFromApi));
    } catch (error) {
      console.error("Erro ao carregar treinamentos:", error);
      setTrainingsError("Erro ao conectar com a API. Confira se o backend Go esta rodando.");
      setTrainingsList([]);
    } finally {
      setIsLoadingTrainings(false);
    }
  };

  useEffect(() => {
    carregarTreinamentos();
  }, []);

  // Dynamic KPI and Data Calculation
  const now = new Date(2026, 4, 4); // May 4, 2026
  let startDate = new Date(2000, 0, 1);
  let endDate = new Date(2100, 0, 1);

  if (dashboardPeriod === "mes") {
    startDate = new Date(now.getFullYear(), now.getMonth(), 1);
    endDate = new Date(now.getFullYear(), now.getMonth() + 1, 0, 23, 59, 59);
  } else if (dashboardPeriod === "trimestre") {
    startDate = new Date(now.getFullYear(), now.getMonth() - 2, 1);
    endDate = new Date(now.getFullYear(), now.getMonth() + 1, 0, 23, 59, 59);
  } else if (dashboardPeriod === "ano") {
    startDate = new Date(now.getFullYear(), 0, 1);
    endDate = new Date(now.getFullYear(), 11, 31, 23, 59, 59);
  }

  const pastTrainings = trainingsList.filter(t => {
    const d = new Date(t.dataHora);
    return d < now && !t.isCancelado;
  });

  const periodTrainings = pastTrainings.filter(t => {
    const d = new Date(t.dataHora);
    return d >= startDate && d <= endDate;
  });

  const periodTrainingIds = new Set(periodTrainings.map(t => t.id));

  let totalInscritos = 0;
  let totalPresentes = 0;
  const storesImpacted = new Set();

  periodTrainings.forEach((t) => {
    if (t.attendanceList) {
      t.attendanceList.forEach((record: any) => {
        totalInscritos++;
        if (record.status === "Presente") {
          totalPresentes++;
          storesImpacted.add(record.id);
        }
      });
    }
  });

  const taxaPresenca = totalInscritos > 0 ? ((totalPresentes / totalInscritos) * 100).toFixed(1) : "0.0";
  const lojasImpactadas = storesImpacted.size;

  const monthsMap: Record<number, string> = { 0: "Jan", 1: "Fev", 2: "Mar", 3: "Abr", 4: "Mai", 5: "Jun", 6: "Jul", 7: "Ago", 8: "Set", 9: "Out", 10: "Nov", 11: "Dez" };
  const chartDataDynamic = [];
  const startYearMonth = startDate.getFullYear() * 12 + startDate.getMonth();
  const endYearMonth = endDate.getFullYear() * 12 + (dashboardPeriod === "ano" ? now.getMonth() : endDate.getMonth());

  for (let ym = startYearMonth; ym <= endYearMonth; ym++) {
    const y = Math.floor(ym / 12);
    const m = ym % 12;

    const monthTrainings = periodTrainings.filter(t => {
      const d = new Date(t.dataHora);
      return d.getMonth() === m && d.getFullYear() === y;
    });
    
    let mInscritos = 0;
    let mPresentes = 0;
    
    monthTrainings.forEach((t) => {
      if (t.attendanceList) {
        t.attendanceList.forEach((record: any) => {
          mInscritos++;
          if (record.status === "Presente") mPresentes++;
        });
      }
    });
    
    if (dashboardPeriod === "ano" || mInscritos > 0 || dashboardPeriod !== "ano") {
       chartDataDynamic.push({
         name: monthsMap[m],
         inscritos: mInscritos,
         presentes: mPresentes,
         monthIndex: m,
         year: y
       });
    }
  }

  const storeExplorerData = mockStoreData.map(store => {
    const periodAttendance = periodTrainings.map(t => {
      if (t.attendanceList) {
        return t.attendanceList.find((a: any) => a.id === store.id);
      }
      return null;
    }).filter(Boolean) as any[];
    
    const total = periodAttendance.length;
    const presentes = periodAttendance.filter(a => a.status === "Presente").length;
    const participation = total > 0 ? Math.round((presentes / total) * 100) : 0;
    
    return {
      ...store,
      totalTrainings: total,
      participation
    };
  });

  const top5Lojas = [...storeExplorerData]
    .sort((a, b) => b.participation - a.participation)
    .slice(0, 5);

  const riskStores = storeExplorerData.map(store => {
    const sortedAttendance = periodTrainings
      .map(t => {
        const record = t.attendanceList?.find((a: any) => a.id === store.id);
        if (record) {
          return { trainingId: t.id, status: record.status, dataHora: t.dataHora };
        }
        return null;
      })
      .filter(Boolean)
      .sort((a: any, b: any) => new Date(b.dataHora).getTime() - new Date(a.dataHora).getTime());
      
    let totalFaltas = 0;
    let ultimaPresenca = "Nunca";
    
    for (const record of sortedAttendance) {
      if (record.status === "Ausente") {
        totalFaltas++;
      }
    }
    
    const lastPresent = sortedAttendance.find((r: any) => r.status === "Presente");
    if (lastPresent) {
      const date = new Date(lastPresent.dataHora);
      ultimaPresenca = `${monthsMap[date.getMonth()]} de ${date.getFullYear()}`;
    }
    
    return {
      ...store,
      faltas: totalFaltas,
      ultimaPresenca
    };
  }).filter(s => s.faltas > 0).sort((a, b) => b.faltas - a.faltas).slice(0, 4);

  const handleDateInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    let value = e.target.value;
    
    if (
      (e.nativeEvent as any).inputType === "deleteContentBackward" &&
      dateInput.endsWith("/") &&
      value === dateInput.slice(0, -1)
    ) {
      value = value.slice(0, -1);
    }

    let digits = value.replace(/\D/g, "");
    if (digits.length > 8) digits = digits.slice(0, 8);
    
    let formatted = digits;
    if (digits.length > 4) {
      formatted = `${digits.slice(0, 2)}/${digits.slice(2, 4)}/${digits.slice(4)}`;
    } else if (digits.length > 2) {
      formatted = `${digits.slice(0, 2)}/${digits.slice(2)}`;
    }
    setDateInput(formatted);

    if (formatted === "") {
      setDateQuery(undefined);
    } else if (formatted.length === 10) {
      const [d, m, y] = formatted.split("/");
      const parsedDate = new Date(parseInt(y), parseInt(m) - 1, parseInt(d));
      if (!isNaN(parsedDate.getTime())) {
        setDateQuery(parsedDate);
      }
    }
  };

  const handleDelete = () => {
    if (trainingToDelete) {
      setTrainingsList((prev) => prev.filter((t) => t.id !== trainingToDelete.id));
      setTrainingToDelete(null);
    }
  };

  const handleSuccess = async () => {
    setIsCreating(false);
    setEditingTraining(null);
    await carregarTreinamentos();
  };

  const filterOptions: { key: "todos" | "tema" | "segmento" | "data" | "conteudo" | "status"; label: string }[] = [
    { key: "todos", label: "Todos" },
    { key: "tema", label: "Tema" },
    { key: "segmento", label: "Segmento" },
    { key: "data", label: "Data" },
    { key: "conteudo", label: "Conteúdo" },
    { key: "status", label: "Status" },
  ];

  if (selectedTraining) {
    return (
      <TrainingDetails 
        training={selectedTraining} 
        onBack={() => setSelectedTraining(null)} 
        onUpdateAttendance={(id, list) => {
          setTrainingsList(prev => prev.map(t => 
            t.id === id ? { ...t, attendanceList: list } : t
          ));
          setSelectedTraining((prev: any) => ({ ...prev, attendanceList: list }));
        }}
      />
    );
  }

  const handlePrevMonth = () => setCalendarMonth(subMonths(calendarMonth, 1));
  const handleNextMonth = () => setCalendarMonth(addMonths(calendarMonth, 1));

  const renderCalendarView = (filteredTrainings: typeof trainingsList) => {
    const monthStart = startOfMonth(calendarMonth);
    const monthEnd = endOfMonth(monthStart);
    const startDate = startOfWeek(monthStart);
    const endDate = endOfWeek(monthEnd);

    const dateFormat = "d";
    const days = eachDayOfInterval({ start: startDate, end: endDate });
    const weekDays = ["Dom", "Seg", "Ter", "Qua", "Qui", "Sex", "Sáb"];

    return (
      <div className="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden flex flex-col mt-4">
        <div className="flex items-center justify-between px-6 py-4 border-b border-gray-200">
          <h2 className="text-lg font-semibold text-gray-900 capitalize">
            {format(calendarMonth, "MMMM yyyy", { locale: ptBR })}
          </h2>
          <div className="flex items-center gap-2">
            <Button variant="outline" size="icon" onClick={handlePrevMonth} className="h-8 w-8 text-gray-600">
              <ChevronLeft className="w-4 h-4" />
            </Button>
            <Button variant="outline" size="icon" onClick={handleNextMonth} className="h-8 w-8 text-gray-600">
              <ChevronRight className="w-4 h-4" />
            </Button>
          </div>
        </div>
        <div className="grid grid-cols-7 border-b border-gray-200 bg-gray-50">
          {weekDays.map(day => (
            <div key={day} className="py-2 text-center text-xs font-semibold text-gray-500 uppercase">
              {day}
            </div>
          ))}
        </div>
        <div className="grid grid-cols-7 auto-rows-[120px] bg-gray-200 gap-px">
          {days.map(day => {
            const isSameMonthDay = isSameMonth(day, monthStart);
            const isToday = isSameDay(day, now);
            const dayTrainings = filteredTrainings.filter(t => isSameDay(new Date(t.dataHora), day));

            return (
              <div
                key={day.toString()}
                className={cn(
                  "bg-white p-2 flex flex-col gap-1 overflow-y-auto",
                  !isSameMonthDay && "bg-gray-50 text-gray-400"
                )}
              >
                <div className="flex justify-end mb-1">
                  <span
                    className={cn(
                      "text-sm font-medium w-6 h-6 flex items-center justify-center rounded-full",
                      isToday ? "bg-[#8B1A1A] text-white" : "text-gray-700"
                    )}
                  >
                    {format(day, dateFormat)}
                  </span>
                </div>
                {dayTrainings.map(t => {
                  const isAgendado = new Date(t.dataHora) > now;
                  const pillClass = t.isCancelado
                    ? "bg-red-50 text-red-700 border-red-200 hover:bg-red-100"
                    : isAgendado
                    ? "bg-amber-50 text-amber-700 border-amber-200 hover:bg-amber-100"
                    : "bg-emerald-50 text-emerald-700 border-emerald-200 hover:bg-emerald-100";
                  
                  const Icon = t.isCancelado ? XCircle : (isAgendado ? Clock : CheckCircle2);

                  return (
                    <button
                      key={t.id}
                      onClick={() => setSelectedTraining(t)}
                      className={cn(
                        "text-left px-2 py-1 text-xs rounded-md border flex items-center gap-1.5 w-full truncate transition-colors",
                        pillClass
                      )}
                      title={t.tema}
                    >
                      <Icon className="w-3 h-3 flex-shrink-0" />
                      <span className="truncate">{t.tema}</span>
                    </button>
                  );
                })}
              </div>
            );
          })}
        </div>
      </div>
    );
  };

  if (selectedStore) {
    return (
      <StoreDetails
        store={selectedStore}
        trainings={periodTrainings}
        onBack={() => {
          setSelectedStore(null);
          setActiveTab("dashboard");
        }}
        onSelectTraining={(training) => {
          setSelectedTraining(training);
        }}
      />
    );
  }

  if (isCreating || editingTraining) {
    const formData = editingTraining ? {
      ...editingTraining,
      data: editingTraining.dataHora ? editingTraining.dataHora.split('T')[0] : "",
      horarioInicio: editingTraining.hora || "",
      segmento: editingTraining.segmento || "Geral (Todos os Segmentos)",
    } : undefined;

    return (
      <TrainingForm
        initialData={formData}
        onBack={() => {
          setIsCreating(false);
          setEditingTraining(null);
        }}
        onSuccess={handleSuccess}
      />
    );
  }

  const filteredTrainings = trainingsList.filter((row) => {
    if (activeFilter === "data") {
      if (!dateInput) return true;
      const rowDate = new Date(row.dataHora);
      if (isNaN(rowDate.getTime())) return true;
      const rowDateStr = format(rowDate, "dd/MM/yyyy");
      return rowDateStr.includes(dateInput);
    }

    if (!searchQuery) return true;
    
    const query = searchQuery.toLowerCase();
    
    const dataTreinamento = new Date(row.dataHora);
    const agora = new Date();
    let status = "agendado";
    if (row.isCancelado) status = "cancelado";
    else if (dataTreinamento < agora) status = "concluido";

    const matchTema = row.tema.toLowerCase().includes(query);
    const matchSegmento = (row.segmento || "Geral (Todos os Segmentos)").toLowerCase().includes(query);
    const matchConteudo = row.conteudo.toLowerCase().includes(query);
    const matchStatus = status.includes(query);
    const matchData = row.data.toLowerCase().includes(query) || row.dataHora.includes(query);

    if (activeFilter === "todos") return matchTema || matchSegmento || matchConteudo || matchStatus || matchData;
    if (activeFilter === "tema") return matchTema;
    if (activeFilter === "segmento") return matchSegmento;
    if (activeFilter === "conteudo") return matchConteudo;
    if (activeFilter === "status") return matchStatus;
    
    return true;
  });

  const totalPages = Math.ceil(filteredTrainings.length / itemsPerPage);
  const currentItems = filteredTrainings.slice(
    (currentPage - 1) * itemsPerPage,
    currentPage * itemsPerPage
  );

  return (
    <main className="flex-1 min-h-screen flex flex-col" style={{ backgroundColor: "#F7F4EF" }}>
      <Tabs value={activeTab} onValueChange={setActiveTab} className="flex-1 flex flex-col w-full h-full">
        {/* Header */}
        <div className="pt-6 px-4 md:px-8 bg-transparent relative z-20">
          <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-4 relative z-20">
            <h1 className="text-xl md:text-[1.5rem] font-semibold" style={{ color: "#1F2937" }}>
              Gestão de Treinamentos
            </h1>
            <div className="flex items-center w-full md:w-auto self-start md:self-auto relative z-30 pointer-events-auto">
              <button
                onClick={() => setIsCreating(true)}
                className="flex items-center justify-center gap-2 px-4 py-2 rounded-lg text-sm font-medium text-white transition-colors hover:opacity-90 w-full md:w-auto shrink-0 min-h-[44px] cursor-pointer relative z-30"
                style={{ backgroundColor: "#D93030" }}
              >
                <Plus className="w-4 h-4" strokeWidth={2} />
                Novo Treinamento
              </button>
            </div>
          </div>
          
          <TabsList className="bg-white border border-gray-200 shadow-sm p-1 rounded-md inline-flex w-full sm:w-auto h-auto">
            <TabsTrigger 
              value="list" 
              className="px-4 py-2 text-sm font-medium text-gray-500 bg-transparent rounded-sm hover:text-red-600 transition-all data-[state=active]:bg-red-50 data-[state=active]:text-red-600 data-[state=active]:shadow-none flex-1 sm:flex-none"
            >
              Agenda e Lista
            </TabsTrigger>
            <TabsTrigger 
              value="dashboard" 
              className="px-4 py-2 text-sm font-medium text-gray-500 bg-transparent rounded-sm hover:text-red-600 transition-all data-[state=active]:bg-red-50 data-[state=active]:text-red-600 data-[state=active]:shadow-none flex-1 sm:flex-none"
            >
              Engajamento das Lojas
            </TabsTrigger>
          </TabsList>
        </div>

        <TabsContent value="list" className="flex-1 flex flex-col m-0 p-0 outline-none data-[state=active]:flex pt-4">
          {/* Action Buttons */}
          <div className="px-4 md:px-8 mb-4 flex flex-col md:flex-row items-center gap-3 w-full justify-between">
            <div className="flex items-center bg-gray-100 p-1 rounded-lg border border-gray-200">
              <button
                onClick={() => setViewMode("list")}
                className={cn(
                  "px-3 py-1.5 text-sm font-medium rounded-md flex items-center gap-2 transition-colors",
                  viewMode === "list" 
                    ? "bg-white shadow-sm text-gray-900" 
                    : "text-gray-500 hover:text-gray-700"
                )}
              >
                <List className="w-4 h-4" />
                Lista
              </button>
              <button
                onClick={() => setViewMode("calendar")}
                className={cn(
                  "px-3 py-1.5 text-sm font-medium rounded-md flex items-center gap-2 transition-colors",
                  viewMode === "calendar" 
                    ? "bg-white shadow-sm text-gray-900" 
                    : "text-gray-500 hover:text-gray-700"
                )}
              >
                <Calendar className="w-4 h-4" />
                Calendário
              </button>
            </div>
            
            {/* Search and Filter */}
            <div className="flex items-center gap-3 w-full md:w-auto">
            <div className="relative flex items-center bg-white border border-gray-200 rounded-lg hover:border-gray-300 focus-within:ring-2 focus-within:ring-[#D93030] focus-within:border-transparent transition-colors group w-full md:w-auto">
              <div className="relative flex items-center flex-1 min-w-[120px]">
                {activeFilter === "data" ? (
                  <Popover>
                    <PopoverAnchor asChild>
                      <div className="flex items-center w-full md:w-48 lg:w-56 pr-2">
                        <Input
                          type="text"
                          value={dateInput}
                          onChange={handleDateInputChange}
                          placeholder="DD/MM/AAAA"
                          className="pl-3 pr-2 py-2 text-sm focus:outline-none focus-visible:ring-0 focus-visible:ring-offset-0 border-0 shadow-none bg-transparent w-full min-h-[44px] text-gray-900"
                          maxLength={10}
                        />
                        <PopoverTrigger asChild>
                          <Button variant="ghost" size="icon" className="h-8 w-8 text-gray-400 hover:text-[#D93030] shrink-0 hover:bg-transparent">
                            <Calendar className="h-4 w-4" />
                          </Button>
                        </PopoverTrigger>
                      </div>
                    </PopoverAnchor>
                    <PopoverContent align="start" sideOffset={4} className="w-auto p-0 bg-white shadow-md z-50">
                      <CalendarComponent
                          mode="single"
                          selected={dateQuery}
                          onSelect={(date) => {
                            setDateQuery(date);
                            if (date) {
                              setDateInput(format(date, "dd/MM/yyyy", { locale: ptBR }));
                            } else {
                              setDateInput("");
                            }
                          }}
                          initialFocus
                        />
                      </PopoverContent>
                  </Popover>
                ) : (
                  <>
                    <Search className="w-4 h-4 text-gray-400 absolute left-3 pointer-events-none group-focus-within:text-[#D93030] transition-colors" />
                    <input
                      type="text"
                      value={searchQuery}
                      onChange={(e) => setSearchQuery(e.target.value)}
                      placeholder="Buscar..."
                      className="pl-9 pr-3 py-2 text-sm focus:outline-none bg-transparent w-full md:w-48 lg:w-56 min-h-[44px]"
                    />
                  </>
                )}
              </div>

              <div className="h-5 w-px bg-gray-200 shrink-0" />

              <div className="relative flex items-center shrink-0">
                <Filter className="w-3.5 h-3.5 text-gray-400 absolute left-3 pointer-events-none group-focus-within:text-[#D93030] transition-colors hidden md:block" />
                <select
                  value={activeFilter}
                  onChange={(e) => {
                    setActiveFilter(e.target.value as any);
                    setSearchQuery("");
                    setDateQuery(undefined);
                    setDateInput("");
                  }}
                  className="appearance-none bg-transparent text-gray-600 text-sm focus:outline-none cursor-pointer w-auto pl-3 md:pl-8 pr-8 py-2 hover:text-gray-900 transition-colors min-h-[44px]"
                >
                  {filterOptions.map((opt) => (
                    <option key={opt.key} value={opt.key}>
                      {opt.label}
                    </option>
                  ))}
                </select>
                <ChevronDown className="w-3.5 h-3.5 text-gray-400 absolute right-3 pointer-events-none group-focus-within:text-[#D93030] transition-colors" />
              </div>
            </div>
            </div>
          </div>

          {/* Content Card */}
          <div className="px-4 md:px-8 pb-8 flex-1 flex flex-col">
            {viewMode === "calendar" ? (
              renderCalendarView(filteredTrainings)
            ) : (
            <div className="bg-white rounded-lg shadow-sm border border-gray-200 flex-1 flex flex-col">
              {/* Card Header */}
              <div className="px-4 md:px-6 py-4 border-b border-gray-200 flex items-center justify-between">
                <h2 className="text-gray-800" style={{ fontSize: "0.9375rem", fontWeight: 600 }}>Próximos e Recentes</h2>
                <span className="text-xs text-gray-500 bg-gray-50 px-2.5 py-1 rounded-full border border-gray-200">
                  {filteredTrainings.length} eventos
                </span>
              </div>

              {/* Table (Desktop) & Cards (Mobile) */}
              <div className="flex-1">
                {/* Desktop Table */}
                <div className="hidden md:block w-full overflow-x-auto">
                  <table className="w-full">
                    <thead className="bg-gray-50 border-b border-gray-200">
                      <tr>
                        <th className="text-center px-6 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wider">
                          <div className="flex items-center justify-center gap-1.5">
                            <Tag className="w-3 h-3" />
                            Tema
                          </div>
                        </th>
                        <th className="text-center px-6 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wider">
                          <div className="flex items-center justify-center gap-1.5">
                            <Filter className="w-3 h-3" />
                            Segmento
                          </div>
                        </th>
                        <th className="text-center px-6 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wider">
                          <div className="flex items-center justify-center gap-1.5">
                            <Calendar className="w-3 h-3" />
                            Data
                          </div>
                        </th>
                        <th className="text-center px-6 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wider">
                          <div className="flex items-center justify-center gap-1.5">
                            <BookOpen className="w-3 h-3" />
                            Conteúdo
                          </div>
                        </th>
                        <th className="text-center px-6 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wider">
                          Status
                        </th>
                        <th className="text-center px-6 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wider">
                          Ações
                        </th>
                      </tr>
                    </thead>
                    <tbody>
                      {(isLoadingTrainings || trainingsError || currentItems.length === 0) && (
                        <tr>
                          <td colSpan={6} className="px-6 py-8 text-center text-sm text-gray-500">
                            {isLoadingTrainings
                              ? "Carregando treinamentos do banco de dados..."
                              : trainingsError || "Nenhum treinamento encontrado."}
                          </td>
                        </tr>
                      )}
                      {currentItems.map((row, idx) => {
                        const dataTreinamento = new Date(row.dataHora);
                        const agora = new Date();
                        
                        let status = "agendado";
                        if (row.isCancelado) {
                          status = "cancelado";
                        } else if (dataTreinamento < agora) {
                          status = "concluido";
                        }

                        return (
                        <tr
                          key={row.id}
                          className={`border-b border-gray-100 hover:bg-gray-50 transition-colors ${
                            idx === currentItems.length - 1 ? "border-b-0" : ""
                          }`}
                        >
                          <td 
                            className="px-6 py-4 cursor-pointer text-center"
                            onClick={() => setSelectedTraining(row)}
                          >
                            <div className="flex items-center justify-center gap-3">
                              <span className="text-sm font-medium text-gray-800">
                                {row.tema}
                              </span>
                            </div>
                          </td>
                          <td 
                            className="px-6 py-4 cursor-pointer text-center"
                            onClick={() => setSelectedTraining(row)}
                          >
                            <span className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-gray-100 text-gray-800">
                              {row.segmento || "Geral (Todos os Segmentos)"}
                            </span>
                          </td>
                          <td 
                            className="px-6 py-4 cursor-pointer text-center"
                            onClick={() => setSelectedTraining(row)}
                          >
                            <span className="inline-flex items-center justify-center gap-1.5 text-sm text-gray-600 bg-gray-50 px-2.5 py-1 rounded-md border border-gray-200 hover:bg-gray-100 transition-colors">
                              <Calendar className="w-3 h-3 text-gray-400" />
                              {row.data} às {row.hora}
                            </span>
                          </td>
                          <td 
                            className="px-6 py-4 cursor-pointer text-center"
                            onClick={() => setSelectedTraining(row)}
                          >
                            <span className="text-sm text-gray-500 hover:text-gray-700 transition-colors line-clamp-2 md:line-clamp-none">{row.conteudo}</span>
                          </td>
                          <td 
                            className="px-6 py-4 text-center cursor-pointer"
                            onClick={() => setSelectedTraining(row)}
                          >
                            <div className="flex justify-center items-center">
                              {status === "concluido" && (
                                <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-50 text-green-700 border border-green-200 whitespace-nowrap">
                                  Concluído
                                </span>
                              )}
                              {status === "agendado" && (
                                <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-orange-50 text-orange-700 border border-orange-200 whitespace-nowrap">
                                  Agendado
                                </span>
                              )}
                              {status === "cancelado" && (
                                <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-red-50 border border-red-200 whitespace-nowrap" style={{ color: "#D93030" }}>
                                  Cancelado
                                </span>
                              )}
                            </div>
                          </td>
                          <td className="px-6 py-4 text-center">
                            <div className="flex items-center justify-center gap-2">
                              <button 
                                onClick={(e) => {
                                  e.stopPropagation();
                                  setEditingTraining(row);
                                }}
                                className="p-1.5 text-gray-400 hover:text-[#D93030] hover:bg-red-50 rounded-md transition-colors min-h-[44px] min-w-[44px] inline-flex items-center justify-center"
                                title="Editar treinamento"
                              >
                                <Edit2 className="w-4 h-4" />
                              </button>
                              <button 
                                onClick={(e) => {
                                  e.stopPropagation();
                                  setTrainingToDelete(row);
                                }}
                                className="p-1.5 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-md transition-colors min-h-[44px] min-w-[44px] inline-flex items-center justify-center"
                                title="Excluir treinamento"
                              >
                                <Trash2 className="w-4 h-4" />
                              </button>
                            </div>
                          </td>
                        </tr>
                        );
                      })}
                    </tbody>
                  </table>
                </div>

                {/* Mobile Cards */}
                <div className="md:hidden flex flex-col">
                  {(isLoadingTrainings || trainingsError || currentItems.length === 0) && (
                    <div className="p-6 text-center text-sm text-gray-500">
                      {isLoadingTrainings
                        ? "Carregando treinamentos do banco de dados..."
                        : trainingsError || "Nenhum treinamento encontrado."}
                    </div>
                  )}
                  {currentItems.map((row, idx) => {
                    const dataTreinamento = new Date(row.dataHora);
                    const agora = new Date();
                    
                    let status = "agendado";
                    if (row.isCancelado) {
                      status = "cancelado";
                    } else if (dataTreinamento < agora) {
                      status = "concluido";
                    }

                    return (
                      <div
                        key={row.id}
                        className={`p-4 border-b border-gray-100 active:bg-gray-50 transition-colors cursor-pointer ${
                          idx === currentItems.length - 1 ? "border-b-0" : ""
                        }`}
                        onClick={() => setSelectedTraining(row)}
                      >
                        <div className="flex justify-between items-start mb-2 gap-3">
                          <div>
                            <h3 className="text-sm font-semibold text-gray-800 line-clamp-2 mb-1">
                              {row.tema}
                            </h3>
                            <span className="inline-flex items-center px-1.5 py-0.5 rounded text-[10px] font-medium bg-gray-100 text-gray-600">
                              {row.segmento || "Geral (Todos os Segmentos)"}
                            </span>
                          </div>
                          <div className="shrink-0 flex items-center gap-2 mt-1">
                            {status === "concluido" && (
                              <span className="inline-flex items-center px-2 py-0.5 rounded-md text-[10px] font-medium bg-green-50 text-green-700 border border-green-200">
                                Concluído
                              </span>
                            )}
                            {status === "agendado" && (
                              <span className="inline-flex items-center px-2 py-0.5 rounded-md text-[10px] font-medium bg-orange-50 text-orange-700 border border-orange-200">
                                Agendado
                              </span>
                            )}
                            {status === "cancelado" && (
                              <span className="inline-flex items-center px-2 py-0.5 rounded-md text-[10px] font-medium bg-red-50 border border-red-200" style={{ color: "#D93030" }}>
                                Cancelado
                              </span>
                            )}
                            <div className="flex items-center gap-1">
                              <button 
                                onClick={(e) => {
                                  e.stopPropagation();
                                  setEditingTraining(row);
                                }}
                                className="p-1 text-gray-400 hover:text-[#D93030] hover:bg-red-50 rounded-md transition-colors min-h-[44px] min-w-[44px] inline-flex items-center justify-center"
                                title="Editar treinamento"
                              >
                                <Edit2 className="w-4 h-4" />
                              </button>
                              <button 
                                onClick={(e) => {
                                  e.stopPropagation();
                                  setTrainingToDelete(row);
                                }}
                                className="p-1 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-md transition-colors min-h-[44px] min-w-[44px] inline-flex items-center justify-center"
                                title="Excluir treinamento"
                              >
                                <Trash2 className="w-4 h-4" />
                              </button>
                            </div>
                          </div>
                        </div>
                        
                        <div className="space-y-1.5 mt-2">
                          <div className="flex items-center gap-1.5 text-xs text-gray-600">
                            <Calendar className="w-3.5 h-3.5 text-gray-400" />
                            {row.data} às {row.hora}
                          </div>
                          <div className="flex items-center gap-1.5 text-xs text-gray-500">
                            <BookOpen className="w-3.5 h-3.5 text-gray-400 shrink-0" />
                            <span className="line-clamp-1">{row.conteudo}</span>
                          </div>
                        </div>
                      </div>
                    );
                  })}
                </div>
              </div>

              {/* Card Footer */}
              <div className="px-4 md:px-6 py-3 bg-gray-50 border-t border-gray-200 flex flex-col sm:flex-row items-center justify-between gap-3 shrink-0">
                <p className="text-xs text-gray-400 text-center sm:text-left">
                  Mostrando {currentItems.length} de {filteredTrainings.length} registros
                </p>
                <div className="flex items-center gap-1 overflow-x-auto w-full sm:w-auto pb-1 sm:pb-0 justify-center sm:justify-start">
                  <button
                    onClick={() => setCurrentPage(prev => Math.max(prev - 1, 1))}
                    disabled={currentPage === 1}
                    className="px-3 py-1.5 text-xs text-gray-500 border border-gray-200 rounded-md bg-white hover:bg-gray-50 disabled:opacity-40 disabled:cursor-not-allowed transition-colors whitespace-nowrap min-h-[44px]"
                  >
                    Anterior
                  </button>
                  
                  {Array.from({ length: totalPages }, (_, i) => i + 1).map((page) => (
                    <button
                      key={page}
                      onClick={() => setCurrentPage(page)}
                      className={`px-3 py-1.5 text-xs rounded-md transition-colors shrink-0 min-h-[44px] ${
                        currentPage === page 
                          ? "text-white shadow-sm" 
                          : "text-gray-600 bg-white border border-gray-200 hover:bg-gray-50"
                      }`}
                      style={currentPage === page ? { backgroundColor: "#D93030", borderColor: "#D93030" } : undefined}
                    >
                      {page}
                    </button>
                  ))}

                  <button
                    onClick={() => setCurrentPage(prev => Math.min(prev + 1, totalPages))}
                    disabled={currentPage === totalPages}
                    className="px-3 py-1.5 text-xs text-gray-500 border border-gray-200 rounded-md bg-white hover:bg-gray-50 disabled:opacity-40 disabled:cursor-not-allowed transition-colors whitespace-nowrap min-h-[44px]"
                  >
                    Próximo
                  </button>
                </div>
              </div>
            </div>
            )}
          </div>
        </TabsContent>

        <TabsContent value="dashboard" className="flex-1 flex flex-col m-0 p-0 outline-none data-[state=active]:flex pt-4">
          <div className="px-4 md:px-8 pb-8 space-y-6">
            <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-6">
              <h2 className="text-lg font-semibold text-gray-800">Visão Geral do Shopping</h2>
              <div className="flex flex-col sm:flex-row items-stretch sm:items-center gap-4 w-full sm:w-auto">
                <Select value={dashboardPeriod} onValueChange={setDashboardPeriod}>
                  <SelectTrigger className="w-full sm:w-[180px] bg-white border-gray-200 text-sm">
                    <SelectValue placeholder="Período" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="mes">Este Mês</SelectItem>
                    <SelectItem value="trimestre">Último Trimestre</SelectItem>
                    <SelectItem value="ano">Este Ano</SelectItem>
                  </SelectContent>
                </Select>

                <Select>
                  <SelectTrigger className="w-full sm:w-[220px] bg-white border-gray-200 text-sm">
                    <SelectValue placeholder="Filtrar por Loja/Setor" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="todas">Todas as Lojas</SelectItem>
                    <SelectItem value="alimentacao">Alimentação</SelectItem>
                    <SelectItem value="vestuario">Vestuário</SelectItem>
                  </SelectContent>
                </Select>
              </div>
            </div>

            {/* Top 4 Cards Grid */}
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
              <Card className="bg-white border-gray-200 shadow-sm">
                <CardHeader className="flex flex-row items-center justify-between pb-2">
                  <CardTitle className="text-sm font-medium text-gray-600">Taxa de Presença</CardTitle>
                  <UserCheck className="h-4 w-4 text-gray-400" />
                </CardHeader>
                <CardContent>
                  <div className="text-2xl font-bold text-gray-900">{taxaPresenca}%</div>
                  <p className="text-xs text-green-600 font-medium mt-1">Geral do período</p>
                  <Progress value={parseFloat(taxaPresenca as string)} className="h-1.5 mt-4 bg-gray-100" indicatorClassName="bg-[#D93030]" />
                </CardContent>
              </Card>

              <Card className="bg-white border-gray-200 shadow-sm">
                <CardHeader className="flex flex-row items-center justify-between pb-2">
                  <CardTitle className="text-sm font-medium text-gray-600">Lojas Impactadas</CardTitle>
                  <Store className="h-4 w-4 text-gray-400" />
                </CardHeader>
                <CardContent>
                  <div className="text-2xl font-bold text-gray-900">{lojasImpactadas}</div>
                  <p className="text-xs text-gray-500 mt-1">Lojas engajadas (de 9 totais)</p>
                  <Progress value={(lojasImpactadas / 9) * 100} className="h-1.5 mt-4 bg-gray-100" indicatorClassName="bg-[#D93030]" />
                </CardContent>
              </Card>

              <Card className="bg-white border-gray-200 shadow-sm">
                <CardHeader className="flex flex-row items-center justify-between pb-2">
                  <CardTitle className="text-sm font-medium text-gray-600">Média Colab. / Loja</CardTitle>
                  <Users className="h-4 w-4 text-gray-400" />
                </CardHeader>
                <CardContent>
                  <div className="text-2xl font-bold text-gray-900">3.5</div>
                  <p className="text-xs text-gray-500 mt-1">Funcionários por loja engajada</p>
                </CardContent>
              </Card>

              <Card className="bg-white border-gray-200 shadow-sm">
                <CardHeader className="flex flex-row items-center justify-between pb-2">
                  <CardTitle className="text-sm font-medium text-gray-600">Total de Participações</CardTitle>
                  <Award className="h-4 w-4 text-gray-400" />
                </CardHeader>
                <CardContent>
                  <div className="text-2xl font-bold text-gray-900">{totalPresentes}</div>
                  <p className="text-xs text-gray-500 mt-1">Capacitações no período</p>
                </CardContent>
              </Card>
            </div>

            {/* Row 2: Top 5 Lojas Engajadas e Radar de Risco */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mt-6">
              <Card className="bg-white border-gray-200 shadow-sm flex flex-col">
                <CardHeader>
                  <CardTitle className="text-base font-semibold text-gray-800 flex items-center gap-2">
                    <Star className="h-4 w-4 text-yellow-500" />
                    Top 5 Lojas Engajadas
                  </CardTitle>
                </CardHeader>
                <CardContent className="flex-1 flex flex-col">
                  <div className="flex flex-col justify-between flex-1 gap-y-3">
                    {top5Lojas.length > 0 ? top5Lojas.map((loja, i) => (
                      <div key={i} className="flex items-center justify-between p-1">
                        <div className="flex items-center gap-4">
                          <span className="flex items-center justify-center w-8 h-8 rounded-full bg-slate-100 text-sm font-bold text-gray-700">
                            {i + 1}
                          </span>
                          <div className="flex flex-col">
                            <span className="text-sm font-semibold text-gray-900">{loja.name}</span>
                            <span className="text-xs text-slate-500">{loja.totalTrainings} {loja.totalTrainings === 1 ? 'presença' : 'presenças'} no período</span>
                          </div>
                        </div>
                        <div className="flex flex-col items-end gap-1.5 w-24">
                          <span className="text-sm font-bold text-emerald-600">{loja.participation}%</span>
                          <div className="w-full h-1.5 bg-slate-100 rounded-full overflow-hidden">
                            <div 
                              className="h-full bg-emerald-500 rounded-full" 
                              style={{ width: `${loja.participation}%` }}
                            />
                          </div>
                        </div>
                      </div>
                    )) : (
                      <p className="text-sm text-gray-500 text-center py-4 my-auto">Nenhum dado no período.</p>
                    )}
                  </div>
                </CardContent>
              </Card>

              <Card className="bg-white border-gray-200 shadow-sm flex flex-col">
                <CardHeader>
                  <CardTitle className="text-base font-semibold text-gray-800 flex items-center gap-2">
                    <AlertTriangle className="h-4 w-4 text-[#D93030]" />
                    Radar de Risco (Lojas Omissas)
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    {riskStores.length > 0 ? riskStores.map((loja, i) => (
                      <div key={i} className="flex items-center justify-between p-3 bg-red-50 rounded-lg border border-red-100">
                        <div>
                          <p className="text-sm font-bold text-gray-900">{loja.name}</p>
                          <p className="text-xs text-gray-500 mt-0.5">Última presença: {loja.ultimaPresenca}</p>
                        </div>
                        <span className="text-sm font-semibold text-red-500">
                          {loja.faltas} {loja.faltas === 1 ? 'falta' : 'faltas'}
                        </span>
                      </div>
                    )) : (
                      <p className="text-sm text-gray-500 text-center py-4">Sem lojas em risco no período.</p>
                    )}
                  </div>
                </CardContent>
              </Card>
            </div>

            {/* Row 3: Evolução de Participação */}
            <Card className="mt-6 bg-white border-gray-200 shadow-sm">
              <CardHeader>
                <CardTitle className="text-base font-semibold text-gray-800 flex items-center gap-2">
                  <Activity className="h-4 w-4 text-[#D93030]" /> Evolução de Participação
                </CardTitle>
                <CardDescription>Comparativo de inscritos x presentes no período selecionado</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="h-[300px] w-full mt-4">
                  {chartDataDynamic.length > 0 ? (
                    <ResponsiveContainer key={`rc-${dashboardPeriod}`} width="100%" height="100%">
                      <BarChart 
                        id="participation-chart" 
                        data={chartDataDynamic} 
                        margin={{ top: 10, right: 10, left: -20, bottom: 0 }}
                        onClick={(data) => {
                          if (data && data.activePayload && data.activePayload.length > 0) {
                            const clickedData = data.activePayload[0].payload;
                            if (clickedData.monthIndex !== undefined && clickedData.year !== undefined) {
                              setCalendarMonth(new Date(clickedData.year, clickedData.monthIndex, 1));
                              setViewMode("calendar");
                              setActiveTab("list");
                            }
                          }
                        }}
                        style={{ cursor: "pointer" }}
                      >
                        <CartesianGrid key="grid" strokeDasharray="3 3" vertical={false} stroke="#E5E7EB" />
                        <XAxis key="xaxis" dataKey="name" axisLine={false} tickLine={false} tick={{ fontSize: 12, fill: '#6B7280' }} dy={10} />
                        <YAxis key="yaxis" axisLine={false} tickLine={false} tick={{ fontSize: 12, fill: '#6B7280' }} />
                        <Tooltip 
                          key="tooltip"
                          cursor={{ fill: '#F3F4F6' }}
                          contentStyle={{ borderRadius: '8px', border: 'none', boxShadow: '0 4px 6px -1px rgb(0 0 0 / 0.1)' }}
                        />
                        <Legend key="legend" iconType="circle" wrapperStyle={{ paddingTop: '20px' }} />
                        <Bar key="bar1" dataKey="inscritos" name="Inscritos" fill="#E5E7EB" radius={[4, 4, 0, 0]} barSize={32} className="hover:opacity-80 transition-opacity" />
                        <Bar key="bar2" dataKey="presentes" name="Presentes" fill="#D93030" radius={[4, 4, 0, 0]} barSize={32} className="hover:opacity-80 transition-opacity" />
                      </BarChart>
                    </ResponsiveContainer>
                  ) : (
                    <div className="w-full h-full flex items-center justify-center text-sm text-gray-500">
                      Sem dados para o período.
                    </div>
                  )}
                </div>
              </CardContent>
            </Card>

            {/* Store Explorer */}
            <Card className="mt-8 bg-white border-gray-200 shadow-sm">
              <CardHeader className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
                <div>
                  <CardTitle className="text-base font-semibold text-gray-800">Explorador de Lojas</CardTitle>
                  <CardDescription>Analise o desempenho e engajamento de cada loja</CardDescription>
                </div>
                <div className="flex flex-col sm:flex-row items-center gap-3 w-full sm:w-auto">
                  <div className="relative w-full sm:w-48">
                    <Hash className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-gray-400" />
                    <Input
                      placeholder="Filtrar por LUC..."
                      value={lucSearchQuery}
                      onChange={(e) => setLucSearchQuery(e.target.value)}
                      className="pl-9 bg-gray-50 border-transparent focus:bg-white"
                    />
                  </div>
                  <div className="relative w-full sm:w-64">
                    <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-gray-400" />
                    <Input
                      placeholder="Buscar loja..."
                      value={storeSearchQuery}
                      onChange={(e) => setStoreSearchQuery(e.target.value)}
                      className="pl-9 bg-gray-50 border-transparent focus:bg-white"
                    />
                  </div>
                </div>
              </CardHeader>
              <CardContent>
                <div className="overflow-x-auto">
                  <Table>
                    <TableHeader>
                      <TableRow className="hover:bg-transparent">
                        <TableHead className="font-semibold text-gray-600">Loja / LUC</TableHead>
                        <TableHead className="font-semibold text-gray-600">Segmento</TableHead>
                        <TableHead className="font-semibold text-gray-600 text-center">Treinamentos Totais</TableHead>
                        <TableHead className="font-semibold text-gray-600 text-center">% Participação</TableHead>
                        <TableHead className="font-semibold text-gray-600 text-right">Ação</TableHead>
                      </TableRow>
                    </TableHeader>
                    <TableBody>
                      {storeExplorerData
                        .filter(store => 
                          store.name.toLowerCase().includes(storeSearchQuery.toLowerCase()) &&
                          store.luc.toLowerCase().includes(lucSearchQuery.toLowerCase())
                        )
                        .map(store => (
                        <TableRow 
                          key={store.id}
                          className="cursor-pointer hover:bg-gray-50 transition-colors group"
                          onClick={() => setSelectedStore(store)}
                        >
                          <TableCell className="font-medium text-gray-900">
                            <div className="flex flex-col">
                              <span>{store.name}</span>
                              <span className="text-xs text-gray-500 font-normal">{store.luc}</span>
                            </div>
                          </TableCell>
                          <TableCell className="text-gray-500">{store.segment}</TableCell>
                          <TableCell className="text-center text-gray-900">{store.totalTrainings}</TableCell>
                          <TableCell className="text-center">
                            <span className={cn(
                              "font-medium",
                              store.participation >= 80 ? "text-green-600" : store.participation >= 50 ? "text-orange-500" : "text-[#D93030]"
                            )}>
                              {store.participation}%
                            </span>
                          </TableCell>
                          <TableCell className="text-right">
                            <Button 
                              variant="outline" 
                              size="sm"
                              className="text-[#D93030] border-gray-200 hover:bg-red-50 hover:text-[#D93030]"
                              onClick={(e) => {
                                e.stopPropagation();
                                setSelectedStore(store);
                              }}
                            >
                              Ver Detalhes
                            </Button>
                          </TableCell>
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                </div>
              </CardContent>
            </Card>
          </div>
        </TabsContent>
      </Tabs>

      {/* Delete Confirmation Modal */}
      {trainingToDelete && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm animate-in fade-in duration-200">
          <div className="bg-white rounded-lg shadow-xl max-w-sm w-full mx-4 overflow-hidden">
            <div className="p-6">
              <h3 className="text-lg font-semibold text-gray-900 mb-2">
                Tem certeza ou não que quer excluir?
              </h3>
              <p className="text-sm text-gray-500 mb-6">
                Você está prestes a excluir o treinamento <strong>{trainingToDelete.tema}</strong>. Esta ação não poderá ser desfeita.
              </p>
              <div className="flex items-center justify-end gap-3">
                <button
                  onClick={() => setTrainingToDelete(null)}
                  className="px-4 py-2 text-sm font-medium text-gray-700 bg-white hover:bg-gray-100 rounded-lg transition-colors focus:outline-none focus:ring-2 focus:ring-gray-200 min-h-[44px]"
                >
                  Cancelar
                </button>
                <button
                  onClick={handleDelete}
                  className="px-4 py-2 text-sm font-medium text-white rounded-lg transition-colors hover:opacity-90 focus:outline-none focus:ring-2 focus:ring-red-200 min-h-[44px]"
                  style={{ backgroundColor: "#D93030" }}
                >
                  Excluir
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
    </main>
  );
}
