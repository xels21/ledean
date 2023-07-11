import { Injectable } from '@angular/core';
import { deepEqual } from 'fast-equals';
import { deepCopy } from 'src/app/lib/deep-copy';
import { WebsocketService } from 'src/app/websocket/websocket.service';
import { Cmd, CmdMode } from 'src/app/websocket/commands';


export interface ModeGradientParameter {
  count: number
  roundTimeMs: number
  brightness: number

}
export interface ModeGradientLimits {
  minCount: number
  maxCount: number
  minRoundTimeMs: number
  maxRoundTimeMs: number
  minBrightness: number
  maxBrightness: number
}

@Injectable({
  providedIn: 'root'
})
export class ModeGradientService {
  public backParameter: ModeGradientParameter
  public parameter: ModeGradientParameter
  public limits: ModeGradientLimits

  private name = "ModeGradient"

  constructor(private websocketService: WebsocketService) { }

  getName() {
    return this.name
  }

  receiveParameter(parm: ModeGradientParameter) {
    if (!deepEqual(this.backParameter, parm)) {
      this.backParameter = parm
      this.parameter = deepCopy(this.backParameter)
    }
  }

  receiveLimits(limits: ModeGradientLimits) {
    this.limits = limits;
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
