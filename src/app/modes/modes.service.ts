import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { REST_GET_MODE_URL, REST_GET_MODE_RESOLVER_URL } from '../config/const';
import { UpdateService, UpdateIntervall } from '../update/update.service';


@Injectable({
  providedIn: 'root'
})
export class ModesService {

  activeMode: number
  public modeResolver: Array<string>
  onModeChange: () => any

  constructor(private httpClient: HttpClient, private updateService: UpdateService) {
    updateService.registerPolling({ cb: () => { this.updateActiveMode() }, timeout: 1000 })
    this.httpClient.get<Array<string>>(REST_GET_MODE_RESOLVER_URL).subscribe((data: Array<string>) => this.modeResolver = data);
  }


  public setOnModeChange(onModeChange: () => any) {
    this.onModeChange = onModeChange
  }

  public updateActiveMode() {
    this.httpClient.get<number>(REST_GET_MODE_URL).subscribe((data: number) => { this.activeMode = data });
  }

  public isActive(mode: string) {
    return mode == this.modeResolver[this.activeMode]
  }

}
