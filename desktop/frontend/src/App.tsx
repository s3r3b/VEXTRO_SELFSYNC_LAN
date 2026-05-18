import { useState, useEffect, useRef } from 'react';
import { GetChatHistory, SendChatMessage, GetSystemStatus, GetActiveNodes, SelectAndSendFile } from '../wailsjs/go/main/App';

interface ChatMessage { timestamp: string; senderId: string; content: string; }

function App() {
    const [localId, setLocalId] = useState<string>("Łączenie...");
    const [messages, setMessages] = useState<ChatMessage[]>([]);
    const [nodes, setNodes] = useState<Record<string, string>>({});
    const [input, setInput] = useState("");
    const chatEndRef = useRef<HTMLDivElement>(null);

    const fetchData = async () => {
        try {
            const status = await GetSystemStatus();
            setLocalId(status);

            const chat = await GetChatHistory();
            if (chat && chat !== "[]" && chat !== "ERROR") setMessages(JSON.parse(chat));

            const activeNodes = await GetActiveNodes();
            if (activeNodes && activeNodes !== "{}") setNodes(JSON.parse(activeNodes));
        } catch (e) {
            console.error("IPC Sync Error", e);
        }
    };

    useEffect(() => {
        fetchData();
        const interval = setInterval(fetchData, 1500);
        return () => clearInterval(interval);
    }, []);

    useEffect(() => {
        chatEndRef.current?.scrollIntoView({ behavior: "smooth" });
    }, [messages]);

    const handleSendMsg = async () => {
        if (!input.trim()) return;
        await SendChatMessage(input);
        setInput("");
        fetchData();
    };

    const handleSendFile = async (targetId: string) => {
        const result = await SelectAndSendFile(targetId);
        if (result === "TRANSFER_STARTED") {
            alert(`Zlecono wysłanie pliku do: ${targetId}`);
        } else if (result !== "CANCELLED") {
            alert(`Błąd transferu: ${result}`);
        }
    };

    return (
        <div style={{ display: 'flex', height: '100vh', backgroundColor: '#0f0f11', color: '#e0e0e0', fontFamily: 'system-ui' }}>
            {/* RADAR / PANEL BOCZNY */}
            <div style={{ width: '280px', borderRight: '1px solid #222', backgroundColor: '#16161a', display: 'flex', flexDirection: 'column' }}>
                <div style={{ padding: '20px', borderBottom: '1px solid #222' }}>
                    <h2 style={{ margin: 0, fontSize: '18px', color: '#00ffcc' }}>RADAR LAN</h2>
                    <div style={{ fontSize: '12px', color: '#888', marginTop: '10px' }}>Twoje ID: <br/><span style={{color: '#fff'}}>{localId}</span></div>
                </div>
                <div style={{ flex: 1, padding: '15px', overflowY: 'auto' }}>
                    <h3 style={{ fontSize: '11px', color: '#666', letterSpacing: '1px' }}>AKTYWNE WĘZŁY ({Object.keys(nodes).length})</h3>
                    {Object.keys(nodes).map(id => (
                        <div key={id} style={{ 
                            backgroundColor: '#222', padding: '12px', borderRadius: '8px', 
                            marginBottom: '10px', display: 'flex', flexDirection: 'column', gap: '8px'
                        }}>
                            <div style={{ fontSize: '12px', fontWeight: 'bold' }}>{id === localId ? `${id} (Ty)` : id}</div>
                            {id !== localId && (
                                <button onClick={() => handleSendFile(id)} style={{
                                    backgroundColor: '#333', color: '#fff', border: 'none', 
                                    padding: '6px', borderRadius: '4px', cursor: 'pointer', fontSize: '11px'
                                }} onMouseOver={(e) => e.currentTarget.style.backgroundColor = '#005f99'}
                                   onMouseOut={(e) => e.currentTarget.style.backgroundColor = '#333'}>
                                    📄 Wyślij plik
                                </button>
                            )}
                        </div>
                    ))}
                </div>
            </div>

            {/* GŁÓWNY PANEL CZATU */}
            <div style={{ flex: 1, display: 'flex', flexDirection: 'column' }}>
                <div style={{ flex: 1, overflowY: 'auto', padding: '20px', display: 'flex', flexDirection: 'column', gap: '15px' }}>
                    {messages.length === 0 ? <div style={{ textAlign: 'center', color: '#555' }}>Brak wiadomości.</div> : 
                        messages.map((msg, idx) => {
                            const isMe = msg.senderId === localId;
                            return (
                                <div key={idx} style={{ 
                                    alignSelf: isMe ? 'flex-end' : 'flex-start', backgroundColor: isMe ? '#005f99' : '#222226', 
                                    padding: '12px 16px', borderRadius: '12px', borderBottomRightRadius: isMe ? '2px' : '12px',
                                    borderBottomLeftRadius: isMe ? '12px' : '2px', maxWidth: '75%'
                                }}>
                                    {!isMe && <div style={{ fontSize: '10px', color: '#00ffcc', marginBottom: '6px' }}>{msg.senderId}</div>}
                                    <div style={{ fontSize: '14px' }}>{msg.content}</div>
                                    <div style={{ fontSize: '9px', color: '#aaa', marginTop: '6px', textAlign: 'right' }}>
                                        {new Date(msg.timestamp).toLocaleTimeString()}
                                    </div>
                                </div>
                            );
                        })
                    }
                    <div ref={chatEndRef} />
                </div>
                
                {/* WPROWADZANIE */}
                <div style={{ padding: '20px', borderTop: '1px solid #222', display: 'flex', gap: '10px' }}>
                    <input value={input} onChange={(e) => setInput(e.target.value)} onKeyDown={(e) => e.key === 'Enter' && handleSendMsg()}
                        placeholder="Napisz..." style={{ flex: 1, padding: '14px', borderRadius: '8px', border: '1px solid #333', backgroundColor: '#0f0f11', color: 'white', outline: 'none' }} />
                    <button onClick={handleSendMsg} style={{ padding: '0 25px', borderRadius: '8px', border: 'none', backgroundColor: '#00ffcc', color: '#000', fontWeight: 'bold', cursor: 'pointer' }}>WYŚLIJ</button>
                </div>
            </div>
        </div>
    );
}

export default App;