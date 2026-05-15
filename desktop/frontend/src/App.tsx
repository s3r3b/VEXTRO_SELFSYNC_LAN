import { useState, useEffect, useRef } from 'react';
import { GetChatHistory, SendChatMessage, GetSystemStatus, GetActiveNodes } from '../wailsjs/go/main/App';

// Interfejs zgodny ze strukturą ChatMessage z Daemona (chat.go)
interface ChatMessage {
    timestamp: string;
    senderId: string;
    content: string;
}

function App() {
    const [localId, setLocalId] = useState<string>("Łączenie z Daemonem...");
    const [messages, setMessages] = useState<ChatMessage[]>([]);
    const [input, setInput] = useState("");
    const chatEndRef = useRef<HTMLDivElement>(null);

    // Inicjalizacja i pobranie własnego DeviceID
    const fetchStatus = async () => {
        try {
            const res = await GetSystemStatus();
            setLocalId(res);
        } catch (e) {
            console.error("Błąd połączenia IPC:", e);
        }
    };

    // Pobieranie historii czatu (JSON -> Array)
    const fetchChat = async () => {
        try {
            const res = await GetChatHistory();
            if (res && res !== "[]" && res !== "ERROR") {
                setMessages(JSON.parse(res));
            } else if (res === "[]") {
                setMessages([]);
            }
        } catch (e) {
            console.error("Błąd parsowania czatu:", e);
        }
    };

    // Efekt uruchamiany przy starcie
    useEffect(() => {
        fetchStatus();
        fetchChat();
        
        // MVP: Polling co 1.5 sekundy dla synchronizacji wiadomości
        const interval = setInterval(fetchChat, 1500);
        return () => clearInterval(interval);
    }, []);

    // Automatyczne przewijanie do najnowszej wiadomości
    useEffect(() => {
        chatEndRef.current?.scrollIntoView({ behavior: "smooth" });
    }, [messages]);

    const handleSend = async () => {
        if (!input.trim()) return;
        await SendChatMessage(input);
        setInput("");
        fetchChat(); // Natychmiastowe odświeżenie po wysłaniu
    };

    return (
        <div style={{ 
            display: 'flex', flexDirection: 'column', height: '100vh', 
            backgroundColor: '#0f0f11', color: '#e0e0e0', fontFamily: 'system-ui, sans-serif' 
        }}>
            {/* GÓRNY PASEK (HEADER) */}
            <header style={{ 
                padding: '15px 20px', borderBottom: '1px solid #222', 
                backgroundColor: '#16161a', display: 'flex', justifyContent: 'space-between', alignItems: 'center'
            }}>
                <h2 style={{ margin: 0, fontSize: '18px', fontWeight: 600, letterSpacing: '1px' }}>VEXTRO <span style={{color: '#00ffcc'}}>LAN</span></h2>
                <div style={{ fontSize: '12px', color: '#888', backgroundColor: '#222', padding: '5px 10px', borderRadius: '6px' }}>
                    ID: <span style={{color: '#fff'}}>{localId}</span>
                </div>
            </header>

            {/* GŁÓWNY PANEL CZATU */}
            <div style={{ 
                flex: 1, overflowY: 'auto', padding: '20px', 
                display: 'flex', flexDirection: 'column', gap: '15px' 
            }}>
                {messages.length === 0 ? (
                    <div style={{ textAlign: 'center', color: '#555', marginTop: '20px' }}>Brak wiadomości w historii LAN.</div>
                ) : (
                    messages.map((msg, idx) => {
                        const isMe = msg.senderId === localId;
                        return (
                            <div key={idx} style={{ 
                                alignSelf: isMe ? 'flex-end' : 'flex-start', 
                                backgroundColor: isMe ? '#005f99' : '#222226', 
                                padding: '12px 16px', borderRadius: '12px', 
                                borderBottomRightRadius: isMe ? '2px' : '12px',
                                borderBottomLeftRadius: isMe ? '12px' : '2px',
                                maxWidth: '75%', boxShadow: '0 4px 6px rgba(0,0,0,0.1)'
                            }}>
                                {!isMe && (
                                    <div style={{ fontSize: '10px', color: '#00ffcc', marginBottom: '6px', fontWeight: 'bold' }}>
                                        {msg.senderId}
                                    </div>
                                )}
                                <div style={{ fontSize: '14px', lineHeight: '1.4' }}>{msg.content}</div>
                                <div style={{ fontSize: '9px', color: isMe ? '#99d6ff' : '#666', marginTop: '6px', textAlign: 'right' }}>
                                    {new Date(msg.timestamp).toLocaleTimeString([], {hour: '2-digit', minute:'2-digit', second:'2-digit'})}
                                </div>
                            </div>
                        );
                    })
                )}
                <div ref={chatEndRef} />
            </div>

            {/* PANEL WPROWADZANIA */}
            <div style={{ 
                padding: '20px', borderTop: '1px solid #222', backgroundColor: '#16161a',
                display: 'flex', gap: '10px'
            }}>
                <input
                    type="text"
                    value={input}
                    onChange={(e) => setInput(e.target.value)}
                    onKeyDown={(e) => e.key === 'Enter' && handleSend()}
                    placeholder="Napisz wiadomość w sieci LAN..."
                    style={{ 
                        flex: 1, padding: '14px 20px', borderRadius: '8px', border: '1px solid #333', 
                        backgroundColor: '#0f0f11', color: 'white', outline: 'none', fontSize: '14px'
                    }}
                />
                <button 
                    onClick={handleSend} 
                    style={{ 
                        padding: '0 25px', borderRadius: '8px', border: 'none', 
                        backgroundColor: '#00ffcc', color: '#000', fontWeight: 'bold', cursor: 'pointer', transition: '0.2s'
                    }}
                    onMouseOver={(e) => e.currentTarget.style.backgroundColor = '#00ccaa'}
                    onMouseOut={(e) => e.currentTarget.style.backgroundColor = '#00ffcc'}
                >
                    WYŚLIJ
                </button>
            </div>
        </div>
    );
}

export default App;