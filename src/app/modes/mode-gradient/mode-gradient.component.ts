import { Component, OnInit } from '@angular/core';
// import { RGB } from '../../color/color';
import { REST_MODE_GRADIENT_URL } from '../../config/const';
import { UpdateService } from '../../update/update.service';
import { HttpClient } from '@angular/common/http';
import { deepCopy } from '../../lib/deep-copy';
import { deepEqual } from 'fast-equals';

// type RunningLedStyle = "linear" | "trigonometric"
// export enum EmitStyle {
//   PULSE = "pulse",
//   DROP = "drop",
// }

interface ModeGradientParameter {
  count: number
  // emitStyle: EmitStyle
  brightness: number
  // minEmitLifetimeMs: number
  // maxEmitLifetimeMs: number
}
interface ModeGradientLimits {
  minCount: number
  maxCount: number
  // minEmitLifetimeMs: number
  // maxEmitLifetimeMs: number
  minBrightness: number
  maxBrightness: number
}


@Component({
  selector: 'app-mode-gradient',
  templateUrl: './mode-gradient.component.html',
  styleUrls: ['./mode-gradient.component.scss']
})
export class ModeGradientComponent implements OnInit {
  public backModeGradientParameter: ModeGradientParameter
  public modeGradientParameter: ModeGradientParameter
  public modeGradientLimits: ModeGradientLimits
  // public brightnessRange: number[] = [0,0]
  // public emitLifetimeMsRange: number[] = [0,0]

  constructor(private updateService: UpdateService, private httpClient: HttpClient) { }

  ngOnInit(): void {
    this.updateModeGradientParameter();
    this.updateModeGradientLimits();
    this.updateService.registerPolling({ cb: () => { this.updateModeGradientParameter() }, timeout: 500 })
  }

  updateModeGradientParameter() {
    this.httpClient.get<ModeGradientParameter>(REST_MODE_GRADIENT_URL).subscribe((data: ModeGradientParameter) => {
      if (!deepEqual(this.backModeGradientParameter, data)) {
        this.backModeGradientParameter = data
        this.modeGradientParameter = deepCopy(this.backModeGradientParameter)
      }
    })
  }

  // setBrightness(){
  //   this.modeGradientParameter.minBrightness = this.brightnessRange[0]
  //   this.modeGradientParameter.maxBrightness = this.brightnessRange[1]
  // }
  // setEmitLifetimeMs(){
  //   this.modeGradientParameter.minEmitLifetimeMs = this.emitLifetimeMsRange[0]
  //   this.modeGradientParameter.maxEmitLifetimeMs = this.emitLifetimeMsRange[1]
  // }

  setModeGradientParameter() {
    this.httpClient.post<ModeGradientParameter>(REST_MODE_GRADIENT_URL, this.modeGradientParameter, {}).subscribe()
  }

  updateModeGradientLimits() {
    this.httpClient.get<ModeGradientLimits>(REST_MODE_GRADIENT_URL + "/limits").subscribe((data: ModeGradientLimits) => { this.modeGradientLimits = data })
  }
  // getAllStyles() {
  //   return new Array<EmitStyle>(
  //     EmitStyle.PULSE,
  //     EmitStyle.DROP,
  //   )
  // }  
}
