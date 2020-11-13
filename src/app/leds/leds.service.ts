import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { RGB } from '../color/color';
import { REST_GET_LEDS_URL } from '../config/const';


@Injectable({
  providedIn: 'root'
})
export class LedsService {

  leds: Array<RGB>
  pollMs: number

  public pollingActive: boolean
  pollingInterval: any

  constructor(private httpClient: HttpClient) {
    this.pollMs = 100;
    this.pollingActive = true
    this.checkPollingInterval()
  }

  public checkPollingInterval() {
    if (this.pollingActive) {
      if (this.pollingInterval == null) {
        this.pollingInterval = setInterval(() => this.updateLeds(), this.pollMs)
      }
    } else {
      if (this.pollingInterval != null) {
        clearInterval(this.pollingInterval);
        this.pollingInterval = null
      }
    }
  }

  public updateLeds() {
    this.httpClient.get<Array<RGB>>(REST_GET_LEDS_URL).subscribe((data: Array<RGB>) => this.leds = data);
  }

}
