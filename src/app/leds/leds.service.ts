import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { RGB } from '../color/color';
import { REST_GET_LEDS_URL } from '../config/const';
import { UpdateService, UpdateIntervall } from '../update/update.service';


@Injectable({
  providedIn: 'root'
})
export class LedsService {

  leds: Array<RGB>

  constructor(private httpClient: HttpClient, private updateService: UpdateService) {
    updateService.registerPolling({ cb: ()=>{this.updateLeds()}, timeout: 100 })
  }

  public updateLeds() {
    this.httpClient.get<Array<RGB>>(REST_GET_LEDS_URL).subscribe((data: Array<RGB>) => this.leds = data);
  }

}
