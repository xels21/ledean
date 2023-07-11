import { Injectable } from '@angular/core';
import { deepCopy } from 'src/app/lib/deep-copy';
import { deepEqual } from 'fast-equals';
import { WebsocketService } from 'src/app/websocket/websocket.service';
import { Cmd, CmdMode } from 'src/app/websocket/commands';


// type RunningLedStyle = "linear" | "trigonometric"
export enum EmitStyle {
  PULSE = "pulse",
  DROP = "drop",
}

export interface ModeEmitterParameter {
  emitCount: number
  emitStyle: EmitStyle
  minBrightness: number
  maxBrightness: number
  minEmitLifetimeMs: number
  maxEmitLifetimeMs: number
}
export interface ModeEmitterLimits {
  minEmitCount: number
  maxEmitCount: number
  minEmitLifetimeMs: number
  maxEmitLifetimeMs: number
  minBrightness: number
  maxBrightness: number
}

@Injectable({
  providedIn: 'root'
})
export class ModeEmitterService {

  public backParameter: ModeEmitterParameter
  public parameter: ModeEmitterParameter
  public limits: ModeEmitterLimits

  public brightnessRange: number[] = [0, 0]
  public emitLifetimeMsRange: number[] = [0, 0]

  private name ="ModeEmitter"

  constructor(private websocketService: WebsocketService) { }

  getName() {
    return this.name
  }

  receiveParameter(parm: ModeEmitterParameter) {
    if (!deepEqual(this.backParameter, parm)) {
      this.backParameter = parm
      this.parameter = deepCopy(this.backParameter)
      this.brightnessRange = [this.parameter.minBrightness, this.parameter.maxBrightness]
      this.emitLifetimeMsRange = [this.parameter.minEmitLifetimeMs, this.parameter.maxEmitLifetimeMs]
    }
  }

  receiveLimits(limits: ModeEmitterLimits) {
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

  
  setBrightness() {
    this.parameter.minBrightness = this.brightnessRange[0]
    this.parameter.maxBrightness = this.brightnessRange[1]
  }

  setEmitLifetimeMs() {
    this.parameter.minEmitLifetimeMs = this.emitLifetimeMsRange[0]
    this.parameter.maxEmitLifetimeMs = this.emitLifetimeMsRange[1]
  }

  getAllStyles() {
    return new Array<EmitStyle>(
      EmitStyle.PULSE,
      EmitStyle.DROP,
    )
  }
}
