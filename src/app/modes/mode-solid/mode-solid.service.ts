import { Injectable } from '@angular/core';
import { deepEqual } from 'fast-equals';
import { deepCopy } from 'src/app/lib/deep-copy';
import { RGB } from 'src/app/color/color'
import { WebsocketService } from 'src/app/websocket/websocket.service';
import { CmdMode, Cmd } from 'src/app/websocket/commands';


export interface ModeSolidParameter {
  rgb: RGB,
  brightness: number,
}
export interface ModeSolidLimits {
  minBrightness: number,
  maxBrightness: number,
}

@Injectable({
  providedIn: 'root'
})
export class ModeSolidService {

  public backParameter: ModeSolidParameter
  public parameter: ModeSolidParameter
  public limits: ModeSolidLimits

  private name = "ModeSolid"

  constructor(private websocketService: WebsocketService) { }

  getName() {
    return this.name
  }

  receiveParameter(parm: ModeSolidParameter) {
    if (!deepEqual(this.backParameter, parm)) {
      this.backParameter = parm
      this.parameter = deepCopy(this.backParameter)
    }
  }


  receiveLimits(limits: ModeSolidLimits) {
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
