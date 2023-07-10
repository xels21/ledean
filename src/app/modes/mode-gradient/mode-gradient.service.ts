import { Injectable } from '@angular/core';
import { UpdateService } from '../../update/update.service';
import { HttpClient } from '@angular/common/http';
import { deepCopy } from '../../lib/deep-copy';
import { deepEqual } from 'fast-equals';
import { REST_MODE_GRADIENT_URL } from '../../config/const';


export interface ModeGradientParameter {
  count: number
  roundTimeMs:number
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
  public backModeGradientParameter: ModeGradientParameter
  public modeGradientParameter: ModeGradientParameter
  public modeGradientLimits: ModeGradientLimits

  constructor( private httpClient: HttpClient) { }

  getName() {
    return "ModeGradient"
  }

  
  updateModeGradientParameter(parm : ModeGradientParameter) {
    // this.httpClient.get<ModeGradientParameter>(REST_MODE_GRADIENT_URL).subscribe((data: ModeGradientParameter) => {
      if (!deepEqual(this.backModeGradientParameter, parm)) {
        this.backModeGradientParameter = parm
        this.modeGradientParameter = deepCopy(this.backModeGradientParameter)
      }
    // })
  }


  setModeGradientParameter() {
    this.httpClient.post<ModeGradientParameter>(REST_MODE_GRADIENT_URL, this.modeGradientParameter, {}).subscribe()
  }

  updateModeGradientLimits(limits: ModeGradientLimits) {
      this.modeGradientLimits = limits;
  }
}
