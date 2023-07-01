import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { RGB } from '../color/color';
import { REST_GET_LEDS_URL, REST_GET_LEDS_COUNT_URL, REST_GET_LEDS_ROWS_URL } from '../config/const';
import { UpdateService, UpdateIntervall } from '../update/update.service';
import { ModesService } from '../modes/modes.service';
import { WebsocketService } from '../websocket/websocket.service';


@Injectable({
  providedIn: 'root'
})
export class LedsService {

  leds: Array<RGB>
  bufferedLeds: Array<Array<RGB>>
  bufferedLedsCount: number
  ledCount: number
  ledRows: number
  public pollingTimeout: number

  constructor(private httpClient: HttpClient, private updateService: UpdateService, public modesService: ModesService, private websocketService: WebsocketService) {
    this.bufferedLedsCount = 16
    this.pollingTimeout = 200
    // this.pollingTimeout=300
    // this.updateLedCount()
    // this.updateLedRows()
    websocketService.ledsCountChanged.subscribe(count => this.updateLedCount(count))
    websocketService.ledsRowsChanged.subscribe(rows => this.updateLedRows(rows))
    websocketService.ledsChanged.subscribe(leds => this.updateLeds(leds))
    // updateService.registerPolling({ cb: () => { this.updateLeds() }, timeout: this.pollingTimeout })
  }

  public updateLedCount(count:number) {
    // this.httpClient.get<number>(REST_GET_LEDS_COUNT_URL).subscribe((data: number) => {
      this.ledCount = count

      this.leds = new Array<RGB>(this.ledCount)
      for (let i = 0; i < this.ledCount; i++) {
        this.leds[i] = { r: 0, g: 0, b: 0 }
      }

      this.bufferedLeds = new Array<Array<RGB>>(this.ledCount)
      for (let i = 0; i < this.ledCount; i++) {
        this.bufferedLeds[i] = new Array<RGB>(this.bufferedLedsCount)
        for (let b = 0; b < this.bufferedLedsCount; b++) {
          this.bufferedLeds[i][b] = { r: 0, g: 0, b: 0 }
        }
      }
    // })
  }
  public updateLedRows(rows:number) {
    // this.httpClient.get<number>(REST_GET_LEDS_ROWS_URL).subscribe((data: number) => {
      this.ledRows = rows
    // })
  }

  public updateLeds(leds:Array<RGB>) {
    // this.httpClient.get<Array<RGB>>(REST_GET_LEDS_URL).subscribe((leds: Array<RGB>) => {
      if (this.ledCount == undefined) {
        console.log("led count was not set before")
        return
        // this.leds = leds
      }
      if (!this.modesService.isPictureMode) {
        for (var i = 0; i < this.leds.length; i++) {
          this.leds[i].r = leds[i].r
          this.leds[i].g = leds[i].g
          this.leds[i].b = leds[i].b
        }
      } else {
        for (var i = 0; i < this.ledCount; i++) {
          for (let b = this.bufferedLedsCount - 1; b > 0; b--) {
            this.bufferedLeds[i][b].r = this.bufferedLeds[i][b - 1].r
            this.bufferedLeds[i][b].g = this.bufferedLeds[i][b - 1].g
            this.bufferedLeds[i][b].b = this.bufferedLeds[i][b - 1].b
          }
        }
        for (var i = 0; i < this.ledCount; i++) {
          this.bufferedLeds[i][0].r = leds[i].r
          this.bufferedLeds[i][0].g = leds[i].g
          this.bufferedLeds[i][0].b = leds[i].b
        }
      }
    // });
  }

}
