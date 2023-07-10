import { Injectable } from '@angular/core';
import { RGB } from '../../color/color';
import { REST_MODE_TRANSITION_RAINBOW_URL } from '../../config/const';
import { UpdateService } from '../../update/update.service';
import { HttpClient } from '@angular/common/http';
import { deepCopy } from '../../lib/deep-copy';
import { deepEqual } from 'fast-equals';

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
  public backModeTransitionRainbowParameter: ModeTransitionRainbowParameter
  public modeTransitionRainbowParameter: ModeTransitionRainbowParameter
  public modeTransitionRainbowLimits: ModeTransitionRainbowLimits

  constructor(private httpClient: HttpClient) { }
  getName() {
    return "ModeTransitionRainbow"
  }


  updateModeTransitionRainbowParameter(parm: ModeTransitionRainbowParameter) {
      if (!deepEqual(this.backModeTransitionRainbowParameter, parm)) {
        this.backModeTransitionRainbowParameter = parm
        this.modeTransitionRainbowParameter = deepCopy(this.backModeTransitionRainbowParameter)
      }
  }

  setModeTransitionRainbowParameter() {
    console.log("set")
    this.httpClient.post<ModeTransitionRainbowParameter>(REST_MODE_TRANSITION_RAINBOW_URL, this.modeTransitionRainbowParameter, {}).subscribe()
  }

  updateModeTransitionRainbowLimits(limits: ModeTransitionRainbowLimits) {
    this.modeTransitionRainbowLimits = limits
  }

}
