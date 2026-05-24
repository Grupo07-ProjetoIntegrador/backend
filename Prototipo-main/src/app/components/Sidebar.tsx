import { ImageWithFallback } from "./figma/ImageWithFallback";
import figmaAsset from "../../imports/logo_2024_(1).png";
import {
  LayoutDashboard,
  Store,
  GraduationCap,
  ShieldCheck,
  Wrench,
  AlertTriangle,
  Megaphone,
  Briefcase,
  Building2,
  BarChart2,
  ChevronLeft,
  LogOut,
} from "lucide-react";

const navItems = [
  { label: "Dashboard", icon: LayoutDashboard, active: false },
  { label: "Lojistas", icon: Store, active: false },
  { label: "Treinamentos", icon: GraduationCap, active: true },
  { label: "Seguros", icon: ShieldCheck, active: false },
  { label: "Manutenção", icon: Wrench, active: false },
  { label: "Sinistros", icon: AlertTriangle, active: false },
  { label: "Marketing", icon: Megaphone, active: false },
  { label: "Comercial", icon: Briefcase, active: false },
  { label: "Institucional", icon: Building2, active: false },
  { label: "Relatórios", icon: BarChart2, active: false },
];

interface SidebarProps {
  isOpen: boolean;
  onClose: () => void;
}

export function Sidebar({ isOpen, onClose }: SidebarProps) {
  return (
    <div className="h-full flex-shrink-0" style={{ backgroundColor: "#8B1A1A" }}>
      <aside className="w-64 h-full flex flex-col flex-shrink-0" style={{ backgroundColor: "#8B1A1A" }}>
        {/* Logo Area */}
        <div className="px-5 py-5 border-b border-white/15">
          <div className="flex items-center gap-3">
            <div className="relative flex items-center justify-center w-12 h-12 flex-shrink-0">
              <div className="absolute inset-0 bg-yellow-400 opacity-60 blur-xl scale-125 rounded-full"></div>
              <ImageWithFallback
                src={figmaAsset}
                className="w-full h-full object-contain relative z-10"
                alt="Flamboyant logo"
              />
            </div>
            <div className="flex-1">
              <p className="text-xs text-white/60 leading-tight">
                Flamboyant
              </p>
              <p className="text-xs font-semibold text-white leading-tight">JP Mall</p>
            </div>
            <button
              onClick={onClose}
              className="p-1 rounded-lg hover:bg-white/10 transition-colors group"
              aria-label="Fechar menu"
            >
              <ChevronLeft className={`w-6 h-6 text-white/70 transition-transform duration-500 ease-in-out group-active:scale-95 ${!isOpen ? '-rotate-180' : 'rotate-0'}`} />
            </button>
          </div>
        </div>

        {/* Navigation */}
        <nav className="flex-1 px-3 py-4 space-y-0.5">
          {navItems.map((item) => {
            const Icon = item.icon;
            return (
              <button
                key={item.label}
                className={`w-full flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm transition-colors cursor-pointer text-left ${
                  item.active
                    ? "bg-white/15 text-white"
                    : "text-white/70 hover:bg-white/10 hover:text-white"
                }`}
              >
                <Icon
                  className="w-4 h-4 flex-shrink-0"
                  strokeWidth={item.active ? 2 : 1.5}
                />
                <span className={item.active ? "font-medium" : ""}>
                  {item.label}
                </span>
              </button>
            );
          })}
        </nav>

        {/* User Profile Card */}
        <div
          className="mx-3 mb-4 p-3 rounded-xl flex items-center justify-between bg-white/10"
        >
          <div className="flex items-center gap-3 min-w-0 pr-2">
            <div
              className="w-9 h-9 rounded-full flex items-center justify-center text-xs font-semibold flex-shrink-0"
              style={{ backgroundColor: "#D93030", color: "#ffffff" }}
            >
              AM
            </div>
            <div className="min-w-0">
              <p className="text-xs font-semibold text-white truncate">
                Admin Master
              </p>
              <p className="text-xs text-white/60 truncate">
                admin@flamboyant.com
              </p>
            </div>
          </div>
          <button
            className="text-white/50 hover:text-white transition-colors flex-shrink-0"
            title="Sair"
          >
            <LogOut className="w-5 h-5" />
          </button>
        </div>
      </aside>
    </div>
  );
}
