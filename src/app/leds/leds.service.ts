import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { RGB } from '../color/color';
import { REST_GET_LEDS_URL, REST_GET_LEDS_COUNT_URL, REST_GET_LEDS_ROWS_URL } from '../config/const';
import { UpdateService, UpdateIntervall } from '../update/update.service';


@Injectable({
  providedIn: 'root'
})
export class LedsService {

  leds: Array<RGB>
  ledCount: number
  ledRows: number
  public pollingTimeout: number

  constructor(private httpClient: HttpClient, private updateService: UpdateService) {
    this.pollingTimeout = 200
    // this.pollingTimeout=300
    this.updateLedCount()
    this.updateLedRows()
    updateService.registerPolling({ cb: () => { this.updateLeds() }, timeout: this.pollingTimeout })
  }

  public updateLedCount() {
    this.httpClient.get<number>(REST_GET_LEDS_COUNT_URL).subscribe((data: number) => {
      this.ledCount = data
    })
  }
  public updateLedRows() {
    this.httpClient.get<number>(REST_GET_LEDS_ROWS_URL).subscribe((data: number) => {
      this.ledRows = data
    })
  }

  public updateLeds() {
    this.httpClient.get<Array<RGB>>(REST_GET_LEDS_URL).subscribe((data: Array<RGB>) => {
      if (this.leds == undefined) {
        this.leds = data
      } else {
        for (var i = 0; i < this.leds.length; i++) {
          this.leds[i].r = data[i].r
          this.leds[i].g = data[i].g
          this.leds[i].b = data[i].b
        }
      }
    });
  }

}
