import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { REST_GET_MODE_URL, REST_GET_MODE_RESOLVER_URL } from '../config/const';
import { UpdateService, UpdateIntervall } from '../update/update.service';
import { ModeSolidService } from './mode-solid/mode-solid.service';
import { WebsocketService } from '../websocket/websocket.service';
import { CmdMode, CmdModeLimits } from '../websocket/commands';


@Injectable({
  providedIn: 'root'
})
export class ModesService {

  activeMode: number
  isPictureMode: boolean
  public modeResolver: Array<string>
  onModeChange: () => any

  constructor(private httpClient: HttpClient, private updateService: UpdateService, private websocketService: WebsocketService, private modeSolidService: ModeSolidService) {
    // updateService.registerPolling({ cb: () => { this.updateActiveMode() }, timeout: 1000 })
    this.httpClient.get<Array<string>>(REST_GET_MODE_RESOLVER_URL).subscribe((data: Array<string>) => this.modeResolver = data);
    this.setOnModeChange(() => { })
    this.websocketService.modeChanged.subscribe((cmdMode: CmdMode) => this.modeChanged(cmdMode))
    this.websocketService.modeLimitChanged.subscribe((cmdModeLimits: CmdModeLimits) => this.modeLimitChanged(cmdModeLimits))
  }

  private modeChanged(cmdMode: CmdMode) {
    console.log(cmdMode)
    switch (cmdMode.id) {
      case this.modeSolidService.getName():
        console.log("_____________")
        this.modeSolidService.updateModeSolidParameter(cmdMode.parm)
        break;
      default:
        console.log("Unknown limit change for mode: ", cmdMode.id)
        break;
    }
  }

  private modeLimitChanged(cmdModeLimits: CmdModeLimits) {
    switch (cmdModeLimits.id) {
      case this.modeSolidService.getName():
        this.modeSolidService.updateModeSolidLimits(cmdModeLimits.limits)
        break;

      default:
        console.log("Unknown limit change for mode: ", cmdModeLimits.id)
        break;
    }
    console.log("needs too be handled: ", cmdModeLimits)
  }


  public setOnModeChange(onModeChange: () => any) {
    this.onModeChange = onModeChange
  }

  public updateActiveMode(id: number) {
    // this.httpClient.get<number>(REST_GET_MODE_URL).subscribe((data: number) => {
    if (this.activeMode != id) {
      this.activeMode = id
      if (this.activeMode >= 0) {
        this.isPictureMode = false
        this.onModeChange()
      } else {
        this.isPictureMode = true
      }
    }
    // });
  }

  public isActive(mode: string) {
    return mode == this.modeResolver[this.activeMode]
  }

  public switchState(mode: string) {
    this.httpClient.get(REST_GET_MODE_URL + "/" + this.modeStrToIdx(mode)).subscribe();
  }

  public modeStrToIdx(mode: string) {
    return this.modeResolver.findIndex(m => { return m == mode })
  }



}
