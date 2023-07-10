import { Injectable } from '@angular/core';

import { REST_MODE_EMITTER_URL } from '../../config/const';
import { deepCopy } from '../../lib/deep-copy';
import { deepEqual } from 'fast-equals';
import { HttpClient } from '@angular/common/http';


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

  public backModeEmitterParameter: ModeEmitterParameter
  public modeEmitterParameter: ModeEmitterParameter
  public modeEmitterLimits: ModeEmitterLimits
  public brightnessRange: number[] = [0, 0]
  public emitLifetimeMsRange: number[] = [0, 0]

  constructor(private httpClient: HttpClient) { }

  getName() {
    return "ModeEmitter"
  }

  updateModeEmitterParameter(parm: ModeEmitterParameter) {
    // this.httpClient.get<ModeEmitterParameter>(REST_MODE_EMITTER_URL).subscribe((data: ModeEmitterParameter) => {
    if (!deepEqual(this.backModeEmitterParameter, parm)) {
      this.backModeEmitterParameter = parm
      this.modeEmitterParameter = deepCopy(this.backModeEmitterParameter)
      this.brightnessRange = [this.modeEmitterParameter.minBrightness, this.modeEmitterParameter.maxBrightness]
      this.emitLifetimeMsRange = [this.modeEmitterParameter.minEmitLifetimeMs, this.modeEmitterParameter.maxEmitLifetimeMs]
    }
    // })
  }

  setBrightness() {
    this.modeEmitterParameter.minBrightness = this.brightnessRange[0]
    this.modeEmitterParameter.maxBrightness = this.brightnessRange[1]
  }
  setEmitLifetimeMs() {
    this.modeEmitterParameter.minEmitLifetimeMs = this.emitLifetimeMsRange[0]
    this.modeEmitterParameter.maxEmitLifetimeMs = this.emitLifetimeMsRange[1]
  }

  setModeEmitterParameter() {
    this.httpClient.post<ModeEmitterParameter>(REST_MODE_EMITTER_URL, this.modeEmitterParameter, {}).subscribe()
  }

  updateModeEmitterLimits(limits: ModeEmitterLimits) {
    this.modeEmitterLimits = limits 
  }
  getAllStyles() {
    return new Array<EmitStyle>(
      EmitStyle.PULSE,
      EmitStyle.DROP,
    )
  }
}
