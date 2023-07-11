import { Injectable } from '@angular/core';
import { deepEqual } from 'fast-equals';
import { deepCopy } from 'src/app/lib/deep-copy';
import { WebsocketService } from 'src/app/websocket/websocket.service';
import { Cmd, CmdMode } from 'src/app/websocket/commands';

export interface ModeSolidRainbowParameter {
  roundTimeMs: number,
  brightness: number,
}
export interface ModeSolidRainbowLimits {
  minRoundTimeMs: number,
  maxRoundTimeMs: number,
  minBrightness: number,
  maxBrightness: number,
}

@Injectable({
  providedIn: 'root'
})
export class ModeSolidRainbowService {

  public backParameter: ModeSolidRainbowParameter
  public parameter: ModeSolidRainbowParameter
  public limits: ModeSolidRainbowLimits

  private name = "ModeSolidRainbow"

  constructor(private websocketService: WebsocketService) { }

  getName() {
    return this.name
  }

  receiveParameter(parm: ModeSolidRainbowParameter) {
    if (!deepEqual(this.backParameter, parm)) {
      this.backParameter = parm
      this.parameter = deepCopy(this.backParameter)
    }
  }

  receiveLimits(limits: ModeSolidRainbowLimits) {
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
