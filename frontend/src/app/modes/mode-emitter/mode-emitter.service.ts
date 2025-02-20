import { Injectable } from '@angular/core';
import { WebsocketService } from 'src/app/websocket/websocket.service';
import { SuperMode } from 'src/app/modes/super-mode';


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
export class ModeEmitterService extends SuperMode {

  public backParameter: ModeEmitterParameter
  public parameter: ModeEmitterParameter
  public limits: ModeEmitterLimits

  public brightnessRange: number[] = [0, 0]
  public emitLifetimeMsRange: number[] = [0, 0]

  constructor(protected websocketService: WebsocketService) {
    super("ModeEmitter", websocketService)
  }

  receiveParameter(parm: ModeEmitterParameter) {
    super.receiveParameter(parm)
    this.brightnessRange = [this.parameter.minBrightness, this.parameter.maxBrightness]
    this.emitLifetimeMsRange = [this.parameter.minEmitLifetimeMs, this.parameter.maxEmitLifetimeMs]
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
