import { Injectable } from '@angular/core';

import { RGB } from '../../color/color';
import { REST_MODE_RUNNING_LED_URL } from '../../config/const';
import { UpdateService } from '../../update/update.service';
import { HttpClient } from '@angular/common/http';
import { deepCopy } from '../../lib/deep-copy';
import { deepEqual } from 'fast-equals';

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
  public backModeRunningLedParameter: ModeRunningLedParameter
  public modeRunningLedParameter: ModeRunningLedParameter
  public modeRunningLedLimits: ModeRunningLedLimits

  constructor(private httpClient: HttpClient) { }
  getName() {
    return "ModeRunningLed"
  }


  updateModeRunningLedParameter(parm:ModeRunningLedParameter) {
      if (!deepEqual(this.backModeRunningLedParameter, parm)) {
        this.backModeRunningLedParameter = parm
        this.modeRunningLedParameter = deepCopy(this.backModeRunningLedParameter)
      }
  }

  setModeRunningLedParameter() {
    this.httpClient.post<ModeRunningLedParameter>(REST_MODE_RUNNING_LED_URL, this.modeRunningLedParameter, {}).subscribe()
  }

  updateModeRunningLedLimits(limits: ModeRunningLedLimits) {
    this.modeRunningLedLimits = limits
  }

  getAllStyles() {
    return new Array<RunningLedStyle>(
      RunningLedStyle.LINEAR,
      RunningLedStyle.TRIGONOMETRIC,
    )
  }
}
