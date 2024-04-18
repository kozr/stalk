import WSTestKit from "./ws_test_kit/index.js";
import InterfaceService from "./services/interface_service.js";

const myInterface = InterfaceService.getInstance();
myInterface.setProgram("Main");
myInterface.write("Welcome to the JS Test Kit for Stalk\n");

myInterface.write("Starting WebSocket Test Kit...\n");
const wsTestKit = new WSTestKit();
wsTestKit.run();
