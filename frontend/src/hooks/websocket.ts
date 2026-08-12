import { useEffect, useState, useRef } from 'react';

export function useWebSocket(url: string, onMessage?: (data: string) => void) {
  const [messages, setMessages] = useState<string[]>([]);
  const [isConnected, setIsConnected] = useState(false);
  const socketRef = useRef<WebSocket | null>(null);
  const onMessageRef = useRef(onMessage);

  useEffect(() => {
    onMessageRef.current = onMessage;
  });

  useEffect(() => {
    const socket = new WebSocket(url);
    socketRef.current = socket;

    socket.onopen = () => setIsConnected(true);
    socket.onclose = () => setIsConnected(false);
    socket.onmessage = (event) => {
      const data = String(event.data);
      setMessages((prev) => [...prev, data]);
      onMessageRef.current?.(data);
    };

    return () => {
      socket.close();
    };
  }, [url]);

  const sendMessage = (data: string) => {
    if (socketRef.current && isConnected) {
      socketRef.current?.send(data);
    }
  };

  return { messages, isConnected, sendMessage };
}
