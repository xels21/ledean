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
  public pollingTimeout:number

  constructor(private httpClient: HttpClient, private updateService: UpdateService) {
    this.pollingTimeout=100
    // this.pollingTimeout=300
    updateService.registerPolling({ cb: ()=>{this.updateLeds()}, timeout: this.pollingTimeout })
  }

  public updateLeds() {
    this.httpClient.get<Array<RGB>>(REST_GET_LEDS_URL).subscribe((data: Array<RGB>) => {
      if(this.leds==undefined){
        this.leds = data
      }else{
        for(var i=0;i<this.leds.length;i++){
          this.leds[i].r=data[i].r
          this.leds[i].g=data[i].g
          this.leds[i].b=data[i].b
        }
      }
    });
  }

}
