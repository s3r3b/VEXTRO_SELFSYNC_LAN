import React, { useEffect, useState } from 'react';
import { GetSystemStatus } from '../wailsjs/go/main/App';

function App() {
  const [daemonStatus, setDaemonStatus] = useState<string>("NAWIĄZYWANIE POŁĄCZENIA...");

  useEffect(() => {
    GetSystemStatus()
      .then((result) => setDaemonStatus(result))
      .catch((err) => setDaemonStatus("BŁĄD IPC: " + err));
  }, []);

  return (
    <div className="flex-1 flex flex-col items-center justify-center p-8 bg-white/10 w-full h-full">
      {/* bg-white/10 to 10% białego filtru na całe okno, który delikatnie rozjaśni ciemne tło Twojego IDE */}
      
      {/* Sekcja Nagłówka */}
      <div className="text-center mb-12 select-none">
        {/* Zmieniamy tekst na biały, aby był w końcu czytelny! */}
        <h1 className="text-5xl font-light tracking-[0.3em] text-white/90 mb-2 drop-shadow-md">VEXTRO</h1>
        <p className="text-[var(--color-neon-orange)] font-mono text-xs uppercase tracking-widest font-bold opacity-100 drop-shadow-sm">
          SelfSync LAN Protocol
        </p>
      </div>
      
      {/* Główny Kontener / Karta Statusu */}
      {/* Zwiększyliśmy przezroczystość (bg-white/15 zamiast 40) i nałożyliśmy mniejszy blur */}
      <div className="w-full max-w-lg bg-white/15 backdrop-blur-lg border border-white/30 rounded-2xl p-6 shadow-[0_8px_32px_rgba(0,0,0,0.2)] flex flex-col gap-6 relative overflow-hidden">
        
        <div className="absolute top-0 left-0 w-full h-[1px] bg-gradient-to-r from-transparent via-white/70 to-transparent"></div>
        
        <div className="flex items-center justify-between border-b border-white/20 pb-4">
          <span className="text-xs text-white/80 uppercase tracking-widest font-bold">Status Środowiska</span>
          <div className="flex items-center gap-2 bg-black/20 px-3 py-1 rounded-full border border-white/20 shadow-sm backdrop-blur-md">
            <span className="w-2 h-2 rounded-full bg-[var(--color-neon-orange)] animate-pulse shadow-[0_0_8px_var(--color-neon-orange)]"></span>
            <span className="text-[10px] font-mono text-white/90 tracking-wider mt-[1px] font-bold">TRYB: IPC BINDING</span>
          </div>
        </div>
        
        <div className="flex flex-col items-center justify-center py-8">
          <svg className="w-8 h-8 text-white/50 mb-4 animate-spin-slow" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <p className="text-xs text-white/80 font-mono text-center leading-relaxed tracking-wide font-medium">
            ODCZYT Z SILNIKA GO:<br/>
            <span className="text-[var(--color-neon-orange)] font-black tracking-widest mt-2 block drop-shadow-md">
              {daemonStatus}
            </span>
          </p>
        </div>

      </div>
    </div>
  )
}

export default App;