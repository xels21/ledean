
import { Component, OnInit } from '@angular/core';
import { RGB } from '../../color/color';
import { REST_MODE_RUNNING_LED_URL } from '../../config/const';
import { UpdateService } from '../../update/update.service';
import { HttpClient } from '@angular/common/http';
import { deepCopy } from '../../lib/deep-copy';
import { deepEqual } from 'fast-equals';

interface ModeRunningLedParameter {
  roundTimeMs: number,
  brightness: number,
}
interface ModeRunningLedLimits {
  minRoundTimeMs: number,
  maxRoundTimeMs: number,
  minBrightness: number,
  maxBrightness: number,
}

type RunningLedStyle = "linear" | "trigonometric"

@Component({
  selector: 'app-mode-running-led',
  templateUrl: './mode-running-led.component.html',
  styleUrls: ['./mode-running-led.component.scss']
})
export class ModeRunningLedComponent implements OnInit {
  public backModeRunningLedParameter: ModeRunningLedParameter
  public modeRunningLedParameter: ModeRunningLedParameter
  public modeRunningLedLimits: ModeRunningLedLimits

  constructor(private httpClient: HttpClient, private updateService: UpdateService) {
  }

  ngOnInit(): void {
    this.updateModeRunningLedParameter();
    this.updateModeRunningLedLimits();
    this.updateService.registerPolling({ cb: () => { this.updateModeRunningLedParameter() }, timeout: 500 })
    setTimeout( ()=>{M.updateTextFields()},100);
  }

  updateModeRunningLedParameter() {
    this.httpClient.get<ModeRunningLedParameter>(REST_MODE_RUNNING_LED_URL).subscribe((data: ModeRunningLedParameter) => {
      if (!deepEqual(this.backModeRunningLedParameter, data)) {
        this.backModeRunningLedParameter = data
        this.modeRunningLedParameter = deepCopy(this.backModeRunningLedParameter)
        console.log(data)
      }
    }
    )
  }

  setModeRunningLedParameter() {
    console.log("set")
    this.httpClient.post<ModeRunningLedParameter>(REST_MODE_RUNNING_LED_URL, this.modeRunningLedParameter, {}).subscribe()
  }

  updateModeRunningLedLimits() {
    this.httpClient.get<ModeRunningLedLimits>(REST_MODE_RUNNING_LED_URL + "/limits").subscribe((data: ModeRunningLedLimits) => { this.modeRunningLedLimits = data })
  }

}
