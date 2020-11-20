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

@Component({
  selector: 'app-mode-solid',
  templateUrl: './mode-solid.component.html',
  styleUrls: ['./mode-solid.component.scss']
})
export class ModeSolidComponent implements OnInit {
  public modeSolidParameter: ModeSolidParameter
  public uiModeSolidParameter: ModeSolidParameter

  constructor(private httpClient: HttpClient, private updateService: UpdateService) {
    this.modeSolidParameter = {
      rgb: { r: 0, g: 0, b: 0 },
      brightness: 1.0
    }
    this.uiModeSolidParameter = this.modeSolidParameter
  }

  ngOnInit(): void {
    this.updateModeSolidParameter();
    this.updateService.registerPolling({ cb: () => { this.updateModeSolidParameter() }, timeout: 500 })
  }

  updateModeSolidParameter() {
    this.httpClient.get<ModeSolidParameter>(REST_MODE_SOLID_URL).subscribe((data: ModeSolidParameter) => {
      if (!deepEqual(this.modeSolidParameter, data)) {
        this.modeSolidParameter = data
        this.uiModeSolidParameter = deepCopy(this.modeSolidParameter)
      }
    }
    )
  }

  setModeSolidParameter() {
    console.log("set")
    this.httpClient.post<ModeSolidParameter>(REST_MODE_SOLID_URL, this.uiModeSolidParameter, {}).subscribe()
  }

}
