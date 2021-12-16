
import { Component, OnInit } from '@angular/core';
import { RGB } from '../../color/color';
import { REST_MODE_SOLID_RAINBOW_URL } from '../../config/const';
import { UpdateService } from '../../update/update.service';
import { HttpClient } from '@angular/common/http';
import { deepCopy } from '../../lib/deep-copy';
import { deepEqual } from 'fast-equals';

interface ModeSolidRainbowParameter {
  roundTimeMs: number,
  brightness: number,
}
interface ModeSolidRainbowLimits {
  minRoundTimeMs: number,
  maxRoundTimeMs: number,
  minBrightness: number,
  maxBrightness: number,
}

@Component({
  selector: 'app-mode-solid-rainbow',
  templateUrl: './mode-solid-rainbow.component.html',
  styleUrls: ['./mode-solid-rainbow.component.scss','../../app.component.scss']
})
export class ModeSolidRainbowComponent implements OnInit {
  public backModeSolidRainbowParameter: ModeSolidRainbowParameter
  public modeSolidRainbowParameter: ModeSolidRainbowParameter
  public modeSolidRainbowLimits: ModeSolidRainbowLimits

  constructor(private httpClient: HttpClient, private updateService: UpdateService) {
  }

  ngOnInit(): void {
    this.updateModeSolidRainbowParameter();
    this.updateModeSolidRainbowLimits();
    this.updateService.registerPolling({ cb: () => { this.updateModeSolidRainbowParameter() }, timeout: 500 })
    setTimeout( ()=>{M.updateTextFields()},100);
  }

  updateModeSolidRainbowParameter() {
    this.httpClient.get<ModeSolidRainbowParameter>(REST_MODE_SOLID_RAINBOW_URL).subscribe((data: ModeSolidRainbowParameter) => {
      if (!deepEqual(this.backModeSolidRainbowParameter, data)) {
        this.backModeSolidRainbowParameter = data
        this.modeSolidRainbowParameter = deepCopy(this.backModeSolidRainbowParameter)
        console.log(data)
      }
    }
    )
  }

  setModeSolidRainbowParameter() {
    console.log("set")
    this.httpClient.post<ModeSolidRainbowParameter>(REST_MODE_SOLID_RAINBOW_URL, this.modeSolidRainbowParameter, {}).subscribe()
  }

  updateModeSolidRainbowLimits() {
    this.httpClient.get<ModeSolidRainbowLimits>(REST_MODE_SOLID_RAINBOW_URL + "/limits").subscribe((data: ModeSolidRainbowLimits) => { this.modeSolidRainbowLimits = data })
  }

}
