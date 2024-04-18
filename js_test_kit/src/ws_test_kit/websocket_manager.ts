import WebSocket from "ws";

const DEFAULT_WS_URL = "ws://localhost:8080/establish-connection";

export default class WebSocketManager {
  private websockets: { [key: string]: WebSocket } = {};

  constructor() {
    this.websockets = {};
  }

  openWebSocket(userId: string): WebSocket {
    if (!userId) throw new Error("userId must be provided");
    try {
      if (this.websockets[userId]) {
        console.log(`WebSocket already open for userId: ${userId}`);
        return this.websockets[userId];
      }
      const ws = new WebSocket(`${DEFAULT_WS_URL}?userId=${userId}`);
      ws.onclose = () => {
        console.log(`WebSocket closed for userId: ${userId}`);
        delete this.websockets[userId];
      };
      ws.onerror = (err) => {
        console.error(`WebSocket error for userId: ${userId}:`, err);
        delete this.websockets[userId];
      };
      this.websockets[userId] = ws;
      return ws;
    } catch (e) {
      console.error(`Error opening WebSocket for userId: ${userId}:`, e);
      throw e;
    }
  }

  getWebSocket(userId: string): WebSocket {
    return this.websockets[userId];
  }

  getAllWebSockets(): { [key: string]: WebSocket } {
    return this.websockets;
  }

  closeWebSocket(userId: string) {
    if (this.websockets[userId]) {
      this.websockets[userId].close();
      delete this.websockets[userId];
    } else {
      console.warn(`No WebSocket found for userId: ${userId}`);
    }
  }

  closeAllWebSockets() {
    Object.keys(this.websockets).forEach((userId) => {
      this.closeWebSocket(userId);
    });
    this.websockets = {};
  }
}
