import { Injectable } from '@angular/core';
import { WebsocketService } from 'src/app/websocket/websocket.service';
import { SuperMode } from 'src/app/modes/super-mode';



export interface ModeSpectrumParameter {
  hueFrom720: number
  hueTo720: number
  brightness: number

}
export interface ModeSpectrumLimits {
  minBrightness: number
  maxBrightness: number
}

@Injectable({
  providedIn: 'root'
})
export class ModeSpectrumService extends SuperMode {
  public backParameter: ModeSpectrumParameter
  public parameter: ModeSpectrumParameter
  public limits: ModeSpectrumLimits
  public hueRange: number[] = [0, 0]

  constructor(protected websocketService: WebsocketService) {
    super("ModeSpectrum", websocketService)
  }

  receiveParameter(parm: ModeSpectrumParameter) {
    super.receiveParameter(parm)
    this.hueRange = [this.parameter.hueFrom720, this.parameter.hueTo720]
    console.log(parm)
  }

  setHue() {
    this.parameter.hueFrom720 = this.hueRange[0]
    this.parameter.hueTo720 = this.hueRange[1]
  }
}
