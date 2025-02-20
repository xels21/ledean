import { Injectable } from '@angular/core';
import { WebsocketService } from 'src/app/websocket/websocket.service';
import { SuperMode } from 'src/app/modes/super-mode';


export interface ModeSpectrumParameterPosition {
  facFrom: number
  facTo: number
  facRoundTimeMs: number
  offFrom: number
  offTo: number
  offRoundTimeMs: number
}

export interface ModeSpectrumParameter {
  hueFrom720: number
  hueTo720: number
  brightness: number
  positions: Array<ModeSpectrumParameterPosition>

}
export interface ModeSpectrumLimits {
  maxRoundTimeMs: number
  minRoundTimeMs: number
  minBrightness: number
  maxBrightness: number
  minFactor: number
  maxFactor: number
  minOffset: number
  maxOffset: number
}

export interface PosRange{
  factor: Array<number>
  offset: Array<number>
}

@Injectable({
  providedIn: 'root'
})
export class ModeSpectrumService extends SuperMode {
  public backParameter: ModeSpectrumParameter
  public parameter: ModeSpectrumParameter
  public limits: ModeSpectrumLimits
  public hueRange: number[] = [0, 0]
  // public posRange: Array<PosRange> = new Array<PosRange>(2)
  public posRange: Array<PosRange> = [{factor:[0,0],offset:[0,0]},{factor:[0,0],offset:[0,0]}]

  constructor(protected websocketService: WebsocketService) {
    super("ModeSpectrum", websocketService)
  }

  receiveParameter(parm: ModeSpectrumParameter) {
    super.receiveParameter(parm)
    this.hueRange = [this.parameter.hueFrom720, this.parameter.hueTo720]
    for(let i in this.parameter.positions){
      this.posRange[i].factor = [this.parameter.positions[i].facFrom, this.parameter.positions[i].facTo]
      this.posRange[i].offset = [this.parameter.positions[i].offFrom, this.parameter.positions[i].offTo]
    }
    console.log(parm)
  }

  setHue() {
    this.parameter.hueFrom720 = this.hueRange[0]
    this.parameter.hueTo720 = this.hueRange[1]
  }


  setPosRangeFactor(i) {
    this.parameter.positions[i].facFrom = this.posRange[i].factor[0]
    this.parameter.positions[i].facTo = this.posRange[i].factor[1]
  }

  setPosRangeOffset(i) {
    this.parameter.positions[i].offFrom = this.posRange[i].offset[0]
    this.parameter.positions[i].offTo = this.posRange[i].offset[1]
  }
}
