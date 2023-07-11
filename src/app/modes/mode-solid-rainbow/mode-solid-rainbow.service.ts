import { Injectable } from '@angular/core';
import { WebsocketService } from 'src/app/websocket/websocket.service';
import { ParentMode } from 'src/app/modes/parent-mode';


export interface ModeSolidRainbowParameter {
  roundTimeMs: number,
  brightness: number,
}
export interface ModeSolidRainbowLimits {
  minRoundTimeMs: number,
  maxRoundTimeMs: number,
  minBrightness: number,
  maxBrightness: number,
}

@Injectable({
  providedIn: 'root'
})
export class ModeSolidRainbowService extends ParentMode {

  public backParameter: ModeSolidRainbowParameter
  public parameter: ModeSolidRainbowParameter
  public limits: ModeSolidRainbowLimits

  constructor(protected websocketService: WebsocketService) {
    super("ModeSolidRainbow", websocketService)
  }
}
