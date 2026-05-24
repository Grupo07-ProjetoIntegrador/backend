import { useState, useCallback, useEffect } from "react";
import { ArrowLeft, UploadCloud, FileSpreadsheet, CheckCircle2, User, Building, Calendar, BookOpen, UserPlus, Trash2, X, XCircle, Users, CheckSquare, Store } from "lucide-react";
import { toast } from "sonner";
import * as XLSX from "xlsx";
import { Card, CardHeader, CardTitle, CardContent } from "./ui/card";

export interface TrainingDetailsProps {
  training: {
    id: number | string;
    tema: string;
    data: string;
    hora: string;
    conteudo: string;
    dataHora?: string;
    attendanceList?: any[];
  };
  onBack: () => void;
  onUpdateAttendance?: (id: number | string, list: any[]) => void;
}

export function TrainingDetails({ training, onBack, onUpdateAttendance }: TrainingDetailsProps) {
  const isConcluido = training.dataHora ? new Date(training.dataHora) < new Date(2026, 4, 4) : false;
  const isAgendado = training.dataHora ? new Date(training.dataHora) > new Date(2026, 4, 4) : false;

  const [isDragging, setIsDragging] = useState(false);
  const [attendees, setAttendees] = useState<{ id: number; luc: string; loja: string; representante: string; status: string }[]>(
    training.attendanceList ? training.attendanceList : []
  );
  const [showAttendanceBlock, setShowAttendanceBlock] = useState(false);
  const [isDataLoaded, setIsDataLoaded] = useState(training.attendanceList && training.attendanceList.length > 0);
  const [isAddModalOpen, setIsAddModalOpen] = useState(false);
  const [newAttendee, setNewAttendee] = useState({ luc: "", loja: "", representante: "", status: "Presente" });

  const updateAttendees = (newList: any[]) => {
    setAttendees(newList);
    if (onUpdateAttendance) {
      onUpdateAttendance(training.id, newList);
    }
  };

  const handleAddAttendee = (e: React.FormEvent) => {
    e.preventDefault();
    if (!newAttendee.luc || !newAttendee.loja || !newAttendee.representante || !newAttendee.status) return;
    
    const newId = attendees.length > 0 ? Math.max(...attendees.map(a => a.id)) + 1 : 1;
    updateAttendees([...attendees, { 
      id: newId, 
      ...newAttendee
    }]);
    
    setNewAttendee({ luc: "", loja: "", representante: "", status: "Presente" });
    setIsAddModalOpen(false);
    toast.success("Participante adicionado!");
  };

  const handleRemoveAttendee = (id: number) => {
    updateAttendees(attendees.filter(a => a.id !== id));
    toast("Participante removido");
  };

  const handleDragOver = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(true);
  }, []);

  const handleDragLeave = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(false);
  }, []);

  const handleDrop = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(false);
    
    if (e.dataTransfer.files && e.dataTransfer.files.length > 0) {
      processFile(e.dataTransfer.files[0]);
    }
  }, []);

  const handleFileInput = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files.length > 0) {
      processFile(e.target.files[0]);
    }
    e.target.value = '';
  };

  const processFile = (file: File) => {
    const isCsv = file.name.toLowerCase().endsWith('.csv');
    const isXlsx = file.name.toLowerCase().endsWith('.xlsx');

    if (!isCsv && !isXlsx) {
      toast.error("Por favor, envie um arquivo no formato .csv ou .xlsx", { id: "upload-toast" });
      return;
    }

    toast.loading("Lendo arquivo...", { id: "upload-toast" });
    
    const reader = new FileReader();
    
    reader.onload = (event) => {
      try {
        const parsedData: { luc: string; loja: string; representante: string; status: string }[] = [];

        if (isXlsx) {
          const data = event.target?.result;
          const workbook = XLSX.read(data, { type: 'array' });
          const firstSheetName = workbook.SheetNames[0];
          const worksheet = workbook.Sheets[firstSheetName];
          const json = XLSX.utils.sheet_to_json(worksheet) as any[];

          if (json.length === 0) {
            toast.error("O arquivo Excel está vazio ou não contém dados suficientes.", { id: "upload-toast" });
            return;
          }

          for (const row of json) {
            // Mapping keys: "LUC", "Nome da Loja", "Nome do Representante", "Status"
            const luc = row["LUC"] != null ? String(row["LUC"]).trim() : "";
            const loja = row["Nome da Loja"] != null ? String(row["Nome da Loja"]).trim() : "";
            const representante = row["Nome do Representante"] != null ? String(row["Nome do Representante"]).trim() : "";
            const status = row["Status"] != null ? String(row["Status"]).trim() : "Presente";

            if (luc || loja || representante) {
              parsedData.push({
                luc,
                loja,
                representante,
                status
              });
            }
          }
        } else if (isCsv) {
          const text = event.target?.result as string;
          const lines = text.split(/\r?\n/).filter(line => line.trim() !== '');
          
          if (lines.length < 2) {
            toast.error("O arquivo CSV está vazio ou não contém dados suficientes.", { id: "upload-toast" });
            return;
          }

          for (let i = 1; i < lines.length; i++) {
            const columns = lines[i].split(',').map(col => col.trim());
            if (columns.length >= 3) {
              const parsedStatus = columns[3] ? columns[3].trim().replace('\r', '') : 'Presente';
              parsedData.push({
                luc: columns[0],
                loja: columns[1],
                representante: columns[2],
                status: parsedStatus
              });
            }
          }
        }

        const newData: { id: number; luc: string; loja: string; representante: string; status: string }[] = [];
        let currentMaxId = attendees.length > 0 ? Math.max(...attendees.map(a => a.id)) : 0;

        for (const newItem of parsedData) {
          const isDuplicate = attendees.some(
            a => a.luc === newItem.luc && a.representante === newItem.representante
          ) || newData.some(
            a => a.luc === newItem.luc && a.representante === newItem.representante
          );

          if (!isDuplicate) {
            newData.push({ ...newItem, id: ++currentMaxId });
          }
        }

        if (newData.length > 0) {
          updateAttendees([...attendees, ...newData]);
          toast.success(`${newData.length} participante(s) adicionado(s) com sucesso!`, { id: "upload-toast" });
          setIsDataLoaded(true);
        } else {
          toast.info("Nenhum novo participante adicionado. Todas as entradas já existiam ou são inválidas.", { id: "upload-toast" });
        }

      } catch (error) {
        toast.error("Ocorreu um erro ao processar o arquivo.", { id: "upload-toast" });
      }
    };

    reader.onerror = () => {
      toast.error("Erro na leitura do arquivo.", { id: "upload-toast" });
    };

    if (isXlsx) {
      reader.readAsArrayBuffer(file);
    } else {
      reader.readAsText(file);
    }
  };

  const ROOM_CAPACITY = 50;
  const isListEmpty = attendees.length === 0;
  const totalInscritos = attendees.length;
  const totalPresentes = attendees.filter(a => a.status === 'Presente').length;
  const comparecimentoPercent = totalInscritos > 0 ? Math.round((totalPresentes / totalInscritos) * 100) : 0;
  const ocupacaoPercent = Math.min(100, Math.round((totalPresentes / ROOM_CAPACITY) * 100));
  const lojasUnicas = new Set(attendees.map(a => a.loja)).size;

  return (
    <div className="flex-1 min-h-screen flex flex-col p-4 md:p-8 overflow-y-auto" style={{ backgroundColor: "#F7F4EF" }}>
      {/* Header */}
      <div className="mb-6 md:mb-8 flex flex-col sm:flex-row sm:items-start gap-4">
        <div className="flex items-start gap-3">
          <button
            onClick={onBack}
            className="mt-0.5 md:mt-1 p-2 -ml-2 sm:ml-0 rounded-lg text-gray-500 hover:bg-white hover:text-[#D93030] transition-colors border border-transparent hover:border-gray-200 hover:shadow-sm shrink-0 min-h-[44px] min-w-[44px] flex items-center justify-center"
            aria-label="Voltar"
          >
            <ArrowLeft className="w-5 h-5" />
          </button>
          <div>
            <h1 className="text-2xl md:text-3xl mb-2 md:mb-3 leading-tight" style={{ color: "#8B1A1A" }}>
              Palestra: {training.tema}
            </h1>
            <div className="flex flex-col sm:flex-row sm:items-center gap-2 sm:gap-4 text-sm text-gray-600">
              <span className="flex items-center gap-1.5">
                <Calendar className="w-4 h-4 shrink-0" style={{ color: "#D93030" }} />
                {training.data} às {training.hora}
              </span>
              <span className="hidden sm:block w-1 h-1 bg-gray-300 rounded-full shrink-0" />
              <span className="flex items-center gap-1.5">
                <BookOpen className="w-4 h-4 shrink-0" style={{ color: "#D93030" }} />
                <span className="line-clamp-2 sm:line-clamp-1">{training.conteudo}</span>
              </span>
            </div>
          </div>
        </div>
      </div>

      <div className="w-full space-y-8">
        {/* KPI Cards */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 w-full">
          <Card className={`bg-white border-gray-200 shadow-sm transition-all duration-700 ease-in-out ${isListEmpty ? 'opacity-60' : 'opacity-100'}`}>
            <CardHeader className="flex flex-row items-center justify-between pb-1 px-4 pt-4">
              <CardTitle className="text-xs font-medium text-gray-600">Comparecimento</CardTitle>
              <Users className="h-3.5 w-3.5 text-gray-400" />
            </CardHeader>
            <CardContent className="px-4 pb-4">
              <div className={`text-xl font-bold text-gray-900 transition-opacity duration-500 delay-100 ${isListEmpty ? 'opacity-50' : 'opacity-100'}`}>
                {!isListEmpty ? `${comparecimentoPercent}%` : '--'}
              </div>
              <p className={`text-[11px] text-muted-foreground mt-0.5 transition-opacity duration-500 delay-200 ${isListEmpty ? 'opacity-50' : 'opacity-100'}`}>
                {!isListEmpty ? `${totalPresentes} presentes / ${totalInscritos} inscritos` : 'Aguardando dados...'}
              </p>
            </CardContent>
          </Card>

          <Card className={`bg-white border-gray-200 shadow-sm transition-all duration-700 ease-in-out ${isListEmpty ? 'opacity-60' : 'opacity-100'}`}>
            <CardHeader className="flex flex-row items-center justify-between pb-1 px-4 pt-4">
              <CardTitle className="text-xs font-medium text-gray-600">Ocupação da Sala</CardTitle>
              <CheckSquare className="h-3.5 w-3.5 text-gray-400" />
            </CardHeader>
            <CardContent className="px-4 pb-4">
              <div className={`text-xl font-bold text-gray-900 transition-opacity duration-500 delay-100 ${isListEmpty ? 'opacity-50' : 'opacity-100'}`}>
                {!isListEmpty ? `${ocupacaoPercent}%` : '--'}
              </div>
              <p className={`text-[11px] text-muted-foreground mt-0.5 transition-opacity duration-500 delay-200 ${isListEmpty ? 'opacity-50' : 'opacity-100'}`}>
                {!isListEmpty ? `${totalPresentes} ocupadas / ${ROOM_CAPACITY} vagas` : 'Aguardando dados...'}
              </p>
            </CardContent>
          </Card>

          <Card className={`bg-white border-gray-200 shadow-sm transition-all duration-700 ease-in-out ${isListEmpty ? 'opacity-60' : 'opacity-100'}`}>
            <CardHeader className="flex flex-row items-center justify-between pb-1 px-4 pt-4">
               <CardTitle className="text-xs font-medium text-gray-600">Lojas Representadas</CardTitle>
               <Store className="h-3.5 w-3.5 text-gray-400" />
            </CardHeader>
            <CardContent className="px-4 pb-4">
              <div className={`text-xl font-bold text-gray-900 transition-opacity duration-500 delay-100 ${isListEmpty ? 'opacity-50' : 'opacity-100'}`}>
                {!isListEmpty ? lojasUnicas : '--'}
              </div>
              <p className={`text-[11px] text-muted-foreground mt-0.5 transition-opacity duration-500 delay-200 ${isListEmpty ? 'opacity-50' : 'opacity-100'}`}>
                {!isListEmpty ? 'Lojas diferentes nesta turma' : 'Aguardando dados...'}
              </p>
            </CardContent>
          </Card>
        </div>

        {isAgendado && attendees.length === 0 && !showAttendanceBlock ? (
          <Card className="bg-gray-50/50 border-dashed border-2 border-gray-200 shadow-none mb-8 w-full flex-1 flex">
            <CardContent className="flex flex-col items-center justify-center py-16 w-full text-center flex-1">
              <div className="w-14 h-14 bg-white rounded-full flex items-center justify-center mb-5 shadow-sm border border-gray-100">
                <Calendar className="h-6 w-6 text-gray-400" />
              </div>
              <h3 className="text-lg font-semibold text-gray-900 mb-2">Treinamento Agendado</h3>
              <p className="text-sm text-gray-500 max-w-md mb-8">
                Este treinamento ainda não ocorreu. A lista de presença será disponibilizada após a sua realização.
              </p>
              <button 
                onClick={() => setShowAttendanceBlock(true)}
                className="inline-flex items-center justify-center px-4 py-2.5 text-sm font-medium border border-gray-300 rounded-lg text-gray-700 bg-white hover:bg-gray-50 transition-colors shadow-sm"
              >
                <FileSpreadsheet className="w-4 h-4 mr-2" />
                Lançar Lista Antecipada
              </button>
            </CardContent>
          </Card>
        ) : (
          <>
            {/* Dropzone */}
            <div 
              className={`border-2 border-dashed rounded-lg p-6 flex flex-col items-center justify-center transition-colors ${
                isDragging ? 'border-[#C4151F] bg-red-50' : 'border-gray-300 hover:border-gray-400 bg-white shadow-sm'
              }`}
              onDragOver={handleDragOver}
              onDragLeave={handleDragLeave}
              onDrop={handleDrop}
            >
              <UploadCloud className={`w-10 h-10 mb-3 ${isDragging ? 'text-[#C4151F]' : 'text-gray-400'}`} />
              <h3 className="text-base font-medium text-gray-900 mb-1">Importar Lista de Presença</h3>
              <p className="text-sm text-gray-500 text-center mb-4 max-w-md">
                Arraste e solte o arquivo da lista de presença (CSV ou XLSX) ou clique para selecionar.
              </p>
              <label className="cursor-pointer inline-flex items-center justify-center px-4 py-2 text-sm font-medium border border-gray-300 rounded-lg text-gray-700 bg-white hover:bg-gray-50 transition-colors shadow-sm">
                <FileSpreadsheet className="w-4 h-4 mr-2" />
                Selecionar Arquivo
                <input type="file" className="hidden" accept=".csv, application/vnd.openxmlformats-officedocument.spreadsheetml.sheet, application/vnd.ms-excel" onChange={handleFileInput} />
              </label>
            </div>

            {/* Upload Section / Lista Consolidada */}
            <div className="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden w-full flex-1 flex flex-col mb-8">
              <div className="px-4 md:px-6 py-4 border-b border-gray-200 flex flex-col sm:flex-row sm:items-center justify-between gap-3">
                <div className="flex flex-col sm:flex-row sm:items-center gap-3">
                  <h2 className="text-gray-800 font-semibold text-lg flex items-center gap-2">
                    <Users className="w-5 h-5" style={{ color: "#D93030" }} />
                    Lista de Presença Consolidada
                  </h2>
                  <span className="text-xs font-medium px-3 py-1 rounded-full self-start sm:self-auto bg-red-50 border border-red-200" style={{ color: "#D93030" }}>
                    {attendees.length} presenças registradas
                  </span>
                </div>
                <div className="flex justify-end pt-4 sm:pt-0">
                  <button
                    onClick={() => setIsAddModalOpen(true)}
                    className="inline-flex items-center justify-center gap-2 px-4 py-2 text-sm font-medium border border-gray-300 rounded-lg text-gray-700 bg-white hover:bg-gray-50 transition-colors"
                  >
                    <UserPlus className="w-4 h-4" />
                    Adicionar Manualmente
                  </button>
                </div>
              </div>
              
              <div className="flex-1 overflow-x-auto p-0">
                <table className="w-full text-sm text-left">
                  <thead className="bg-gray-50/80 text-gray-500 font-medium border-b border-gray-200 sticky top-0 z-10">
                    <tr>
                      <th className="px-6 py-3 whitespace-nowrap font-medium">LUC</th>
                      <th className="px-6 py-3 whitespace-nowrap font-medium">Loja</th>
                      <th className="px-6 py-3 whitespace-nowrap font-medium">Representante</th>
                      <th className="px-6 py-3 text-center whitespace-nowrap font-medium">Status</th>
                      <th className="px-6 py-3 text-right whitespace-nowrap font-medium">Ações</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-gray-100 bg-white">
                    {attendees.map((attendee) => (
                      <tr key={attendee.id} className="hover:bg-gray-50/50 transition-colors group">
                        <td className="px-6 py-4 text-gray-500 whitespace-nowrap font-mono">
                          {attendee.luc ? attendee.luc.replace(/\D/g, '').slice(0, 3) : ''}
                        </td>
                        <td className="px-6 py-4 text-gray-900 font-medium whitespace-nowrap">
                          {attendee.loja}
                        </td>
                        <td className="px-6 py-4 text-gray-600 whitespace-nowrap">
                          {attendee.representante}
                        </td>
                        <td className="px-6 py-4 text-center whitespace-nowrap">
                          <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                            attendee.status === 'Presente'
                              ? 'bg-green-50 text-green-700'
                              : 'bg-red-50 text-red-700'
                          }`}>
                            {attendee.status === 'Presente' ? "Presente" : "Ausente"}
                          </span>
                        </td>
                        <td className="px-6 py-4 text-right whitespace-nowrap">
                          <button
                            onClick={() => handleRemoveAttendee(attendee.id)}
                            className="text-gray-500 hover:text-red-600 hover:bg-red-50 p-1.5 rounded-md transition-colors"
                            title="Excluir"
                          >
                            <Trash2 className="w-4 h-4" />
                          </button>
                        </td>
                      </tr>
                    ))}
                    {attendees.length === 0 && (
                      <tr>
                        <td colSpan={5} className="px-6 py-8 text-center text-gray-500">
                          Nenhum participante registrado.
                        </td>
                      </tr>
                    )}
                  </tbody>
                </table>
              </div>
            </div>
          </>
        )}
      </div>

      {/* Add Manual Attendee Modal */}
      {isAddModalOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm animate-in fade-in duration-200">
          <div className="bg-white rounded-lg shadow-xl max-w-md w-full mx-4 overflow-hidden animate-in zoom-in-95 duration-200">
            <div className="flex items-center justify-between px-6 py-4 border-b border-gray-100">
              <div>
                <h3 className="text-lg font-semibold text-gray-900">Adicionar Participante</h3>
                <p className="text-sm text-gray-500 mt-1">Insira os dados do lojista manualmente.</p>
              </div>
              <button
                onClick={() => setIsAddModalOpen(false)}
                className="p-2 text-gray-400 hover:text-gray-600 hover:bg-gray-100 rounded-full transition-colors"
              >
                <X className="w-5 h-5" />
              </button>
            </div>
            
            <form onSubmit={handleAddAttendee} className="p-6">
              <div className="space-y-4">
                <div className="space-y-2">
                  <label htmlFor="luc" className="text-sm font-medium" style={{ color: "#2A1A1A" }}>LUC</label>
                  <input
                    id="luc"
                    required
                    value={newAttendee.luc}
                    onChange={(e) => setNewAttendee({...newAttendee, luc: e.target.value})}
                    placeholder="Ex: LUC-101"
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#C4151F]/20 focus:border-[#C4151F] transition-colors text-sm text-gray-900"
                  />
                </div>
                
                <div className="space-y-2">
                  <label htmlFor="loja" className="text-sm font-medium" style={{ color: "#2A1A1A" }}>Nome da Loja</label>
                  <input
                    id="loja"
                    required
                    value={newAttendee.loja}
                    onChange={(e) => setNewAttendee({...newAttendee, loja: e.target.value})}
                    placeholder="Ex: O Boticário"
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#C4151F]/20 focus:border-[#C4151F] transition-colors text-sm text-gray-900"
                  />
                </div>
                
                <div className="space-y-2">
                  <label htmlFor="representante" className="text-sm font-medium" style={{ color: "#2A1A1A" }}>Nome do Representante</label>
                  <input
                    id="representante"
                    required
                    value={newAttendee.representante}
                    onChange={(e) => setNewAttendee({...newAttendee, representante: e.target.value})}
                    placeholder="Ex: Ana Silva"
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#C4151F]/20 focus:border-[#C4151F] transition-colors text-sm text-gray-900"
                  />
                </div>

                <div className="space-y-2">
                  <label htmlFor="status" className="text-sm font-medium" style={{ color: "#2A1A1A" }}>Status</label>
                  <select
                    id="status"
                    required
                    value={newAttendee.status}
                    onChange={(e) => setNewAttendee({...newAttendee, status: e.target.value})}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#C4151F]/20 focus:border-[#C4151F] transition-colors text-sm text-gray-900 bg-white"
                  >
                    <option value="Presente">Presente</option>
                    <option value="Não Presente">Não Presente</option>
                  </select>
                </div>
              </div>
              
              <div className="mt-8 flex items-center justify-end gap-3">
                <button
                  type="button"
                  onClick={() => setIsAddModalOpen(false)}
                  className="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 hover:bg-gray-50 rounded-lg transition-colors focus:outline-none focus:ring-2 focus:ring-gray-200"
                >
                  Cancelar
                </button>
                <button
                  type="submit"
                  className="px-4 py-2 text-sm font-medium text-white rounded-lg transition-colors hover:opacity-90 focus:outline-none focus:ring-2 focus:ring-red-200"
                  style={{ backgroundColor: "#C4151F" }}
                >
                  Salvar
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
