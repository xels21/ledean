import { Component, OnInit } from '@angular/core';
import { REST_MODE_GRADIENT_URL } from '../../config/const';
import { UpdateService } from '../../update/update.service';
import { HttpClient } from '@angular/common/http';
import { deepCopy } from '../../lib/deep-copy';
import { deepEqual } from 'fast-equals';


interface ModeGradientParameter {
  count: number
  roundTimeMs:number
  brightness: number

}
interface ModeGradientLimits {
  minCount: number
  maxCount: number
  minRoundTimeMs: number
  maxRoundTimeMs: number
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


  setModeGradientParameter() {
    this.httpClient.post<ModeGradientParameter>(REST_MODE_GRADIENT_URL, this.modeGradientParameter, {}).subscribe()
  }

  updateModeGradientLimits() {
    this.httpClient.get<ModeGradientLimits>(REST_MODE_GRADIENT_URL + "/limits").subscribe((data: ModeGradientLimits) => { this.modeGradientLimits = data })
  }

}
