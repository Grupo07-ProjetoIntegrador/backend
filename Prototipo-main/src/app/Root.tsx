import { useState } from "react";
import { Sidebar } from "./components/Sidebar";
import { Menu } from "lucide-react";
import { Toaster } from "sonner";
import { ImageWithFallback } from "./components/figma/ImageWithFallback";
import figmaAsset from "../imports/logo_2024-1.png";
import { Outlet, useLocation } from "react-router";

export default function Root() {
  const [sidebarOpen, setSidebarOpen] = useState(true);
  const location = useLocation();

  return (
    <div className="flex min-h-screen w-full overflow-hidden relative" style={{ backgroundColor: "#F7F4EF" }}>
      {/* Mobile overlay */}
      <div 
        className={`fixed inset-0 bg-black/50 z-20 md:hidden transition-opacity duration-500 ease-in-out ${
          sidebarOpen ? "opacity-100" : "opacity-0 pointer-events-none"
        }`} 
        onClick={() => setSidebarOpen(false)}
      />

      {/* Sidebar container */}
      <div 
        className={`fixed inset-y-0 left-0 z-30 transform md:relative md:translate-x-0 transition-[transform,width] duration-500 ease-in-out ${
          sidebarOpen ? "translate-x-0 w-64 md:w-64" : "-translate-x-full w-64 md:w-0"
        } overflow-hidden shadow-xl md:shadow-none`}
        style={{ backgroundColor: "#8B1A1A" }}
      >
        <Sidebar isOpen={sidebarOpen} onClose={() => setSidebarOpen(false)} />
      </div>

      <div className="flex-1 flex flex-col h-screen overflow-y-auto w-full relative">
        <Toaster position="top-right" richColors />
        
        {/* Mobile Header */}
        <div className="md:hidden p-4 bg-white border-b border-gray-200 flex items-center sticky top-0 z-10 shadow-sm shrink-0">
          <button
            onClick={() => setSidebarOpen(true)}
            className="p-2 -ml-2 rounded-lg hover:bg-gray-100 transition-colors min-h-[44px] min-w-[44px] flex items-center justify-center group"
            aria-label="Abrir menu"
          >
            <Menu className={`w-5 h-5 text-gray-600 transition-transform duration-500 ease-in-out group-active:scale-95 ${sidebarOpen ? 'rotate-180' : 'rotate-0'}`} />
          </button>
          <div className="ml-2 relative flex items-center justify-center h-8 w-auto">
            <ImageWithFallback src={figmaAsset} alt="Flamboyant" className="h-full w-auto object-contain relative z-10 grayscale brightness-0 opacity-80 hover:opacity-100 transition-opacity" />
          </div>
        </div>

        {/* Desktop Menu Button (when sidebar closed) */}
        <div 
          className={`hidden md:flex absolute top-0 left-0 w-full items-center shrink-0 transition-opacity duration-500 ease-in-out z-10 ${
            sidebarOpen ? "opacity-0 pointer-events-none" : "opacity-100 pointer-events-auto"
          }`}
        >
          <div className="p-4 flex items-center w-full h-[76px]">
            <button
              onClick={() => setSidebarOpen(true)}
              className="p-2 rounded-lg hover:bg-gray-100 transition-colors bg-white shadow-sm border border-gray-200 min-h-[44px] min-w-[44px] flex items-center justify-center mr-4 group"
              aria-label="Abrir menu"
            >
              <Menu className={`w-5 h-5 text-gray-600 transition-transform duration-500 ease-in-out group-active:scale-95 ${sidebarOpen ? 'rotate-180' : 'rotate-0'}`} />
            </button>
            <div className="relative flex items-center justify-center h-10 w-auto">
              <ImageWithFallback src={figmaAsset} alt="Flamboyant" className="h-full w-auto object-contain relative z-10 grayscale brightness-0 opacity-80 hover:opacity-100 transition-opacity" />
            </div>
          </div>
        </div>
        
        <div className="md:pt-[76px] flex-1 flex flex-col w-full h-full">
          <Outlet />
        </div>
      </div>
    </div>
  );
}