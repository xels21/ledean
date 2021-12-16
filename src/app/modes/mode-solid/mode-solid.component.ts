import { Component, OnInit } from '@angular/core';
import { RGB } from '../../color/color';
import { REST_MODE_SOLID_URL } from '../../config/const';
import { UpdateService } from '../../update/update.service';
import { HttpClient } from '@angular/common/http';
import { deepCopy } from '../../lib/deep-copy';
import { deepEqual } from 'fast-equals';


interface ModeSolidParameter {
  rgb: RGB,
  brightness: number,
}
interface ModeSolidLimits {
  minBrightness: number,
  maxBrightness: number,
}

@Component({
  selector: 'app-mode-solid',
  templateUrl: './mode-solid.component.html',
  styleUrls: ['./mode-solid.component.scss','../../app.component.scss']
})
export class ModeSolidComponent implements OnInit {
  public backModeSolidParameter: ModeSolidParameter
  public modeSolidParameter: ModeSolidParameter
  public modeSolidLimits: ModeSolidLimits


  constructor(private httpClient: HttpClient, private updateService: UpdateService) {
  }

  ngOnInit(): void {
    this.updateModeSolidParameter();
    this.updateModeSolidLimits();
    this.updateService.registerPolling({ cb: () => { this.updateModeSolidParameter() }, timeout: 500 })
  }

  updateModeSolidParameter() {
    this.httpClient.get<ModeSolidParameter>(REST_MODE_SOLID_URL).subscribe((data: ModeSolidParameter) => {
      if (!deepEqual(this.backModeSolidParameter, data)) {
        this.backModeSolidParameter = data
        this.modeSolidParameter = deepCopy(this.backModeSolidParameter)
      }
    }
    )
  }

  setModeSolidParameter() {
    console.log("set")
    this.httpClient.post<ModeSolidParameter>(REST_MODE_SOLID_URL, this.modeSolidParameter, {}).subscribe()
  }

  updateModeSolidLimits() {
    this.httpClient.get<ModeSolidLimits>(REST_MODE_SOLID_URL + "/limits").subscribe((data: ModeSolidLimits) => { this.modeSolidLimits = data })
  }

}
