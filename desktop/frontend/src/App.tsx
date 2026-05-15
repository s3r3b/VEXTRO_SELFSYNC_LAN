import { useEffect, useState } from 'react';
import { GetSystemStatus, GetActiveNodes } from '../wailsjs/go/main/App';

// Deterministyczny generator awatarów na podstawie ID urządzenia
const generateAvatar = (id: string) => {
  if (!id || !id.startsWith('VXT-')) return { color: '#333', initials: '??' };
  
  const hex = id.replace('VXT-', '');
  const color = `#${hex.slice(0, 6).padEnd(6, '8')}`;
  const initials = hex.slice(0, 2).toUpperCase();
  
  return { color, initials };
};

// Struktura węzła odebrana z demona Go
interface NodeInfo {
  device_id: string;
  status: string;
  port: string;
  ip: string;
  last_seen: number;
}

function App() {
  const [deviceID, setDeviceID] = useState<string>("ŁĄCZENIE Z RDZENIEM...");
  const [activeNodes, setActiveNodes] = useState<NodeInfo[]>([]);

  useEffect(() => {
    const fetchData = async () => {
      try {
        // Odpytanie o własny status
        const resultStatus = await GetSystemStatus();
        setDeviceID(resultStatus);

        // Odpytanie o inne węzły (Radar)
        const nodesJsonStr = await GetActiveNodes();
        try {
          const parsedNodes = JSON.parse(nodesJsonStr);
          setActiveNodes(Array.isArray(parsedNodes) ? parsedNodes : []);
        } catch (e) {
          console.error("Błąd dekodowania radaru JSON:", e);
        }
      } catch (err) {
        setDeviceID("BŁĄD IPC");
      }
    };
    
    fetchData();
    // Odświeżanie całego interfejsu (i radaru) co 3 sekundy
    const interval = setInterval(fetchData, 3000);
    return () => clearInterval(interval);
  }, []);

  const avatar = generateAvatar(deviceID);
  const isConnected = deviceID.startsWith('VXT-');

  return (
    <div className="flex-1 flex flex-row items-stretch justify-center p-8 gap-8 w-full h-full select-none">
      
      {/* Lewa Sekcja: Osobista tożsamość */}
      <div className="flex flex-col items-center justify-center w-full max-w-md gap-10">
        
        {/* Nagłówek Typograficzny */}
        <div className="text-center">
          <h1 className="text-4xl font-light tracking-[0.4em] text-white/90 mb-2 drop-shadow-lg">VEXTRO</h1>
          <p className="text-[var(--color-neon-orange)] font-mono text-[10px] uppercase tracking-[0.3em] font-bold">
            SelfSync LAN Protocol
          </p>
        </div>
        
        {/* Karta Identyfikacji (Smoked Glass) */}
        <div className="w-full bg-white/10 backdrop-blur-xl border border-white/20 rounded-3xl p-6 shadow-[0_8px_32px_rgba(0,0,0,0.3)] flex flex-col gap-6 relative overflow-hidden">
          
          <div className="absolute top-0 left-0 w-full h-[1px] bg-gradient-to-r from-transparent via-white/50 to-transparent"></div>
          
          <div className="flex items-center justify-between border-b border-white/10 pb-4">
            <span className="text-[10px] text-white/60 uppercase tracking-widest font-bold">Tożsamość Węzła</span>
            <div className="flex items-center gap-2 bg-black/30 px-3 py-1 rounded-full border border-white/10 shadow-inner">
              <span className={`w-2 h-2 rounded-full ${isConnected ? 'bg-[var(--color-neon-orange)] animate-pulse shadow-[0_0_8px_var(--color-neon-orange)]' : 'bg-red-500'}`}></span>
              <span className="text-[9px] font-mono text-white/80 tracking-wider mt-[1px] font-bold">
                {isConnected ? 'mDNS BROADCAST' : 'OFFLINE'}
              </span>
            </div>
          </div>
          
          {/* Blok Awatara i ID */}
          <div className="flex items-center gap-5 bg-black/20 p-5 rounded-2xl border border-white/5">
            <div 
              className="w-14 h-14 rounded-xl flex items-center justify-center text-xl font-black text-white/90 shadow-lg border border-white/20"
              style={{ backgroundColor: avatar.color }}
            >
              {avatar.initials}
            </div>
            
            <div className="flex flex-col">
              <span className="text-[9px] font-mono text-white/40 uppercase tracking-widest mb-1">Local Device ID</span>
              <span className="text-xl font-mono font-bold tracking-widest text-white/90 drop-shadow-md">
                {deviceID}
              </span>
            </div>
          </div>
        </div>
      </div>

      {/* Prawa Sekcja: Radar mDNS (Topologia) */}
      <div className="w-full max-w-xs bg-white/5 backdrop-blur-lg border border-white/10 rounded-3xl p-5 shadow-[0_8px_32px_rgba(0,0,0,0.2)] flex flex-col relative overflow-hidden h-[80vh] self-center">
        
        <div className="flex items-center justify-between border-b border-white/10 pb-4 mb-4">
          <span className="text-[10px] text-white/60 uppercase tracking-widest font-bold">Radar LAN (mDNS)</span>
          <span className="text-[10px] font-mono text-[var(--color-neon-orange)] tracking-wider font-bold">
            WĘZŁY: {activeNodes.length}
          </span>
        </div>

        {/* Lista Kontaktów */}
        <div className="flex-1 overflow-y-auto flex flex-col gap-3 custom-scrollbar pr-1">
          {activeNodes.length === 0 ? (
            <div className="flex-1 flex flex-col items-center justify-center text-center opacity-40">
              <div className="w-6 h-6 rounded-full border-2 border-[var(--color-neon-orange)] border-t-transparent animate-spin mb-3"></div>
              <span className="text-[9px] font-mono text-white uppercase tracking-widest">Skanowanie Sieci...</span>
            </div>
          ) : (
            activeNodes.map((node) => {
              const nodeAvatar = generateAvatar(node.device_id);
              return (
                <div key={node.device_id} className="flex items-center justify-between bg-black/40 p-3 rounded-2xl border border-white/5 hover:bg-black/60 transition-colors cursor-pointer group">
                  <div className="flex items-center gap-3">
                    <div 
                      className="w-10 h-10 rounded-lg flex items-center justify-center text-sm font-black text-white/90 shadow-md border border-white/10 group-hover:scale-105 transition-transform"
                      style={{ backgroundColor: nodeAvatar.color }}
                    >
                      {nodeAvatar.initials}
                    </div>
                    <div className="flex flex-col">
                      <span className="text-[10px] font-mono font-bold tracking-widest text-white/90">{node.device_id}</span>
                      <span className="text-[8px] font-mono text-[var(--color-neon-orange)] tracking-widest">{node.ip}</span>
                    </div>
                  </div>
                  <div className="w-2 h-2 rounded-full bg-[var(--color-neon-orange)] shadow-[0_0_8px_var(--color-neon-orange)]"></div>
                </div>
              )
            })
          )}
        </div>

      </div>

    </div>
  )
}

export default App;