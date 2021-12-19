import { Component, OnInit } from '@angular/core';
// import { RGB } from '../../color/color';
import { REST_MODE_EMITTER_URL } from '../../config/const';
import { UpdateService } from '../../update/update.service';
import { HttpClient } from '@angular/common/http';
import { deepCopy } from '../../lib/deep-copy';
import { deepEqual } from 'fast-equals';

// type RunningLedStyle = "linear" | "trigonometric"
export enum EmitStyle {
  PULSE = "pulse",
  DROP = "drop",
}

interface ModeEmitterParameter {
  emitCount: number
  emitStyle: EmitStyle
  minBrightness: number
  maxBrightness: number
  minEmitLifetimeMs: number
  maxEmitLifetimeMs: number
}
interface ModeEmitterLimits {
  minEmitCount: number
  maxEmitCount: number
  minEmitLifetimeMs: number
  maxEmitLifetimeMs: number
  minBrightness: number
  maxBrightness: number
}


@Component({
  selector: 'app-mode-emitter',
  templateUrl: './mode-emitter.component.html',
  styleUrls: ['./mode-emitter.component.scss','../../app.component.scss']
})
export class ModeEmitterComponent implements OnInit {
  public backModeEmitterParameter: ModeEmitterParameter
  public modeEmitterParameter: ModeEmitterParameter
  public modeEmitterLimits: ModeEmitterLimits
  public brightnessRange: number[] = [0,0]
  public emitLifetimeMsRange: number[] = [0,0]

  constructor(private updateService: UpdateService, private httpClient: HttpClient) { }

  ngOnInit(): void {
    this.updateModeEmitterParameter();
    this.updateModeEmitterLimits();
    this.updateService.registerPolling({ cb: () => { this.updateModeEmitterParameter() }, timeout: 500 })
  }

  updateModeEmitterParameter() {
    this.httpClient.get<ModeEmitterParameter>(REST_MODE_EMITTER_URL).subscribe((data: ModeEmitterParameter) => {
      if (!deepEqual(this.backModeEmitterParameter, data)) {
        this.backModeEmitterParameter = data
        this.modeEmitterParameter = deepCopy(this.backModeEmitterParameter)
        this.brightnessRange =  [this.modeEmitterParameter.minBrightness, this.modeEmitterParameter.maxBrightness]
        this.emitLifetimeMsRange =  [this.modeEmitterParameter.minEmitLifetimeMs, this.modeEmitterParameter.maxEmitLifetimeMs]
      }
    })
  }

  setBrightness(){
    this.modeEmitterParameter.minBrightness = this.brightnessRange[0]
    this.modeEmitterParameter.maxBrightness = this.brightnessRange[1]
  }
  setEmitLifetimeMs(){
    this.modeEmitterParameter.minEmitLifetimeMs = this.emitLifetimeMsRange[0]
    this.modeEmitterParameter.maxEmitLifetimeMs = this.emitLifetimeMsRange[1]
  }

  setModeEmitterParameter() {
    this.httpClient.post<ModeEmitterParameter>(REST_MODE_EMITTER_URL, this.modeEmitterParameter, {}).subscribe()
  }

  updateModeEmitterLimits() {
    this.httpClient.get<ModeEmitterLimits>(REST_MODE_EMITTER_URL + "/limits").subscribe((data: ModeEmitterLimits) => { this.modeEmitterLimits = data })
  }
  getAllStyles() {
    return new Array<EmitStyle>(
      EmitStyle.PULSE,
      EmitStyle.DROP,
    )
  }  
}
