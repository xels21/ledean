import { Injectable } from '@angular/core';
import { deepEqual } from 'fast-equals';
import { deepCopy } from 'src/app/lib/deep-copy';
import { WebsocketService } from 'src/app/websocket/websocket.service';
import { Cmd, CmdMode } from 'src/app/websocket/commands';


// type RunningLedStyle = "linear" | "trigonometric"
export enum RunningLedStyle {
  LINEAR = "linear",
  TRIGONOMETRIC = "trigonometric",
}

export interface ModeRunningLedParameter {
  brightness: number,
  fadePct: number,
  roundTimeMs: number,
  hueFrom: number,
  huerTo: number,
  style: RunningLedStyle,
}

export interface ModeRunningLedLimits {
  minRoundTimeMs: number,
  maxRoundTimeMs: number,
  minBrightness: number,
  maxBrightness: number,
}


@Injectable({
  providedIn: 'root'
})
export class ModeRunningLedService {
  public backParameter: ModeRunningLedParameter
  public parameter: ModeRunningLedParameter
  public limits: ModeRunningLedLimits

  private name = "ModeRunningLed"

  constructor(private websocketService: WebsocketService) { }

  getName() {
    return this.name
  }


  receiveParameter(parm: ModeRunningLedParameter) {
    if (!deepEqual(this.backParameter, parm)) {
      this.backParameter = parm
      this.parameter = deepCopy(this.backParameter)
    }
  }

  receiveLimits(limits: ModeRunningLedLimits) {
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


  getAllStyles() {
    return new Array<RunningLedStyle>(
      RunningLedStyle.LINEAR,
      RunningLedStyle.TRIGONOMETRIC,
    )
  }
}
