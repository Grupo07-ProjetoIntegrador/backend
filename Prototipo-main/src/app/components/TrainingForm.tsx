import { useState } from "react";
import { useForm, Controller } from "react-hook-form";
import { toast } from "sonner";
import {
  ArrowLeft,
  Calendar as CalendarIcon,
  Clock,
  MapPin,
  FileText,
  UploadCloud,
  X,
  Plus,
  Info
} from "lucide-react";
import { Popover, PopoverContent, PopoverTrigger } from "./ui/popover";
import { Calendar } from "./ui/calendar";
import { ScrollArea } from "./ui/scroll-area";

interface FormValues {
  tema: string;
  categoria: string;
  segmento: string;
  data: string;
  horarioInicio: string;
  horarioFim: string;
  status: string;
  capacidadeMaxima: string;
  local: string;
  modalidade: string;
  conteudo: string;
  objetivo: string;
  observacoes: string;
  responsavel: string;
  areaResponsavel: string;
  publicoAlvo: string;
  tags: string;
  recorrente: boolean;
  frequencia?: string;
}

interface TrainingFormProps {
  onBack: () => void;
  onSuccess: (data: any) => void;
  initialData?: FormValues | any;
}

const API_BASE_URL = "http://localhost:8080";

const toBrazilianDate = (date: string) => {
  if (!date) return "";
  const [year, month, day] = date.split("-");
  return `${day}/${month}/${year}`;
};

const toApiPayload = (data: FormValues) => ({
  tema: data.tema,
  descricao: data.conteudo || data.objetivo || "",
  categoria: data.categoria || "Geral",
  data: toBrazilianDate(data.data),
  horario_inicio: data.horarioInicio || "00:00",
  horario_fim: data.horarioFim || data.horarioInicio || "00:00",
  local: data.local || "",
  modalidade: data.modalidade || "Presencial",
  conteudo: data.conteudo,
  capacidade_maxima: Number(data.capacidadeMaxima) || 0,
  segmento_alvo: data.segmento || data.publicoAlvo || "Geral",
  status: data.status || "agendado",
  objetivo: data.objetivo || "",
  observacoes: data.observacoes || "",
  material_apoio: "",
  responsavel: data.responsavel || "Admin",
  area_responsavel: data.areaResponsavel || "",
  tags: data.tags || "",
  recorrente: !!data.recorrente,
});

const TimePickerContent = ({ value, onChange }: { value: string, onChange: (val: string) => void }) => {
  const hours = Array.from({ length: 24 }, (_, i) => i.toString().padStart(2, '0'));
  const minutes = Array.from({ length: 60 }, (_, i) => i.toString().padStart(2, '0'));
  
  const [h, m] = (value || "08:00").split(":");
  
  return (
    <div className="flex gap-2 p-3 h-[200px] w-auto">
      <ScrollArea className="h-full w-16">
        <div className="flex flex-col gap-1 pr-3">
          {hours.map((hour) => (
            <button
              key={`h-${hour}`}
              type="button"
              className={`p-2 text-sm rounded-md text-center hover:bg-gray-100 ${hour === h ? 'bg-red-50 text-red-600 font-medium' : 'text-gray-700'}`}
              onClick={() => onChange(`${hour}:${m || "00"}`)}
            >
              {hour}
            </button>
          ))}
        </div>
      </ScrollArea>
      <div className="flex items-center justify-center font-bold text-gray-400">:</div>
      <ScrollArea className="h-full w-16">
        <div className="flex flex-col gap-1 pr-3">
          {minutes.map((minute) => (
            <button
              key={`m-${minute}`}
              type="button"
              className={`p-2 text-sm rounded-md text-center hover:bg-gray-100 ${minute === m ? 'bg-red-50 text-red-600 font-medium' : 'text-gray-700'}`}
              onClick={() => onChange(`${h || "08"}:${minute}`)}
            >
              {minute}
            </button>
          ))}
        </div>
      </ScrollArea>
    </div>
  );
};

