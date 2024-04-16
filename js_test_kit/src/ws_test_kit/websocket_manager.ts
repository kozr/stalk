import InterfaceService from "../services/interface_service";

interface Command {
  method: () => void;
  description: string;
}

export default class WSTestKit {
  private myInterface: InterfaceService;
  private commands: { [key: string]: Command };

  constructor() {
    this.myInterface = InterfaceService.getInstance();
    this.myInterface.setProgram("Websocket");
    this.commands = {
      add: {
        method: this.addWebSocket,
        description: "Add a new WebSocket",
      },
      help: {
        method: this.displayHelp,
        description: "Display help information",
      },
      exit: {
        method: this.exit,
        description: "Exit the application",
      },
    };
  }

  run(): void {
    this.myInterface.question(
      "Enter a command, or 'help' for a list of commands\n> ",
      (command: string) => {
        const cmd = this.commands[command.trim()];
        if (cmd && typeof cmd.method === "function") {
          cmd.method.call(this);
        } else {
          this.myInterface.write("Invalid command\n");
          this.run();
        }
      }
    );
  }

  private addWebSocket(): void {
    this.myInterface.write("Adding a new WebSocket...\n");
  }

  private displayHelp(): void {
    this.myInterface.write("Available commands:\n");
    Object.entries(this.commands).forEach(([command, { description }]) => {
      this.myInterface.write(`${command}: ${description}\n`);
    });
    this.run();
  }

  private exit(): void {
    this.myInterface.write("Exiting...\n");
    this.myInterface.close();
  }

  private test(): void {
    // private method for internal class use
  }
}
