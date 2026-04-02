// ในตัวแปรสถานะ
const [randomCode, setRandomCode] = useState("0000");

useEffect(() => {
    if (status === 'F-16 ONLINE') {
        const interval = setInterval(() => {
            setRandomCode(Math.random().toString(16).substring(2, 6).toUpperCase());
        }, 100);
        return () => clearInterval(interval);
    }
}, [status]);

// ในส่วนการแสดงผล Cockpit
<div className="text-center">
    {status === 'F-16 ONLINE' && (
        <div className="text-[10px] text-cyan-500 mb-2 tracking-[0.5em] animate-pulse">
            DECRYPTING: {randomCode}
        </div>
    )}
    <div className={`text-4xl font-black tracking-widest uppercase italic ${status === 'STANDBY' ? 'text-slate-800' : 'text-cyan-400 glow-cyan'}`}>
        {status}
    </div>
</div>