export function TrainingForm({ onBack, onSuccess, initialData }: TrainingFormProps) {
  const {
    register,
    handleSubmit,
    watch,
    control,
    setValue,
    formState: { errors, isSubmitting },
  } = useForm<FormValues>({
    defaultValues: initialData || {
      categoria: "",
      segmento: "Geral",
      status: "agendado",
      capacidadeMaxima: "",
      modalidade: "",
      responsavel: "Admin Master",
      areaResponsavel: "",
      recorrente: false,
    },
    mode: "onBlur"
  });

  const isRecorrente = watch("recorrente");
  const [showCancelModal, setShowCancelModal] = useState(false);
  
  const isEditing = !!initialData;

  const onSubmit = async (data: FormValues) => {
    if (isEditing) {
      toast.info("A API atual ainda nao possui rota de edicao para treinamentos.");
      return;
    }

    try {
      const response = await fetch(`${API_BASE_URL}/api/treinamentos/cadastrar`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(toApiPayload(data)),
      });

      if (!response.ok) {
        const message = await response.text();
        throw new Error(message || "Erro ao salvar treinamento");
      }

      toast.success("Treinamento cadastrado com sucesso");
      onSuccess(data);
    } catch (error) {
      console.error("Erro ao cadastrar treinamento:", error);
      toast.error("Nao foi possivel salvar. Confira se o backend Go esta rodando.");
    }
  };

  const handleCancelClick = () => {
    setShowCancelModal(true);
  };

  const confirmCancel = () => {
    setShowCancelModal(false);
    onBack();
  };

  const inputBase = "w-full px-3 py-2.5 bg-gray-50 border border-gray-200 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-[#D93030] focus:border-transparent focus:bg-white transition-colors text-sm text-gray-900 min-h-[44px]";
  const inputError = "w-full px-3 py-2.5 bg-gray-50 border border-red-500 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-red-500 focus:border-transparent focus:bg-white transition-colors text-sm text-gray-900 min-h-[44px]";

  return (
    <div className="flex-1 flex flex-col h-full overflow-hidden" style={{ backgroundColor: "#F7F4EF" }}>
      {/* Header */}
      <header className="bg-white border-b border-gray-200 sticky top-0 z-10">
        <div className="px-6 py-4 md:px-8 max-w-5xl mx-auto flex flex-col md:flex-row md:items-center justify-between gap-4">
          <div>
            <nav className="flex items-center text-xs text-gray-500 mb-2 font-medium">
              <button
                onClick={handleCancelClick}
                className="hover:text-gray-800 transition-colors flex items-center"
              >
                <ArrowLeft className="w-3.5 h-3.5 mr-1" />
                Treinamentos
              </button>
              <span className="mx-2">/</span>
              <span className="text-gray-900">{isEditing ? "Editar treinamento" : "Novo treinamento"}</span>
            </nav>
            <h1 className="text-2xl" style={{ color: "#1F2937" }}>
              {isEditing ? "Editar Treinamento" : "Cadastrar Novo Treinamento"}
            </h1>
            <p className="text-sm text-gray-500 mt-1">
              {isEditing ? "Atualize as informações do evento." : "Cadastre uma nova ação de capacitação para os lojistas."}
            </p>
          </div>

          <div className="flex items-center gap-3 self-start md:self-auto">
            <button
              type="button"
              onClick={handleCancelClick}
              className="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors focus:ring-2 focus:ring-offset-2 focus:ring-gray-200 focus:outline-none min-h-[44px]"
            >
              Cancelar
            </button>
            <button
              type="button"
              onClick={handleSubmit(onSubmit)}
              disabled={isSubmitting}
              className="px-4 py-2 text-sm font-medium text-white rounded-lg transition-all focus:ring-2 focus:ring-offset-2 focus:ring-[#D93030] focus:outline-none disabled:opacity-50 disabled:cursor-not-allowed flex items-center shadow-sm min-h-[44px]"
              style={{ backgroundColor: "#D93030" }}
            >
              {isSubmitting ? "Salvando..." : (isEditing ? "Salvar alterações" : "Salvar treinamento")}
            </button>
          </div>
        </div>
      </header>

      {/* Main Form Content */}
      <main className="flex-1 overflow-y-auto p-4 md:p-8">
        <form
          id="training-form"
          className="max-w-5xl mx-auto space-y-6"
          onSubmit={handleSubmit(onSubmit)}
        >
          {/* Seção 1: Informações principais */}
          <section className="bg-white p-6 md:p-8 rounded-lg shadow-sm border border-gray-200">
            <h2 className="text-lg font-semibold text-gray-900 mb-5 pb-3 border-b border-gray-200 flex items-center gap-2">
              <span className="w-1.5 h-5 rounded-full inline-block" style={{ backgroundColor: "#8B1A1A" }}></span>
              Informações principais
            </h2>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-x-6 gap-y-5">
              <div className="md:col-span-2">
                <label className="block text-sm font-medium text-gray-700 mb-1.5">
                  Tema do treinamento <span style={{ color: "#D93030" }}>*</span>
                </label>
                <input
                  type="text"
                  placeholder="Ex.: Treinamento de atendimento ao lojista"
                  {...register("tema", {
                    required: "Preencha este campo para continuar",
                  })}
                  className={errors.tema ? inputError : inputBase}
                />
                {errors.tema && (
                  <p className="mt-1.5 text-xs flex items-center" style={{ color: "#D93030" }}>
                    <Info className="w-3.5 h-3.5 mr-1" /> {errors.tema.message}
                  </p>
                )}
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1.5">
                  Categoria
                </label>
                <select
                  {...register("categoria")}
                  className={inputBase}
                >
                  <option value="" disabled>Selecione uma categoria</option>
                  <option value="palestra">Palestra</option>
                  <option value="treinamento">Treinamento</option>
                  <option value="workshop">Workshop</option>
                  <option value="capacitacao">Capacitação</option>
                  <option value="outro">Outro</option>
                </select>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1.5">
                  Segmento Alvo
                </label>
                <select
                  {...register("segmento")}
                  className={inputBase}
                >
                  <option value="Geral">Geral</option>
                  <option value="Vestuário">Vestuário</option>
                  <option value="Calçados">Calçados</option>
                  <option value="Cosméticos e Perfumaria">Cosméticos e Perfumaria</option>
                  <option value="Artigos Esportivos">Artigos Esportivos</option>
                  <option value="Tecnologia e Eletro">Tecnologia e Eletro</option>
                  <option value="Óticas e Joalheria">Óticas e Joalheria</option>
                  <option value="Alimentação">Alimentação</option>
                  <option value="Casa e Decoração">Casa e Decoração</option>
                  <option value="Serviços e Conveniência">Serviços e Conveniência</option>
                  <option value="Entretenimento">Entretenimento</option>
                </select>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1.5">
                  Modalidade
                </label>
                <select
                  {...register("modalidade")}
                  className={inputBase}
                >
                  <option value="" disabled>Selecione a modalidade</option>
                  <option value="presencial">Presencial</option>
                  <option value="online">Online</option>
                  <option value="hibrido">Híbrido</option>
                </select>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1.5">
                  Data do treinamento <span style={{ color: "#D93030" }}>*</span>
                </label>
                <div className="relative">
                  <Popover>
                    <PopoverTrigger asChild>
                      <button type="button" className="absolute inset-y-0 left-0 pl-3 flex items-center z-10 cursor-pointer group">
                        <CalendarIcon className="h-4 w-4 text-gray-400 group-hover:text-[#D93030] transition-colors" />
                      </button>
                    </PopoverTrigger>
                    <PopoverContent className="w-auto p-0" align="start">
                      <Calendar
                        mode="single"
                        selected={watch("data") ? new Date(watch("data") + "T12:00:00") : undefined}
                        onSelect={(d) => {
                          if (d) {
                            const dateStr = d.toISOString().split('T')[0];
                            setValue("data", dateStr, { shouldValidate: true });
                          }
                        }}
                      />
                    </PopoverContent>
                  </Popover>
                  <input
                    type="date"
                    {...register("data", {
                      required: "Preencha este campo para continuar",
                    })}
                    className={`pl-10 pr-3 ${errors.data ? inputError : inputBase} [&::-webkit-calendar-picker-indicator]:hidden`}
                  />
                </div>
                {errors.data && (
                  <p className="mt-1.5 text-xs flex items-center" style={{ color: "#D93030" }}>
                    <Info className="w-3.5 h-3.5 mr-1" /> {errors.data.message}
                  </p>
                )}
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1.5">
                    Horário de início
                  </label>
                  <div className="relative">
                    <Popover>
                      <PopoverTrigger asChild>
                        <button type="button" className="absolute inset-y-0 left-0 pl-3 flex items-center z-10 cursor-pointer group">
                          <Clock className="h-4 w-4 text-gray-400 group-hover:text-[#D93030] transition-colors" />
                        </button>
                      </PopoverTrigger>
                      <PopoverContent className="w-auto p-0" align="start">
                        <TimePickerContent 
                          value={watch("horarioInicio")} 
                          onChange={(val) => setValue("horarioInicio", val, { shouldValidate: true })} 
                        />
                      </PopoverContent>
                    </Popover>
                    <input
                      type="time"
                      {...register("horarioInicio")}
                      className={`pl-10 pr-3 ${inputBase} [&::-webkit-calendar-picker-indicator]:hidden`}
                    />
                  </div>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1.5">
                    Horário de fim
                  </label>
                  <div className="relative">
                    <Popover>
                      <PopoverTrigger asChild>
                        <button type="button" className="absolute inset-y-0 left-0 pl-3 flex items-center z-10 cursor-pointer group">
                          <Clock className="h-4 w-4 text-gray-400 group-hover:text-[#D93030] transition-colors" />
                        </button>
                      </PopoverTrigger>
                      <PopoverContent className="w-auto p-0" align="start">
                        <TimePickerContent 
                          value={watch("horarioFim")} 
                          onChange={(val) => setValue("horarioFim", val, { shouldValidate: true })} 
                        />
                      </PopoverContent>
                    </Popover>
                    <input
                      type="time"
                      {...register("horarioFim")}
                      className={`pl-10 pr-3 ${inputBase} [&::-webkit-calendar-picker-indicator]:hidden`}
                    />
                  </div>
                </div>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1.5">
                  Local
                </label>
                <div className="relative">
                  <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                    <MapPin className="h-4 w-4 text-gray-400" />
                  </div>
                  <input
                    type="text"
                    placeholder="Ex.: Sala de Reuniões 1, Auditório..."
                    {...register("local")}
                    className={`pl-10 pr-3 ${inputBase}`}
                  />
                </div>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1.5">
                  Status inicial
                </label>
                <select
                  {...register("status")}
                  className={inputBase}
                >
                  <option value="agendado">Agendado</option>
                  <option value="realizado">Realizado</option>
                  <option value="cancelado">Cancelado</option>
                </select>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1.5">
                  Capacidade máxima
                </label>
                <input
                  type="number"
                  min="1"
                  placeholder="Ex.: 50"
                  {...register("capacidadeMaxima")}
                  className={inputBase}
                />
              </div>
            </div>
          </section>

          {/* Seção 2: Conteúdo e detalhamento */}
          <section className="bg-white p-6 md:p-8 rounded-lg shadow-sm border border-gray-200">
            <h2 className="text-lg font-semibold text-gray-900 mb-5 pb-3 border-b border-gray-200 flex items-center gap-2">
              <span className="w-1.5 h-5 rounded-full inline-block" style={{ backgroundColor: "#8B1A1A" }}></span>
              Conteúdo e detalhamento
            </h2>

            <div className="space-y-5">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1.5">
                  Conteúdo / descrição <span style={{ color: "#D93030" }}>*</span>
                </label>
                <textarea
                  rows={4}
                  placeholder="Descreva os principais tópicos, orientações ou informações da capacitação"
                  {...register("conteudo", {
                    required: "Preencha este campo para continuar",
                  })}
                  className={`resize-none ${errors.conteudo ? inputError : inputBase}`}
                />
                {errors.conteudo && (
                  <p className="mt-1.5 text-xs flex items-center" style={{ color: "#D93030" }}>
                    <Info className="w-3.5 h-3.5 mr-1" /> {errors.conteudo.message}
                  </p>
                )}
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1.5">
                  Objetivo do treinamento
                </label>
                <textarea
                  rows={2}
                  placeholder="Qual o resultado esperado ao final deste treinamento?"
                  {...register("objetivo")}
                  className={`resize-none ${inputBase}`}
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1.5">
                  Observações internas
                </label>
                <textarea
                  rows={2}
                  placeholder="Informações visíveis apenas para a equipe administrativa"
                  {...register("observacoes")}
                  className={`resize-none ${inputBase}`}
                />
              </div>

              {/* Anexos */}
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1.5">
                  Material de apoio (Opcional)
                </label>
                <div className="mt-1 flex justify-center px-6 pt-5 pb-6 border-2 border-gray-200 border-dashed rounded-lg hover:bg-gray-50 transition-colors cursor-pointer group">
                  <div className="space-y-2 text-center">
                    <div className="mx-auto h-10 w-10 bg-gray-100 rounded-full flex items-center justify-center group-hover:bg-red-50 transition-colors">
                      <UploadCloud className="h-5 w-5 text-gray-400 group-hover:text-[#D93030]" />
                    </div>
                    <div className="flex text-sm text-gray-600 justify-center">
                      <span className="relative cursor-pointer bg-transparent rounded-md font-medium hover:opacity-80 focus-within:outline-none" style={{ color: "#D93030" }}>
                        Faça upload de um arquivo
                      </span>
                      <p className="pl-1">ou arraste e solte</p>
                    </div>
                    <p className="text-xs text-gray-500">
                      PDF, PPTX, DOCX até 10MB (Apenas materiais do curso)
                    </p>
                  </div>
                </div>
              </div>
            </div>
          </section>

          {/* Seção 3: Organização */}
          <section className="bg-white p-6 md:p-8 rounded-lg shadow-sm border border-gray-200">
            <h2 className="text-lg font-semibold text-gray-900 mb-5 pb-3 border-b border-gray-200 flex items-center gap-2">
              <span className="w-1.5 h-5 rounded-full inline-block" style={{ backgroundColor: "#8B1A1A" }}></span>
              Organização
            </h2>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-x-6 gap-y-5">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1.5">
                  Responsável pelo cadastro
                </label>
                <input
                  type="text"
                  readOnly
                  {...register("responsavel")}
                  className="w-full px-3 py-2.5 bg-gray-100 border border-gray-200 text-gray-500 rounded-lg shadow-sm focus:outline-none text-sm cursor-not-allowed min-h-[44px]"
                />
                <p className="mt-1 text-xs text-gray-400">
                  Preenchido automaticamente com seu usuário logado.
                </p>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1.5">
                  Área responsável
                </label>
                <select
                  {...register("areaResponsavel")}
                  className={inputBase}
                >
                  <option value="" disabled>Selecione a área</option>
                  <option value="relacionamento">Relacionamento</option>
                  <option value="marketing">Marketing</option>
                  <option value="comercial">Comercial</option>
                  <option value="rh">Recursos Humanos</option>
                  <option value="operacoes">Operações</option>
                </select>
              </div>

              <div className="md:col-span-2">
                <label className="block text-sm font-medium text-gray-700 mb-1.5">
                  Público-alvo
                </label>
                <input
                  type="text"
                  placeholder="Ex.: Gerentes, Supervisores de Venda, Atendentes..."
                  {...register("publicoAlvo")}
                  className={inputBase}
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1.5">
                  Palavras-chave / Tags (Opcional)
                </label>
                <input
                  type="text"
                  placeholder="Ex.: vendas, liderança, sistema"
                  {...register("tags")}
                  className={inputBase}
                />
              </div>

              <div>
                <div className="flex items-center justify-between mt-6 p-3 bg-gray-50 rounded-lg border border-gray-200">
                  <div>
                    <label className="text-sm font-medium text-gray-800">
                      Treinamento recorrente?
                    </label>
                    <p className="text-xs text-gray-500">
                      Ative se esta capacitação acontecer com frequência.
                    </p>
                  </div>
                  <label className="relative inline-flex items-center cursor-pointer">
                    <input
                      type="checkbox"
                      className="sr-only peer"
                      {...register("recorrente")}
                    />
                    <div className="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-2 peer-focus:ring-red-100 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-[#D93030]"></div>
                  </label>
                </div>

                {isRecorrente && (
                  <div className="mt-3 animate-in fade-in slide-in-from-top-2 duration-200">
                    <label className="block text-sm font-medium text-gray-700 mb-1.5">
                      Frequência
                    </label>
                    <select
                      {...register("frequencia")}
                      className={inputBase}
                    >
                      <option value="" disabled>Selecione a frequência</option>
                      <option value="diaria">Diária</option>
                      <option value="semanal">Semanal</option>
                      <option value="quinzenal">Quinzenal</option>
                      <option value="mensal">Mensal</option>
                      <option value="semestral">Semestral</option>
                      <option value="anual">Anual</option>
                    </select>
                  </div>
                )}
              </div>
            </div>
          </section>

          <div className="pb-10"></div>
        </form>
      </main>

      {/* Cancel Confirmation Modal */}
      {showCancelModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm animate-in fade-in duration-200">
          <div className="bg-white rounded-lg shadow-xl max-w-sm w-full mx-4 overflow-hidden">
            <div className="p-6">
              <h3 className="text-lg font-semibold text-gray-900 mb-2">
                Cancelar cadastro?
              </h3>
              <p className="text-sm text-gray-500 mb-6">
                Todas as informações preenchidas serão perdidas. Deseja realmente voltar para a lista de treinamentos?
              </p>
              <div className="flex items-center justify-end gap-3">
                <button
                  onClick={() => setShowCancelModal(false)}
                  className="px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-100 rounded-lg transition-colors focus:outline-none focus:ring-2 focus:ring-gray-200 min-h-[44px]"
                >
                  Continuar editando
                </button>
                <button
                  onClick={confirmCancel}
                  className="px-4 py-2 text-sm font-medium text-white rounded-lg transition-colors hover:opacity-90 focus:outline-none focus:ring-2 focus:ring-red-200 min-h-[44px]"
                  style={{ backgroundColor: "#D93030" }}
                >
                  Sim, cancelar
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
