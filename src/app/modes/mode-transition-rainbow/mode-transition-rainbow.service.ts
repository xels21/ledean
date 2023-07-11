import { Injectable } from '@angular/core';
import { deepCopy } from 'src/app/lib/deep-copy';
import { deepEqual } from 'fast-equals';
import { WebsocketService } from 'src/app/websocket/websocket.service';
import { Cmd, CmdMode } from 'src/app/websocket/commands';


export interface ModeTransitionRainbowParameter {
  roundTimeMs: number,
  brightness: number,
  spectrum: number,
  reverse: boolean,
}
export interface ModeTransitionRainbowLimits {
  minRoundTimeMs: number,
  maxRoundTimeMs: number,
  minBrightness: number,
  maxBrightness: number,
}

@Injectable({
  providedIn: 'root'
})
export class ModeTransitionRainbowService {
  public backParameter: ModeTransitionRainbowParameter
  public parameter: ModeTransitionRainbowParameter
  public limits: ModeTransitionRainbowLimits

  private name = "ModeTransitionRainbow"

  constructor(private websocketService: WebsocketService) { }

  getName() {
    return this.name
  }

  receiveParameter(parm: ModeTransitionRainbowParameter) {
    if (!deepEqual(this.backParameter, parm)) {
      this.backParameter = parm
      this.parameter = deepCopy(this.backParameter)
    }
  }

  receiveLimits(limits: ModeTransitionRainbowLimits) {
    this.limits = limits
  }

  sendParameter() {
    this.websocketService.send({
      cmd: "mode",
      parm: {
        id: this.name,
        parm: this.parameter
      } as CmdMode
    } as Cmd)
  }
}
