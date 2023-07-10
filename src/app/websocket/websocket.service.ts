import { Injectable, EventEmitter } from '@angular/core';
import { webSocket, WebSocketSubject } from "rxjs/webSocket";
import { Cmd, CmdLeds, CmdLedsParameter, CmdMode, CmdModeLimits, CmdModeResolver } from "./commands"
import { RGB } from '../color/color';
// import { RGB } from '../led-display/led-display.component';
// import { LedsService} from "../leds/leds.service" //Circular dependency

// import { CookieService } from 'ngx-cookie-service';
// import { Cmd, Cmd2sGet, Cmd2sSet, Cmd2sSubscribe, Cmd2cUpdate, Cmd2sConsole, Cmd2cConsole, Cmd2sUnsubscribe, Cmd2sPing, Cmd2cLogin, Cmd2sLogin, Cmd2sLogout, Cmd2cLogout, Cmd2cLog } from './commands';

const WEBSOCKET_PROTOCOL = "ws";
// const WEBSOCKET_HOST = window.location.host //address:port
const WEBSOCKET_HOST = window.location.hostname + ":" + "2211" //needed for debugging with angular
// const WEBSOCKET_HOST = "127.0.0.1" + ":" + "2211" //address:port
const WEBSOCKET_PATH = "ws";
const WEBSOCKET_COMPLETE_URL = WEBSOCKET_PROTOCOL + "://"
  + WEBSOCKET_HOST + "/"
  + WEBSOCKET_PATH;
const WEBSOCKET_RECONNECT_TIMEOUT = 1000;

@Injectable({
  providedIn: 'root'
})
export class WebsocketService {

  subject: WebSocketSubject<any>
  connected: boolean
  ledsChanged: EventEmitter<Array<RGB>>
  ledsParameterChanged: EventEmitter<CmdLedsParameter>
  modeChanged: EventEmitter<CmdMode>
  modeLimitChanged: EventEmitter<CmdModeLimits>
  modeResolverChanged: EventEmitter<CmdModeResolver>
  // connectedChangeCnt: number

  // subscriptions: Map<string, SubscribeElement>

  constructor() {
    this.connected = false;
    this.ledsChanged = new EventEmitter();
    this.ledsParameterChanged = new EventEmitter();
    this.modeChanged = new EventEmitter();
    this.modeLimitChanged = new EventEmitter();
    this.modeResolverChanged = new EventEmitter();

    // this.subscriptions = new Map<string, SubscribeElement>()
    this.subject = webSocket(
      WEBSOCKET_COMPLETE_URL
    );
    this.run()
  }

  run() {
    this.subject.subscribe(
      msg => {
        if (!this.connected) {
          this.connected = true
          // if(this.authService.log)
        }
        if (msg.hasOwnProperty("cmd") && msg.hasOwnProperty("parm")) {
          let cmd = msg as Cmd
          switch (cmd.cmd) {
            case "leds":
              var cmd2cLeds = cmd.parm as CmdLeds
              this.ledsChanged.emit(cmd2cLeds.leds)
              break;
            case "ledsParameter":
              var cmdLedsParameter = cmd.parm as CmdLedsParameter
              this.ledsParameterChanged.emit(cmdLedsParameter)
              break;
            case "mode":
              var cmdMode = cmd.parm as CmdMode
              this.modeChanged.emit(cmdMode)
              break;
            case "modeLimits":
              var cmdModeLimits = cmd.parm as CmdModeLimits
              this.modeLimitChanged.emit(cmdModeLimits)
              break;
            case "modeResolver":
              var cmdModeResolver = cmd.parm as CmdModeResolver
              this.modeResolverChanged.emit(cmdModeResolver)
              break;
            default:
              console.log("something went wrong with message: ", msg)
          }
        }
        else {
          console.log("unknown websocket message: ", msg)
        }
      },
      (err) => this.reRun(), // Called if at any point WebSocket API signals some kind of error.
      () => this.reRun() // Called when connection is closed (for whatever reason).
    );
    this.subject.next({ message: "hello" })
  }



  private reRun() {
    this.connected = false
    setTimeout(() => this.run(), WEBSOCKET_RECONNECT_TIMEOUT)
  }

  public send(cmd: Cmd) {
    this.subject.next(cmd)
    // this.subject.next(JSON.stringify(cmd))
  }

}


