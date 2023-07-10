import { Injectable } from '@angular/core';
import { RGB } from '../../color/color';
import { REST_MODE_SOLID_RAINBOW_URL } from '../../config/const';
import { UpdateService } from '../../update/update.service';
import { HttpClient } from '@angular/common/http';
import { deepCopy } from '../../lib/deep-copy';
import { deepEqual } from 'fast-equals';

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

  public backModeSolidRainbowParameter: ModeSolidRainbowParameter
  public modeSolidRainbowParameter: ModeSolidRainbowParameter
  public modeSolidRainbowLimits: ModeSolidRainbowLimits

  constructor(private httpClient: HttpClient,) { }

  getName() {
    return "ModeSolidRainbow"
  }

  
  updateModeSolidRainbowParameter(parm: ModeSolidRainbowParameter) {
      if (!deepEqual(this.backModeSolidRainbowParameter, parm)) {
        this.backModeSolidRainbowParameter = parm
        this.modeSolidRainbowParameter = deepCopy(this.backModeSolidRainbowParameter)
        // console.log(data)
      }
  }

  setModeSolidRainbowParameter() {
    console.log("set")
    this.httpClient.post<ModeSolidRainbowParameter>(REST_MODE_SOLID_RAINBOW_URL, this.modeSolidRainbowParameter, {}).subscribe()
  }

  updateModeSolidRainbowLimits(limits: ModeSolidRainbowLimits) {
     this.modeSolidRainbowLimits = limits
  }

}
