import WebSocketManager from "./websocket_manager.js";
import InterfaceService from "../services/interface_service.js";
import WebSocket from "ws";

interface Command {
  description: string;
  args: string[];
  method: (...args: string[]) => void;
}

export default class WSTestKit {
  private myInterface: InterfaceService;
  private wsManager: WebSocketManager;
  private commands: { [key: string]: Command };

  constructor() {
    this.myInterface = InterfaceService.getInstance();
    this.myInterface.setProgram("Websocket");
    this.wsManager = new WebSocketManager();
    this.commands = {
      add: {
        method: this.addWebSocket,
        args: ["user_id"],
        description: "Add a new WebSocket",
      },
      help: {
        method: this.displayHelp,
        args: [],
        description: "Display help information",
      },
      exit: {
        method: this.exit,
        args: [],
        description: "Exit the application",
      },
    };

    // Bind methods
    this.addWebSocket = this.addWebSocket.bind(this);
    this.displayHelp = this.displayHelp.bind(this);
    this.exit = this.exit.bind(this);
  }

  run(): void {
    this.myInterface.question(
      "Enter a command, or 'help' for a list of commands",
      (command: string) => {
        const [commandName, ...args] = command.trim().split(" ");
        const cmd = this.commands[commandName];

        if (cmd && typeof cmd.method === "function") {
          try {
            cmd.method.apply(this, args);
          } catch (e) {
            this.myInterface.write(`Error: ${e}\n`);
            this.run();
          }
        } else {
          // New logic to handle <user_id> send <message>
          this.processDynamicCommands(commandName, args);
        }
      }
    );
  }

  private addWebSocket(userId: string): void {
    if (!userId) {
      this.myInterface.write("Invalid user ID\n");
      this.run();
      return;
    }

    this.wsManager.openWebSocket(userId);
    this.myInterface.write(`WebSocket opened for user ${userId}\n`);
    this.run();
  }

  private processDynamicCommands(userId: string, args: string[]): void {
    const ws = this.wsManager.getWebSocket(userId);
    if (!ws) {
      this.myInterface.write("Invalid user ID or command\n");
      this.run();
      return;
    }

    const [action, ...messageParts] = args;
    if (action === "send" && messageParts.length > 0) {
      this.sendMessage(userId, messageParts.join(" "));
    } else if (action === "close") {
      this.closeWebSocketCommand(userId);
    } else if (action === "status") {
      this.websocketStatus(userId);
    } else {
      this.myInterface.write("Invalid action or missing parameters\n");
      this.run();
    }
  }

  private sendMessage(userId: string, message: string): void {
    const ws = this.wsManager.getWebSocket(userId);
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(message);
      this.myInterface.write(`Message sent to user ${userId}\n`);
    } else {
      this.myInterface.write(
        `WebSocket is not open or does not exist for user ${userId}\n`
      );
    }
    this.run();
  }

  private closeWebSocketCommand(userId: string): void {
    this.wsManager.closeWebSocket(userId);
    this.myInterface.write(`WebSocket closed for user ${userId}\n`);
    this.run();
  }

  private websocketStatus(userId: string): void {
    const ws = this.wsManager.getWebSocket(userId);
    if (ws) {
      const status = `WebSocket for ${userId} is currently ${
        ws.readyState === WebSocket.OPEN ? "open" : "closed"
      }\n`;
      this.myInterface.write(status);
    } else {
      this.myInterface.write(`No WebSocket found for user ${userId}\n`);
    }
    this.run();
  }

  private displayHelp(): void {
    this.myInterface.write("Available commands:\n");
    Object.entries(this.commands).forEach(
      ([command, { args, description }]) => {
        this.myInterface.write(
          `${command} ${args
            .map((arg) => `<${arg}>`)
            .join(" ")} - ${description}\n`
        );
      }
    );
    this.run();
  }

  private exit(): void {
    this.myInterface.write("Exiting...\n");
    this.wsManager.closeAllWebSockets();
    this.myInterface.close();
  }
}
