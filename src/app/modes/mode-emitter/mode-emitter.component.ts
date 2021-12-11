import { Component, OnInit } from '@angular/core';
// import { RGB } from '../../color/color';
import { REST_MODE_EMITTER_URL } from '../../config/const';
import { UpdateService } from '../../update/update.service';
import { HttpClient } from '@angular/common/http';
import { deepCopy } from '../../lib/deep-copy';
import { deepEqual } from 'fast-equals';

interface ModeEmitterParameter {
  emitCount:number
}
interface ModeEmitterLimits {
  minEmitCount: number
  maxEmitCount: number
}





@Component({
  selector: 'app-mode-emitter',
  templateUrl: './mode-emitter.component.html',
  styleUrls: ['./mode-emitter.component.scss']
})
export class ModeEmitterComponent implements OnInit {
  public backModeEmitterParameter: ModeEmitterParameter
  public modeEmitterParameter: ModeEmitterParameter
  public modeEmitterLimits: ModeEmitterLimits

  constructor(private updateService:UpdateService, private httpClient:HttpClient) { }

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
      }
    }
    )
  }

  setModeEmitterParameter() {
    console.log("set")
    this.httpClient.post<ModeEmitterParameter>(REST_MODE_EMITTER_URL, this.modeEmitterParameter, {}).subscribe()
  }

  updateModeEmitterLimits() {
    this.httpClient.get<ModeEmitterLimits>(REST_MODE_EMITTER_URL + "/limits").subscribe((data: ModeEmitterLimits) => { this.modeEmitterLimits = data })
  }
}
