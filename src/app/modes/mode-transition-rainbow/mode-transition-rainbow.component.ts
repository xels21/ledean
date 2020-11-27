
import { Component, OnInit } from '@angular/core';
import { RGB } from '../../color/color';
import { REST_MODE_TRANSITION_RAINBOW_URL } from '../../config/const';
import { UpdateService } from '../../update/update.service';
import { HttpClient } from '@angular/common/http';
import { deepCopy } from '../../lib/deep-copy';
import { deepEqual } from 'fast-equals';

interface ModeTransitionRainbowParameter {
  roundTimeMs: number,
  brightness: number,
  spectrum: number,
  reverse: boolean,
}
interface ModeTransitionRainbowLimits {
  minRoundTimeMs: number,
  maxRoundTimeMs: number,
  minBrightness: number,
  maxBrightness: number,
}

@Component({
  selector: 'app-mode-transition-rainbow',
  templateUrl: './mode-transition-rainbow.component.html',
  styleUrls: ['./mode-transition-rainbow.component.scss']
})
export class ModeTransitionRainbowComponent implements OnInit {
  public backModeTransitionRainbowParameter: ModeTransitionRainbowParameter
  public modeTransitionRainbowParameter: ModeTransitionRainbowParameter
  public modeTransitionRainbowLimits: ModeTransitionRainbowLimits

  constructor(private httpClient: HttpClient, private updateService: UpdateService) {
  }

  ngOnInit(): void {
    this.updateModeTransitionRainbowParameter();
    this.updateModeTransitionRainbowLimits();
    this.updateService.registerPolling({ cb: () => { this.updateModeTransitionRainbowParameter() }, timeout: 500 })
    setTimeout( ()=>{M.updateTextFields()},100);
  }

  updateModeTransitionRainbowParameter() {
    this.httpClient.get<ModeTransitionRainbowParameter>(REST_MODE_TRANSITION_RAINBOW_URL).subscribe((data: ModeTransitionRainbowParameter) => {
      if (!deepEqual(this.backModeTransitionRainbowParameter, data)) {
        this.backModeTransitionRainbowParameter = data
        this.modeTransitionRainbowParameter = deepCopy(this.backModeTransitionRainbowParameter)
        console.log(data)
      }
    }
    )
  }

  setModeTransitionRainbowParameter() {
    console.log("set")
    this.httpClient.post<ModeTransitionRainbowParameter>(REST_MODE_TRANSITION_RAINBOW_URL, this.modeTransitionRainbowParameter, {}).subscribe()
  }

  updateModeTransitionRainbowLimits() {
    this.httpClient.get<ModeTransitionRainbowLimits>(REST_MODE_TRANSITION_RAINBOW_URL + "/limits").subscribe((data: ModeTransitionRainbowLimits) => { this.modeTransitionRainbowLimits = data })
  }

}
