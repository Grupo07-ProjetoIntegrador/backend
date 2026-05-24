import { ArrowLeft, Calendar, Building, Store, TrendingUp, TrendingDown, Users, Award, ExternalLink, BookOpen } from "lucide-react";
import { Badge } from "./ui/badge";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "./ui/table";
import { Card, CardHeader, CardTitle, CardContent } from "./ui/card";
import { cn } from "./ui/utils";

export interface StoreDetailsProps {
  store: {
    id: number;
    name: string;
    luc: string;
    segment: string;
    manager: string;
  };
  trainings: any[];
  onBack: () => void;
  onSelectTraining: (training: any) => void;
}

export function StoreDetails({ store, trainings, onBack, onSelectTraining }: StoreDetailsProps) {
  const attendanceMap = new Map();
  trainings.forEach(t => {
    if (t.attendanceList) {
      const record = t.attendanceList.find((a: any) => a.id === store.id);
      if (record) {
        attendanceMap.set(t.id, record.status);
      }
    }
  });

  const invitedTrainings = trainings
    .filter(t => attendanceMap.has(t.id))
    .sort((a, b) => new Date(b.dataHora).getTime() - new Date(a.dataHora).getTime());

  const totalConvites = invitedTrainings.length;
  const totalPresentes = invitedTrainings.filter(t => attendanceMap.get(t.id) === "Presente").length;
  const taxaPresenca = totalConvites > 0 ? Math.round((totalPresentes / totalConvites) * 100) : 0;
  
  // Mocking engajamento único based on presences
  const engajamentoUnico = Math.max(1, Math.min(totalPresentes, 4));

  const formatLuc = (lucString: string) => {
    const digits = lucString.replace(/\D/g, '').slice(0, 3);
    return digits ? `LUC ${digits}` : lucString;
  };

  return (
    <main className="flex-1 min-h-screen flex flex-col" style={{ backgroundColor: "#F7F4EF" }}>
      {/* Header */}
      <div className="bg-white border-b border-gray-200 px-4 md:px-8 py-6 shrink-0 relative z-20 shadow-sm">
        <button 
          onClick={onBack}
          className="flex items-center text-sm font-medium text-gray-500 hover:text-[#C4151F] transition-colors mb-6 w-fit border border-transparent hover:border-gray-200 py-1.5 px-3 -ml-3 rounded-md"
        >
          <ArrowLeft className="w-4 h-4 mr-2" />
          Voltar para gestão
        </button>

        <div className="flex flex-col md:flex-row md:items-start justify-between gap-6">
          <div>
            <div className="flex items-center gap-4 mb-2">
              <div className="w-14 h-14 rounded-2xl flex items-center justify-center shrink-0 shadow-sm border border-red-100" style={{ backgroundColor: "#C4151F" }}>
                <Store className="w-7 h-7 text-white" />
              </div>
              <div>
                <h1 className="text-2xl font-bold text-gray-900 tracking-tight">{store.name}</h1>
                <p className="text-sm font-medium text-gray-500 mt-1 flex items-center gap-2">
                  <span className="font-semibold text-gray-700 bg-gray-100 px-2 py-0.5 rounded-md">{formatLuc(store.luc)}</span>
                  <span>•</span>
                  {store.segment}
                  <span>•</span>
                  Gerente: {store.manager}
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Content */}
      <div className="p-4 md:p-8 flex-1 overflow-y-auto space-y-8">
        
        {/* KPI Cards */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 w-full">
          <Card className="bg-white border-gray-200 shadow-sm">
            <CardHeader className="flex flex-row items-center justify-between pb-1 px-5 pt-5">
              <CardTitle className="text-xs font-medium text-gray-600">Taxa de Presença</CardTitle>
              <Award className="h-4 w-4 text-gray-400" />
            </CardHeader>
            <CardContent className="px-5 pb-5">
              <div className="flex items-baseline gap-3">
                <div className="text-2xl font-bold text-gray-900">
                  {taxaPresenca}%
                </div>
                {taxaPresenca >= 70 ? (
                  <span className="flex items-center text-xs font-medium text-green-600 bg-green-50 px-1.5 py-0.5 rounded-md">
                    <TrendingUp className="w-3 h-3 mr-1" />
                    Alta
                  </span>
                ) : (
                  <span className="flex items-center text-xs font-medium text-red-600 bg-red-50 px-1.5 py-0.5 rounded-md">
                    <TrendingDown className="w-3 h-3 mr-1" />
                    Baixa
                  </span>
                )}
              </div>
              <p className="text-[11px] text-muted-foreground mt-1">
                {totalPresentes} presenças registradas
              </p>
            </CardContent>
          </Card>

          <Card className="bg-white border-gray-200 shadow-sm">
            <CardHeader className="flex flex-row items-center justify-between pb-1 px-5 pt-5">
              <CardTitle className="text-xs font-medium text-gray-600">Total de Convites</CardTitle>
              <Calendar className="h-4 w-4 text-gray-400" />
            </CardHeader>
            <CardContent className="px-5 pb-5">
              <div className="text-2xl font-bold text-gray-900">
                {totalConvites}
              </div>
              <p className="text-[11px] text-muted-foreground mt-1">
                Treinamentos agendados/concluídos
              </p>
            </CardContent>
          </Card>

          <Card className="bg-white border-gray-200 shadow-sm">
            <CardHeader className="flex flex-row items-center justify-between pb-1 px-5 pt-5">
              <CardTitle className="text-xs font-medium text-gray-600">Engajamento Único</CardTitle>
              <Users className="h-4 w-4 text-gray-400" />
            </CardHeader>
            <CardContent className="px-5 pb-5">
              <div className="text-2xl font-bold text-gray-900">
                {totalPresentes > 0 ? engajamentoUnico : 0}
              </div>
              <p className="text-[11px] text-muted-foreground mt-1">
                Colaboradores diferentes treinados
              </p>
            </CardContent>
          </Card>
        </div>

        {/* Training History Table */}
        <div className="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
          <div className="px-6 py-5 border-b border-gray-200 flex items-center justify-between">
            <h2 className="text-base font-semibold text-gray-900 flex items-center gap-2">
              <BookOpen className="w-5 h-5 text-[#C4151F]" />
              Histórico de Treinamentos
            </h2>
          </div>
          
          <div className="overflow-x-auto">
            <Table>
              <TableHeader>
                <TableRow className="bg-gray-50/80 hover:bg-gray-50/80 border-b border-gray-200">
                  <TableHead className="font-semibold text-gray-600 w-[140px] px-6">Data</TableHead>
                  <TableHead className="font-semibold text-gray-600 px-6">Tema</TableHead>
                  <TableHead className="font-semibold text-gray-600 w-[160px] text-center px-6">Status</TableHead>
                  <TableHead className="font-semibold text-gray-600 w-[140px] text-right px-6">Ações</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody className="divide-y divide-gray-100">
                {invitedTrainings.length > 0 ? invitedTrainings.map((training) => {
                  const rawStatus = attendanceMap.get(training.id);
                  const isAgendado = new Date(training.dataHora) > new Date();
                  
                  let badgeVariant = "bg-gray-50 text-gray-700 border-gray-200";
                  let badgeText = "Desconhecido";

                  if (isAgendado) {
                    badgeVariant = "bg-amber-50 text-amber-700 border-amber-200";
                    badgeText = "Agendado";
                  } else if (rawStatus === "Presente") {
                    badgeVariant = "bg-green-50 text-green-700 border-green-200";
                    badgeText = "Presente";
                  } else {
                    badgeVariant = "bg-red-50 text-red-700 border-red-200";
                    badgeText = "Ausente";
                  }
                  
                  return (
                    <TableRow 
                      key={training.id} 
                      className="group cursor-pointer hover:bg-slate-50 transition-colors"
                      onClick={() => onSelectTraining(training)}
                    >
                      <TableCell className="text-gray-600 font-medium whitespace-nowrap px-6">
                        {training.data}
                        <div className="text-xs text-gray-400 font-normal mt-0.5">{training.hora}</div>
                      </TableCell>
                      <TableCell className="px-6">
                        <div className="flex flex-col gap-1">
                          <span className="font-semibold text-gray-900 group-hover:text-[#C4151F] transition-colors">
                            {training.tema}
                          </span>
                          <span className="text-sm text-gray-500 line-clamp-1">
                            {training.conteudo}
                          </span>
                        </div>
                      </TableCell>
                      <TableCell className="text-center px-6">
                        <span className={cn(
                          "inline-flex items-center justify-center px-2.5 py-0.5 rounded-full text-xs font-medium border whitespace-nowrap",
                          badgeVariant
                        )}>
                          {badgeText}
                        </span>
                      </TableCell>
                      <TableCell className="text-right px-6">
                        <div
                          className="inline-flex items-center justify-center gap-1.5 px-3 py-1.5 text-xs font-medium text-gray-700 bg-white border border-gray-300 rounded-md transition-colors focus:outline-none group-hover:bg-gray-100 group-hover:text-gray-900 group-hover:border-gray-400"
                        >
                          Ver Detalhes
                          <ExternalLink className="w-3.5 h-3.5" />
                        </div>
                      </TableCell>
                    </TableRow>
                  );
                }) : (
                  <TableRow>
                    <TableCell colSpan={4} className="text-center py-10 text-gray-500">
                      Nenhum treinamento registrado para esta loja.
                    </TableCell>
                  </TableRow>
                )}
              </TableBody>
            </Table>
          </div>
        </div>
      </div>
    </main>
  );
}