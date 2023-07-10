import { REST_MODE_SOLID_URL } from '../../config/const';
import { UpdateService } from '../../update/update.service';
import { HttpClient } from '@angular/common/http';
import { deepCopy } from '../../lib/deep-copy';
import { deepEqual } from 'fast-equals';

import { Injectable } from '@angular/core';
import { RGB } from '../../color/color'

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

  public backModeSolidParameter: ModeSolidParameter
  public modeSolidParameter: ModeSolidParameter
  public modeSolidLimits: ModeSolidLimits

  constructor(private httpClient: HttpClient) { }


  updateModeSolidParameter(parm: ModeSolidParameter) {
    console.log(parm)
    if (!deepEqual(this.backModeSolidParameter, parm)) {
      this.backModeSolidParameter = parm
      this.modeSolidParameter = deepCopy(this.backModeSolidParameter)
    }
  }

  getName() {
    return "ModeSolid"
  }

  setModeSolidParameter() {
    this.httpClient.post<ModeSolidParameter>(REST_MODE_SOLID_URL, this.modeSolidParameter, {}).subscribe()
  }

  updateModeSolidLimits(limits: ModeSolidLimits) {
    this.modeSolidLimits = limits
  }


  // updateModeSolidLimits() {
  // this.httpClient.get<ModeSolidLimits>(REST_MODE_SOLID_URL + "/limits").subscribe((data: ModeSolidLimits) => { this.modeSolidLimits = data })
  // }

}
