import { Injectable } from '@angular/core';
import { WebsocketService } from 'src/app/websocket/websocket.service';
import { SuperMode } from 'src/app/modes/super-mode';



export interface ModeGradientParameter {
  count: number
  roundTimeMs: number
  brightness: number

}
export interface ModeGradientLimits {
  minCount: number
  maxCount: number
  minRoundTimeMs: number
  maxRoundTimeMs: number
  minBrightness: number
  maxBrightness: number
}

@Injectable({
  providedIn: 'root'
})
export class ModeGradientService extends SuperMode{
  public backParameter: ModeGradientParameter
  public parameter: ModeGradientParameter
  public limits: ModeGradientLimits

  constructor(protected websocketService: WebsocketService) {
    super("ModeGradient", websocketService)
   }
}
