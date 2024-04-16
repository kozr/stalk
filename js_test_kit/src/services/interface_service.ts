import * as readline from "readline";

class InterfaceService {
  private static instance: InterfaceService;
  private myInterface: readline.Interface  = {} as readline.Interface;
  private programName: string = "";

  private constructor() {
    if (InterfaceService.instance) {
      return InterfaceService.instance;
    }

    this.myInterface = readline.createInterface({
      input: process.stdin,
      output: process.stdout,
    });
    this.programName = "Main";
    InterfaceService.instance = this;
  }

  public setProgram(name: string): void {
    this.programName = name;
    this.myInterface.write("-----------------------------------------\n")
  }

  public write(message: string, level: string = "info"): void {
    switch (level) {
      case "error":
        this.myInterface.write(`ERROR: ${message}\n`);
        break;
      case "warning":
        this.myInterface.write(`WARNING: ${message}\n`);
        break;
      default:
        this.myInterface.write(`${message}\n`);
    }
  }

  public question(prompt: string, callback: (input: string) => void): void {
    const fullPrompt = `${this.programName}: ${prompt}\n> `;
    this.myInterface.question(fullPrompt, (input) => {
      callback(input);
    });
  }

  public close(): void {
    this.myInterface.close();
  }

  public static getInstance(): InterfaceService {
    if (!InterfaceService.instance) {
      InterfaceService.instance = new InterfaceService();
    }
    return InterfaceService.instance;
  }
}

export default InterfaceService;
