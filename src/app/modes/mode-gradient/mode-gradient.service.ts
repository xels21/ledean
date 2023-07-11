import { Injectable } from '@angular/core';
import { WebsocketService } from 'src/app/websocket/websocket.service';
import { ParentMode } from 'src/app/modes/parent-mode';



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
export class ModeGradientService extends ParentMode{
  public backParameter: ModeGradientParameter
  public parameter: ModeGradientParameter
  public limits: ModeGradientLimits

  constructor(protected websocketService: WebsocketService) {
    super("ModeGradient", websocketService)
   }
}
