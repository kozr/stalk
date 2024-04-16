import WSTestKit from "./ws_test_kit/index.js";
import InterfaceService from "./services/interface_service.js";

const myInterface = InterfaceService.getInstance();
myInterface.setProgram("Main");
myInterface.write("Welcome to the JS Test Kit for Stalk\n");

interface Command {
  description: string;
  function: (callback: () => void) => void;
}

const commands: Record<string, Command> = {
  help: {
    description: "Display this help message",
    function: displayHelp,
  },
  ws: {
    description: "WebSocket",
    function: startWebSocket,
  },
  exit: {
    description: "Exit the program",
    function: exitProgram,
  },
};

function prompt(): void {
  myInterface.question("Enter a command, or 'help' for a list of commands", (command: string) => {
    const cmd = commands[command.trim()];
    if (cmd && typeof cmd.function === "function") {
      cmd.function(() => {
        prompt(); // Recurse after the command completes
      });
    } else {
      myInterface.write("Invalid command\n");
      prompt(); // Recurse to continue prompting for command
    }
  });
}

function displayHelp(callback: () => void): void {
  for (let command in commands) {
    myInterface.write(`${command} - ${commands[command].description}\n`);
  }
  callback();
}

function startWebSocket(callback: () => void): void {
  myInterface.write("Starting WebSocket Test Kit...\n");
  const wsTestKit = new WSTestKit();
  wsTestKit.run();
  callback();
}

function exitProgram(callback: () => void): void {
  myInterface.write("Exiting...\n");
  process.exit(0);
}

prompt();  // Start the interactive prompt
