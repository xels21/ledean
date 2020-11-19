import { Component, OnInit } from '@angular/core';
import { RGB } from '../../color/color';
import { REST_MODE_SOLID_URL } from '../../config/const';
import { UpdateService } from '../../update/update.service';
import { HttpClient } from '@angular/common/http';


interface ModeSolidParameter {
  rgb: RGB,
}

@Component({
  selector: 'app-mode-solid',
  templateUrl: './mode-solid.component.html',
  styleUrls: ['./mode-solid.component.scss']
})
export class ModeSolidComponent implements OnInit {
  public modeSolidParameter: ModeSolidParameter

  constructor(private httpClient: HttpClient, private updateService:UpdateService) {
    this.modeSolidParameter = {
      rgb: { r: 0, g: 0, b: 0 }
    }
  }

  ngOnInit(): void {
    console.log("this")
    this.updateModeSolidParameter();
    this.updateService.registerPolling({ cb: ()=>{this.updateModeSolidParameter()}, timeout: 500 })
  }

  updateModeSolidParameter() {
    this.httpClient.get<ModeSolidParameter>(REST_MODE_SOLID_URL).subscribe((data: ModeSolidParameter) => {
      this.modeSolidParameter = data
    });
  }

  setModeSolidParameter() {
    this.httpClient.post<ModeSolidParameter>(REST_MODE_SOLID_URL,this.modeSolidParameter,{}).subscribe()
  }

}
