import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { RGB } from '../color/color';


@Injectable({
  providedIn: 'root'
})
export class LedsService {

  leds: Array<RGB>
  pollMs: number
  addr: string
  port: number
  path: string
  url: string

  public pollingActive: boolean
  pollingInterval: any

  constructor(private httpClient: HttpClient) {
    this.pollMs = 1000;
    this.addr = window.location.hostname//"localhost"
    // this.addr = "localhost"
    this.port = 2211
    this.path = "leds"
    this.url = "http://" + this.addr + ":" + this.port + "/" + this.path
    this.pollingActive = true
    this.checkPollingInterval()
  }

  public checkPollingInterval() {
    if (this.pollingActive) {
      if (this.pollingInterval == null) {
        this.pollingInterval = setInterval(() => this.updateLeds(), 1000)
      }
    } else {
      if (this.pollingInterval != null) {
        clearInterval(this.pollingInterval);
        this.pollingInterval = null
      }
    }
  }

  public updateLeds() {
    const headers = new HttpHeaders()
    this.httpClient.get<Array<RGB>>(this.url, { headers }).subscribe((data: Array<RGB>) => this.leds = data);
  }

}
